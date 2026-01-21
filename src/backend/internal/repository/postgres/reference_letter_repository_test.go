package postgres_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestReferenceLetterRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:        "letteruser@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID:       user.ID,
		Title:        strPtr("Reference for John Doe"),
		AuthorName:   strPtr("Jane Smith"),
		AuthorTitle:  strPtr("Senior Manager"),
		Organization: strPtr("Acme Corp"),
		Status:       domain.ReferenceLetterStatusPending,
	}

	err := letterRepo.Create(ctx, letter)
	if err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	if letter.ID == uuid.Nil {
		t.Error("expected letter ID to be set after create")
	}
}

func TestReferenceLetterRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	ctx := context.Background()

	// Create user and letter
	user := &domain.User{
		Email:        "lettergetbyid@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID:     user.ID,
		Title:      strPtr("Reference Letter"),
		AuthorName: strPtr("Author"),
		Status:     domain.ReferenceLetterStatusPending,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Retrieve by ID
	found, err := letterRepo.GetByID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find letter, got nil")
	}

	if *found.Title != *letter.Title {
		t.Errorf("title mismatch: got %q, want %q", *found.Title, *letter.Title)
	}

	// Test not found
	notFound, err := letterRepo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID for non-existent letter failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent letter")
	}
}

func TestReferenceLetterRepository_GetByUserID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Email:        "lettergetbyuser@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create multiple letters
	for i := 0; i < 3; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Title:  strPtr("Letter " + string(rune('A'+i))),
			Status: domain.ReferenceLetterStatusPending,
		}
		if err := letterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("Create letter %d failed: %v", i, err)
		}
	}

	// Retrieve by user ID
	letters, err := letterRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}

	if len(letters) != 3 {
		t.Errorf("expected 3 letters, got %d", len(letters))
	}

	// Test empty result for non-existent user
	empty, err := letterRepo.GetByUserID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByUserID for non-existent user failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected 0 letters, got %d", len(empty))
	}
}

func TestReferenceLetterRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	ctx := context.Background()

	// Create user and letter
	user := &domain.User{
		Email:        "letterupdate@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Title:  strPtr("Original Title"),
		Status: domain.ReferenceLetterStatusPending,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Update the letter
	letter.Status = domain.ReferenceLetterStatusCompleted
	letter.RawText = strPtr("Extracted text content")
	extractedData := map[string]string{"skill": "Go programming"}
	var jsonErr error
	letter.ExtractedData, jsonErr = json.Marshal(extractedData)
	if jsonErr != nil {
		t.Fatalf("failed to marshal extracted data: %v", jsonErr)
	}

	if err := letterRepo.Update(ctx, letter); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the update
	found, err := letterRepo.GetByID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("status not updated: got %q, want %q", found.Status, domain.ReferenceLetterStatusCompleted)
	}

	if found.RawText == nil || *found.RawText != "Extracted text content" {
		t.Errorf("raw text not updated: got %v", found.RawText)
	}
}

func TestReferenceLetterRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	ctx := context.Background()

	// Create user and letter
	user := &domain.User{
		Email:        "letterdelete@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Title:  strPtr("To Delete"),
		Status: domain.ReferenceLetterStatusPending,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Delete the letter
	if err := letterRepo.Delete(ctx, letter.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	found, err := letterRepo.GetByID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected letter to be deleted")
	}
}

func TestReferenceLetterRepository_WithFile(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Email:        "letterwithfile@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create file
	file := &domain.File{
		UserID:      user.ID,
		Filename:    "reference.pdf",
		ContentType: "application/pdf",
		SizeBytes:   54321,
		StorageKey:  "users/" + user.ID.String() + "/reference.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	// Create letter with file reference
	dateWritten := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	letter := &domain.ReferenceLetter{
		UserID:      user.ID,
		FileID:      &file.ID,
		Title:       strPtr("Reference from Prof. Smith"),
		AuthorName:  strPtr("Prof. John Smith"),
		DateWritten: &dateWritten,
		Status:      domain.ReferenceLetterStatusPending,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Retrieve and verify file reference
	found, err := letterRepo.GetByID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.FileID == nil {
		t.Error("expected file ID to be set")
	} else if *found.FileID != file.ID {
		t.Errorf("file ID mismatch: got %s, want %s", *found.FileID, file.ID)
	}
}
