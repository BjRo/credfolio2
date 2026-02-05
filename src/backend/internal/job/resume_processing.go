package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"backend/internal/domain"
	"backend/internal/logger"
)

// ResumeProcessingArgs contains the arguments for a resume processing job.
type ResumeProcessingArgs struct { //nolint:govet // Field ordering prioritizes JSON consistency over memory alignment
	StorageKey  string    `json:"storage_key"`
	ResumeID    uuid.UUID `json:"resume_id"`
	FileID      uuid.UUID `json:"file_id"`
	ContentType string    `json:"content_type"` // Added to avoid extra DB lookup
}

// Kind returns the job type identifier for River.
func (ResumeProcessingArgs) Kind() string {
	return "resume_processing"
}

// InsertOpts returns default insert options — limits retries to avoid
// repeatedly hammering LLM providers on persistent failures.
func (ResumeProcessingArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 2}
}

// ResumeProcessingWorker processes uploaded resumes to extract profile data.
type ResumeProcessingWorker struct {
	river.WorkerDefaults[ResumeProcessingArgs]
	resumeRepo     domain.ResumeRepository
	fileRepo       domain.FileRepository
	profileRepo    domain.ProfileRepository
	profileExpRepo   domain.ProfileExperienceRepository
	profileEduRepo   domain.ProfileEducationRepository
	profileSkillRepo domain.ProfileSkillRepository
	storage          domain.Storage
	extractor      domain.DocumentExtractor
	log            logger.Logger
}

// NewResumeProcessingWorker creates a new resume processing worker.
func NewResumeProcessingWorker(
	resumeRepo domain.ResumeRepository,
	fileRepo domain.FileRepository,
	profileRepo domain.ProfileRepository,
	profileExpRepo domain.ProfileExperienceRepository,
	profileEduRepo domain.ProfileEducationRepository,
	profileSkillRepo domain.ProfileSkillRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *ResumeProcessingWorker {
	return &ResumeProcessingWorker{
		resumeRepo:       resumeRepo,
		fileRepo:         fileRepo,
		profileRepo:      profileRepo,
		profileExpRepo:   profileExpRepo,
		profileEduRepo:   profileEduRepo,
		profileSkillRepo: profileSkillRepo,
		storage:          storage,
		extractor:        extractor,
		log:              log,
	}
}

// Timeout disables River's default 60s job timeout. LLM extraction can take
// several minutes for structured output; timeout is handled by the resilient
// provider layer instead (300s).
func (w *ResumeProcessingWorker) Timeout(*river.Job[ResumeProcessingArgs]) time.Duration {
	return -1 // no River-level timeout
}

// Work processes a resume and extracts profile data using LLM.
func (w *ResumeProcessingWorker) Work(ctx context.Context, job *river.Job[ResumeProcessingArgs]) error {
	args := job.Args
	w.log.Info("Processing resume",
		logger.Feature("jobs"),
		logger.String("resume_id", args.ResumeID.String()),
		logger.String("file_id", args.FileID.String()),
		logger.String("storage_key", args.StorageKey),
	)

	// Update status to processing
	if err := w.updateStatus(ctx, args.ResumeID, domain.ResumeStatusProcessing, nil); err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Get content type from file record (in case it wasn't passed in args)
	contentType := args.ContentType
	if contentType == "" {
		file, err := w.fileRepo.GetByID(ctx, args.FileID)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get file record: %v", err)
			w.log.Error("Failed to get file record",
				logger.Feature("jobs"),
				logger.String("resume_id", args.ResumeID.String()),
				logger.String("file_id", args.FileID.String()),
				logger.Err(err),
			)
			w.updateStatusFailed(ctx, args.ResumeID, errMsg)
			return fmt.Errorf("failed to get file record: %w", err)
		}
		if file == nil {
			errMsg := "file record not found"
			w.updateStatusFailed(ctx, args.ResumeID, errMsg)
			return fmt.Errorf("file record not found: %s", args.FileID)
		}
		contentType = file.ContentType
	}

	// Download file from storage
	reader, err := w.storage.Download(ctx, args.StorageKey)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download file: %v", err)
		w.log.Error("Storage download failed",
			logger.Feature("jobs"),
			logger.String("resume_id", args.ResumeID.String()),
			logger.String("storage_key", args.StorageKey),
			logger.Err(err),
		)
		w.updateStatusFailed(ctx, args.ResumeID, errMsg)
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer reader.Close() //nolint:errcheck // Best effort cleanup

	// Read file data
	data, err := io.ReadAll(reader)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read file data: %v", err)
		w.log.Error("Failed to read file data",
			logger.Feature("jobs"),
			logger.String("resume_id", args.ResumeID.String()),
			logger.Err(err),
		)
		w.updateStatusFailed(ctx, args.ResumeID, errMsg)
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Extract profile data using LLM
	extractedData, err := w.extractResumeData(ctx, data, contentType)
	if err != nil {
		errMsg := fmt.Sprintf("failed to extract resume data: %v", err)
		w.log.Error("Resume extraction failed",
			logger.Feature("jobs"),
			logger.String("resume_id", args.ResumeID.String()),
			logger.Err(err),
		)
		w.updateStatusFailed(ctx, args.ResumeID, errMsg)
		return fmt.Errorf("failed to extract resume data: %w", err)
	}

	// Save extracted data
	if saveErr := w.saveExtractedData(ctx, args.ResumeID, extractedData); saveErr != nil {
		errMsg := fmt.Sprintf("failed to save extracted data: %v", saveErr)
		w.log.Error("Failed to save extracted data",
			logger.Feature("jobs"),
			logger.String("resume_id", args.ResumeID.String()),
			logger.Err(saveErr),
		)
		w.updateStatusFailed(ctx, args.ResumeID, errMsg)
		return fmt.Errorf("failed to save extracted data: %w", saveErr)
	}

	// Materialize extracted data into profile tables
	resume, err := w.resumeRepo.GetByID(ctx, args.ResumeID)
	if err != nil {
		w.log.Error("Failed to get resume for materialization",
			logger.Feature("jobs"),
			logger.String("resume_id", args.ResumeID.String()),
			logger.Err(err),
		)
	} else if resume != nil {
		if matErr := w.materializeExtractedData(ctx, args.ResumeID, resume.UserID, extractedData); matErr != nil {
			w.log.Error("Failed to materialize extracted data into profile",
				logger.Feature("jobs"),
				logger.String("resume_id", args.ResumeID.String()),
				logger.Err(matErr),
			)
			// Log but don't fail — extraction data is still saved in JSONB
		} else {
			w.log.Info("Materialized extracted data into profile tables",
				logger.Feature("jobs"),
				logger.String("resume_id", args.ResumeID.String()),
				logger.Int("experience_count", len(extractedData.Experience)),
				logger.Int("education_count", len(extractedData.Education)),
				logger.Int("skills_count", len(extractedData.Skills)),
			)
		}
	}

	// Mark as completed AFTER materialization so the frontend doesn't
	// redirect to the profile page before profile data is ready.
	if err := w.updateStatus(ctx, args.ResumeID, domain.ResumeStatusCompleted, nil); err != nil {
		return fmt.Errorf("failed to update status to completed: %w", err)
	}

	w.log.Info("Resume processing completed",
		logger.Feature("jobs"),
		logger.String("resume_id", args.ResumeID.String()),
		logger.String("name", extractedData.Name),
		logger.Int("skills_count", len(extractedData.Skills)),
		logger.Int("experience_count", len(extractedData.Experience)),
	)

	return nil
}

// extractResumeData uses the LLM to extract structured data from the resume.
func (w *ResumeProcessingWorker) extractResumeData(ctx context.Context, data []byte, contentType string) (*domain.ResumeExtractedData, error) {
	ctx, span := otel.Tracer("credfolio").Start(ctx, "resume_extraction")
	defer span.End()

	// First, extract text from the document
	text, err := w.extractor.ExtractText(ctx, data, contentType)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Then, use LLM to extract structured profile data from the text
	extractedData, err := w.extractor.ExtractResumeData(ctx, text)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to extract resume data: %w", err)
	}

	// Set extraction timestamp
	extractedData.ExtractedAt = time.Now()

	return extractedData, nil
}

// saveExtractedData saves the extracted data JSONB to the resume record without changing status.
// Status is updated separately after materialization to avoid a race condition where
// the frontend sees COMPLETED before profile data is materialized.
func (w *ResumeProcessingWorker) saveExtractedData(ctx context.Context, resumeID uuid.UUID, data *domain.ResumeExtractedData) error {
	resume, err := w.resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		return fmt.Errorf("failed to get resume: %w", err)
	}
	if resume == nil {
		return fmt.Errorf("resume not found: %s", resumeID)
	}

	// Marshal extracted data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal extracted data: %w", err)
	}

	// Save extracted data but keep status as processing — status is set to
	// completed only after materialization finishes so the frontend won't
	// redirect before profile tables are populated.
	resume.ExtractedData = jsonData
	resume.ErrorMessage = nil

	if err := w.resumeRepo.Update(ctx, resume); err != nil {
		return fmt.Errorf("failed to update resume: %w", err)
	}

	return nil
}

// updateStatus updates the resume status.
func (w *ResumeProcessingWorker) updateStatus(ctx context.Context, id uuid.UUID, status domain.ResumeStatus, errMsg *string) error {
	resume, err := w.resumeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get resume: %w", err)
	}
	if resume == nil {
		return fmt.Errorf("resume not found: %s", id)
	}

	resume.Status = status
	resume.ErrorMessage = errMsg

	if err := w.resumeRepo.Update(ctx, resume); err != nil {
		return fmt.Errorf("failed to update resume: %w", err)
	}

	return nil
}

// updateStatusFailed updates resume status to failed and logs any DB errors.
// This is best-effort — if the DB update fails, we log it but don't propagate
// since the caller is already returning the original error.
func (w *ResumeProcessingWorker) updateStatusFailed(ctx context.Context, id uuid.UUID, errMsg string) {
	if err := w.updateStatus(ctx, id, domain.ResumeStatusFailed, &errMsg); err != nil {
		w.log.Error("Failed to update resume status to failed",
			logger.Feature("jobs"),
			logger.String("resume_id", id.String()),
			logger.Err(err),
		)
	}
}

// materializeExtractedData creates profile education, experience, and skill rows from extracted resume data.
// This makes profile tables the single source of truth for display.
// Each category (experiences, education, skills) is processed independently so that failures in one
// category don't prevent the others from being saved.
func (w *ResumeProcessingWorker) materializeExtractedData(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID, data *domain.ResumeExtractedData) error {
	// Get or create the user's profile
	profile, err := w.profileRepo.GetOrCreateByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get or create profile: %w", err)
	}

	// Delete any existing entries from this resume (idempotent re-processing)
	if delErr := w.profileExpRepo.DeleteBySourceResumeID(ctx, resumeID); delErr != nil {
		return fmt.Errorf("failed to delete existing experiences for resume: %w", delErr)
	}
	if delErr := w.profileEduRepo.DeleteBySourceResumeID(ctx, resumeID); delErr != nil {
		return fmt.Errorf("failed to delete existing education for resume: %w", delErr)
	}
	if delErr := w.profileSkillRepo.DeleteBySourceResumeID(ctx, resumeID); delErr != nil {
		return fmt.Errorf("failed to delete existing skills for resume: %w", delErr)
	}

	// Process each category independently to ensure partial success
	var errors []error

	if expErr := w.materializeExperiences(ctx, resumeID, profile.ID, data.Experience); expErr != nil {
		errors = append(errors, fmt.Errorf("experiences: %w", expErr))
	}
	if eduErr := w.materializeEducation(ctx, resumeID, profile.ID, data.Education); eduErr != nil {
		errors = append(errors, fmt.Errorf("education: %w", eduErr))
	}
	if skillErr := w.materializeSkills(ctx, resumeID, profile.ID, data.Skills); skillErr != nil {
		errors = append(errors, fmt.Errorf("skills: %w", skillErr))
	}

	// Return aggregated error if any category failed
	if len(errors) > 0 {
		return fmt.Errorf("materialization errors: %v", errors)
	}

	return nil
}

func (w *ResumeProcessingWorker) materializeExperiences(ctx context.Context, resumeID, profileID uuid.UUID, experiences []domain.WorkExperience) error {
	displayOrder, err := w.profileExpRepo.GetNextDisplayOrder(ctx, profileID)
	if err != nil {
		return fmt.Errorf("failed to get next experience display order: %w", err)
	}

	for i, exp := range experiences {
		originalJSON, marshalErr := json.Marshal(exp)
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal experience original data: %w", marshalErr)
		}

		profileExp := &domain.ProfileExperience{
			ID:             uuid.New(),
			ProfileID:      profileID,
			Company:        exp.Company,
			Title:          exp.Title,
			Location:       exp.Location,
			StartDate:      exp.StartDate,
			EndDate:        exp.EndDate,
			IsCurrent:      exp.IsCurrent,
			Description:    exp.Description,
			DisplayOrder:   displayOrder + i,
			Source:         domain.ExperienceSourceResumeExtracted,
			SourceResumeID: &resumeID,
			OriginalData:   originalJSON,
		}
		if createErr := w.profileExpRepo.Create(ctx, profileExp); createErr != nil {
			return fmt.Errorf("failed to create experience for %s at %s: %w", exp.Title, exp.Company, createErr)
		}
	}
	return nil
}

func (w *ResumeProcessingWorker) materializeEducation(ctx context.Context, resumeID, profileID uuid.UUID, educations []domain.Education) error {
	displayOrder, err := w.profileEduRepo.GetNextDisplayOrder(ctx, profileID)
	if err != nil {
		return fmt.Errorf("failed to get next education display order: %w", err)
	}

	for i, edu := range educations {
		originalJSON, marshalErr := json.Marshal(edu)
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal education original data: %w", marshalErr)
		}

		degree := "Degree"
		if edu.Degree != nil && *edu.Degree != "" {
			degree = *edu.Degree
		}

		profileEdu := &domain.ProfileEducation{
			ID:             uuid.New(),
			ProfileID:      profileID,
			Institution:    edu.Institution,
			Degree:         degree,
			Field:          edu.Field,
			StartDate:      edu.StartDate,
			EndDate:        edu.EndDate,
			Description:    edu.Achievements, // Map achievements -> description
			GPA:            edu.GPA,
			DisplayOrder:   displayOrder + i,
			Source:         domain.ExperienceSourceResumeExtracted,
			SourceResumeID: &resumeID,
			OriginalData:   originalJSON,
		}
		if createErr := w.profileEduRepo.Create(ctx, profileEdu); createErr != nil {
			return fmt.Errorf("failed to create education for %s: %w", edu.Institution, createErr)
		}
	}
	return nil
}

func (w *ResumeProcessingWorker) materializeSkills(ctx context.Context, resumeID, profileID uuid.UUID, skills []string) error {
	displayOrder, err := w.profileSkillRepo.GetNextDisplayOrder(ctx, profileID)
	if err != nil {
		return fmt.Errorf("failed to get next skill display order: %w", err)
	}

	// Deduplicate skills before inserting to avoid unique constraint violations
	dedupedSkills := deduplicateSkills(skills)

	for i, skillName := range dedupedSkills {
		profileSkill := &domain.ProfileSkill{
			ID:             uuid.New(),
			ProfileID:      profileID,
			Name:           skillName,
			NormalizedName: strings.ToLower(skillName),
			Category:       "TECHNICAL",
			DisplayOrder:   displayOrder + i,
			Source:         domain.ExperienceSourceResumeExtracted,
			SourceResumeID: &resumeID,
		}
		// Use CreateIgnoreDuplicate to silently skip skills that already exist
		// (e.g., from manual entry or previous extraction from another resume)
		if createErr := w.profileSkillRepo.CreateIgnoreDuplicate(ctx, profileSkill); createErr != nil {
			return fmt.Errorf("failed to create skill %q: %w", skillName, createErr)
		}
	}
	return nil
}

// deduplicateSkills removes duplicate skills by normalized name, preserving the first occurrence.
// It also trims whitespace and filters out empty strings.
func deduplicateSkills(skills []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(skills))

	for _, skill := range skills {
		trimmed := strings.TrimSpace(skill)
		if trimmed == "" {
			continue
		}
		normalized := strings.ToLower(trimmed)
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, trimmed)
		}
	}
	return result
}
