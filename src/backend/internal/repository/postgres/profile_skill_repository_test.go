package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestProfileSkillRepository_GetByIDs(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	skillRepo := postgres.NewProfileSkillRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "skillbatch@example.com",
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

	// Create 3 skills
	skill1 := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Go",
		NormalizedName: "go",
		Category:       "TECHNICAL",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill1); err != nil {
		t.Fatalf("Create skill1 failed: %v", err)
	}

	skill2 := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Python",
		NormalizedName: "python",
		Category:       "TECHNICAL",
		DisplayOrder:   1,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill2); err != nil {
		t.Fatalf("Create skill2 failed: %v", err)
	}

	skill3 := &domain.ProfileSkill{
		ProfileID:      profile.ID,
		Name:           "Leadership",
		NormalizedName: "leadership",
		Category:       "SOFT",
		DisplayOrder:   2,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill3); err != nil {
		t.Fatalf("Create skill3 failed: %v", err)
	}

	// Batch get all 3 skills plus one non-existent ID
	nonExistentID := uuid.New()
	skills, err := skillRepo.GetByIDs(ctx, []uuid.UUID{skill1.ID, skill2.ID, skill3.ID, nonExistentID})
	if err != nil {
		t.Fatalf("GetByIDs failed: %v", err)
	}

	// Verify we got 3 skills (non-existent ID should not be in map)
	if len(skills) != 3 {
		t.Errorf("expected 3 skills, got %d", len(skills))
	}

	// Verify each skill is in the map with correct data
	if s, ok := skills[skill1.ID]; !ok {
		t.Errorf("skill1 not found in results")
	} else if s.Name != "Go" {
		t.Errorf("skill1 name mismatch: got %s, want Go", s.Name)
	}

	if s, ok := skills[skill2.ID]; !ok {
		t.Errorf("skill2 not found in results")
	} else if s.Name != "Python" {
		t.Errorf("skill2 name mismatch: got %s, want Python", s.Name)
	}

	if s, ok := skills[skill3.ID]; !ok {
		t.Errorf("skill3 not found in results")
	} else if s.Name != "Leadership" {
		t.Errorf("skill3 name mismatch: got %s, want Leadership", s.Name)
	}

	// Verify non-existent ID is not in map
	if _, ok := skills[nonExistentID]; ok {
		t.Errorf("non-existent ID should not be in results")
	}
}

func TestProfileSkillRepository_GetByIDs_EmptyInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	skillRepo := postgres.NewProfileSkillRepository(db)
	ctx := context.Background()

	// Empty input should return empty map
	skills, err := skillRepo.GetByIDs(ctx, []uuid.UUID{})
	if err != nil {
		t.Fatalf("GetByIDs with empty input failed: %v", err)
	}

	if len(skills) != 0 {
		t.Errorf("expected empty map, got %d entries", len(skills))
	}
}
