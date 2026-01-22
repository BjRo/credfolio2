package job

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/domain"
	"backend/internal/logger"
)

// DocumentProcessingArgs contains the arguments for a document processing job.
type DocumentProcessingArgs struct {
	StorageKey        string    `json:"storage_key"`
	ReferenceLetterID uuid.UUID `json:"reference_letter_id"`
	FileID            uuid.UUID `json:"file_id"`
}

// Kind returns the job type identifier for River.
func (DocumentProcessingArgs) Kind() string {
	return "document_processing"
}

// DocumentProcessingWorker processes uploaded documents to extract reference letter data.
type DocumentProcessingWorker struct {
	river.WorkerDefaults[DocumentProcessingArgs]
	refLetterRepo domain.ReferenceLetterRepository
	storage       domain.Storage
	log           logger.Logger
}

// NewDocumentProcessingWorker creates a new document processing worker.
func NewDocumentProcessingWorker(
	refLetterRepo domain.ReferenceLetterRepository,
	storage domain.Storage,
	log logger.Logger,
) *DocumentProcessingWorker {
	return &DocumentProcessingWorker{
		refLetterRepo: refLetterRepo,
		storage:       storage,
		log:           log,
	}
}

// Work processes a document and updates the reference letter status.
func (w *DocumentProcessingWorker) Work(ctx context.Context, job *river.Job[DocumentProcessingArgs]) error {
	args := job.Args
	w.log.Info("Processing document",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
		logger.String("file_id", args.FileID.String()),
		logger.String("storage_key", args.StorageKey),
	)

	// Update status to processing
	if err := w.updateStatus(ctx, args.ReferenceLetterID, domain.ReferenceLetterStatusProcessing); err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Verify file exists in storage
	exists, err := w.storage.Exists(ctx, args.StorageKey)
	if err != nil {
		w.log.Error("Storage check failed",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.String("storage_key", args.StorageKey),
			logger.Err(err),
		)
		if statusErr := w.updateStatus(ctx, args.ReferenceLetterID, domain.ReferenceLetterStatusFailed); statusErr != nil {
			w.log.Warning("Failed to update status after storage check error",
				logger.Feature("jobs"),
				logger.String("reference_letter_id", args.ReferenceLetterID.String()),
				logger.Err(statusErr),
			)
		}
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	if !exists {
		w.log.Error("File not found in storage",
			logger.Feature("jobs"),
			logger.String("reference_letter_id", args.ReferenceLetterID.String()),
			logger.String("storage_key", args.StorageKey),
		)
		if statusErr := w.updateStatus(ctx, args.ReferenceLetterID, domain.ReferenceLetterStatusFailed); statusErr != nil {
			w.log.Warning("Failed to update status after file not found",
				logger.Feature("jobs"),
				logger.String("reference_letter_id", args.ReferenceLetterID.String()),
				logger.Err(statusErr),
			)
		}
		return fmt.Errorf("file not found in storage: %s", args.StorageKey)
	}

	// TODO: Future implementation will:
	// 1. Download file from storage
	// 2. Extract text content (PDF, DOCX, TXT)
	// 3. Send to LLM for analysis
	// 4. Store extracted data in reference letter
	//
	// For now, we just mark it as completed to establish the job flow

	// Mark as completed
	if err := w.updateStatus(ctx, args.ReferenceLetterID, domain.ReferenceLetterStatusCompleted); err != nil {
		return fmt.Errorf("failed to update status to completed: %w", err)
	}

	w.log.Info("Document processing completed",
		logger.Feature("jobs"),
		logger.String("reference_letter_id", args.ReferenceLetterID.String()),
	)
	return nil
}

// updateStatus updates the reference letter status.
func (w *DocumentProcessingWorker) updateStatus(ctx context.Context, id uuid.UUID, status domain.ReferenceLetterStatus) error {
	letter, err := w.refLetterRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get reference letter: %w", err)
	}

	if letter == nil {
		return fmt.Errorf("reference letter not found: %s", id)
	}

	letter.Status = status
	if err := w.refLetterRepo.Update(ctx, letter); err != nil {
		return fmt.Errorf("failed to update reference letter: %w", err)
	}

	return nil
}
