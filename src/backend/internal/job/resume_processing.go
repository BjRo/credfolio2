package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

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

// ResumeProcessingWorker processes uploaded resumes to extract profile data.
type ResumeProcessingWorker struct {
	river.WorkerDefaults[ResumeProcessingArgs]
	resumeRepo domain.ResumeRepository
	fileRepo   domain.FileRepository
	storage    domain.Storage
	extractor  domain.DocumentExtractor
	log        logger.Logger
}

// NewResumeProcessingWorker creates a new resume processing worker.
func NewResumeProcessingWorker(
	resumeRepo domain.ResumeRepository,
	fileRepo domain.FileRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *ResumeProcessingWorker {
	return &ResumeProcessingWorker{
		resumeRepo: resumeRepo,
		fileRepo:   fileRepo,
		storage:    storage,
		extractor:  extractor,
		log:        log,
	}
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
			_ = w.updateStatusFailed(ctx, args.ResumeID, errMsg) //nolint:errcheck
			return fmt.Errorf("failed to get file record: %w", err)
		}
		if file == nil {
			errMsg := "file record not found"
			_ = w.updateStatusFailed(ctx, args.ResumeID, errMsg) //nolint:errcheck
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
		_ = w.updateStatusFailed(ctx, args.ResumeID, errMsg) //nolint:errcheck
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
		_ = w.updateStatusFailed(ctx, args.ResumeID, errMsg) //nolint:errcheck
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
		_ = w.updateStatusFailed(ctx, args.ResumeID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to extract resume data: %w", err)
	}

	// Save extracted data
	if err := w.saveExtractedData(ctx, args.ResumeID, extractedData); err != nil {
		errMsg := fmt.Sprintf("failed to save extracted data: %v", err)
		w.log.Error("Failed to save extracted data",
			logger.Feature("jobs"),
			logger.String("resume_id", args.ResumeID.String()),
			logger.Err(err),
		)
		_ = w.updateStatusFailed(ctx, args.ResumeID, errMsg) //nolint:errcheck
		return fmt.Errorf("failed to save extracted data: %w", err)
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
	// First, extract text from the document
	text, err := w.extractor.ExtractText(ctx, data, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Then, use LLM to extract structured profile data from the text
	extractedData, err := w.extractor.ExtractResumeData(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract resume data: %w", err)
	}

	// Set extraction timestamp
	extractedData.ExtractedAt = time.Now()

	return extractedData, nil
}

// saveExtractedData saves the extracted data to the resume record.
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

	resume.Status = domain.ResumeStatusCompleted
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

// updateStatusFailed is a helper to update status to failed with an error message.
func (w *ResumeProcessingWorker) updateStatusFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return w.updateStatus(ctx, id, domain.ResumeStatusFailed, &errMsg)
}
