package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestSkillValidationRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	skillRepo := postgres.NewProfileSkillRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewSkillValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, skill, and reference letter
	user := &domain.User{
		Email:        "skillvalidation@example.com",
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

	skill := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Go",
		NormalizedName: "go",
		Category:       "TECHNICAL",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill); err != nil {
		t.Fatalf("Create skill failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create validation
	validation := &domain.SkillValidation{
		ProfileSkillID:    skill.ID,
		ReferenceLetterID: letter.ID,
		QuoteSnippet:      strPtr("Her expertise in Go was exceptional..."),
	}

	err := validationRepo.Create(ctx, validation)
	if err != nil {
		t.Fatalf("Create validation failed: %v", err)
	}

	if validation.ID == uuid.Nil {
		t.Error("expected validation ID to be set after create")
	}
}

func TestSkillValidationRepository_GetByProfileSkillID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	skillRepo := postgres.NewProfileSkillRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewSkillValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, skill
	user := &domain.User{
		Email:        "skillvalidationget@example.com",
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

	skill := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Leadership",
		NormalizedName: "leadership",
		Category:       "SOFT",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill); err != nil {
		t.Fatalf("Create skill failed: %v", err)
	}

	// Create multiple reference letters with validations
	for i := 0; i < 3; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if err := letterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("Create letter %d failed: %v", i, err)
		}

		validation := &domain.SkillValidation{
			ProfileSkillID:    skill.ID,
			ReferenceLetterID: letter.ID,
			QuoteSnippet:      strPtr("Leadership quote " + string(rune('A'+i))),
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation %d failed: %v", i, err)
		}
	}

	// Retrieve by skill ID
	validations, err := validationRepo.GetByProfileSkillID(ctx, skill.ID)
	if err != nil {
		t.Fatalf("GetByProfileSkillID failed: %v", err)
	}

	if len(validations) != 3 {
		t.Errorf("expected 3 validations, got %d", len(validations))
	}
}

func TestSkillValidationRepository_CountByProfileSkillID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	skillRepo := postgres.NewProfileSkillRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewSkillValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, skill
	user := &domain.User{
		Email:        "skillvalidationcount@example.com",
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

	skill := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Python",
		NormalizedName: "python",
		Category:       "TECHNICAL",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill); err != nil {
		t.Fatalf("Create skill failed: %v", err)
	}

	// Verify count is 0
	count, err := validationRepo.CountByProfileSkillID(ctx, skill.ID)
	if err != nil {
		t.Fatalf("CountByProfileSkillID failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 validations, got %d", count)
	}

	// Add validations
	for i := 0; i < 2; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if createErr := letterRepo.Create(ctx, letter); createErr != nil {
			t.Fatalf("Create letter %d failed: %v", i, createErr)
		}

		validation := &domain.SkillValidation{
			ProfileSkillID:    skill.ID,
			ReferenceLetterID: letter.ID,
		}
		if createErr := validationRepo.Create(ctx, validation); createErr != nil {
			t.Fatalf("Create validation %d failed: %v", i, createErr)
		}
	}

	// Verify count is 2
	count, err = validationRepo.CountByProfileSkillID(ctx, skill.ID)
	if err != nil {
		t.Fatalf("CountByProfileSkillID failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 validations, got %d", count)
	}
}

func TestSkillValidationRepository_UniqueConstraint(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	skillRepo := postgres.NewProfileSkillRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewSkillValidationRepository(db)
	ctx := context.Background()

	// Create user, profile, skill, letter
	user := &domain.User{
		Email:        "skillvalidationunique@example.com",
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

	skill := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Unique Test",
		NormalizedName: "unique-test",
		Category:       "TECHNICAL",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill); err != nil {
		t.Fatalf("Create skill failed: %v", err)
	}

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create first validation
	validation1 := &domain.SkillValidation{
		ProfileSkillID:    skill.ID,
		ReferenceLetterID: letter.ID,
	}
	if err := validationRepo.Create(ctx, validation1); err != nil {
		t.Fatalf("Create first validation failed: %v", err)
	}

	// Attempt to create duplicate validation - should fail
	validation2 := &domain.SkillValidation{
		ProfileSkillID:    skill.ID,
		ReferenceLetterID: letter.ID,
	}
	err := validationRepo.Create(ctx, validation2)
	if err == nil {
		t.Error("expected error for duplicate validation, got nil")
	}
}

func TestSkillValidationRepository_BatchCountByProfileSkillIDs(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	skillRepo := postgres.NewProfileSkillRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	validationRepo := postgres.NewSkillValidationRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "skillvalidationbatch@example.com",
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

	// Create 3 skills with different validation counts
	skill1 := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Batch Skill 1",
		NormalizedName: "batch-skill-1",
		Category:       "TECHNICAL",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill1); err != nil {
		t.Fatalf("Create skill1 failed: %v", err)
	}

	skill2 := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Batch Skill 2",
		NormalizedName: "batch-skill-2",
		Category:       "TECHNICAL",
		DisplayOrder:   1,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill2); err != nil {
		t.Fatalf("Create skill2 failed: %v", err)
	}

	skill3 := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Batch Skill 3",
		NormalizedName: "batch-skill-3",
		Category:       "TECHNICAL",
		DisplayOrder:   2,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill3); err != nil {
		t.Fatalf("Create skill3 failed: %v", err)
	}

	// Create validations: skill1=2, skill2=0, skill3=3
	for i := 0; i < 2; i++ {
		letter := &domain.ReferenceLetter{
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		if err := letterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("Create letter failed: %v", err)
		}
		validation := &domain.SkillValidation{
			ProfileSkillID:    skill1.ID,
			ReferenceLetterID: letter.ID,
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation for skill1 failed: %v", err)
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
		validation := &domain.SkillValidation{
			ProfileSkillID:    skill3.ID,
			ReferenceLetterID: letter.ID,
		}
		if err := validationRepo.Create(ctx, validation); err != nil {
			t.Fatalf("Create validation for skill3 failed: %v", err)
		}
	}

	// Batch count all skills in one query
	counts, err := validationRepo.BatchCountByProfileSkillIDs(ctx, []uuid.UUID{skill1.ID, skill2.ID, skill3.ID})
	if err != nil {
		t.Fatalf("BatchCountByProfileSkillIDs failed: %v", err)
	}

	// Verify counts
	if counts[skill1.ID] != 2 {
		t.Errorf("expected skill1 count=2, got %d", counts[skill1.ID])
	}
	if counts[skill2.ID] != 0 {
		t.Errorf("expected skill2 count=0, got %d", counts[skill2.ID])
	}
	if counts[skill3.ID] != 3 {
		t.Errorf("expected skill3 count=3, got %d", counts[skill3.ID])
	}
}

func TestSkillValidationRepository_BatchCountByProfileSkillIDs_EmptyInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	validationRepo := postgres.NewSkillValidationRepository(db)
	ctx := context.Background()

	// Empty input should return empty map
	counts, err := validationRepo.BatchCountByProfileSkillIDs(ctx, []uuid.UUID{})
	if err != nil {
		t.Fatalf("BatchCountByProfileSkillIDs failed: %v", err)
	}

	if len(counts) != 0 {
		t.Errorf("expected empty map, got %d entries", len(counts))
	}
}
