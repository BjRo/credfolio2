// Package queue provides background job processing using River.
package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/riverqueue/river/rivertype"

	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/job"
)

// Client wraps River client and pgx pool for job queue operations.
type Client struct {
	riverClient *river.Client[pgx.Tx]
	pool        *pgxpool.Pool
}

// NewClient creates a new River queue client.
// It establishes a pgx connection pool and runs River migrations.
func NewClient(ctx context.Context, cfg config.DatabaseConfig, queueCfg config.QueueConfig, workers *river.Workers) (*Client, error) {
	// Create pgx pool for River (separate from bun connection)
	pool, err := pgxpool.New(ctx, cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Run River migrations
	if migrateErr := runMigrations(ctx, pool); migrateErr != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run river migrations: %w", migrateErr)
	}

	// Create River client
	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: queueCfg.MaxWorkers},
		},
		Workers: workers,
	})
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	return &Client{
		riverClient: riverClient,
		pool:        pool,
	}, nil
}

// runMigrations runs River's schema migrations.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		return fmt.Errorf("failed to create river migrator: %w", err)
	}

	res, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	for _, v := range res.Versions {
		log.Printf("River migration applied: version %d", v.Version)
	}

	return nil
}

// Start begins processing jobs. This is non-blocking.
func (c *Client) Start(ctx context.Context) error {
	return c.riverClient.Start(ctx)
}

// Stop gracefully shuts down the job processor.
func (c *Client) Stop(ctx context.Context) error {
	return c.riverClient.Stop(ctx)
}

// Close releases all resources.
func (c *Client) Close() {
	c.pool.Close()
}

// Insert adds a job to the queue.
func (c *Client) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return c.riverClient.Insert(ctx, args, opts)
}

// River returns the underlying River client for advanced operations.
func (c *Client) River() *river.Client[pgx.Tx] {
	return c.riverClient
}

// EnqueueDocumentProcessing adds a reference letter processing job to the queue.
func (c *Client) EnqueueDocumentProcessing(ctx context.Context, req domain.DocumentProcessingRequest) error {
	args := job.ReferenceLetterProcessingArgs{
		StorageKey:        req.StorageKey,
		ReferenceLetterID: req.ReferenceLetterID,
		FileID:            req.FileID,
		ContentType:       req.ContentType,
	}

	_, err := c.riverClient.Insert(ctx, args, nil)
	if err != nil {
		return fmt.Errorf("failed to enqueue reference letter processing job: %w", err)
	}

	return nil
}

// EnqueueResumeProcessing adds a resume processing job to the queue.
func (c *Client) EnqueueResumeProcessing(ctx context.Context, req domain.ResumeProcessingRequest) error {
	args := job.ResumeProcessingArgs{
		StorageKey:  req.StorageKey,
		ResumeID:    req.ResumeID,
		FileID:      req.FileID,
		ContentType: req.ContentType,
	}

	_, err := c.riverClient.Insert(ctx, args, nil)
	if err != nil {
		return fmt.Errorf("failed to enqueue resume processing job: %w", err)
	}

	return nil
}

// EnqueueUnifiedDocumentProcessing adds a unified document processing job to the queue.
// The unified worker extracts text once and runs the selected extractors (resume, letter, or both).
func (c *Client) EnqueueUnifiedDocumentProcessing(ctx context.Context, req domain.UnifiedDocumentProcessingRequest) error {
	args := job.DocumentProcessingArgs{
		StorageKey:        req.StorageKey,
		FileID:            req.FileID,
		ContentType:       req.ContentType,
		UserID:            req.UserID,
		ResumeID:          req.ResumeID,
		ReferenceLetterID: req.ReferenceLetterID,
	}

	_, err := c.riverClient.Insert(ctx, args, nil)
	if err != nil {
		return fmt.Errorf("failed to enqueue unified document processing job: %w", err)
	}

	return nil
}

// EnqueueDocumentDetection adds a document detection job to the queue.
func (c *Client) EnqueueDocumentDetection(ctx context.Context, req domain.DocumentDetectionRequest) error {
	args := job.DocumentDetectionArgs{
		StorageKey:  req.StorageKey,
		FileID:      req.FileID,
		ContentType: req.ContentType,
		UserID:      req.UserID,
	}

	_, err := c.riverClient.Insert(ctx, args, nil)
	if err != nil {
		return fmt.Errorf("failed to enqueue document detection job: %w", err)
	}

	return nil
}

// Verify Client implements domain.JobEnqueuer.
var _ domain.JobEnqueuer = (*Client)(nil)
