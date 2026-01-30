package job_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/domain"
	"backend/internal/job"
	"backend/internal/logger"
)

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

// mockStorage is a mock implementation of domain.Storage.
type mockStorage struct {
	existingKeys map[string]bool
}

func newMockStorage() *mockStorage {
	return &mockStorage{existingKeys: make(map[string]bool)}
}

func (s *mockStorage) Upload(_ context.Context, _ string, _ io.Reader, _ int64, _ string) (*domain.StorageObject, error) {
	return &domain.StorageObject{}, nil
}

func (s *mockStorage) Download(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, nil
}

func (s *mockStorage) Delete(_ context.Context, _ string) error {
	return nil
}

func (s *mockStorage) GetPresignedURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", nil
}

func (s *mockStorage) GetPublicURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", nil
}

func (s *mockStorage) Exists(_ context.Context, key string) (bool, error) {
	return s.existingKeys[key], nil
}

func (s *mockStorage) addKey(key string) {
	s.existingKeys[key] = true
}

// errorStorage always returns errors for testing error handling.
type errorStorage struct{}

func (s *errorStorage) Upload(_ context.Context, _ string, _ io.Reader, _ int64, _ string) (*domain.StorageObject, error) {
	return nil, errors.New("storage error")
}

func (s *errorStorage) Download(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, errors.New("storage error")
}

func (s *errorStorage) Delete(_ context.Context, _ string) error {
	return errors.New("storage error")
}

func (s *errorStorage) GetPresignedURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", errors.New("storage error")
}

func (s *errorStorage) GetPublicURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", errors.New("storage error")
}

func (s *errorStorage) Exists(_ context.Context, _ string) (bool, error) {
	return false, errors.New("storage error")
}

// testLogger returns a logger that discards all output (for tests).
func testLogger() logger.Logger {
	return logger.NewStdoutLogger(logger.WithMinLevel(logger.Severity(100))) // level 100 = discard all
}

func TestDocumentProcessingArgs_Kind(t *testing.T) {
	args := job.DocumentProcessingArgs{}
	if got := args.Kind(); got != "document_processing" {
		t.Errorf("Kind() = %q, want %q", got, "document_processing")
	}
}

func TestDocumentProcessingWorker_Work(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully processes document", func(t *testing.T) {
		refLetterRepo := newMockReferenceLetterRepository()
		storage := newMockStorage()

		// Create test data
		letterID := uuid.New()
		fileID := uuid.New()
		storageKey := "test/storage/key.pdf"

		letter := &domain.ReferenceLetter{
			ID:     letterID,
			UserID: uuid.New(),
			FileID: &fileID,
			Status: domain.ReferenceLetterStatusPending,
		}
		if err := refLetterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("failed to create letter: %v", err)
		}

		// Add file to storage
		storage.addKey(storageKey)

		worker := job.NewDocumentProcessingWorker(refLetterRepo, storage, testLogger())

		args := job.DocumentProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
		}

		riverJob := &river.Job[job.DocumentProcessingArgs]{
			Args: args,
		}

		err := worker.Work(ctx, riverJob)
		if err != nil {
			t.Fatalf("Work() error = %v", err)
		}

		// Verify status was updated to completed
		updatedLetter, getErr := refLetterRepo.GetByID(ctx, letterID)
		if getErr != nil {
			t.Fatalf("failed to get updated letter: %v", getErr)
		}
		if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
			t.Errorf("letter status = %s, want %s", updatedLetter.Status, domain.ReferenceLetterStatusCompleted)
		}
	})

	t.Run("fails when file not found in storage", func(t *testing.T) {
		refLetterRepo := newMockReferenceLetterRepository()
		storage := newMockStorage()

		letterID := uuid.New()
		fileID := uuid.New()
		storageKey := "nonexistent/key.pdf"

		letter := &domain.ReferenceLetter{
			ID:     letterID,
			UserID: uuid.New(),
			FileID: &fileID,
			Status: domain.ReferenceLetterStatusPending,
		}
		if err := refLetterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("failed to create letter: %v", err)
		}

		// Do NOT add file to storage

		worker := job.NewDocumentProcessingWorker(refLetterRepo, storage, testLogger())

		args := job.DocumentProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
		}

		riverJob := &river.Job[job.DocumentProcessingArgs]{
			Args: args,
		}

		err := worker.Work(ctx, riverJob)
		if err == nil {
			t.Fatal("Work() expected error for missing file")
		}

		// Verify status was updated to failed
		updatedLetter, getErr := refLetterRepo.GetByID(ctx, letterID)
		if getErr != nil {
			t.Fatalf("failed to get updated letter: %v", getErr)
		}
		if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
			t.Errorf("letter status = %s, want %s", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
		}
	})

	t.Run("fails when storage check fails", func(t *testing.T) {
		refLetterRepo := newMockReferenceLetterRepository()
		storage := &errorStorage{}

		letterID := uuid.New()
		fileID := uuid.New()
		storageKey := "test/key.pdf"

		letter := &domain.ReferenceLetter{
			ID:     letterID,
			UserID: uuid.New(),
			FileID: &fileID,
			Status: domain.ReferenceLetterStatusPending,
		}
		if err := refLetterRepo.Create(ctx, letter); err != nil {
			t.Fatalf("failed to create letter: %v", err)
		}

		worker := job.NewDocumentProcessingWorker(refLetterRepo, storage, testLogger())

		args := job.DocumentProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
		}

		riverJob := &river.Job[job.DocumentProcessingArgs]{
			Args: args,
		}

		err := worker.Work(ctx, riverJob)
		if err == nil {
			t.Fatal("Work() expected error for storage failure")
		}

		// Verify status was updated to failed
		updatedLetter, getErr := refLetterRepo.GetByID(ctx, letterID)
		if getErr != nil {
			t.Fatalf("failed to get updated letter: %v", getErr)
		}
		if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
			t.Errorf("letter status = %s, want %s", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
		}
	})

	t.Run("fails when reference letter not found", func(t *testing.T) {
		refLetterRepo := newMockReferenceLetterRepository()
		storage := newMockStorage()

		nonExistentID := uuid.New()
		fileID := uuid.New()
		storageKey := "test/key.pdf"

		storage.addKey(storageKey)

		worker := job.NewDocumentProcessingWorker(refLetterRepo, storage, testLogger())

		args := job.DocumentProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: nonExistentID,
			FileID:            fileID,
		}

		riverJob := &river.Job[job.DocumentProcessingArgs]{
			Args: args,
		}

		err := worker.Work(ctx, riverJob)
		if err == nil {
			t.Fatal("Work() expected error for missing reference letter")
		}
	})
}
