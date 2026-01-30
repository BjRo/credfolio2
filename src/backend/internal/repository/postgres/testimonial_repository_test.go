package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestTestimonialRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	testimonialRepo := postgres.NewTestimonialRepository(db)
	ctx := context.Background()

	// Create user, profile, and reference letter
	user := &domain.User{
		Email:        "testimonialuser@example.com",
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

	letter := &domain.ReferenceLetter{
		UserID:     user.ID,
		Title:      strPtr("Reference Letter"),
		AuthorName: strPtr("John Smith"),
		Status:     domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create testimonial
	testimonial := &domain.Testimonial{
		ProfileID:         profile.ID,
		ReferenceLetterID: letter.ID,
		Quote:             "Jane is an exceptional engineer with outstanding leadership skills.",
		AuthorName:        "John Smith",
		AuthorTitle:       strPtr("Engineering Manager"),
		AuthorCompany:     strPtr("Acme Corp"),
		Relationship:      domain.TestimonialRelationshipManager,
	}

	err := testimonialRepo.Create(ctx, testimonial)
	if err != nil {
		t.Fatalf("Create testimonial failed: %v", err)
	}

	if testimonial.ID == uuid.Nil {
		t.Error("expected testimonial ID to be set after create")
	}
}

func TestTestimonialRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	testimonialRepo := postgres.NewTestimonialRepository(db)
	ctx := context.Background()

	// Create user, profile, and reference letter
	user := &domain.User{
		Email:        "testimonialgetbyid@example.com",
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

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create testimonial
	testimonial := &domain.Testimonial{
		ProfileID:         profile.ID,
		ReferenceLetterID: letter.ID,
		Quote:             "Great team player.",
		AuthorName:        "Jane Doe",
		Relationship:      domain.TestimonialRelationshipPeer,
	}
	if err := testimonialRepo.Create(ctx, testimonial); err != nil {
		t.Fatalf("Create testimonial failed: %v", err)
	}

	// Retrieve by ID
	found, err := testimonialRepo.GetByID(ctx, testimonial.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find testimonial, got nil")
	}

	if found.Quote != testimonial.Quote {
		t.Errorf("quote mismatch: got %q, want %q", found.Quote, testimonial.Quote)
	}

	// Test not found
	notFound, err := testimonialRepo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID for non-existent testimonial failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent testimonial")
	}
}

func TestTestimonialRepository_GetByProfileID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	testimonialRepo := postgres.NewTestimonialRepository(db)
	ctx := context.Background()

	// Create user, profile, and reference letter
	user := &domain.User{
		Email:        "testimonialbyprofile@example.com",
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

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create multiple testimonials
	for i := 0; i < 3; i++ {
		testimonial := &domain.Testimonial{
			ProfileID:         profile.ID,
			ReferenceLetterID: letter.ID,
			Quote:             "Testimonial " + string(rune('A'+i)),
			AuthorName:        "Author " + string(rune('A'+i)),
			Relationship:      domain.TestimonialRelationshipOther,
		}
		if err := testimonialRepo.Create(ctx, testimonial); err != nil {
			t.Fatalf("Create testimonial %d failed: %v", i, err)
		}
	}

	// Retrieve by profile ID
	testimonials, err := testimonialRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		t.Fatalf("GetByProfileID failed: %v", err)
	}

	if len(testimonials) != 3 {
		t.Errorf("expected 3 testimonials, got %d", len(testimonials))
	}

	// Test empty result for non-existent profile
	empty, err := testimonialRepo.GetByProfileID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByProfileID for non-existent profile failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected 0 testimonials, got %d", len(empty))
	}
}

func TestTestimonialRepository_DeleteByReferenceLetterID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	letterRepo := postgres.NewReferenceLetterRepository(db)
	testimonialRepo := postgres.NewTestimonialRepository(db)
	ctx := context.Background()

	// Create user, profile, and reference letter
	user := &domain.User{
		Email:        "testimonialdelete@example.com",
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

	letter := &domain.ReferenceLetter{
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	if err := letterRepo.Create(ctx, letter); err != nil {
		t.Fatalf("Create letter failed: %v", err)
	}

	// Create multiple testimonials
	for i := 0; i < 3; i++ {
		testimonial := &domain.Testimonial{
			ProfileID:         profile.ID,
			ReferenceLetterID: letter.ID,
			Quote:             "Delete test " + string(rune('A'+i)),
			AuthorName:        "Author",
			Relationship:      domain.TestimonialRelationshipOther,
		}
		if err := testimonialRepo.Create(ctx, testimonial); err != nil {
			t.Fatalf("Create testimonial %d failed: %v", i, err)
		}
	}

	// Delete all testimonials by reference letter ID
	if err := testimonialRepo.DeleteByReferenceLetterID(ctx, letter.ID); err != nil {
		t.Fatalf("DeleteByReferenceLetterID failed: %v", err)
	}

	// Verify deletion
	remaining, err := testimonialRepo.GetByReferenceLetterID(ctx, letter.ID)
	if err != nil {
		t.Fatalf("GetByReferenceLetterID failed: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("expected 0 testimonials after delete, got %d", len(remaining))
	}
}
