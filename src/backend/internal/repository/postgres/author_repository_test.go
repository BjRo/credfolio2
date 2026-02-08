package postgres_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"backend/internal/domain"
	"backend/internal/repository/postgres"
)

func TestAuthorRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authoruser@example.com",
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

	// Create author
	author := &domain.Author{
		ProfileID:   profile.ID,
		Name:        "John Smith",
		Title:       strPtr("Engineering Manager"),
		Company:     strPtr("Acme Corp"),
		LinkedInURL: strPtr("https://linkedin.com/in/johnsmith"),
	}

	err := authorRepo.Create(ctx, author)
	if err != nil {
		t.Fatalf("Create author failed: %v", err)
	}

	if author.ID == uuid.Nil {
		t.Error("expected author ID to be set after create")
	}
}

func TestAuthorRepository_Create_DuplicatePrevention(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authordup@example.com",
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

	// Create first author
	author1 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "John Smith",
		Company:   strPtr("Acme Corp"),
	}
	if err := authorRepo.Create(ctx, author1); err != nil {
		t.Fatalf("Create first author failed: %v", err)
	}

	// Try to create duplicate author with same name and company
	author2 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "John Smith",
		Company:   strPtr("Acme Corp"),
	}
	err := authorRepo.Create(ctx, author2)
	if err == nil {
		t.Fatal("expected error when creating duplicate author, got nil")
	}

	// Verify it's a unique constraint violation (PostgreSQL error code 23505)
	// Note: errors.As unwraps through Bun's pgdriver.Error to get the underlying pgconn.PgError
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code != "23505" {
			t.Errorf("expected unique violation error code 23505, got %s", pgErr.Code)
		}
		if pgErr.ConstraintName != "idx_authors_profile_name_company" {
			t.Errorf("expected constraint name %q, got %q", "idx_authors_profile_name_company", pgErr.ConstraintName)
		}
	} else {
		// If we can't unwrap to pgconn.PgError, at least verify the error mentions the constraint
		errMsg := err.Error()
		if !strings.Contains(errMsg, "duplicate key") && !strings.Contains(errMsg, "23505") {
			t.Errorf("error doesn't appear to be a unique constraint violation: %v", err)
		}
	}
}

func TestAuthorRepository_Create_AllowSameNameDifferentCompany(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authordiffcompany@example.com",
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

	// Create first author at Company A
	author1 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "John Smith",
		Company:   strPtr("Acme Corp"),
	}
	if err := authorRepo.Create(ctx, author1); err != nil {
		t.Fatalf("Create first author failed: %v", err)
	}

	// Create second author with same name but different company - should work
	author2 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "John Smith",
		Company:   strPtr("Other Inc"),
	}
	if err := authorRepo.Create(ctx, author2); err != nil {
		t.Fatalf("Create second author with different company failed: %v", err)
	}

	if author1.ID == author2.ID {
		t.Error("expected different IDs for different authors")
	}
}

func TestAuthorRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authorgetbyid@example.com",
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

	// Create author
	author := &domain.Author{
		ProfileID:   profile.ID,
		Name:        "Jane Doe",
		Title:       strPtr("CTO"),
		Company:     strPtr("Tech Corp"),
		LinkedInURL: strPtr("https://linkedin.com/in/janedoe"),
	}
	if err := authorRepo.Create(ctx, author); err != nil {
		t.Fatalf("Create author failed: %v", err)
	}

	// Retrieve by ID
	found, err := authorRepo.GetByID(ctx, author.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected to find author, got nil")
	}

	if found.Name != author.Name {
		t.Errorf("name mismatch: got %q, want %q", found.Name, author.Name)
	}

	if *found.Title != *author.Title {
		t.Errorf("title mismatch: got %q, want %q", *found.Title, *author.Title)
	}

	if *found.LinkedInURL != *author.LinkedInURL {
		t.Errorf("linkedInUrl mismatch: got %q, want %q", *found.LinkedInURL, *author.LinkedInURL)
	}

	// Test not found
	notFound, err := authorRepo.GetByID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByID for non-existent author failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent author")
	}
}

func TestAuthorRepository_GetByProfileID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authorbyprofile@example.com",
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

	// Create multiple authors
	for i := 0; i < 3; i++ {
		author := &domain.Author{
			ProfileID: profile.ID,
			Name:      "Author " + string(rune('A'+i)),
			Company:   strPtr("Company " + string(rune('A'+i))),
		}
		if err := authorRepo.Create(ctx, author); err != nil {
			t.Fatalf("Create author %d failed: %v", i, err)
		}
	}

	// Retrieve by profile ID
	authors, err := authorRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		t.Fatalf("GetByProfileID failed: %v", err)
	}

	if len(authors) != 3 {
		t.Errorf("expected 3 authors, got %d", len(authors))
	}

	// Test empty result for non-existent profile
	empty, err := authorRepo.GetByProfileID(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetByProfileID for non-existent profile failed: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected 0 authors, got %d", len(empty))
	}
}

func TestAuthorRepository_FindByNameAndCompany(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authorfind@example.com",
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

	// Create author with company
	authorWithCompany := &domain.Author{
		ProfileID: profile.ID,
		Name:      "John Smith",
		Company:   strPtr("Acme Corp"),
	}
	if err := authorRepo.Create(ctx, authorWithCompany); err != nil {
		t.Fatalf("Create author with company failed: %v", err)
	}

	// Create author without company
	authorNoCompany := &domain.Author{
		ProfileID: profile.ID,
		Name:      "Jane Doe",
		Company:   nil,
	}
	if err := authorRepo.Create(ctx, authorNoCompany); err != nil {
		t.Fatalf("Create author without company failed: %v", err)
	}

	// Find by name and company
	found, err := authorRepo.FindByNameAndCompany(ctx, profile.ID, "John Smith", strPtr("Acme Corp"))
	if err != nil {
		t.Fatalf("FindByNameAndCompany failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find author, got nil")
	}
	if found.ID != authorWithCompany.ID {
		t.Errorf("found wrong author: got ID %s, want %s", found.ID, authorWithCompany.ID)
	}

	// Find by name with nil company
	foundNil, err := authorRepo.FindByNameAndCompany(ctx, profile.ID, "Jane Doe", nil)
	if err != nil {
		t.Fatalf("FindByNameAndCompany with nil company failed: %v", err)
	}
	if foundNil == nil {
		t.Fatal("expected to find author with nil company, got nil")
	}
	if foundNil.ID != authorNoCompany.ID {
		t.Errorf("found wrong author: got ID %s, want %s", foundNil.ID, authorNoCompany.ID)
	}

	// Not found - different company
	notFound, err := authorRepo.FindByNameAndCompany(ctx, profile.ID, "John Smith", strPtr("Other Corp"))
	if err != nil {
		t.Fatalf("FindByNameAndCompany for non-existent author failed: %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent author")
	}

	// Not found - different name
	notFound2, err := authorRepo.FindByNameAndCompany(ctx, profile.ID, "Bob Jones", strPtr("Acme Corp"))
	if err != nil {
		t.Fatalf("FindByNameAndCompany for non-existent name failed: %v", err)
	}
	if notFound2 != nil {
		t.Error("expected nil for non-existent author name")
	}
}

func TestAuthorRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authorupdate@example.com",
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

	// Create author
	author := &domain.Author{
		ProfileID: profile.ID,
		Name:      "John Smith",
		Title:     strPtr("Manager"),
		Company:   strPtr("Old Corp"),
	}
	if err := authorRepo.Create(ctx, author); err != nil {
		t.Fatalf("Create author failed: %v", err)
	}

	// Update author
	author.Name = "John A. Smith"
	author.Title = strPtr("Senior Manager")
	author.Company = strPtr("New Corp")
	author.LinkedInURL = strPtr("https://linkedin.com/in/johnsmith")

	if err := authorRepo.Update(ctx, author); err != nil {
		t.Fatalf("Update author failed: %v", err)
	}

	// Verify update
	updated, err := authorRepo.GetByID(ctx, author.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}

	if updated.Name != "John A. Smith" {
		t.Errorf("name not updated: got %q, want %q", updated.Name, "John A. Smith")
	}
	if *updated.Title != "Senior Manager" {
		t.Errorf("title not updated: got %q, want %q", *updated.Title, "Senior Manager")
	}
	if *updated.Company != "New Corp" {
		t.Errorf("company not updated: got %q, want %q", *updated.Company, "New Corp")
	}
	if *updated.LinkedInURL != "https://linkedin.com/in/johnsmith" {
		t.Errorf("linkedInUrl not updated: got %q, want %q", *updated.LinkedInURL, "https://linkedin.com/in/johnsmith")
	}
}

func TestAuthorRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authordelete@example.com",
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

	// Create author
	author := &domain.Author{
		ProfileID: profile.ID,
		Name:      "Delete Me",
	}
	if err := authorRepo.Create(ctx, author); err != nil {
		t.Fatalf("Create author failed: %v", err)
	}

	// Delete author
	if err := authorRepo.Delete(ctx, author.ID); err != nil {
		t.Fatalf("Delete author failed: %v", err)
	}

	// Verify deletion
	deleted, err := authorRepo.GetByID(ctx, author.ID)
	if err != nil {
		t.Fatalf("GetByID after delete failed: %v", err)
	}
	if deleted != nil {
		t.Error("expected author to be deleted, but found it")
	}
}

func TestAuthorRepository_Upsert_Concurrent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authorupsert@example.com",
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

	// Simulate concurrent upserts (same name + company)
	authorData := &domain.Author{
		ID:        uuid.New(),
		ProfileID: profile.ID,
		Name:      "Jane Concurrent",
		Company:   strPtr("ConcurrentCo"),
	}

	// Clone for second goroutine
	authorData2 := &domain.Author{
		ID:        uuid.New(),
		ProfileID: profile.ID,
		Name:      "Jane Concurrent",
		Company:   strPtr("ConcurrentCo"),
	}

	var author1, author2 *domain.Author
	var err1, err2 error
	var wg sync.WaitGroup
	wg.Add(2)

	// Launch two goroutines that both try to upsert the same author
	go func() {
		defer wg.Done()
		author1, err1 = authorRepo.Upsert(ctx, authorData)
	}()

	go func() {
		defer wg.Done()
		author2, err2 = authorRepo.Upsert(ctx, authorData2)
	}()

	wg.Wait()

	// Both should succeed
	if err1 != nil {
		t.Fatalf("first upsert failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("second upsert failed: %v", err2)
	}

	// Both should return the same author ID (no duplicates)
	if author1.ID != author2.ID {
		t.Errorf("concurrent upserts created duplicates: %s vs %s", author1.ID, author2.ID)
	}

	// Verify only one author exists in database
	authors, err := authorRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		t.Fatalf("GetByProfileID failed: %v", err)
	}
	if len(authors) != 1 {
		t.Errorf("expected 1 author, got %d (duplicate created)", len(authors))
	}
}

func TestAuthorRepository_GetByIDs(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	userRepo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Create user and profile
	user := &domain.User{
		Email:        "authorbatch@example.com",
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

	// Create 3 authors
	author1 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "Author One",
		Company:   strPtr("Company A"),
	}
	if err := authorRepo.Create(ctx, author1); err != nil {
		t.Fatalf("Create author1 failed: %v", err)
	}

	author2 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "Author Two",
		Company:   strPtr("Company B"),
	}
	if err := authorRepo.Create(ctx, author2); err != nil {
		t.Fatalf("Create author2 failed: %v", err)
	}

	author3 := &domain.Author{
		ProfileID: profile.ID,
		Name:      "Author Three",
		Company:   strPtr("Company C"),
	}
	if err := authorRepo.Create(ctx, author3); err != nil {
		t.Fatalf("Create author3 failed: %v", err)
	}

	// Batch get all 3 authors plus one non-existent ID
	nonExistentID := uuid.New()
	authors, err := authorRepo.GetByIDs(ctx, []uuid.UUID{author1.ID, author2.ID, author3.ID, nonExistentID})
	if err != nil {
		t.Fatalf("GetByIDs failed: %v", err)
	}

	// Verify we got 3 authors (non-existent ID should not be in map)
	if len(authors) != 3 {
		t.Errorf("expected 3 authors, got %d", len(authors))
	}

	// Verify each author is in the map with correct data
	if a, ok := authors[author1.ID]; !ok {
		t.Errorf("author1 not found in results")
	} else if a.Name != "Author One" {
		t.Errorf("author1 name mismatch: got %s, want Author One", a.Name)
	}

	if a, ok := authors[author2.ID]; !ok {
		t.Errorf("author2 not found in results")
	} else if a.Name != "Author Two" {
		t.Errorf("author2 name mismatch: got %s, want Author Two", a.Name)
	}

	if a, ok := authors[author3.ID]; !ok {
		t.Errorf("author3 not found in results")
	} else if a.Name != "Author Three" {
		t.Errorf("author3 name mismatch: got %s, want Author Three", a.Name)
	}

	// Verify non-existent ID is not in map
	if _, ok := authors[nonExistentID]; ok {
		t.Errorf("non-existent ID should not be in results")
	}
}

func TestAuthorRepository_GetByIDs_EmptyInput(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	cleanupTestData(t, db)

	authorRepo := postgres.NewAuthorRepository(db)
	ctx := context.Background()

	// Empty input should return empty map
	authors, err := authorRepo.GetByIDs(ctx, []uuid.UUID{})
	if err != nil {
		t.Fatalf("GetByIDs with empty input failed: %v", err)
	}

	if len(authors) != 0 {
		t.Errorf("expected empty map, got %d entries", len(authors))
	}
}
