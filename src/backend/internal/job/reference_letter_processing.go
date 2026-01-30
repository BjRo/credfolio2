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
	letterRepo domain.ReferenceLetterRepository
	fileRepo   domain.FileRepository
	storage    domain.Storage
	extractor  domain.DocumentExtractor
	log        logger.Logger
}

// NewReferenceLetterProcessingWorker creates a new reference letter processing worker.
func NewReferenceLetterProcessingWorker(
	letterRepo domain.ReferenceLetterRepository,
	fileRepo domain.FileRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *ReferenceLetterProcessingWorker {
	return &ReferenceLetterProcessingWorker{
		letterRepo: letterRepo,
		fileRepo:   fileRepo,
		storage:    storage,
		extractor:  extractor,
		log:        log,
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

	// Extract credibility data using LLM
	extractedData, err := w.extractLetterData(ctx, data, contentType)
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

	w.log.Info("Reference letter processing completed",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
		logger.String("author_name", extractedData.Author.Name),
		logger.Int("testimonials_count", len(extractedData.Testimonials)),
		logger.Int("skill_mentions_count", len(extractedData.SkillMentions)),
		logger.Int("experience_mentions_count", len(extractedData.ExperienceMentions)),
	)

	return nil
}

// extractLetterData uses the LLM to extract structured credibility data from the reference letter.
func (w *ReferenceLetterProcessingWorker) extractLetterData(ctx context.Context, data []byte, contentType string) (*domain.ExtractedLetterData, error) {
	// First, extract text from the document
	text, err := w.extractor.ExtractText(ctx, data, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Then, use LLM to extract structured credibility data from the text
	extractedData, err := w.extractor.ExtractLetterData(ctx, text)
	if err != nil {
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
