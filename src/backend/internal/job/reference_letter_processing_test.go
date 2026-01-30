//nolint:goconst // Test file - string constants are fine inline
package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/domain"
	"backend/internal/logger"
)

// mockReferenceLetterRepository implements domain.ReferenceLetterRepository for testing.
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
	if _, exists := r.letters[letter.ID]; !exists {
		return fmt.Errorf("letter not found: %s", letter.ID)
	}
	r.letters[letter.ID] = letter
	return nil
}

func (r *mockReferenceLetterRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.letters, id)
	return nil
}

// mockFileRepository implements domain.FileRepository for testing.
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

func (r *mockFileRepository) GetByStorageKey(_ context.Context, _ string) (*domain.File, error) {
	return nil, nil
}

func (r *mockFileRepository) Update(_ context.Context, file *domain.File) error {
	r.files[file.ID] = file
	return nil
}

func (r *mockFileRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.files, id)
	return nil
}

func (r *mockFileRepository) GetByUserID(_ context.Context, _ uuid.UUID) ([]*domain.File, error) {
	return nil, nil
}

// mockStorage implements domain.Storage for testing.
type mockStorage struct {
	data map[string][]byte
}

func newMockStorage() *mockStorage {
	return &mockStorage{data: make(map[string][]byte)}
}

func (s *mockStorage) Upload(_ context.Context, key string, reader io.Reader, _ int64, _ string) (*domain.StorageObject, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	s.data[key] = data
	return &domain.StorageObject{Key: key}, nil
}

func (s *mockStorage) Download(_ context.Context, key string) (io.ReadCloser, error) {
	data, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", key)
	}
	return io.NopCloser(strings.NewReader(string(data))), nil
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
	_, ok := s.data[key]
	return ok, nil
}

// mockDocumentExtractor implements domain.DocumentExtractor for testing.
type mockDocumentExtractor struct {
	extractTextResult  string
	extractTextError   error
	extractLetterData  *domain.ExtractedLetterData
	extractLetterError error
}

func (e *mockDocumentExtractor) ExtractText(_ context.Context, _ []byte, _ string) (string, error) {
	return e.extractTextResult, e.extractTextError
}

func (e *mockDocumentExtractor) ExtractResumeData(_ context.Context, _ string) (*domain.ResumeExtractedData, error) {
	return nil, nil
}

func (e *mockDocumentExtractor) ExtractLetterData(_ context.Context, _ string) (*domain.ExtractedLetterData, error) {
	return e.extractLetterData, e.extractLetterError
}

// mockLogger implements logger.Logger for testing.
type mockLogger struct{}

func (l *mockLogger) Debug(_ string, _ ...logger.Attr)    {}
func (l *mockLogger) Info(_ string, _ ...logger.Attr)     {}
func (l *mockLogger) Warning(_ string, _ ...logger.Attr)  {}
func (l *mockLogger) Error(_ string, _ ...logger.Attr)    {}
func (l *mockLogger) Critical(_ string, _ ...logger.Attr) {}

func newTestLetterWorker() (*ReferenceLetterProcessingWorker, *mockReferenceLetterRepository, *mockFileRepository, *mockStorage, *mockDocumentExtractor) {
	letterRepo := newMockReferenceLetterRepository()
	fileRepo := newMockFileRepository()
	storage := newMockStorage()
	extractor := &mockDocumentExtractor{}

	worker := NewReferenceLetterProcessingWorker(
		letterRepo,
		fileRepo,
		storage,
		extractor,
		&mockLogger{},
	)

	return worker, letterRepo, fileRepo, storage, extractor
}

func TestReferenceLetterProcessingArgs_Kind(t *testing.T) {
	args := ReferenceLetterProcessingArgs{}
	if args.Kind() != "reference_letter_processing" {
		t.Errorf("Kind() = %q, want %q", args.Kind(), "reference_letter_processing")
	}
}

func TestReferenceLetterProcessingWorker_Work_Success(t *testing.T) {
	worker, letterRepo, fileRepo, storage, extractor := newTestLetterWorker()

	// Set up test data
	letterID := uuid.New()
	fileID := uuid.New()
	userID := uuid.New()
	storageKey := "test/letter.pdf"

	// Create reference letter
	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: userID,
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	letterRepo.letters[letterID] = letter

	// Create file record
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	fileRepo.files[fileID] = file

	// Upload test file
	storage.data[storageKey] = []byte("fake pdf content")

	// Set up extractor response
	authorTitle := "Engineering Manager"
	authorCompany := "Acme Corp"
	skillContext := "technical skills"
	extractor.extractTextResult = "This is a reference letter for Jane..."
	extractor.extractLetterData = &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "John Smith",
			Title:        &authorTitle,
			Company:      &authorCompany,
			Relationship: domain.AuthorRelationshipManager,
		},
		Testimonials: []domain.ExtractedTestimonial{
			{Quote: "Jane's leadership was exceptional.", SkillsMentioned: []string{"leadership"}},
		},
		SkillMentions: []domain.ExtractedSkillMention{
			{Skill: "Go", Quote: "Her Go expertise helped us...", Context: &skillContext},
		},
		ExperienceMentions: []domain.ExtractedExperienceMention{
			{Company: "Acme Corp", Role: "Senior Engineer", Quote: "During her time as Senior Engineer..."},
		},
		DiscoveredSkills: []string{"mentoring", "system design"},
	}

	// Create job
	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	// Execute
	err := worker.Work(context.Background(), job)
	if err != nil {
		t.Fatalf("Work() returned error: %v", err)
	}

	// Verify status is completed
	updatedLetter := letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusCompleted)
	}

	// Verify extracted data is saved
	if updatedLetter.ExtractedData == nil {
		t.Fatal("ExtractedData should not be nil")
	}

	var savedData domain.ExtractedLetterData
	if err := json.Unmarshal(updatedLetter.ExtractedData, &savedData); err != nil {
		t.Fatalf("Failed to unmarshal ExtractedData: %v", err)
	}

	if savedData.Author.Name != "John Smith" {
		t.Errorf("Author.Name = %q, want %q", savedData.Author.Name, "John Smith")
	}

	// Verify author fields are populated on the letter
	if updatedLetter.AuthorName == nil || *updatedLetter.AuthorName != "John Smith" {
		t.Errorf("AuthorName = %v, want %q", updatedLetter.AuthorName, "John Smith")
	}
	if updatedLetter.AuthorTitle == nil || *updatedLetter.AuthorTitle != "Engineering Manager" {
		t.Errorf("AuthorTitle = %v, want %q", updatedLetter.AuthorTitle, "Engineering Manager")
	}
	if updatedLetter.Organization == nil || *updatedLetter.Organization != "Acme Corp" {
		t.Errorf("Organization = %v, want %q", updatedLetter.Organization, "Acme Corp")
	}
}

func TestReferenceLetterProcessingWorker_Work_ExtractTextFails(t *testing.T) {
	worker, letterRepo, fileRepo, storage, extractor := newTestLetterWorker()

	// Set up test data
	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	letterRepo.letters[letterID] = letter

	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	fileRepo.files[fileID] = file

	storage.data[storageKey] = []byte("fake pdf content")

	// Set extractor to fail
	extractor.extractTextError = fmt.Errorf("OCR failed")

	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err == nil {
		t.Fatal("expected error when text extraction fails")
	}

	// Verify status is failed
	updatedLetter := letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
	if updatedLetter.ErrorMessage == nil {
		t.Error("ErrorMessage should be set")
	}
}

func TestReferenceLetterProcessingWorker_Work_LetterExtractionFails(t *testing.T) {
	worker, letterRepo, fileRepo, storage, extractor := newTestLetterWorker()

	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	letterRepo.letters[letterID] = letter

	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	fileRepo.files[fileID] = file

	storage.data[storageKey] = []byte("fake pdf content")

	// Text extraction succeeds but structured extraction fails
	extractor.extractTextResult = "This is a reference letter..."
	extractor.extractLetterError = fmt.Errorf("LLM parsing failed")

	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err == nil {
		t.Fatal("expected error when letter extraction fails")
	}

	updatedLetter := letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
}

func TestReferenceLetterProcessingWorker_Work_StorageDownloadFails(t *testing.T) {
	worker, letterRepo, fileRepo, _, extractor := newTestLetterWorker()

	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	letterRepo.letters[letterID] = letter

	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	fileRepo.files[fileID] = file

	// Don't upload any file to storage - download will fail
	extractor.extractTextResult = "unused"

	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err == nil {
		t.Fatal("expected error when storage download fails")
	}

	updatedLetter := letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
}

func TestReferenceLetterProcessingWorker_Work_LetterNotFound(t *testing.T) {
	worker, _, _, _, _ := newTestLetterWorker()

	// Don't create the letter - it won't be found
	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        "test/letter.pdf",
			ReferenceLetterID: uuid.New(),
			FileID:            uuid.New(),
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err == nil {
		t.Fatal("expected error when letter not found")
	}
}

func TestReferenceLetterProcessingWorker_Work_FallbackToFileContentType(t *testing.T) {
	worker, letterRepo, fileRepo, storage, extractor := newTestLetterWorker()

	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	letterRepo.letters[letterID] = letter

	// File has the content type
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	fileRepo.files[fileID] = file

	storage.data[storageKey] = []byte("fake pdf content")

	extractor.extractTextResult = "Letter text"
	extractor.extractLetterData = &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "Test Author",
			Relationship: domain.AuthorRelationshipOther,
		},
	}

	// Don't pass content type in args - should fall back to file record
	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "", // Empty - should use file record
		},
	}

	err := worker.Work(context.Background(), job)
	if err != nil {
		t.Fatalf("Work() returned error: %v", err)
	}

	// Verify it completed successfully
	updatedLetter := letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusCompleted)
	}
}

func TestReferenceLetterProcessingWorker_updateStatus(t *testing.T) {
	worker, letterRepo, _, _, _ := newTestLetterWorker()

	letterID := uuid.New()
	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		Status: domain.ReferenceLetterStatusPending,
	}
	letterRepo.letters[letterID] = letter

	// Test updating to processing
	err := worker.updateStatus(context.Background(), letterID, domain.ReferenceLetterStatusProcessing, nil)
	if err != nil {
		t.Fatalf("updateStatus() returned error: %v", err)
	}

	if letterRepo.letters[letterID].Status != domain.ReferenceLetterStatusProcessing {
		t.Errorf("Status = %q, want %q", letterRepo.letters[letterID].Status, domain.ReferenceLetterStatusProcessing)
	}

	// Test updating to failed with error message
	errMsg := "Something went wrong"
	err = worker.updateStatus(context.Background(), letterID, domain.ReferenceLetterStatusFailed, &errMsg)
	if err != nil {
		t.Fatalf("updateStatus() returned error: %v", err)
	}

	if letterRepo.letters[letterID].Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", letterRepo.letters[letterID].Status, domain.ReferenceLetterStatusFailed)
	}
	if letterRepo.letters[letterID].ErrorMessage == nil || *letterRepo.letters[letterID].ErrorMessage != errMsg {
		t.Errorf("ErrorMessage = %v, want %q", letterRepo.letters[letterID].ErrorMessage, errMsg)
	}
}

func TestReferenceLetterProcessingWorker_updateStatus_NotFound(t *testing.T) {
	worker, _, _, _, _ := newTestLetterWorker()

	// Try to update non-existent letter
	err := worker.updateStatus(context.Background(), uuid.New(), domain.ReferenceLetterStatusProcessing, nil)
	if err == nil {
		t.Fatal("expected error when letter not found")
	}
}
