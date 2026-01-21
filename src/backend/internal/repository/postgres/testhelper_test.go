package postgres_test

import (
	"context"
	"testing"

	"github.com/uptrace/bun"

	"backend/internal/config"
	"backend/internal/infrastructure/database"
)

// setupTestDB creates a database connection for testing.
// It returns the database and a cleanup function.
func setupTestDB(t *testing.T) (*bun.DB, func()) {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	cleanup := func() {
		db.Close() //nolint:errcheck,gosec // Best effort cleanup in test
	}

	return db, cleanup
}

// cleanupTestData removes test data from all tables.
// Call this at the start of tests to ensure a clean state.
func cleanupTestData(t *testing.T, db *bun.DB) {
	t.Helper()

	ctx := context.Background()

	// Delete in reverse order of dependencies
	_, err := db.NewDelete().TableExpr("reference_letters").Where("1=1").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to clean reference_letters: %v", err)
	}

	_, err = db.NewDelete().TableExpr("files").Where("1=1").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to clean files: %v", err)
	}

	_, err = db.NewDelete().TableExpr("users").Where("1=1").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to clean users: %v", err)
	}
}
