//nolint:goconst // Test file - string constants are fine inline
package job_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/domain"
	"backend/internal/job"
)

// mockDetectionExtractor implements domain.DocumentExtractor for detection tests.
type mockDetectionExtractor struct {
	extractTextResult string
	extractTextErr    error
	detectionResult   *domain.DocumentDetectionResult
	detectionErr      error
}

func (e *mockDetectionExtractor) ExtractText(_ context.Context, _ []byte, _ string) (string, error) {
	return e.extractTextResult, e.extractTextErr
}

func (e *mockDetectionExtractor) ExtractResumeData(_ context.Context, _ string) (*domain.ResumeExtractedData, error) {
	return nil, nil
}

func (e *mockDetectionExtractor) ExtractLetterData(_ context.Context, _ string, _ []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	return nil, nil
}

func (e *mockDetectionExtractor) DetectDocumentContent(_ context.Context, _ string) (*domain.DocumentDetectionResult, error) {
	return e.detectionResult, e.detectionErr
}

func TestDocumentDetectionArgs_Kind(t *testing.T) {
	args := job.DocumentDetectionArgs{}
	if args.Kind() != "document_detection" {
		t.Errorf("Kind() = %q, want %q", args.Kind(), "document_detection")
	}
}

func TestDocumentDetectionArgs_InsertOpts(t *testing.T) {
	args := job.DocumentDetectionArgs{}
	opts := args.InsertOpts()
	if opts.MaxAttempts != 2 {
		t.Errorf("MaxAttempts = %d, want 2", opts.MaxAttempts)
	}
}

func TestDocumentDetectionWorker_Timeout(t *testing.T) {
	worker := job.NewDocumentDetectionWorker(
		newMockFileRepository(), newMockDownloadStorage(nil),
		&mockDetectionExtractor{}, testLogger(),
	)
	timeout := worker.Timeout(nil)
	if timeout != 5*time.Minute {
		t.Errorf("Timeout = %v, want 5m", timeout)
	}
}

func TestDocumentDetectionWorker_Success(t *testing.T) {
	ctx := context.Background()
	fileRepo := newMockFileRepository()
	fileID := uuid.New()
	userID := uuid.New()

	pending := domain.DetectionStatusPending
	if err := fileRepo.Create(ctx, &domain.File{
		ID:              fileID,
		UserID:          userID,
		Filename:        "resume.pdf",
		ContentType:     "application/pdf",
		StorageKey:      "test/key.pdf",
		DetectionStatus: &pending,
	}); err != nil {
		t.Fatalf("Create file: %v", err)
	}

	detection := &domain.DocumentDetectionResult{
		HasCareerInfo:    true,
		HasTestimonial:   false,
		Confidence:       0.95,
		Summary:          "A resume document",
		DocumentTypeHint: domain.DocumentTypeResume,
	}

	extractor := &mockDetectionExtractor{
		extractTextResult: "John Doe, Software Engineer",
		detectionResult:   detection,
	}

	worker := job.NewDocumentDetectionWorker(
		fileRepo, newMockDownloadStorage([]byte("pdf data")),
		extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentDetectionArgs]{
		Args: job.DocumentDetectionArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
			UserID:      userID,
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	// Verify file detection fields
	file, err := fileRepo.GetByID(ctx, fileID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if file.DetectionStatus == nil || *file.DetectionStatus != domain.DetectionStatusCompleted {
		t.Errorf("detection_status = %v, want %q", file.DetectionStatus, domain.DetectionStatusCompleted)
	}
	if file.DetectionResult == nil {
		t.Fatal("detection_result is nil, want non-nil")
	}
	if file.DetectionError != nil {
		t.Errorf("detection_error = %v, want nil", file.DetectionError)
	}

	// Verify the stored JSON can be deserialized
	var stored domain.DocumentDetectionResult
	if err := json.Unmarshal(file.DetectionResult, &stored); err != nil {
		t.Fatalf("Failed to unmarshal detection_result: %v", err)
	}
	if !stored.HasCareerInfo {
		t.Error("stored HasCareerInfo = false, want true")
	}
	if stored.DocumentTypeHint != domain.DocumentTypeResume {
		t.Errorf("stored DocumentTypeHint = %q, want %q", stored.DocumentTypeHint, domain.DocumentTypeResume)
	}
}

func TestDocumentDetectionWorker_StorageDownloadFailure(t *testing.T) {
	ctx := context.Background()
	fileRepo := newMockFileRepository()
	fileID := uuid.New()

	pending := domain.DetectionStatusPending
	if err := fileRepo.Create(ctx, &domain.File{
		ID:              fileID,
		UserID:          uuid.New(),
		Filename:        "doc.pdf",
		ContentType:     "application/pdf",
		StorageKey:      "test/key.pdf",
		DetectionStatus: &pending,
	}); err != nil {
		t.Fatalf("Create file: %v", err)
	}

	storage := newMockDownloadStorage(nil)
	storage.dlErr = errors.New("storage unavailable")

	worker := job.NewDocumentDetectionWorker(
		fileRepo, storage,
		&mockDetectionExtractor{}, testLogger(),
	)

	riverJob := &river.Job[job.DocumentDetectionArgs]{
		Args: job.DocumentDetectionArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
		},
	}

	err := worker.Work(ctx, riverJob)
	if err == nil {
		t.Fatal("Work() expected error, got nil")
	}

	file, err := fileRepo.GetByID(ctx, fileID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if file.DetectionStatus == nil || *file.DetectionStatus != domain.DetectionStatusFailed {
		t.Errorf("detection_status = %v, want %q", file.DetectionStatus, domain.DetectionStatusFailed)
	}
	if file.DetectionError == nil {
		t.Error("detection_error is nil, want non-nil")
	}
}

func TestDocumentDetectionWorker_ExtractTextFailure(t *testing.T) {
	ctx := context.Background()
	fileRepo := newMockFileRepository()
	fileID := uuid.New()

	pending := domain.DetectionStatusPending
	if err := fileRepo.Create(ctx, &domain.File{
		ID:              fileID,
		UserID:          uuid.New(),
		Filename:        "doc.pdf",
		ContentType:     "application/pdf",
		StorageKey:      "test/key.pdf",
		DetectionStatus: &pending,
	}); err != nil {
		t.Fatalf("Create file: %v", err)
	}

	extractor := &mockDetectionExtractor{
		extractTextErr: errors.New("LLM unavailable"),
	}

	worker := job.NewDocumentDetectionWorker(
		fileRepo, newMockDownloadStorage([]byte("pdf data")),
		extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentDetectionArgs]{
		Args: job.DocumentDetectionArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
		},
	}

	err := worker.Work(ctx, riverJob)
	if err == nil {
		t.Fatal("Work() expected error, got nil")
	}

	file, err := fileRepo.GetByID(ctx, fileID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if file.DetectionStatus == nil || *file.DetectionStatus != domain.DetectionStatusFailed {
		t.Errorf("detection_status = %v, want %q", file.DetectionStatus, domain.DetectionStatusFailed)
	}
}

func TestDocumentDetectionWorker_DetectionFailure(t *testing.T) {
	ctx := context.Background()
	fileRepo := newMockFileRepository()
	fileID := uuid.New()

	pending := domain.DetectionStatusPending
	if err := fileRepo.Create(ctx, &domain.File{
		ID:              fileID,
		UserID:          uuid.New(),
		Filename:        "doc.pdf",
		ContentType:     "application/pdf",
		StorageKey:      "test/key.pdf",
		DetectionStatus: &pending,
	}); err != nil {
		t.Fatalf("Create file: %v", err)
	}

	extractor := &mockDetectionExtractor{
		extractTextResult: "some text",
		detectionErr:      errors.New("classification failed"),
	}

	worker := job.NewDocumentDetectionWorker(
		fileRepo, newMockDownloadStorage([]byte("pdf data")),
		extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentDetectionArgs]{
		Args: job.DocumentDetectionArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
		},
	}

	err := worker.Work(ctx, riverJob)
	if err == nil {
		t.Fatal("Work() expected error, got nil")
	}

	file, err := fileRepo.GetByID(ctx, fileID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if file.DetectionStatus == nil || *file.DetectionStatus != domain.DetectionStatusFailed {
		t.Errorf("detection_status = %v, want %q", file.DetectionStatus, domain.DetectionStatusFailed)
	}
}
