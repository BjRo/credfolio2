package database_test

import (
	"context"
	"testing"
	"time"

	"backend/internal/config"
	"backend/internal/infrastructure/database"
)

func TestConnect_Success(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close() //nolint:errcheck // Best effort cleanup in test

	// Verify connection works by pinging
	if err := db.PingContext(context.Background()); err != nil {
		t.Errorf("ping failed: %v", err)
	}
}

func TestConnect_InvalidConfig(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "invalid-host-that-does-not-exist",
		Port:     5432,
		User:     "invalid",
		Password: "invalid",
		Name:     "invalid",
		SSLMode:  "disable",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Short timeout
	defer cancel()

	_, err := database.Connect(ctx, cfg)
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	}
}
