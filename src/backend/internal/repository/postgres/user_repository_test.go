package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         strPtr("Test User"),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == uuid.Nil {
		t.Error("expected user ID to be set after create")
	}

	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set after create")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:        "getbyid@example.com",
		PasswordHash: "hashed_password",
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Retrieve by ID
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find user, got nil")
	}

	if found.Email != user.Email {
		t.Errorf("email mismatch: got %q, want %q", found.Email, user.Email)
	}

	// Test not found
	notFound, err := repo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID for non-existent user failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent user")
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:        "getbyemail@example.com",
		PasswordHash: "hashed_password",
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Retrieve by email
	found, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find user, got nil")
	}

	if found.ID != user.ID {
		t.Errorf("ID mismatch: got %s, want %s", found.ID, user.ID)
	}

	// Test not found
	notFound, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	if err != nil {
		t.Fatalf("GetByEmail for non-existent user failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent email")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:        "update@example.com",
		PasswordHash: "hashed_password",
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update the user
	user.Name = strPtr("Updated Name")
	if err := repo.Update(ctx, user); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the update
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.Name == nil || *found.Name != "Updated Name" {
		t.Errorf("Name not updated: got %v", found.Name)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	repo := postgres.NewUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:        "delete@example.com",
		PasswordHash: "hashed_password",
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete the user
	if err := repo.Delete(ctx, user.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected user to be deleted")
	}
}

func strPtr(s string) *string {
	return &s
}
