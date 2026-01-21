package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestFileRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	ctx := context.Background()

	// Create a user first (files require a user)
	user := &domain.User{
		Email:        "fileuser@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/test.pdf",
	}

	err := fileRepo.Create(ctx, file)
	if err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	if file.ID == uuid.Nil {
		t.Error("expected file ID to be set after create")
	}
}

func TestFileRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	ctx := context.Background()

	// Create user and file
	user := &domain.User{
		Email:        "filegetbyid@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/getbyid.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	// Retrieve by ID
	found, err := fileRepo.GetByID(ctx, file.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find file, got nil")
	}

	if found.Filename != file.Filename {
		t.Errorf("filename mismatch: got %q, want %q", found.Filename, file.Filename)
	}

	// Test not found
	notFound, err := fileRepo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID for non-existent file failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent file")
	}
}

func TestFileRepository_GetByUserID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Email:        "filegetbyuser@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create multiple files
	for i := 0; i < 3; i++ {
		file := &domain.File{
			UserID:      user.ID,
			Filename:    "test.pdf",
			ContentType: "application/pdf",
			SizeBytes:   int64(12345 + i),
			StorageKey:  uuid.New().String(), // Unique storage key
		}
		if err := fileRepo.Create(ctx, file); err != nil {
			t.Fatalf("Create file %d failed: %v", i, err)
		}
	}

	// Retrieve by user ID
	files, err := fileRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d", len(files))
	}

	// Test empty result for non-existent user
	empty, err := fileRepo.GetByUserID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByUserID for non-existent user failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected 0 files, got %d", len(empty))
	}
}

func TestFileRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	ctx := context.Background()

	// Create user and file
	user := &domain.User{
		Email:        "filedelete@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/delete.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	// Delete the file
	if err := fileRepo.Delete(ctx, file.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	found, err := fileRepo.GetByID(ctx, file.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected file to be deleted")
	}
}
