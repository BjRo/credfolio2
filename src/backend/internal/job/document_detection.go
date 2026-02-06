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

// DocumentDetectionArgs contains the arguments for a document detection job.
type DocumentDetectionArgs struct { //nolint:govet // Field ordering prioritizes readability
	FileID      uuid.UUID `json:"file_id"`
	UserID      uuid.UUID `json:"user_id"`
	StorageKey  string    `json:"storage_key"`
	ContentType string    `json:"content_type"`
}

// Kind returns the job type identifier for River.
func (DocumentDetectionArgs) Kind() string {
	return "document_detection"
}

// InsertOpts returns default insert options.
func (DocumentDetectionArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 2}
}

// DocumentDetectionWorker runs lightweight document content detection in the background.
// It downloads the file, extracts text, runs classification, and stores the result.
type DocumentDetectionWorker struct {
	river.WorkerDefaults[DocumentDetectionArgs]
	fileRepo  domain.FileRepository
	storage   domain.Storage
	extractor domain.DocumentExtractor
	log       logger.Logger
}

// NewDocumentDetectionWorker creates a new document detection worker.
func NewDocumentDetectionWorker(
	fileRepo domain.FileRepository,
	storage domain.Storage,
	extractor domain.DocumentExtractor,
	log logger.Logger,
) *DocumentDetectionWorker {
	return &DocumentDetectionWorker{
		fileRepo:  fileRepo,
		storage:   storage,
		extractor: extractor,
		log:       log,
	}
}

// Timeout returns the maximum duration for a detection job.
// Detection involves text extraction + classification, which is lighter than full extraction.
func (w *DocumentDetectionWorker) Timeout(*river.Job[DocumentDetectionArgs]) time.Duration {
	return 5 * time.Minute
}

// Work processes a document for content detection.
func (w *DocumentDetectionWorker) Work(ctx context.Context, job *river.Job[DocumentDetectionArgs]) error {
	args := job.Args
	w.log.Info("Starting document detection",
		logger.Feature("jobs"),
		logger.String("file_id", args.FileID.String()),
	)

	// Mark as processing
	if err := w.updateDetectionStatus(ctx, args.FileID, domain.DetectionStatusProcessing, nil, nil, nil); err != nil {
		return fmt.Errorf("failed to mark as processing: %w", err)
	}

	// Download file from storage
	reader, err := w.storage.Download(ctx, args.StorageKey)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download file: %v", err)
		w.markFailed(ctx, args.FileID, errMsg)
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer reader.Close() //nolint:errcheck // Best effort cleanup

	data, err := io.ReadAll(reader)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read file data: %v", err)
		w.markFailed(ctx, args.FileID, errMsg)
		return fmt.Errorf("failed to read file data: %w", err)
	}

	// Extract text from the document
	text, err := w.extractor.ExtractText(ctx, data, args.ContentType)
	if err != nil {
		errMsg := fmt.Sprintf("failed to extract text: %v", err)
		w.markFailed(ctx, args.FileID, errMsg)
		return fmt.Errorf("failed to extract text: %w", err)
	}

	// Run lightweight detection
	detection, err := w.extractor.DetectDocumentContent(ctx, text)
	if err != nil {
		errMsg := fmt.Sprintf("failed to detect content: %v", err)
		w.markFailed(ctx, args.FileID, errMsg)
		return fmt.Errorf("failed to detect content: %w", err)
	}

	// Marshal detection result to JSON
	resultJSON, err := json.Marshal(detection)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal detection result: %v", err)
		w.markFailed(ctx, args.FileID, errMsg)
		return fmt.Errorf("failed to marshal detection result: %w", err)
	}

	// Save completed result with extracted text for reuse by processing worker
	if err := w.updateDetectionStatus(ctx, args.FileID, domain.DetectionStatusCompleted, resultJSON, nil, &text); err != nil {
		return fmt.Errorf("failed to save detection result: %w", err)
	}

	w.log.Info("Document detection completed",
		logger.Feature("jobs"),
		logger.String("file_id", args.FileID.String()),
	)
	return nil
}

// updateDetectionStatus updates the detection fields on the file record.
// When extractedText is non-nil, the text is stored for reuse by the processing worker,
// eliminating a redundant LLM text extraction call.
func (w *DocumentDetectionWorker) updateDetectionStatus(ctx context.Context, fileID uuid.UUID, status domain.DetectionStatus, result json.RawMessage, errMsg *string, extractedText *string) error {
	file, err := w.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}
	if file == nil {
		return fmt.Errorf("file not found: %s", fileID)
	}

	file.DetectionStatus = &status
	file.DetectionResult = result
	file.DetectionError = errMsg
	file.ExtractedText = extractedText

	return w.fileRepo.Update(ctx, file)
}

// markFailed marks the detection as failed with an error message.
func (w *DocumentDetectionWorker) markFailed(ctx context.Context, fileID uuid.UUID, errMsg string) {
	w.log.Error("Document detection failed",
		logger.Feature("jobs"),
		logger.String("file_id", fileID.String()),
		logger.String("error", errMsg),
	)

	if updateErr := w.updateDetectionStatus(ctx, fileID, domain.DetectionStatusFailed, nil, &errMsg, nil); updateErr != nil {
		w.log.Error("Failed to update detection status",
			logger.Feature("jobs"),
			logger.Err(updateErr),
		)
	}
}
