package domain

import (
	"context"

	"github.com/google/uuid"
)

// DocumentProcessingRequest contains the data needed to enqueue a document processing job.
type DocumentProcessingRequest struct {
	StorageKey        string
	ReferenceLetterID uuid.UUID
	FileID            uuid.UUID
}

// JobEnqueuer defines the interface for enqueueing background jobs.
type JobEnqueuer interface {
	// EnqueueDocumentProcessing adds a document processing job to the queue.
	EnqueueDocumentProcessing(ctx context.Context, req DocumentProcessingRequest) error

	// EnqueueResumeProcessing adds a resume processing job to the queue.
	EnqueueResumeProcessing(ctx context.Context, req ResumeProcessingRequest) error
}
