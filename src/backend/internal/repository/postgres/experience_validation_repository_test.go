package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestExperienceValidationRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	experienceRepo := postgres.NewProfileExperienceRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewExperienceValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, experience, and reference letter
	user := &domain.User{
		Email:        "expvalidation@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	profile := &domain.Profile{
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("Create profile failed: %v", err)
	}

	experience := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Acme Corp",
		Title:        "Senior Engineer",
		DisplayOrder: 0,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, experience); err != nil {
		t.Fatalf("Create experience failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create validation
	validation := &domain.ExperienceValidation{
		ProfileExperienceID: experience.ID,
		ReferenceLetterID:   letter.ID,
		QuoteSnippet:        strPtr("During her time as Senior Engineer at Acme Corp..."),
	}

	err := validationRepo.Create(ctx, validation)
	if err != nil {
		t.Fatalf("Create validation failed: %v", err)
	}

	if validation.ID == uuid.Nil {
		t.Error("expected validation ID to be set after create")
	}
}

func TestExperienceValidationRepository_GetByProfileExperienceID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	experienceRepo := postgres.NewProfileExperienceRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewExperienceValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, experience
	user := &domain.User{
		Email:        "expvalidationget@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	profile := &domain.Profile{
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("Create profile failed: %v", err)
	}

	experience := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Test Corp",
		Title:        "Manager",
		DisplayOrder: 0,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, experience); err != nil {
		t.Fatalf("Create experience failed: %v", err)
	}

	// Create multiple reference letters with validations
	for i := 0; i < 2; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if err := letterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("Create letter %d failed: %v", i, err)
		}

		validation := &domain.ExperienceValidation{
			ProfileExperienceID: experience.ID,
			ReferenceLetterID:   letter.ID,
			QuoteSnippet:        strPtr("Experience quote " + string(rune('A'+i))),
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation %d failed: %v", i, err)
		}
	}

	// Retrieve by experience ID
	validations, err := validationRepo.GetByProfileExperienceID(ctx, experience.ID)
	if err != nil {
		t.Fatalf("GetByProfileExperienceID failed: %v", err)
	}

	if len(validations) != 2 {
		t.Errorf("expected 2 validations, got %d", len(validations))
	}
}

func TestExperienceValidationRepository_CountByProfileExperienceID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	experienceRepo := postgres.NewProfileExperienceRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewExperienceValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, experience
	user := &domain.User{
		Email:        "expvalidationcount@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	profile := &domain.Profile{
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("Create profile failed: %v", err)
	}

	experience := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Count Corp",
		Title:        "Developer",
		DisplayOrder: 0,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, experience); err != nil {
		t.Fatalf("Create experience failed: %v", err)
	}

	// Verify count is 0
	count, err := validationRepo.CountByProfileExperienceID(ctx, experience.ID)
	if err != nil {
		t.Fatalf("CountByProfileExperienceID failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 validations, got %d", count)
	}

	// Add validations
	for i := 0; i < 3; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if createErr := letterRepo.Create(ctx, letter); createErr != nil {
			t.Fatalf("Create letter %d failed: %v", i, createErr)
		}

		validation := &domain.ExperienceValidation{
			ProfileExperienceID: experience.ID,
			ReferenceLetterID:   letter.ID,
		}
		if createErr := validationRepo.Create(ctx, validation); createErr != nil {
			t.Fatalf("Create validation %d failed: %v", i, createErr)
		}
	}

	// Verify count is 3
	count, err = validationRepo.CountByProfileExperienceID(ctx, experience.ID)
	if err != nil {
		t.Fatalf("CountByProfileExperienceID failed: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 validations, got %d", count)
	}
}

func TestExperienceValidationRepository_DeleteByReferenceLetterID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	experienceRepo := postgres.NewProfileExperienceRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewExperienceValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, experience
	user := &domain.User{
		Email:        "expvalidationdelete@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	profile := &domain.Profile{
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("Create profile failed: %v", err)
	}

	experience := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Delete Corp",
		Title:        "Tester",
		DisplayOrder: 0,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, experience); err != nil {
		t.Fatalf("Create experience failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create validations
	for i := 0; i < 2; i++ {
		exp := &domain.ProfileExperience{
			ProfileID:    profile.ID,
			Company:      "Company " + string(rune('A'+i)),
			Title:        "Title",
			DisplayOrder: i + 1,
			Source:       domain.ExperienceSourceManual,
		}
		if err := experienceRepo.Create(ctx, exp); err != nil {
			t.Fatalf("Create experience %d failed: %v", i, err)
		}

		validation := &domain.ExperienceValidation{
			ProfileExperienceID: exp.ID,
			ReferenceLetterID:   letter.ID,
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation %d failed: %v", i, err)
		}
	}

	// Verify validations exist
	validations, err := validationRepo.GetByReferenceLetterID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByReferenceLetterID failed: %v", err)
	}
	if len(validations) != 2 {
		t.Fatalf("expected 2 validations, got %d", len(validations))
	}

	// Delete all validations by reference letter ID
	if deleteErr := validationRepo.DeleteByReferenceLetterID(ctx, letter.ID); deleteErr != nil {
		t.Fatalf("DeleteByReferenceLetterID failed: %v", deleteErr)
	}

	// Verify deletion
	remaining, err := validationRepo.GetByReferenceLetterID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByReferenceLetterID failed: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("expected 0 validations after delete, got %d", len(remaining))
	}
}

func TestExperienceValidationRepository_BatchCountByProfileExperienceIDs(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	experienceRepo := postgres.NewProfileExperienceRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewExperienceValidationRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "expvalidationbatch@example.com",
		PasswordHash: "hashed_password",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	profile := &domain.Profile{
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("Create profile failed: %v", err)
	}

	// Create 3 experiences with different validation counts
	exp1 := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Batch Corp 1",
		Title:        "Engineer",
		DisplayOrder: 0,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, exp1); err != nil {
		t.Fatalf("Create exp1 failed: %v", err)
	}

	exp2 := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Batch Corp 2",
		Title:        "Manager",
		DisplayOrder: 1,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, exp2); err != nil {
		t.Fatalf("Create exp2 failed: %v", err)
	}

	exp3 := &domain.ProfileExperience{
		ProfileID:    profile.ID,
		Company:      "Batch Corp 3",
		Title:        "Director",
		DisplayOrder: 2,
		Source:       domain.ExperienceSourceManual,
	}
	if err := experienceRepo.Create(ctx, exp3); err != nil {
		t.Fatalf("Create exp3 failed: %v", err)
	}

	// Create validations: exp1=2, exp2=0, exp3=3
	for i := 0; i < 2; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if err := letterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("Create letter failed: %v", err)
		}
		validation := &domain.ExperienceValidation{
			ProfileExperienceID: exp1.ID,
			ReferenceLetterID:   letter.ID,
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation for exp1 failed: %v", err)
		}
	}

	for i := 0; i < 3; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if err := letterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("Create letter failed: %v", err)
		}
		validation := &domain.ExperienceValidation{
			ProfileExperienceID: exp3.ID,
			ReferenceLetterID:   letter.ID,
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation for exp3 failed: %v", err)
		}
	}

	// Batch count all experiences in one query
	counts, err := validationRepo.BatchCountByProfileExperienceIDs(ctx, []uuid.UUID{exp1.ID, exp2.ID, exp3.ID})
	if err != nil {
		t.Fatalf("BatchCountByProfileExperienceIDs failed: %v", err)
	}

	// Verify counts
	if counts[exp1.ID] != 2 {
		t.Errorf("expected exp1 count=2, got %d", counts[exp1.ID])
	}
	if counts[exp2.ID] != 0 {
		t.Errorf("expected exp2 count=0, got %d", counts[exp2.ID])
	}
	if counts[exp3.ID] != 3 {
		t.Errorf("expected exp3 count=3, got %d", counts[exp3.ID])
	}
}

func TestExperienceValidationRepository_BatchCountByProfileExperienceIDs_EmptyInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	validationRepo := postgres.NewExperienceValidationRepository(db)
	ctx := context.Background()

	// Empty input should return empty map
	counts, err := validationRepo.BatchCountByProfileExperienceIDs(ctx, []uuid.UUID{})
	if err != nil {
		t.Fatalf("BatchCountByProfileExperienceIDs failed: %v", err)
	}

	if len(counts) != 0 {
		t.Errorf("expected empty map, got %d entries", len(counts))
	}
}
