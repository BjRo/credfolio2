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
	authorRepo := postgres.NewAuthorRepository(db)
	testimonialRepo := postgres.NewTestimonialRepository(db)
	skillValidationRepo := postgres.NewSkillValidationRepository(db)
	expValidationRepo := postgres.NewExperienceValidationRepository(db)

	// Ensure demo user exists (development convenience)
	if seedErr := ensureDemoUser(context.Background(), userRepo, log); seedErr != nil {
		log.Warning("Failed to ensure demo user exists", logger.Feature("seed"), logger.Err(seedErr))
	}

	// Create LLM extractor with provider registry for per-operation chains
	extractor, extractHandler, btTracing := createLLMExtractor(cfg, log)
	if btTracing != nil {
		defer func() {
			if shutdownErr := btTracing.Shutdown(context.Background()); shutdownErr != nil {
				log.Error("Failed to shutdown Braintrust tracing", logger.Feature("llm"), logger.Err(shutdownErr))
			}
		}()
	}

	// Create job queue workers
	workers := river.NewWorkers()

	// Register processing workers only if LLM is configured
	if extractor != nil {
		river.AddWorker(workers, job.NewResumeProcessingWorker(resumeRepo, fileRepo, profileRepo, profileExpRepo, profileEduRepo, profileSkillRepo, fileStorage, extractor, log))
		log.Info("Resume processing worker registered", logger.Feature("jobs"))

		river.AddWorker(workers, job.NewReferenceLetterProcessingWorker(refLetterRepo, fileRepo, profileRepo, profileSkillRepo, authorRepo, testimonialRepo, skillValidationRepo, fileStorage, extractor, log))
		log.Info("Reference letter processing worker registered", logger.Feature("jobs"))
	} else {
		log.Warning("Processing workers not registered (LLM not configured)", logger.Feature("jobs"))
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
	r.Handle("/graphql", graphql.NewHandler(userRepo, fileRepo, refLetterRepo, resumeRepo, profileRepo, profileExpRepo, profileEduRepo, profileSkillRepo, authorRepo, testimonialRepo, skillValidationRepo, expValidationRepo, fileStorage, queueClient, log))
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
// Returns the registry, list of provider names, and Braintrust tracing.
func createProviderRegistry(cfg *config.Config, log logger.Logger) (*llm.ProviderRegistry, []string, *llm.BraintrustTracing) {
	registry := llm.NewProviderRegistry()
	var providerNames []string

	// Initialize Braintrust tracing if configured
	btTracing, err := llm.NewBraintrustTracing(llm.BraintrustConfig{
		APIKey:  cfg.Braintrust.APIKey,
		Project: cfg.Braintrust.Project,
	}, log)
	if err != nil {
		log.Warning("Failed to initialize Braintrust tracing", logger.Feature("llm"), logger.Err(err))
	} else if btTracing != nil {
		log.Info("Braintrust tracing enabled", logger.Feature("llm"), logger.String("project", btTracing.Project()))
	}

	// Register Anthropic if API key is available
	if cfg.Anthropic.APIKey != "" {
		providerConfig := llm.AnthropicConfig{
			APIKey: cfg.Anthropic.APIKey,
		}
		// Add Braintrust middleware if tracing is enabled
		if btTracing != nil {
			providerConfig.Middleware = btTracing.AnthropicMiddleware() //nolint:bodyclose // middleware, not response
		}
		provider := llm.NewAnthropicProvider(providerConfig)
		registry.Register("anthropic", provider)
		providerNames = append(providerNames, "anthropic")
		log.Debug("Registered Anthropic provider", logger.Feature("llm"))
	}

	// Register OpenAI if API key is available
	if cfg.OpenAI.APIKey != "" {
		providerConfig := llm.OpenAIConfig{
			APIKey: cfg.OpenAI.APIKey,
		}
		// Add Braintrust middleware if tracing is enabled
		if btTracing != nil {
			providerConfig.Middleware = btTracing.OpenAIMiddleware() //nolint:bodyclose // middleware, not response
		}
		provider := llm.NewOpenAIProvider(providerConfig)
		registry.Register("openai", provider)
		providerNames = append(providerNames, "openai")
		log.Debug("Registered OpenAI provider", logger.Feature("llm"))
	}

	return registry, providerNames, btTracing
}

// createLLMExtractor creates the document extractor with per-operation provider chains.
// Returns the extractor (nil if no providers available), the HTTP handler, and the Braintrust tracing instance.
func createLLMExtractor(cfg *config.Config, log logger.Logger) (*llm.DocumentExtractor, http.Handler, *llm.BraintrustTracing) {
	registry, providerNames, btTracing := createProviderRegistry(cfg, log)

	// Parse per-use-case model configuration
	docProvider, docModel := cfg.LLM.ParseDocumentExtractionModel()
	resumeProvider, resumeModel := cfg.LLM.ParseResumeExtractionModel()

	// Determine a default provider from the document extraction chain config.
	// This serves as the fallback when a chain references an unregistered provider.
	var defaultProvider domain.LLMProvider
	if p, ok := registry.Get(docProvider); ok {
		defaultProvider = p
	} else if len(providerNames) > 0 {
		defaultProvider, _ = registry.Get(providerNames[0])
		log.Warning("Document extraction provider not available, falling back",
			logger.Feature("llm"),
			logger.String("configured", docProvider),
			logger.String("fallback", providerNames[0]))
	}

	if defaultProvider == nil {
		log.Warning("LLM extraction disabled (no API key configured)", logger.Feature("llm"))
		return nil, handler.NewExtractUnavailableHandler(), btTracing
	}

	// Wrap default provider with resilience
	resilientProvider := llm.NewResilientProvider(defaultProvider, llm.ResilientConfig{
		RequestTimeout: 120 * time.Second, // Extraction can be slow for large docs
	})

	// Configure provider chains for each operation
	docChain := llm.ProviderChain{{Provider: docProvider, Model: docModel}}
	resumeChain := llm.ProviderChain{{Provider: resumeProvider, Model: resumeModel}}

	// Validate that configured chains reference registered providers
	if _, ok := registry.Get(docProvider); !ok {
		log.Warning("Document extraction chain references unregistered provider — will fall back to default",
			logger.Feature("llm"),
			logger.String("provider", docProvider),
			logger.String("registered", fmt.Sprintf("%v", providerNames)),
		)
	}
	if _, ok := registry.Get(resumeProvider); !ok {
		log.Warning("Resume extraction chain references unregistered provider — will fall back to default",
			logger.Feature("llm"),
			logger.String("provider", resumeProvider),
			logger.String("registered", fmt.Sprintf("%v", providerNames)),
		)
	}

	// Log which providers are being used for each operation
	log.Info("Configured extraction providers",
		logger.Feature("llm"),
		logger.String("document_extraction", fmt.Sprintf("%s/%s", docProvider, docModel)),
		logger.String("resume_extraction", fmt.Sprintf("%s/%s", resumeProvider, resumeModel)),
	)

	extractor := llm.NewDocumentExtractor(resilientProvider, llm.DocumentExtractorConfig{
		ProviderRegistry:        registry,
		DocumentExtractionChain: docChain,
		ResumeExtractionChain:   resumeChain,
		Logger:                  log,
	})
	extractHandler := handler.NewExtractHandler(extractor, log)
	log.Info("LLM extraction enabled", logger.Feature("llm"), logger.String("providers", fmt.Sprintf("%v", providerNames)))

	return extractor, extractHandler, btTracing
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
