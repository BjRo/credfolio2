package resolver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/graphql/resolver"
	"backend/internal/infrastructure/storage"
	"backend/internal/logger"
)

// mustCreateUser creates a user in the mock repository. Panics on error (should never happen with mocks).
func mustCreateUser(repo *mockUserRepository, user *domain.User) {
	if err := repo.Create(context.Background(), user); err != nil {
		panic("unexpected error creating user: " + err.Error())
	}
}

// mustCreateFile creates a file in the mock repository. Panics on error (should never happen with mocks).
func mustCreateFile(repo *mockFileRepository, file *domain.File) {
	if err := repo.Create(context.Background(), file); err != nil {
		panic("unexpected error creating file: " + err.Error())
	}
}

// mustCreateReferenceLetter creates a reference letter in the mock repository.
// Panics on error (should never happen with mocks).
func mustCreateReferenceLetter(repo *mockReferenceLetterRepository, letter *domain.ReferenceLetter) {
	if err := repo.Create(context.Background(), letter); err != nil {
		panic("unexpected error creating reference letter: " + err.Error())
	}
}

// mockUserRepository is a mock implementation of domain.UserRepository.
type mockUserRepository struct {
	users map[uuid.UUID]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[uuid.UUID]*domain.User)}
}

func (r *mockUserRepository) Create(_ context.Context, user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	r.users[user.ID] = user
	return nil
}

func (r *mockUserRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (r *mockUserRepository) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (r *mockUserRepository) Update(_ context.Context, user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *mockUserRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.users, id)
	return nil
}

// mockFileRepository is a mock implementation of domain.FileRepository.
type mockFileRepository struct {
	files map[uuid.UUID]*domain.File
}

func newMockFileRepository() *mockFileRepository {
	return &mockFileRepository{files: make(map[uuid.UUID]*domain.File)}
}

func (r *mockFileRepository) Create(_ context.Context, file *domain.File) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	r.files[file.ID] = file
	return nil
}

func (r *mockFileRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.File, error) {
	file, ok := r.files[id]
	if !ok {
		return nil, nil
	}
	return file, nil
}

func (r *mockFileRepository) GetByUserID(_ context.Context, userID uuid.UUID) ([]*domain.File, error) {
	var result []*domain.File
	for _, file := range r.files {
		if file.UserID == userID {
			result = append(result, file)
		}
	}
	return result, nil
}

func (r *mockFileRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.files, id)
	return nil
}

// mockReferenceLetterRepository is a mock implementation of domain.ReferenceLetterRepository.
type mockReferenceLetterRepository struct {
	letters map[uuid.UUID]*domain.ReferenceLetter
}

func newMockReferenceLetterRepository() *mockReferenceLetterRepository {
	return &mockReferenceLetterRepository{letters: make(map[uuid.UUID]*domain.ReferenceLetter)}
}

func (r *mockReferenceLetterRepository) Create(_ context.Context, letter *domain.ReferenceLetter) error {
	if letter.ID == uuid.Nil {
		letter.ID = uuid.New()
	}
	r.letters[letter.ID] = letter
	return nil
}

func (r *mockReferenceLetterRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	letter, ok := r.letters[id]
	if !ok {
		return nil, nil
	}
	return letter, nil
}

func (r *mockReferenceLetterRepository) GetByUserID(_ context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	var result []*domain.ReferenceLetter
	for _, letter := range r.letters {
		if letter.UserID == userID {
			result = append(result, letter)
		}
	}
	return result, nil
}

func (r *mockReferenceLetterRepository) Update(_ context.Context, letter *domain.ReferenceLetter) error {
	r.letters[letter.ID] = letter
	return nil
}

func (r *mockReferenceLetterRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.letters, id)
	return nil
}

// errorUserRepository returns errors for testing error handling.
type errorUserRepository struct{}

func (r *errorUserRepository) Create(_ context.Context, _ *domain.User) error {
	return errors.New("database error")
}

func (r *errorUserRepository) GetByID(_ context.Context, _ uuid.UUID) (*domain.User, error) {
	return nil, errors.New("database error")
}

func (r *errorUserRepository) GetByEmail(_ context.Context, _ string) (*domain.User, error) {
	return nil, errors.New("database error")
}

func (r *errorUserRepository) Update(_ context.Context, _ *domain.User) error {
	return errors.New("database error")
}

func (r *errorUserRepository) Delete(_ context.Context, _ uuid.UUID) error {
	return errors.New("database error")
}

// mockResumeRepository is a mock implementation of domain.ResumeRepository.
type mockResumeRepository struct {
	resumes map[uuid.UUID]*domain.Resume
}

func newMockResumeRepository() *mockResumeRepository {
	return &mockResumeRepository{resumes: make(map[uuid.UUID]*domain.Resume)}
}

func (r *mockResumeRepository) Create(_ context.Context, resume *domain.Resume) error {
	if resume.ID == uuid.Nil {
		resume.ID = uuid.New()
	}
	r.resumes[resume.ID] = resume
	return nil
}

func (r *mockResumeRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Resume, error) {
	resume, ok := r.resumes[id]
	if !ok {
		return nil, nil
	}
	return resume, nil
}

func (r *mockResumeRepository) GetByUserID(_ context.Context, userID uuid.UUID) ([]*domain.Resume, error) {
	var result []*domain.Resume
	for _, resume := range r.resumes {
		if resume.UserID == userID {
			result = append(result, resume)
		}
	}
	return result, nil
}

func (r *mockResumeRepository) Update(_ context.Context, resume *domain.Resume) error {
	r.resumes[resume.ID] = resume
	return nil
}

func (r *mockResumeRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.resumes, id)
	return nil
}

// mockJobEnqueuer is a mock implementation of domain.JobEnqueuer.
type mockJobEnqueuer struct {
	enqueuedDocJobs    []domain.DocumentProcessingRequest
	enqueuedResumeJobs []domain.ResumeProcessingRequest
}

func newMockJobEnqueuer() *mockJobEnqueuer {
	return &mockJobEnqueuer{
		enqueuedDocJobs:    make([]domain.DocumentProcessingRequest, 0),
		enqueuedResumeJobs: make([]domain.ResumeProcessingRequest, 0),
	}
}

func (e *mockJobEnqueuer) EnqueueDocumentProcessing(_ context.Context, req domain.DocumentProcessingRequest) error {
	e.enqueuedDocJobs = append(e.enqueuedDocJobs, req)
	return nil
}

func (e *mockJobEnqueuer) EnqueueResumeProcessing(_ context.Context, req domain.ResumeProcessingRequest) error {
	e.enqueuedResumeJobs = append(e.enqueuedResumeJobs, req)
	return nil
}

// testLogger returns a logger that discards all output (for tests).
func testLogger() logger.Logger {
	return logger.NewStdoutLogger(logger.WithMinLevel(logger.Severity(100))) // level 100 = discard all
}

func TestUserQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	fileRepo := newMockFileRepository()
	refLetterRepo := newMockReferenceLetterRepository()

	ctx := context.Background()

	// Create a test user
	name := "Test User"
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns user when found", func(t *testing.T) {
		result, err := query.User(ctx, user.ID.String())
		if err != nil {
			t.Fatalf("User query failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected user, got nil")
		}

		if result.ID != user.ID.String() {
			t.Errorf("ID mismatch: got %s, want %s", result.ID, user.ID.String())
		}

		if result.Email != user.Email {
			t.Errorf("Email mismatch: got %s, want %s", result.Email, user.Email)
		}

		if result.Name == nil || *result.Name != *user.Name {
			t.Errorf("Name mismatch: got %v, want %v", result.Name, user.Name)
		}
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		result, err := query.User(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("User query failed: %v", err)
		}

		if result != nil {
			t.Error("expected nil for non-existent user")
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.User(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		errorR := resolver.NewResolver(&errorUserRepository{}, fileRepo, refLetterRepo, newMockResumeRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
		errorQuery := errorR.Query()

		_, err := errorQuery.User(ctx, uuid.New().String())
		if err == nil {
			t.Error("expected error when repository fails")
		}
	})
}

func TestFileQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	fileRepo := newMockFileRepository()
	refLetterRepo := newMockReferenceLetterRepository()

	ctx := context.Background()

	// Create test user and file
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "file-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	file := &domain.File{
		ID:          uuid.New(),
		UserID:      user.ID,
		Filename:    "test-document.pdf",
		ContentType: "application/pdf",
		SizeBytes:   1024,
		StorageKey:  "test-key",
	}
	mustCreateFile(fileRepo, file)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns file when found", func(t *testing.T) {
		result, err := query.File(ctx, file.ID.String())
		if err != nil {
			t.Fatalf("File query failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected file, got nil")
		}

		if result.ID != file.ID.String() {
			t.Errorf("ID mismatch: got %s, want %s", result.ID, file.ID.String())
		}

		if result.Filename != file.Filename {
			t.Errorf("Filename mismatch: got %s, want %s", result.Filename, file.Filename)
		}

		// Verify user relation is populated
		if result.User == nil {
			t.Fatal("expected user relation to be populated")
		}

		if result.User.ID != user.ID.String() {
			t.Errorf("User ID mismatch: got %s, want %s", result.User.ID, user.ID.String())
		}
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		result, err := query.File(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("File query failed: %v", err)
		}

		if result != nil {
			t.Error("expected nil for non-existent file")
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.File(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestFilesQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	fileRepo := newMockFileRepository()
	refLetterRepo := newMockReferenceLetterRepository()

	ctx := context.Background()

	// Create test user and files
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "files-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	file1 := &domain.File{
		ID:          uuid.New(),
		UserID:      user.ID,
		Filename:    "document1.pdf",
		ContentType: "application/pdf",
		SizeBytes:   1024,
		StorageKey:  "test-key-1",
	}
	mustCreateFile(fileRepo, file1)

	file2 := &domain.File{
		ID:          uuid.New(),
		UserID:      user.ID,
		Filename:    "document2.pdf",
		ContentType: "application/pdf",
		SizeBytes:   2048,
		StorageKey:  "test-key-2",
	}
	mustCreateFile(fileRepo, file2)

	// Create another user with files (to verify filtering)
	otherUser := &domain.User{
		ID:           uuid.New(),
		Email:        "other@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, otherUser)

	otherFile := &domain.File{
		ID:          uuid.New(),
		UserID:      otherUser.ID,
		Filename:    "other-doc.pdf",
		ContentType: "application/pdf",
		SizeBytes:   512,
		StorageKey:  "other-key",
	}
	mustCreateFile(fileRepo, otherFile)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns files for user", func(t *testing.T) {
		results, err := query.Files(ctx, user.ID.String())
		if err != nil {
			t.Fatalf("Files query failed: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("expected 2 files, got %d", len(results))
		}

		// Verify the files belong to the user
		ids := make(map[string]bool)
		for _, f := range results {
			ids[f.ID] = true
			// Verify user relation is populated
			if f.User == nil {
				t.Error("expected user relation to be populated")
			} else if f.User.ID != user.ID.String() {
				t.Errorf("User ID mismatch: got %s, want %s", f.User.ID, user.ID.String())
			}
		}

		if !ids[file1.ID.String()] || !ids[file2.ID.String()] {
			t.Error("returned files don't match expected files")
		}
	})

	t.Run("returns empty slice for user with no files", func(t *testing.T) {
		noFilesUser := &domain.User{
			ID:           uuid.New(),
			Email:        "nofiles@example.com",
			PasswordHash: "hashed",
		}
		mustCreateUser(userRepo, noFilesUser)

		results, err := query.Files(ctx, noFilesUser.ID.String())
		if err != nil {
			t.Fatalf("Files query failed: %v", err)
		}

		if results == nil {
			t.Error("expected empty slice, got nil")
		}

		if len(results) != 0 {
			t.Errorf("expected 0 files, got %d", len(results))
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.Files(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

//nolint:gocyclo // Test function with multiple subtests
func TestReferenceLetterQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	fileRepo := newMockFileRepository()
	refLetterRepo := newMockReferenceLetterRepository()

	ctx := context.Background()

	// Create test user, file, and reference letter
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "letter-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	file := &domain.File{
		ID:          uuid.New(),
		UserID:      user.ID,
		Filename:    "letter.pdf",
		ContentType: "application/pdf",
		SizeBytes:   1024,
		StorageKey:  "letter-key",
	}
	mustCreateFile(fileRepo, file)

	title := "Test Letter"
	letter := &domain.ReferenceLetter{
		ID:     uuid.New(),
		UserID: user.ID,
		FileID: &file.ID,
		Title:  &title,
		Status: domain.ReferenceLetterStatusPending,
	}
	mustCreateReferenceLetter(refLetterRepo, letter)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns reference letter when found", func(t *testing.T) {
		result, err := query.ReferenceLetter(ctx, letter.ID.String())
		if err != nil {
			t.Fatalf("ReferenceLetter query failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected reference letter, got nil")
		}

		if result.ID != letter.ID.String() {
			t.Errorf("ID mismatch: got %s, want %s", result.ID, letter.ID.String())
		}

		if result.Title == nil || *result.Title != *letter.Title {
			t.Errorf("Title mismatch: got %v, want %v", result.Title, letter.Title)
		}

		// Verify user relation is populated
		if result.User == nil {
			t.Fatal("expected user relation to be populated")
		}

		if result.User.ID != user.ID.String() {
			t.Errorf("User ID mismatch: got %s, want %s", result.User.ID, user.ID.String())
		}

		// Verify file relation is populated
		if result.File == nil {
			t.Fatal("expected file relation to be populated")
		}

		if result.File.ID != file.ID.String() {
			t.Errorf("File ID mismatch: got %s, want %s", result.File.ID, file.ID.String())
		}
	})

	t.Run("returns nil for file when no file associated", func(t *testing.T) {
		titleNoFile := "Letter Without File"
		letterNoFile := &domain.ReferenceLetter{
			ID:     uuid.New(),
			UserID: user.ID,
			FileID: nil,
			Title:  &titleNoFile,
			Status: domain.ReferenceLetterStatusPending,
		}
		mustCreateReferenceLetter(refLetterRepo, letterNoFile)

		result, err := query.ReferenceLetter(ctx, letterNoFile.ID.String())
		if err != nil {
			t.Fatalf("ReferenceLetter query failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected reference letter, got nil")
		}

		if result.File != nil {
			t.Error("expected nil file for letter without file association")
		}
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		result, err := query.ReferenceLetter(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("ReferenceLetter query failed: %v", err)
		}

		if result != nil {
			t.Error("expected nil for non-existent reference letter")
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.ReferenceLetter(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestReferenceLettersQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	fileRepo := newMockFileRepository()
	refLetterRepo := newMockReferenceLetterRepository()

	ctx := context.Background()

	// Create test user and reference letters
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "letters-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	file := &domain.File{
		ID:          uuid.New(),
		UserID:      user.ID,
		Filename:    "letter1.pdf",
		ContentType: "application/pdf",
		SizeBytes:   1024,
		StorageKey:  "letter1-key",
	}
	mustCreateFile(fileRepo, file)

	title1 := "Letter 1"
	letter1 := &domain.ReferenceLetter{
		ID:     uuid.New(),
		UserID: user.ID,
		FileID: &file.ID,
		Title:  &title1,
		Status: domain.ReferenceLetterStatusPending,
	}
	mustCreateReferenceLetter(refLetterRepo, letter1)

	title2 := "Letter 2"
	letter2 := &domain.ReferenceLetter{
		ID:     uuid.New(),
		UserID: user.ID,
		FileID: nil,
		Title:  &title2,
		Status: domain.ReferenceLetterStatusPending,
	}
	mustCreateReferenceLetter(refLetterRepo, letter2)

	// Create another user with letters
	otherUser := &domain.User{
		ID:           uuid.New(),
		Email:        "other-letters@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, otherUser)

	otherTitle := "Other Letter"
	otherLetter := &domain.ReferenceLetter{
		ID:     uuid.New(),
		UserID: otherUser.ID,
		FileID: nil,
		Title:  &otherTitle,
		Status: domain.ReferenceLetterStatusPending,
	}
	mustCreateReferenceLetter(refLetterRepo, otherLetter)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns reference letters for user", func(t *testing.T) {
		results, err := query.ReferenceLetters(ctx, user.ID.String())
		if err != nil {
			t.Fatalf("ReferenceLetters query failed: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("expected 2 letters, got %d", len(results))
		}

		// Verify the letters belong to the user
		ids := make(map[string]bool)
		for _, l := range results {
			ids[l.ID] = true
			// Verify user relation is populated
			if l.User == nil {
				t.Error("expected user relation to be populated")
			} else if l.User.ID != user.ID.String() {
				t.Errorf("User ID mismatch: got %s, want %s", l.User.ID, user.ID.String())
			}
		}

		if !ids[letter1.ID.String()] || !ids[letter2.ID.String()] {
			t.Error("returned letters don't match expected letters")
		}
	})

	t.Run("returns empty slice for user with no letters", func(t *testing.T) {
		noLettersUser := &domain.User{
			ID:           uuid.New(),
			Email:        "noletters@example.com",
			PasswordHash: "hashed",
		}
		mustCreateUser(userRepo, noLettersUser)

		results, err := query.ReferenceLetters(ctx, noLettersUser.ID.String())
		if err != nil {
			t.Fatalf("ReferenceLetters query failed: %v", err)
		}

		if results == nil {
			t.Error("expected empty slice, got nil")
		}

		if len(results) != 0 {
			t.Errorf("expected 0 letters, got %d", len(results))
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.ReferenceLetters(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}
