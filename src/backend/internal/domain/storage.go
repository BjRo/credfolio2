package domain

import (
	"context"
	"io"
	"time"
)

// StorageObject represents metadata about a stored object.
type StorageObject struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	Key         string
	Size        int64
	ContentType string
	ETag        string
}

// Storage defines operations for object storage (MinIO/S3).
type Storage interface {
	// Upload stores data with the given key.
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*StorageObject, error)

	// Download retrieves data by key.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes an object by key.
	Delete(ctx context.Context, key string) error

	// GetPresignedURL generates a time-limited URL for direct access.
	GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	// Exists checks if an object exists.
	Exists(ctx context.Context, key string) (bool, error)
}
