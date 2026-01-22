// Package main provides the entry point for the backend server.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riverqueue/river"

	"backend/internal/config"
	"backend/internal/graphql"
	"backend/internal/handler"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/queue"
	"backend/internal/infrastructure/storage"
	"backend/internal/job"
	"backend/internal/repository/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to database
	db, err := database.Connect(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close() //nolint:errcheck // Best effort cleanup on shutdown

	log.Printf("Connected to database: %s", cfg.Database.Name)

	// Create storage
	fileStorage, err := storage.NewMinIOStorage(cfg.MinIO)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}

	// Ensure bucket exists
	if bucketErr := fileStorage.EnsureBucket(context.Background()); bucketErr != nil {
		return fmt.Errorf("failed to ensure storage bucket: %w", bucketErr)
	}

	log.Printf("Connected to storage: %s/%s", cfg.MinIO.Endpoint, cfg.MinIO.Bucket)

	// Create repositories (needed for workers)
	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	refLetterRepo := postgres.NewReferenceLetterRepository(db)

	// Create job queue workers
	workers := river.NewWorkers()
	river.AddWorker(workers, job.NewDocumentProcessingWorker(refLetterRepo, fileStorage))

	// Create job queue client
	queueClient, err := queue.NewClient(context.Background(), cfg.Database, cfg.Queue, workers)
	if err != nil {
		return fmt.Errorf("failed to create queue client: %w", err)
	}
	defer queueClient.Close()

	// Start job queue processing
	if err := queueClient.Start(context.Background()); err != nil {
		return fmt.Errorf("failed to start queue client: %w", err)
	}

	log.Printf("Job queue started with %d max workers", cfg.Queue.MaxWorkers)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Routes
	r.Get("/", handler.NewRootHandler().ServeHTTP)
	r.Get("/health", handler.NewHealthHandler(db).ServeHTTP)

	// GraphQL API
	r.Handle("/graphql", graphql.NewHandler(userRepo, fileRepo, refLetterRepo, fileStorage, queueClient))
	r.Get("/playground", graphql.NewPlaygroundHandler("/graphql").ServeHTTP)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s\n", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		log.Printf("Received signal %s, shutting down...", sig)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.WriteTimeout)
	defer cancel()

	// Stop queue processing
	if err := queueClient.Stop(ctx); err != nil {
		log.Printf("Error stopping queue client: %v", err)
	}

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
