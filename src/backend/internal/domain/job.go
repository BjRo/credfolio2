package domain

import (
	"context"

	"github.com/google/uuid"
)

// DocumentProcessingRequest contains the data needed to enqueue a document processing job.
type DocumentProcessingRequest struct { //nolint:govet // Field ordering matches job args convention
	StorageKey        string
	ReferenceLetterID uuid.UUID
	FileID            uuid.UUID
	ContentType       string
}

// UnifiedDocumentProcessingRequest contains the data needed to enqueue a unified document processing job.
// The unified worker downloads the file once, extracts text once, then runs the selected extractors.
type UnifiedDocumentProcessingRequest struct { //nolint:govet // Field ordering matches job args convention
	StorageKey  string
	FileID      uuid.UUID
	ContentType string
	UserID      uuid.UUID

	// At least one of these must be set. When set, the worker runs the corresponding extractor
	// and stores results in the linked entity.
	ResumeID          *uuid.UUID
	ReferenceLetterID *uuid.UUID
}

// DocumentDetectionRequest contains the data needed to enqueue a document detection job.
type DocumentDetectionRequest struct { //nolint:govet // Field ordering prioritizes readability
	FileID      uuid.UUID
	UserID      uuid.UUID
	StorageKey  string
	ContentType string
}

// JobEnqueuer defines the interface for enqueueing background jobs.
type JobEnqueuer interface {
	// EnqueueDocumentProcessing adds a document processing job to the queue.
	EnqueueDocumentProcessing(ctx context.Context, req DocumentProcessingRequest) error

	// EnqueueResumeProcessing adds a resume processing job to the queue.
	EnqueueResumeProcessing(ctx context.Context, req ResumeProcessingRequest) error

	// EnqueueUnifiedDocumentProcessing adds a unified document processing job to the queue.
	// The unified worker extracts text once and runs the selected extractors (resume, letter, or both).
	EnqueueUnifiedDocumentProcessing(ctx context.Context, req UnifiedDocumentProcessingRequest) error

	// EnqueueDocumentDetection adds a document detection job to the queue.
	// The worker extracts text and runs lightweight content classification.
	EnqueueDocumentDetection(ctx context.Context, req DocumentDetectionRequest) error
}
