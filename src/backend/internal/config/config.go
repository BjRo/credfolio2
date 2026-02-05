// Package config handles application configuration loading.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct { //nolint:govet // Field order prioritizes readability
	Environment string
	Database    DatabaseConfig
	MinIO       MinIOConfig
	Server      ServerConfig
	Queue       QueueConfig
	LLM         LLMConfig
	Anthropic   AnthropicConfig
	OpenAI      OpenAIConfig
	Braintrust  BraintrustConfig
}

// LLMConfig holds general LLM configuration.
type LLMConfig struct {
	// Provider specifies which LLM provider to use: "anthropic" or "openai".
	// Defaults to "anthropic" if not specified.
	Provider string

	// ResumeExtractionModel specifies the provider and model for resume data extraction.
	// Format: "provider/model" (e.g., "openai/gpt-4o").
	// Defaults to "openai/gpt-4o" for best structured output support.
	ResumeExtractionModel string
}

// ParseResumeExtractionModel parses the ResumeExtractionModel into provider and model parts.
// Returns (provider, model). If no "/" is present, treats the whole string as provider.
func (c *LLMConfig) ParseResumeExtractionModel() (provider, model string) {
	value := c.ResumeExtractionModel
	if value == "" {
		return "openai", "gpt-4o" // Default
	}

	// Split on first "/"
	for i, ch := range value {
		if ch == '/' {
			return value[:i], value[i+1:]
		}
	}

	// No "/" found, treat as provider only
	return value, ""
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	SSLMode  string
	url      string // Explicit URL from DATABASE_URL env var
	Port     int
}

// URL returns the PostgreSQL connection string.
// If DATABASE_URL was set, it takes precedence over individual settings.
func (c *DatabaseConfig) URL() string {
	if c.url != "" {
		return c.url
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

// MinIOConfig holds MinIO/S3 connection settings.
type MinIOConfig struct {
	Endpoint        string
	PublicEndpoint  string // External endpoint for presigned URLs (defaults to Endpoint if not set)
	StorageProxyURL string // If set, use proxy URLs instead of presigned URLs (e.g., "/api/storage")
	AccessKey       string
	SecretKey       string
	Bucket          string
	UseSSL          bool
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// QueueConfig holds River job queue settings.
type QueueConfig struct {
	MaxWorkers int
}

// AnthropicConfig holds Anthropic API settings.
// Note: Model selection is done per-request, not globally configured.
type AnthropicConfig struct {
	APIKey string
}

// OpenAIConfig holds OpenAI API settings.
// Note: Model selection is done per-request, not globally configured.
type OpenAIConfig struct {
	APIKey string
}

// BraintrustConfig holds Braintrust observability settings.
type BraintrustConfig struct {
	// APIKey is the Braintrust API key for sending traces.
	// If empty, Braintrust tracing is disabled.
	APIKey string

	// Project is the Braintrust project name for organizing traces.
	// Defaults to "credfolio" if not specified.
	Project string
}

// Load reads configuration from environment variables.
// It uses sensible defaults matching docker-compose.yml for local development.
func Load() (*Config, error) {
	// Load .env file if present (silently ignore if not found).
	// Try to find .env by walking up from current directory to project root.
	loadEnvFiles()

	dbPort, err := getEnvInt("POSTGRES_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}

	serverPort, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	useSSL, err := getEnvBool("MINIO_USE_SSL", false)
	if err != nil {
		return nil, fmt.Errorf("invalid MINIO_USE_SSL: %w", err)
	}

	queueMaxWorkers, err := getEnvInt("QUEUE_MAX_WORKERS", 10)
	if err != nil {
		return nil, fmt.Errorf("invalid QUEUE_MAX_WORKERS: %w", err)
	}

	// Environment determines database name: credfolio_dev or credfolio_test
	env := getEnv("CREDFOLIO_ENV", "dev")
	dbName := "credfolio_" + env

	// Default hosts use docker container names for devcontainer environment
	cfg := &Config{
		Environment: env,
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "credfolio2-postgres"),
			Port:     dbPort,
			User:     getEnv("POSTGRES_USER", "credfolio"),
			Password: getEnv("POSTGRES_PASSWORD", "credfolio_dev"),
			Name:     dbName,
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
			url:      os.Getenv("DATABASE_URL"),
		},
		MinIO: MinIOConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "credfolio2-minio:9000"),
			PublicEndpoint:  getEnv("MINIO_PUBLIC_ENDPOINT", "localhost:9000"),
			StorageProxyURL: getEnv("STORAGE_PROXY_URL", "/api/storage"), // Default to Next.js proxy for devcontainer
			AccessKey:       getEnv("MINIO_ROOT_USER", "minioadmin"),
			SecretKey:       getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
			UseSSL:          useSSL,
			Bucket:          getEnv("MINIO_BUCKET", "credfolio"),
		},
		Server: ServerConfig{
			Port:         serverPort,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 120 * time.Second, // Long timeout for LLM extraction requests
			IdleTimeout:  120 * time.Second,
		},
		Queue: QueueConfig{
			MaxWorkers: queueMaxWorkers,
		},
		LLM: LLMConfig{
			Provider:                 getEnv("LLM_PROVIDER", "anthropic"),
			ResumeExtractionModel: os.Getenv("RESUME_EXTRACTION_MODEL"), // Default: openai/gpt-4o
		},
		Anthropic: AnthropicConfig{
			APIKey: os.Getenv("ANTHROPIC_API_KEY"),
		},
		OpenAI: OpenAIConfig{
			APIKey: os.Getenv("OPENAI_API_KEY"),
		},
		Braintrust: BraintrustConfig{
			APIKey:  os.Getenv("BRAINTRUST_API_KEY"),
			Project: getEnv("BRAINTRUST_PROJECT", "credfolio"),
		},
	}

	return cfg, nil
}

// loadEnvFiles attempts to load .env files from the current directory
// and by walking up the directory tree to find the project root.
func loadEnvFiles() {
	var envPaths []string

	// Walk up the directory tree to find .env files
	if cwd, err := os.Getwd(); err == nil {
		dir := cwd
		for i := 0; i < 10; i++ { // Limit search depth
			envPath := filepath.Join(dir, ".env")
			if _, err := os.Stat(envPath); err == nil {
				envPaths = append(envPaths, envPath)
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break // Reached filesystem root
			}
			dir = parent
		}
	}

	// Load all found .env files (first one wins for each variable)
	if len(envPaths) > 0 {
		_ = godotenv.Load(envPaths...) //nolint:errcheck // Best effort, env vars may come from system
	}
}

// getEnv returns the environment variable value or a default.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns the environment variable as an int or a default.
func getEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}

// getEnvBool returns the environment variable as a bool or a default.
func getEnvBool(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseBool(value)
}
