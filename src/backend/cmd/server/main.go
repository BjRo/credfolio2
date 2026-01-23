// Package main provides the entry point for the backend server.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riverqueue/river"

	"backend/internal/config"
	"backend/internal/graphql"
	"backend/internal/handler"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/llm"
	"backend/internal/infrastructure/queue"
	"backend/internal/infrastructure/storage"
	"backend/internal/job"
	"backend/internal/logger"
	"backend/internal/repository/postgres"
)

func main() {
	log := logger.NewStdoutLogger()

	if err := run(log); err != nil {
		log.Critical("Server error", logger.Feature("server"), logger.Err(err))
		os.Exit(1)
	}
}

func run(log logger.Logger) error {
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

	log.Info("Connected to database", logger.Feature("database"), logger.String("database", cfg.Database.Name))

	// Create storage
	fileStorage, err := storage.NewMinIOStorage(cfg.MinIO)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}

	// Ensure bucket exists
	if bucketErr := fileStorage.EnsureBucket(context.Background()); bucketErr != nil {
		return fmt.Errorf("failed to ensure storage bucket: %w", bucketErr)
	}

	log.Info("Connected to storage", logger.Feature("storage"), logger.String("endpoint", cfg.MinIO.Endpoint), logger.String("bucket", cfg.MinIO.Bucket))

	// Create repositories (needed for workers)
	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	refLetterRepo := postgres.NewReferenceLetterRepository(db)

	// Create job queue workers
	workers := river.NewWorkers()
	river.AddWorker(workers, job.NewDocumentProcessingWorker(refLetterRepo, fileStorage, log))

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

	log.Info("Job queue started", logger.Feature("queue"), logger.Int("max_workers", cfg.Queue.MaxWorkers))

	// Create LLM extraction handler (optional - only if API key configured)
	var extractHandler http.Handler
	if cfg.Anthropic.APIKey != "" {
		anthropicProvider := llm.NewAnthropicProvider(llm.AnthropicConfig{
			APIKey: cfg.Anthropic.APIKey,
		})
		resilientProvider := llm.NewResilientProvider(anthropicProvider, llm.ResilientConfig{
			RequestTimeout: 120 * time.Second, // Extraction can be slow for large docs
		})
		extractor := llm.NewDocumentExtractor(resilientProvider, llm.DocumentExtractorConfig{})
		extractHandler = handler.NewExtractHandler(extractor, log)
		log.Info("LLM extraction enabled", logger.Feature("llm"))
	} else {
		extractHandler = handler.NewExtractUnavailableHandler()
		log.Warning("LLM extraction disabled (ANTHROPIC_API_KEY not set)", logger.Feature("llm"))
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/", handler.NewRootHandler().ServeHTTP)
	r.Get("/health", handler.NewHealthHandler(db).ServeHTTP)

	// Document extraction API (for testing)
	r.Post("/api/extract", extractHandler.ServeHTTP)

	// GraphQL API
	r.Handle("/graphql", graphql.NewHandler(userRepo, fileRepo, refLetterRepo, fileStorage, queueClient, log))
	r.Get("/playground", graphql.NewPlaygroundHandler("/graphql").ServeHTTP)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("Server starting", logger.Feature("server"), logger.String("address", fmt.Sprintf("http://localhost%s", addr)), logger.Int("port", cfg.Server.Port))

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
		log.Info("Received shutdown signal", logger.Feature("server"), logger.String("signal", sig.String()))
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.WriteTimeout)
	defer cancel()

	// Stop queue processing
	if err := queueClient.Stop(ctx); err != nil {
		log.Error("Error stopping queue client", logger.Feature("queue"), logger.Err(err))
	}

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Info("Server stopped gracefully", logger.Feature("server"))
	return nil
}
