// Package storage provides object storage implementations.
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"backend/internal/config"
	"backend/internal/domain"
)

// MinIOStorage implements domain.Storage using MinIO/S3.
type MinIOStorage struct {
	client           *minio.Client // Client for internal operations (upload, download, delete)
	publicClient     *minio.Client // Client for generating presigned URLs with public endpoint
	bucket           string
	publicEndpoint   string
	internalEndpoint string
	storageProxyURL  string // If set, use proxy URLs instead of presigned URLs
}

// NewMinIOStorage creates a new MinIO storage client.
func NewMinIOStorage(cfg config.MinIOConfig) (*MinIOStorage, error) {
	// Internal client for backend operations
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Public client for generating presigned URLs
	// Uses public endpoint so signatures are valid when accessed from browser
	// Set region explicitly to skip bucket location lookup (which would fail
	// because the public endpoint isn't accessible from inside the container)
	publicEndpoint := cfg.PublicEndpoint
	if publicEndpoint == "" {
		publicEndpoint = cfg.Endpoint
	}

	publicClient, err := minio.New(publicEndpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure:       cfg.UseSSL,
		Region:       "us-east-1",      // Default region for MinIO
		BucketLookup: minio.BucketLookupPath, // Use path-style URLs, skip bucket location lookup
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create public MinIO client: %w", err)
	}

	return &MinIOStorage{
		client:           client,
		publicClient:     publicClient,
		bucket:           cfg.Bucket,
		publicEndpoint:   publicEndpoint,
		internalEndpoint: cfg.Endpoint,
		storageProxyURL:  cfg.StorageProxyURL,
	}, nil
}

// EnsureBucket creates the bucket if it doesn't exist and sets up bucket policies.
func (s *MinIOStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Set bucket policy to allow public read for profile photos
	// This enables the storage proxy to serve images without authentication
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/profile-photos/*"]
			}
		]
	}`, s.bucket)

	err = s.client.SetBucketPolicy(ctx, s.bucket, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

// Upload stores data with the given key.
func (s *MinIOStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (*domain.StorageObject, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := s.client.PutObject(ctx, s.bucket, key, reader, size, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	return &domain.StorageObject{
		Key:         info.Key,
		Size:        info.Size,
		ContentType: contentType,
		ETag:        info.ETag,
	}, nil
}

// Download retrieves data by key.
func (s *MinIOStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	// Verify the object exists by checking stat
	_, err = obj.Stat()
	if err != nil {
		obj.Close() //nolint:errcheck,gosec // Best effort cleanup on error path
		var errResponse minio.ErrorResponse
		if errors.As(err, &errResponse) && errResponse.Code == "NoSuchKey" {
			return nil, fmt.Errorf("object not found: %s", key)
		}
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}

	return obj, nil
}

// Delete removes an object by key.
func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// GetPresignedURL generates a time-limited URL for direct access.
// Uses the public client so the signature is valid for the public endpoint.
func (s *MinIOStorage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	url, err := s.publicClient.PresignedGetObject(ctx, s.bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// GetPublicURL returns a publicly accessible URL for the object.
// If a storage proxy is configured, returns a proxy URL (e.g., "/api/storage/path/to/file").
// Otherwise, falls back to a presigned URL.
func (s *MinIOStorage) GetPublicURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if s.storageProxyURL != "" {
		// Use proxy URL - the frontend proxy will handle authentication
		return fmt.Sprintf("%s/%s", s.storageProxyURL, key), nil
	}
	// Fall back to presigned URL
	return s.GetPresignedURL(ctx, key, expiry)
}

// Exists checks if an object exists.
func (s *MinIOStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		var errResponse minio.ErrorResponse
		if errors.As(err, &errResponse) && errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}
	return true, nil
}

// Verify MinIOStorage implements domain.Storage.
var _ domain.Storage = (*MinIOStorage)(nil)
