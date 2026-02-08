package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"backend/internal/domain"
	"backend/internal/logger"
	"backend/internal/service"
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
	resumeRepo        domain.ResumeRepository
	fileRepo          domain.FileRepository
	storage           domain.Storage
	extractor         domain.DocumentExtractor
	materializationSvc *service.MaterializationService
	log               logger.Logger
}

// NewResumeProcessingWorker creates a new resume processing worker.
func NewResumeProcessingWorker(
	resumeRepo domain.ResumeRepository,
	fileRepo domain.FileRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	materializationSvc *service.MaterializationService,
	log logger.Logger,
) *ResumeProcessingWorker {
	return &ResumeProcessingWorker{
		resumeRepo:        resumeRepo,
		fileRepo:          fileRepo,
		storage:           storage,
		extractor:         extractor,
		materializationSvc: materializationSvc,
		log:               log,
	}
}

// Timeout overrides River's default 60s job timeout with a 10-minute safety net.
// LLM extraction can take several minutes for structured output; the primary
// timeout is handled by the resilient provider layer (300s). This River timeout
// serves as an outer safety net to prevent worker pool exhaustion.
func (w *ResumeProcessingWorker) Timeout(*river.Job[ResumeProcessingArgs]) time.Duration {
	return 10 * time.Minute
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
	extractedData, err := w.extractResumeData(ctx, args.FileID, data, contentType)
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
		if _, matErr := w.materializationSvc.MaterializeResumeData(ctx, args.ResumeID, resume.UserID, extractedData); matErr != nil {
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
func (w *ResumeProcessingWorker) extractResumeData(ctx context.Context, fileID uuid.UUID, data []byte, contentType string) (*domain.ResumeExtractedData, error) {
	ctx, span := otel.Tracer("credfolio").Start(ctx, "resume_extraction")
	defer span.End()

	var text string

	// Check if we already have extracted text from the detection phase
	file, err := w.fileRepo.GetByID(ctx, fileID)
	if err == nil && file != nil && file.ExtractedText != nil && *file.ExtractedText != "" {
		w.log.Info("Reusing extracted text from detection phase",
			logger.Feature("jobs"),
			logger.String("file_id", fileID.String()),
		)
		text = *file.ExtractedText
	} else {
		// Extract text from the document
		text, err = w.extractor.ExtractText(ctx, data, contentType)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, fmt.Errorf("failed to extract text: %w", err)
		}
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
