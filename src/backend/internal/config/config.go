// Package config handles application configuration loading.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Database DatabaseConfig
	MinIO    MinIOConfig
	Server   ServerConfig
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
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Load reads configuration from environment variables.
// It uses sensible defaults matching docker-compose.yml for local development.
func Load() (*Config, error) {
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

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("POSTGRES_USER", "credfolio"),
			Password: getEnv("POSTGRES_PASSWORD", "credfolio_dev"),
			Name:     getEnv("POSTGRES_DB", "credfolio"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
			url:      os.Getenv("DATABASE_URL"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ROOT_USER", "minioadmin"),
			SecretKey: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
			UseSSL:    useSSL,
			Bucket:    getEnv("MINIO_BUCKET", "credfolio"),
		},
		Server: ServerConfig{
			Port:         serverPort,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	return cfg, nil
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
