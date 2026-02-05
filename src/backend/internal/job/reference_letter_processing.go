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

// ReferenceLetterProcessingArgs contains the arguments for a reference letter processing job.
type ReferenceLetterProcessingArgs struct { //nolint:govet // Field ordering prioritizes JSON consistency over memory alignment
	StorageKey        string    `json:"storage_key"`
	ReferenceLetterID uuid.UUID `json:"reference_letter_id"`
	FileID            uuid.UUID `json:"file_id"`
	ContentType       string    `json:"content_type"`
}

// Kind returns the job type identifier for River.
func (ReferenceLetterProcessingArgs) Kind() string {
	return "reference_letter_processing"
}

// ReferenceLetterProcessingWorker processes uploaded reference letters to extract credibility data.
type ReferenceLetterProcessingWorker struct {
	river.WorkerDefaults[ReferenceLetterProcessingArgs]
	letterRepo          domain.ReferenceLetterRepository
	fileRepo            domain.FileRepository
	profileRepo         domain.ProfileRepository
	profileSkillRepo    domain.ProfileSkillRepository
	authorRepo          domain.AuthorRepository
	testimonialRepo     domain.TestimonialRepository
	skillValidationRepo domain.SkillValidationRepository
	storage             domain.Storage
	extractor           domain.DocumentExtractor
	log                 logger.Logger
}

// NewReferenceLetterProcessingWorker creates a new reference letter processing worker.
func NewReferenceLetterProcessingWorker(
	letterRepo domain.ReferenceLetterRepository,
	fileRepo domain.FileRepository,
	profileRepo domain.ProfileRepository,
	profileSkillRepo domain.ProfileSkillRepository,
	authorRepo domain.AuthorRepository,
	testimonialRepo domain.TestimonialRepository,
	skillValidationRepo domain.SkillValidationRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *ReferenceLetterProcessingWorker {
	return &ReferenceLetterProcessingWorker{
		letterRepo:          letterRepo,
		fileRepo:            fileRepo,
		profileRepo:         profileRepo,
		profileSkillRepo:    profileSkillRepo,
		authorRepo:          authorRepo,
		testimonialRepo:     testimonialRepo,
		skillValidationRepo: skillValidationRepo,
		storage:             storage,
		extractor:           extractor,
		log:                 log,
	}
}

// Work processes a reference letter and extracts credibility data using LLM.
func (w *ReferenceLetterProcessingWorker) Work(ctx context.Context, job *river.Job[ReferenceLetterProcessingArgs]) error {
	args := job.Args
	w.log.Info("Processing reference letter",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
		logger.String("file_id", args.FileID.String()),
		logger.String("storage_key", args.StorageKey),
	)

	// Update status to processing
	if err := w.updateStatus(ctx, args.ReferenceLetterID, domain.ReferenceLetterStatusProcessing, nil); err != nil {
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
				logger.String("reference_letter_id", args.ReferenceLetterID.String()),
				logger.String("file_id", args.FileID.String()),
				logger.Err(err),
			)
			_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
			return fmt.Errorf("failed to get file record: %w", err)
		}
		if file == nil {
			errMsg := "file record not found"
			_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
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
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.String("storage_key", args.StorageKey),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer reader.Close() //nolint:errcheck // Best effort cleanup

	// Read file data
	data, err := io.ReadAll(reader)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read file data: %v", err)
		w.log.Error("Failed to read file data",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Get reference letter to find user ID
	letter, err := w.letterRepo.GetByID(ctx, args.ReferenceLetterID)
	if err != nil || letter == nil {
		errMsg := "reference letter not found"
		if err != nil {
			errMsg = fmt.Sprintf("failed to get reference letter: %v", err)
		}
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("reference letter not found: %s", args.ReferenceLetterID)
	}

	// Get the user's profile and existing skills for context
	// Note: existingSkillsMap is no longer needed since processExtractedData was removed
	profileSkills, _ := w.getProfileSkillsContext(ctx, letter.UserID)

	// Extract credibility data using LLM with profile skills context
	extractedData, err := w.extractLetterData(ctx, data, contentType, profileSkills)
	if err != nil {
		errMsg := fmt.Sprintf("failed to extract letter data: %v", err)
		w.log.Error("Letter extraction failed",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to extract letter data: %w", err)
	}

	// Save extracted data and mark as completed
	if saveErr := w.saveExtractedData(ctx, args.ReferenceLetterID, extractedData); saveErr != nil {
		errMsg := fmt.Sprintf("failed to save extracted data: %v", saveErr)
		w.log.Error("Failed to save extracted data",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.Err(saveErr),
		)
		_ = w.updateStatusFailed(ctx, args.ReferenceLetterID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to save extracted data: %w", saveErr)
	}

	// Note: We don't create testimonials/validations here. The extracted data is stored
	// in JSON format for the user to preview. When they click "Apply Selected" on the
	// validation preview page, the applyReferenceLetterValidations mutation creates
	// the selected records. This gives users control over what gets imported.

	w.log.Info("Reference letter processing completed",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
		logger.String("author_name", extractedData.Author.Name),
		logger.Int("testimonials_count", len(extractedData.Testimonials)),
		logger.Int("skill_mentions_count", len(extractedData.SkillMentions)),
		logger.Int("experience_mentions_count", len(extractedData.ExperienceMentions)),
		logger.Int("discovered_skills_count", len(extractedData.DiscoveredSkills)),
	)

	return nil
}

// extractLetterData uses the LLM to extract structured credibility data from the reference letter.
// The profileSkills parameter provides context about existing profile skills, enabling the LLM to
// distinguish between mentions of existing skills (for validation) and newly discovered skills.
func (w *ReferenceLetterProcessingWorker) extractLetterData(ctx context.Context, data []byte, contentType string, profileSkills []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	ctx, span := otel.Tracer("credfolio").Start(ctx, "reference_letter_extraction")
	defer span.End()

	// First, extract text from the document
	text, err := w.extractor.ExtractText(ctx, data, contentType)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Then, use LLM to extract structured credibility data from the text
	extractedData, err := w.extractor.ExtractLetterData(ctx, text, profileSkills)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to extract letter data: %w", err)
	}

	// Set extraction timestamp
	extractedData.Metadata = domain.ExtractionMetadata{
		ExtractedAt:  time.Now(),
		ModelVersion: "claude-sonnet-4-20250514", // TODO: Get from extractor config
	}

	return extractedData, nil
}

// saveExtractedData saves the extracted data to the reference letter record and marks as completed.
func (w *ReferenceLetterProcessingWorker) saveExtractedData(ctx context.Context, letterID uuid.UUID, data *domain.ExtractedLetterData) error {
	letter, err := w.letterRepo.GetByID(ctx, letterID)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}
	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", letterID)
	}

	// Marshal extracted data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal extracted data: %w", err)
	}

	// Update letter with extracted data
	letter.ExtractedData = jsonData
	letter.Status = domain.ReferenceLetterStatusCompleted
	letter.ErrorMessage = nil

	// Also populate the author fields for easier querying
	if data.Author.Name != "" {
		letter.AuthorName = &data.Author.Name
	}
	if data.Author.Title != nil {
		letter.AuthorTitle = data.Author.Title
	}
	if data.Author.Company != nil {
		letter.Organization = data.Author.Company
	}

	if err := w.letterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}

	return nil
}

// updateStatus updates the reference letter status.
func (w *ReferenceLetterProcessingWorker) updateStatus(ctx context.Context, id uuid.UUID, status domain.ReferenceLetterStatus, errMsg *string) error {
	letter, err := w.letterRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}
	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", id)
	}

	letter.Status = status
	letter.ErrorMessage = errMsg

	if err := w.letterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}

	return nil
}

// updateStatusFailed is a helper to update status to failed with an error message.
func (w *ReferenceLetterProcessingWorker) updateStatusFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return w.updateStatus(ctx, id, domain.ReferenceLetterStatusFailed, &errMsg)
}

// normalizeSkillName produces a lowercase, trimmed version of a skill name for matching.
func normalizeSkillName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// getProfileSkillsContext fetches the user's profile skills and returns:
// 1. A slice of ProfileSkillContext for the LLM
// 2. A map of normalized skill names to ProfileSkill records for matching
func (w *ReferenceLetterProcessingWorker) getProfileSkillsContext(ctx context.Context, userID uuid.UUID) ([]domain.ProfileSkillContext, map[string]*domain.ProfileSkill) {
	// Get user's profile
	profile, err := w.profileRepo.GetByUserID(ctx, userID)
	if err != nil || profile == nil {
		w.log.Warning("Could not get profile for skill context",
			logger.Feature("jobs"),
			logger.String("user_id", userID.String()),
			logger.Err(err),
		)
		return nil, nil
	}

	// Get profile skills
	skills, err := w.profileSkillRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		w.log.Warning("Could not get profile skills for context",
			logger.Feature("jobs"),
			logger.String("profile_id", profile.ID.String()),
			logger.Err(err),
		)
		return nil, nil
	}

	// Build context and lookup map
	skillContext := make([]domain.ProfileSkillContext, 0, len(skills))
	skillsMap := make(map[string]*domain.ProfileSkill, len(skills))
	for _, skill := range skills {
		skillContext = append(skillContext, domain.ProfileSkillContext{
			Name:           skill.Name,
			NormalizedName: skill.NormalizedName,
			Category:       skill.Category,
		})
		skillsMap[normalizeSkillName(skill.Name)] = skill
		// Also add by normalized name for better matching
		if skill.NormalizedName != "" {
			skillsMap[skill.NormalizedName] = skill
		}
	}

	return skillContext, skillsMap
}

// NOTE: processExtractedData, findOrCreateAuthor, and related functions were removed.
// The job now only extracts data and stores it as JSON. Record creation (testimonials,
// skill validations, discovered skills) happens in the applyReferenceLetterValidations
// mutation based on user selection, not automatically during processing.

// mapAuthorToTestimonialRelationship maps an AuthorRelationship to a TestimonialRelationship.
// NOTE: This function is also defined in converter.go for use by the resolver.
func mapAuthorToTestimonialRelationship(ar domain.AuthorRelationship) domain.TestimonialRelationship {
	switch ar {
	case domain.AuthorRelationshipManager:
		return domain.TestimonialRelationshipManager
	case domain.AuthorRelationshipPeer, domain.AuthorRelationshipColleague:
		return domain.TestimonialRelationshipPeer
	case domain.AuthorRelationshipDirectReport:
		return domain.TestimonialRelationshipDirectReport
	case domain.AuthorRelationshipClient:
		return domain.TestimonialRelationshipClient
	default:
		return domain.TestimonialRelationshipOther
	}
}
