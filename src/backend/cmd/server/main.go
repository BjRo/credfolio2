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
	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/config"
	"backend/internal/domain"
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
	resumeRepo := postgres.NewResumeRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	profileExpRepo := postgres.NewProfileExperienceRepository(db)
	profileEduRepo := postgres.NewProfileEducationRepository(db)
	profileSkillRepo := postgres.NewProfileSkillRepository(db)

	// Ensure demo user exists (development convenience)
	if seedErr := ensureDemoUser(context.Background(), userRepo, log); seedErr != nil {
		log.Warning("Failed to ensure demo user exists", logger.Feature("seed"), logger.Err(seedErr))
	}

	// Create LLM extractor with provider registry for per-operation chains
	extractor, extractHandler := createLLMExtractor(cfg, log)

	// Create job queue workers
	workers := river.NewWorkers()
	river.AddWorker(workers, job.NewDocumentProcessingWorker(refLetterRepo, fileStorage, log))

	// Register resume processing worker only if LLM is configured
	if extractor != nil {
		river.AddWorker(workers, job.NewResumeProcessingWorker(resumeRepo, fileRepo, profileRepo, profileExpRepo, profileEduRepo, profileSkillRepo, fileStorage, extractor, log))
		log.Info("Resume processing worker registered", logger.Feature("jobs"))
	} else {
		log.Warning("Resume processing worker not registered (LLM not configured)", logger.Feature("jobs"))
	}

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
	r.Handle("/graphql", graphql.NewHandler(userRepo, fileRepo, refLetterRepo, resumeRepo, profileRepo, profileExpRepo, profileEduRepo, profileSkillRepo, fileStorage, queueClient, log))
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

// demoUserID is the well-known ID for the demo user used in development.
var demoUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

// userCreator is the interface needed for creating users.
type userCreator interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

// createProviderRegistry creates all available LLM providers and returns a registry.
// Returns the registry, the default provider (based on LLM_PROVIDER config), and list of provider names.
func createProviderRegistry(cfg *config.Config, log logger.Logger) (*llm.ProviderRegistry, domain.LLMProvider, []string) {
	registry := llm.NewProviderRegistry()
	var providerNames []string
	var defaultProvider domain.LLMProvider

	// Register Anthropic if API key is available
	if cfg.Anthropic.APIKey != "" {
		provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
			APIKey: cfg.Anthropic.APIKey,
		})
		registry.Register("anthropic", provider)
		providerNames = append(providerNames, "anthropic")
		log.Debug("Registered Anthropic provider", logger.Feature("llm"))
	}

	// Register OpenAI if API key is available
	if cfg.OpenAI.APIKey != "" {
		provider := llm.NewOpenAIProvider(llm.OpenAIConfig{
			APIKey: cfg.OpenAI.APIKey,
		})
		registry.Register("openai", provider)
		providerNames = append(providerNames, "openai")
		log.Debug("Registered OpenAI provider", logger.Feature("llm"))
	}

	// Determine default provider based on config
	selectedProvider := cfg.LLM.Provider
	if selectedProvider == "" {
		selectedProvider = "anthropic" // Default
	}

	if p, ok := registry.Get(selectedProvider); ok {
		defaultProvider = p
	} else if len(providerNames) > 0 {
		// Fall back to first available provider
		defaultProvider, _ = registry.Get(providerNames[0])
		log.Warning("Configured provider not available, falling back",
			logger.Feature("llm"),
			logger.String("configured", selectedProvider),
			logger.String("fallback", providerNames[0]))
	}

	return registry, defaultProvider, providerNames
}

// createLLMExtractor creates the document extractor with per-operation provider chains.
// Returns the extractor (nil if no providers available) and the HTTP handler.
func createLLMExtractor(cfg *config.Config, log logger.Logger) (*llm.DocumentExtractor, http.Handler) {
	registry, defaultProvider, providerNames := createProviderRegistry(cfg, log)
	if defaultProvider == nil {
		log.Warning("LLM extraction disabled (no API key configured)", logger.Feature("llm"))
		return nil, handler.NewExtractUnavailableHandler()
	}

	// Wrap default provider with resilience
	resilientProvider := llm.NewResilientProvider(defaultProvider, llm.ResilientConfig{
		RequestTimeout: 120 * time.Second, // Extraction can be slow for large docs
	})

	// Configure provider chains for each operation
	// Document extraction uses the default provider (good for vision/PDF)
	// Resume extraction can use a different provider (better for structured output)
	docProvider := cfg.LLM.Provider
	if docProvider == "" {
		docProvider = "anthropic"
	}
	docChain := llm.ProviderChain{{Provider: docProvider}}

	// Resume extraction uses dedicated provider if configured, otherwise falls back to default
	resumeProvider := cfg.LLM.ResumeExtractionProvider
	if resumeProvider == "" {
		resumeProvider = docProvider
	}
	resumeChain := llm.ProviderChain{{Provider: resumeProvider}}

	// Log which providers are being used for each operation
	log.Info("Configured extraction providers",
		logger.Feature("llm"),
		logger.String("document_extraction", docProvider),
		logger.String("resume_extraction", resumeProvider),
	)

	extractor := llm.NewDocumentExtractor(resilientProvider, llm.DocumentExtractorConfig{
		ProviderRegistry:        registry,
		DocumentExtractionChain: docChain,
		ResumeExtractionChain:   resumeChain,
	})
	extractHandler := handler.NewExtractHandler(extractor, log)
	log.Info("LLM extraction enabled", logger.Feature("llm"), logger.String("providers", fmt.Sprintf("%v", providerNames)))

	return extractor, extractHandler
}

// ensureDemoUser creates the demo user if it doesn't exist.
// This provides a reliable way to have a demo user for development/testing
// that doesn't depend on migrations being in a specific state.
func ensureDemoUser(ctx context.Context, repo userCreator, log logger.Logger) error {
	// Check if demo user already exists
	existing, err := repo.GetByID(ctx, demoUserID)
	if err != nil {
		return fmt.Errorf("failed to check for demo user: %w", err)
	}
	if existing != nil {
		log.Debug("Demo user already exists", logger.Feature("seed"))
		return nil
	}

	// Create demo user
	name := "Demo User"
	demoUser := &domain.User{
		ID:           demoUserID,
		Email:        "demo@example.com",
		PasswordHash: "demo_hash", // Not a real hash - demo user only
		Name:         &name,
	}

	if err := repo.Create(ctx, demoUser); err != nil {
		return fmt.Errorf("failed to create demo user: %w", err)
	}

	log.Info("Demo user created", logger.Feature("seed"), logger.String("user_id", demoUserID.String()))
	return nil
}
