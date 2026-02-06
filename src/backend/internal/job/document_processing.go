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

// DocumentProcessingArgs contains the arguments for a unified document processing job.
// The worker downloads the file once, extracts text once, then runs the selected extractors.
type DocumentProcessingArgs struct { //nolint:govet // Field ordering prioritizes JSON consistency over memory alignment
	StorageKey  string    `json:"storage_key"`
	FileID      uuid.UUID `json:"file_id"`
	ContentType string    `json:"content_type"`
	UserID      uuid.UUID `json:"user_id"`

	// At least one of these must be set. When set, the worker runs the corresponding extractor.
	ResumeID          *uuid.UUID `json:"resume_id,omitempty"`
	ReferenceLetterID *uuid.UUID `json:"reference_letter_id,omitempty"`
}

// Kind returns the job type identifier for River.
func (DocumentProcessingArgs) Kind() string {
	return "document_processing"
}

// InsertOpts returns default insert options — limits retries to avoid
// repeatedly hammering LLM providers on persistent failures.
func (DocumentProcessingArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 2}
}

// DocumentProcessingWorker processes uploaded documents using the unified extraction pipeline.
// It downloads the file once, extracts text once, then runs the selected extractors
// (resume and/or reference letter) and stores the results as JSON.
type DocumentProcessingWorker struct {
	river.WorkerDefaults[DocumentProcessingArgs]
	resumeRepo       domain.ResumeRepository
	refLetterRepo    domain.ReferenceLetterRepository
	fileRepo         domain.FileRepository
	profileRepo      domain.ProfileRepository
	profileSkillRepo domain.ProfileSkillRepository
	storage          domain.Storage
	extractor        domain.DocumentExtractor
	log              logger.Logger
}

// NewDocumentProcessingWorker creates a new unified document processing worker.
func NewDocumentProcessingWorker(
	resumeRepo domain.ResumeRepository,
	refLetterRepo domain.ReferenceLetterRepository,
	fileRepo domain.FileRepository,
	profileRepo domain.ProfileRepository,
	profileSkillRepo domain.ProfileSkillRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *DocumentProcessingWorker {
	return &DocumentProcessingWorker{
		resumeRepo:       resumeRepo,
		refLetterRepo:    refLetterRepo,
		fileRepo:         fileRepo,
		profileRepo:      profileRepo,
		profileSkillRepo: profileSkillRepo,
		storage:          storage,
		extractor:        extractor,
		log:              log,
	}
}

// Timeout overrides River's default 60s job timeout with a 10-minute safety net.
// LLM extraction can take several minutes for structured output; the primary
// timeout is handled by the resilient provider layer (300s). This River timeout
// serves as an outer safety net to prevent worker pool exhaustion.
func (w *DocumentProcessingWorker) Timeout(*river.Job[DocumentProcessingArgs]) time.Duration {
	return 10 * time.Minute
}

// Work processes a document through the unified extraction pipeline.
func (w *DocumentProcessingWorker) Work(ctx context.Context, job *river.Job[DocumentProcessingArgs]) error {
	args := job.Args
	w.log.Info("Processing document (unified)",
		logger.Feature("jobs"),
		logger.String("file_id", args.FileID.String()),
		logger.String("storage_key", args.StorageKey),
		logger.String("resume_id", uuidPtrStr(args.ResumeID)),
		logger.String("reference_letter_id", uuidPtrStr(args.ReferenceLetterID)),
	)

	// Mark entities as processing
	if err := w.markProcessing(ctx, args); err != nil {
		return err
	}

	// Load file record to check for stored extracted text and resolve content type
	file, err := w.fileRepo.GetByID(ctx, args.FileID)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get file record: %v", err)
		w.markAllFailed(ctx, args, errMsg)
		return fmt.Errorf("failed to get file record: %w", err)
	}
	if file == nil {
		errMsg := "file record not found" //nolint:goconst // same string in different workers, not worth extracting
		w.markAllFailed(ctx, args, errMsg)
		return fmt.Errorf("file record not found: %s", args.FileID) //nolint:goconst // see above
	}

	contentType := args.ContentType
	if contentType == "" {
		contentType = file.ContentType
	}

	// Use stored extracted text if available (set by detection worker), otherwise download and extract
	var text string
	if file.ExtractedText != nil && *file.ExtractedText != "" {
		text = *file.ExtractedText
		w.log.Info("Using stored extracted text, skipping download and LLM extraction",
			logger.Feature("jobs"),
			logger.String("file_id", args.FileID.String()),
		)
	} else {
		// Download file from storage (once for all extractors)
		reader, dlErr := w.storage.Download(ctx, args.StorageKey)
		if dlErr != nil {
			errMsg := fmt.Sprintf("failed to download file: %v", dlErr)
			w.log.Error("Storage download failed",
				logger.Feature("jobs"),
				logger.String("storage_key", args.StorageKey),
				logger.Err(dlErr),
			)
			w.markAllFailed(ctx, args, errMsg)
			return fmt.Errorf("failed to download file: %w", dlErr)
		}
		defer reader.Close() //nolint:errcheck // Best effort cleanup

		// Read file data
		data, readErr := io.ReadAll(reader)
		if readErr != nil {
			errMsg := fmt.Sprintf("failed to read file data: %v", readErr)
			w.markAllFailed(ctx, args, errMsg)
			return fmt.Errorf("failed to read file data: %w", readErr)
		}

		// Extract text via LLM
		extracted, extractErr := w.extractText(ctx, data, contentType)
		if extractErr != nil {
			errMsg := fmt.Sprintf("failed to extract text: %v", extractErr)
			w.markAllFailed(ctx, args, errMsg)
			return fmt.Errorf("failed to extract text: %w", extractErr)
		}
		text = extracted
	}

	// Run extractors sequentially — resume first so we can pass its skills to the letter extractor
	var extractionErrors []string
	var resumeSkills []domain.ProfileSkillContext

	if args.ResumeID != nil {
		skills, resumeErr := w.processResumeExtraction(ctx, *args.ResumeID, text)
		if resumeErr != nil {
			extractionErrors = append(extractionErrors, fmt.Sprintf("resume: %v", resumeErr))
		} else {
			resumeSkills = skills
		}
	}

	if args.ReferenceLetterID != nil {
		if letterErr := w.processLetterExtraction(ctx, *args.ReferenceLetterID, args.UserID, text, resumeSkills); letterErr != nil {
			extractionErrors = append(extractionErrors, fmt.Sprintf("letter: %v", letterErr))
		}
	}

	if len(extractionErrors) > 0 {
		return fmt.Errorf("extraction errors: %s", strings.Join(extractionErrors, "; "))
	}

	w.log.Info("Document processing completed (unified)",
		logger.Feature("jobs"),
		logger.String("file_id", args.FileID.String()),
	)
	return nil
}

// extractText uses the document extractor to extract raw text from the document.
func (w *DocumentProcessingWorker) extractText(ctx context.Context, data []byte, contentType string) (string, error) {
	ctx, span := otel.Tracer("credfolio").Start(ctx, "unified_text_extraction")
	defer span.End()

	text, err := w.extractor.ExtractText(ctx, data, contentType)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	return text, nil
}

// processResumeExtraction runs the resume extractor and saves results.
// Returns the extracted skills as ProfileSkillContext so they can be passed
// to the reference letter extractor (which needs them before materialization).
func (w *DocumentProcessingWorker) processResumeExtraction(ctx context.Context, resumeID uuid.UUID, text string) ([]domain.ProfileSkillContext, error) {
	ctx, span := otel.Tracer("credfolio").Start(ctx, "unified_resume_extraction")
	defer span.End()

	extractedData, err := w.extractor.ExtractResumeData(ctx, text)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		errMsg := fmt.Sprintf("failed to extract resume data: %v", err)
		w.log.Error("Resume extraction failed",
			logger.Feature("jobs"),
			logger.String("resume_id", resumeID.String()),
			logger.Err(err),
		)
		w.updateResumeStatusFailed(ctx, resumeID, errMsg)
		return nil, fmt.Errorf("resume extraction failed: %w", err)
	}

	extractedData.ExtractedAt = time.Now()

	// Save extracted data to resume record
	if saveErr := w.saveResumeExtractedData(ctx, resumeID, extractedData); saveErr != nil {
		errMsg := fmt.Sprintf("failed to save resume extracted data: %v", saveErr)
		w.updateResumeStatusFailed(ctx, resumeID, errMsg)
		return nil, fmt.Errorf("failed to save resume data: %w", saveErr)
	}

	// Mark resume as completed (no auto-materialization in unified flow)
	if err := w.updateResumeStatus(ctx, resumeID, domain.ResumeStatusCompleted, nil); err != nil {
		return nil, fmt.Errorf("failed to mark resume completed: %w", err)
	}

	// Build skill context from extracted data for the letter extractor
	skillCtx := resumeSkillsToContext(extractedData.Skills)

	w.log.Info("Resume extraction completed (unified)",
		logger.Feature("jobs"),
		logger.String("resume_id", resumeID.String()),
		logger.String("name", extractedData.Name),
		logger.Int("experience_count", len(extractedData.Experience)),
		logger.Int("education_count", len(extractedData.Education)),
		logger.Int("skills_count", len(extractedData.Skills)),
	)
	return skillCtx, nil
}

// resumeSkillsToContext converts extracted resume skill names to ProfileSkillContext
// for use in the reference letter extraction prompt.
func resumeSkillsToContext(skills []string) []domain.ProfileSkillContext {
	result := make([]domain.ProfileSkillContext, 0, len(skills))
	for _, name := range skills {
		result = append(result, domain.ProfileSkillContext{
			Name:           name,
			NormalizedName: strings.ToLower(strings.TrimSpace(name)),
		})
	}
	return result
}

// mergeSkillContexts combines existing profile skills with resume-extracted skills,
// deduplicating by normalized name.
func mergeSkillContexts(existing, additional []domain.ProfileSkillContext) []domain.ProfileSkillContext {
	if len(additional) == 0 {
		return existing
	}
	seen := make(map[string]bool, len(existing))
	for _, s := range existing {
		seen[s.NormalizedName] = true
	}
	merged := make([]domain.ProfileSkillContext, len(existing))
	copy(merged, existing)
	for _, s := range additional {
		if !seen[s.NormalizedName] {
			merged = append(merged, s)
			seen[s.NormalizedName] = true
		}
	}
	return merged
}

// processLetterExtraction runs the reference letter extractor and saves results.
// resumeSkills are skills just extracted from a co-uploaded resume (not yet materialized).
func (w *DocumentProcessingWorker) processLetterExtraction(ctx context.Context, letterID uuid.UUID, userID uuid.UUID, text string, resumeSkills []domain.ProfileSkillContext) error {
	ctx, span := otel.Tracer("credfolio").Start(ctx, "unified_letter_extraction")
	defer span.End()

	// Merge existing profile skills with just-extracted resume skills
	profileSkills := w.getProfileSkillsContext(ctx, userID)
	profileSkills = mergeSkillContexts(profileSkills, resumeSkills)

	extractedData, err := w.extractor.ExtractLetterData(ctx, text, profileSkills)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		errMsg := fmt.Sprintf("failed to extract letter data: %v", err)
		w.log.Error("Letter extraction failed",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", letterID.String()),
			logger.Err(err),
		)
		w.updateLetterStatusFailed(ctx, letterID, errMsg)
		return fmt.Errorf("letter extraction failed: %w", err)
	}

	extractedData.Metadata = domain.ExtractionMetadata{
		ExtractedAt:  time.Now(),
		ModelVersion: "claude-sonnet-4-20250514", // TODO: Get from extractor config
	}

	// Save extracted data and mark as completed
	if saveErr := w.saveLetterExtractedData(ctx, letterID, extractedData); saveErr != nil {
		errMsg := fmt.Sprintf("failed to save letter extracted data: %v", saveErr)
		w.updateLetterStatusFailed(ctx, letterID, errMsg)
		return fmt.Errorf("failed to save letter data: %w", saveErr)
	}

	w.log.Info("Letter extraction completed (unified)",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", letterID.String()),
		logger.String("author_name", extractedData.Author.Name),
		logger.Int("testimonials_count", len(extractedData.Testimonials)),
	)
	return nil
}

// saveResumeExtractedData saves the extracted data JSON to the resume record.
func (w *DocumentProcessingWorker) saveResumeExtractedData(ctx context.Context, resumeID uuid.UUID, data *domain.ResumeExtractedData) error {
	resume, err := w.resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		return fmt.Errorf("failed to get resume: %w", err)
	}
	if resume == nil {
		return fmt.Errorf("resume not found: %s", resumeID)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal extracted data: %w", err)
	}

	resume.ExtractedData = jsonData
	resume.ErrorMessage = nil

	if err := w.resumeRepo.Update(ctx, resume); err != nil {
		return fmt.Errorf("failed to update resume: %w", err)
	}
	return nil
}

// saveLetterExtractedData saves the extracted data and marks the letter as completed.
func (w *DocumentProcessingWorker) saveLetterExtractedData(ctx context.Context, letterID uuid.UUID, data *domain.ExtractedLetterData) error {
	letter, err := w.refLetterRepo.GetByID(ctx, letterID)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}
	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", letterID)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal extracted data: %w", err)
	}

	letter.ExtractedData = jsonData
	letter.Status = domain.ReferenceLetterStatusCompleted
	letter.ErrorMessage = nil

	if data.Author.Name != "" {
		letter.AuthorName = &data.Author.Name
	}
	if data.Author.Title != nil {
		letter.AuthorTitle = data.Author.Title
	}
	if data.Author.Company != nil {
		letter.Organization = data.Author.Company
	}

	if err := w.refLetterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}
	return nil
}

// getProfileSkillsContext fetches the user's profile skills for LLM context.
func (w *DocumentProcessingWorker) getProfileSkillsContext(ctx context.Context, userID uuid.UUID) []domain.ProfileSkillContext {
	profile, err := w.profileRepo.GetByUserID(ctx, userID)
	if err != nil || profile == nil {
		return nil
	}

	skills, err := w.profileSkillRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		return nil
	}

	result := make([]domain.ProfileSkillContext, 0, len(skills))
	for _, skill := range skills {
		result = append(result, domain.ProfileSkillContext{
			Name:           skill.Name,
			NormalizedName: skill.NormalizedName,
			Category:       skill.Category,
		})
	}
	return result
}

// markProcessing marks all linked entities as processing.
func (w *DocumentProcessingWorker) markProcessing(ctx context.Context, args DocumentProcessingArgs) error {
	if args.ResumeID != nil {
		if err := w.updateResumeStatus(ctx, *args.ResumeID, domain.ResumeStatusProcessing, nil); err != nil {
			return fmt.Errorf("failed to mark resume as processing: %w", err)
		}
	}
	if args.ReferenceLetterID != nil {
		if err := w.updateLetterStatus(ctx, *args.ReferenceLetterID, domain.ReferenceLetterStatusProcessing, nil); err != nil {
			return fmt.Errorf("failed to mark letter as processing: %w", err)
		}
	}
	return nil
}

// markAllFailed marks all linked entities as failed with the given error message.
func (w *DocumentProcessingWorker) markAllFailed(ctx context.Context, args DocumentProcessingArgs, errMsg string) {
	if args.ResumeID != nil {
		w.updateResumeStatusFailed(ctx, *args.ResumeID, errMsg)
	}
	if args.ReferenceLetterID != nil {
		w.updateLetterStatusFailed(ctx, *args.ReferenceLetterID, errMsg)
	}
}

// updateResumeStatus updates the resume status.
func (w *DocumentProcessingWorker) updateResumeStatus(ctx context.Context, id uuid.UUID, status domain.ResumeStatus, errMsg *string) error {
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

// updateResumeStatusFailed is a best-effort helper to mark resume as failed.
func (w *DocumentProcessingWorker) updateResumeStatusFailed(ctx context.Context, id uuid.UUID, errMsg string) {
	if err := w.updateResumeStatus(ctx, id, domain.ResumeStatusFailed, &errMsg); err != nil {
		w.log.Error("Failed to update resume status to failed",
			logger.Feature("jobs"),
			logger.String("resume_id", id.String()),
			logger.Err(err),
		)
	}
}

// updateLetterStatus updates the reference letter status.
func (w *DocumentProcessingWorker) updateLetterStatus(ctx context.Context, id uuid.UUID, status domain.ReferenceLetterStatus, errMsg *string) error {
	letter, err := w.refLetterRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}
	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", id)
	}
	letter.Status = status
	letter.ErrorMessage = errMsg
	if err := w.refLetterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}
	return nil
}

// updateLetterStatusFailed is a best-effort helper to mark letter as failed.
func (w *DocumentProcessingWorker) updateLetterStatusFailed(ctx context.Context, id uuid.UUID, errMsg string) {
	if err := w.updateLetterStatus(ctx, id, domain.ReferenceLetterStatusFailed, &errMsg); err != nil {
		w.log.Error("Failed to update letter status to failed",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", id.String()),
			logger.Err(err),
		)
	}
}

// uuidPtrStr returns the string representation of a UUID pointer, or "nil" if nil.
func uuidPtrStr(u *uuid.UUID) string {
	if u == nil {
		return "nil"
	}
	return u.String()
}
