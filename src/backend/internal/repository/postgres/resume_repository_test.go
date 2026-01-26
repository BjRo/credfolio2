package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestResumeRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	resumeRepo := postgres.NewResumeRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:        "resumeuser@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create a file
	file := &domain.File{
		UserID:      user.ID,
		Filename:    "resume.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/resume.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	resume := &domain.Resume{
		UserID: user.ID,
		FileID: file.ID,
		Status: domain.ResumeStatusPending,
	}

	err := resumeRepo.Create(ctx, resume)
	if err != nil {
		t.Fatalf("Create resume failed: %v", err)
	}

	if resume.ID == uuid.Nil {
		t.Error("expected resume ID to be set after create")
	}
}

func TestResumeRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	resumeRepo := postgres.NewResumeRepository(db)
	ctx := context.Background()

	// Create user and file
	user := &domain.User{
		Email:        "resumegetbyid@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "resume.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/resume_getbyid.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	resume := &domain.Resume{
		UserID: user.ID,
		FileID: file.ID,
		Status: domain.ResumeStatusPending,
	}
	if err := resumeRepo.Create(ctx, resume); err != nil {
		t.Fatalf("Create resume failed: %v", err)
	}

	// Retrieve by ID
	found, err := resumeRepo.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find resume, got nil")
	}

	if found.Status != resume.Status {
		t.Errorf("status mismatch: got %q, want %q", found.Status, resume.Status)
	}

	// Test not found
	notFound, err := resumeRepo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID for non-existent resume failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent resume")
	}
}

func TestResumeRepository_GetByUserID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	resumeRepo := postgres.NewResumeRepository(db)
	ctx := context.Background()

	// Create user
	user := &domain.User{
		Email:        "resumegetbyuser@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create multiple resumes
	for i := 0; i < 3; i++ {
		file := &domain.File{
			UserID:      user.ID,
			Filename:    "resume.pdf",
			ContentType: "application/pdf",
			SizeBytes:   12345,
			StorageKey:  "users/" + user.ID.String() + "/resume_" + string(rune('a'+i)) + ".pdf",
		}
		if err := fileRepo.Create(ctx, file); err != nil {
			t.Fatalf("Create file %d failed: %v", i, err)
		}

		resume := &domain.Resume{
			UserID: user.ID,
			FileID: file.ID,
			Status: domain.ResumeStatusPending,
		}
		if err := resumeRepo.Create(ctx, resume); err != nil {
			t.Fatalf("Create resume %d failed: %v", i, err)
		}
	}

	// Retrieve by user ID
	resumes, err := resumeRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByUserID failed: %v", err)
	}

	if len(resumes) != 3 {
		t.Errorf("expected 3 resumes, got %d", len(resumes))
	}

	// Test empty result for non-existent user
	empty, err := resumeRepo.GetByUserID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByUserID for non-existent user failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected 0 resumes, got %d", len(empty))
	}
}

func TestResumeRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	resumeRepo := postgres.NewResumeRepository(db)
	ctx := context.Background()

	// Create user and file
	user := &domain.User{
		Email:        "resumeupdate@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "resume.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/resume_update.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	resume := &domain.Resume{
		UserID: user.ID,
		FileID: file.ID,
		Status: domain.ResumeStatusPending,
	}
	if err := resumeRepo.Create(ctx, resume); err != nil {
		t.Fatalf("Create resume failed: %v", err)
	}

	// Update the resume with extracted data
	resume.Status = domain.ResumeStatusCompleted
	extractedData := domain.ResumeExtractedData{
		Name:   "John Doe",
		Skills: []string{"Go", "PostgreSQL", "Docker"},
	}
	var jsonErr error
	resume.ExtractedData, jsonErr = json.Marshal(extractedData)
	if jsonErr != nil {
		t.Fatalf("failed to marshal extracted data: %v", jsonErr)
	}

	if err := resumeRepo.Update(ctx, resume); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the update
	found, err := resumeRepo.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.Status != domain.ResumeStatusCompleted {
		t.Errorf("status not updated: got %q, want %q", found.Status, domain.ResumeStatusCompleted)
	}

	// Verify extracted data
	var foundData domain.ResumeExtractedData
	if err := json.Unmarshal(found.ExtractedData, &foundData); err != nil {
		t.Fatalf("failed to unmarshal extracted data: %v", err)
	}

	if foundData.Name != "John Doe" {
		t.Errorf("name not updated: got %q, want %q", foundData.Name, "John Doe")
	}

	if len(foundData.Skills) != 3 {
		t.Errorf("skills count mismatch: got %d, want 3", len(foundData.Skills))
	}
}

func TestResumeRepository_UpdateWithError(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	resumeRepo := postgres.NewResumeRepository(db)
	ctx := context.Background()

	// Create user and file
	user := &domain.User{
		Email:        "resumeerror@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "resume.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/resume_error.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	resume := &domain.Resume{
		UserID: user.ID,
		FileID: file.ID,
		Status: domain.ResumeStatusPending,
	}
	if err := resumeRepo.Create(ctx, resume); err != nil {
		t.Fatalf("Create resume failed: %v", err)
	}

	// Update with error
	resume.Status = domain.ResumeStatusFailed
	errMsg := "LLM extraction failed: timeout"
	resume.ErrorMessage = &errMsg

	if err := resumeRepo.Update(ctx, resume); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify error message
	found, err := resumeRepo.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.Status != domain.ResumeStatusFailed {
		t.Errorf("status not updated: got %q, want %q", found.Status, domain.ResumeStatusFailed)
	}

	if found.ErrorMessage == nil || *found.ErrorMessage != errMsg {
		t.Errorf("error message not updated: got %v, want %q", found.ErrorMessage, errMsg)
	}
}

func TestResumeRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	resumeRepo := postgres.NewResumeRepository(db)
	ctx := context.Background()

	// Create user and file
	user := &domain.User{
		Email:        "resumedelete@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	file := &domain.File{
		UserID:      user.ID,
		Filename:    "resume.pdf",
		ContentType: "application/pdf",
		SizeBytes:   12345,
		StorageKey:  "users/" + user.ID.String() + "/resume_delete.pdf",
	}
	if err := fileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create file failed: %v", err)
	}

	resume := &domain.Resume{
		UserID: user.ID,
		FileID: file.ID,
		Status: domain.ResumeStatusPending,
	}
	if err := resumeRepo.Create(ctx, resume); err != nil {
		t.Fatalf("Create resume failed: %v", err)
	}

	// Delete the resume
	if err := resumeRepo.Delete(ctx, resume.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	found, err := resumeRepo.GetByID(ctx, resume.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found != nil {
		t.Error("expected resume to be deleted")
	}
}
