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
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

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

func (r *mockFileRepository) GetByUserIDAndContentHash(_ context.Context, userID uuid.UUID, contentHash string) (*domain.File, error) {
	for _, file := range r.files {
		if file.UserID == userID && file.ContentHash != nil && *file.ContentHash == contentHash {
			return file, nil
		}
	}
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
	extractTextResult    string
	extractTextError     error
	extractLetterData    *domain.ExtractedLetterData
	extractLetterError   error
	lastProfileSkillsCtx []domain.ProfileSkillContext // captures what was passed to ExtractLetterData
}

func (e *mockDocumentExtractor) ExtractText(_ context.Context, _ []byte, _ string) (string, error) {
	return e.extractTextResult, e.extractTextError
}

func (e *mockDocumentExtractor) ExtractResumeData(_ context.Context, _ string) (*domain.ResumeExtractedData, error) {
	return nil, nil
}

func (e *mockDocumentExtractor) ExtractLetterData(_ context.Context, _ string, profileSkills []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	e.lastProfileSkillsCtx = profileSkills
	return e.extractLetterData, e.extractLetterError
}

// mockAuthorRepository implements domain.AuthorRepository for testing.
// Note: mockProfileRepository and mockProfileSkillRepository are defined in resume_processing_test.go
type mockAuthorRepository struct {
	authors map[uuid.UUID]*domain.Author
}

func newMockAuthorRepository() *mockAuthorRepository {
	return &mockAuthorRepository{authors: make(map[uuid.UUID]*domain.Author)}
}

func (r *mockAuthorRepository) Create(_ context.Context, author *domain.Author) error {
	if author.ID == uuid.Nil {
		author.ID = uuid.New()
	}
	r.authors[author.ID] = author
	return nil
}

func (r *mockAuthorRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Author, error) {
	return r.authors[id], nil
}

func (r *mockAuthorRepository) GetByProfileID(_ context.Context, profileID uuid.UUID) ([]*domain.Author, error) {
	var result []*domain.Author
	for _, a := range r.authors {
		if a.ProfileID == profileID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *mockAuthorRepository) FindByNameAndCompany(_ context.Context, profileID uuid.UUID, name string, company *string) (*domain.Author, error) {
	for _, a := range r.authors {
		if a.ProfileID == profileID && a.Name == name {
			if company == nil && a.Company == nil {
				return a, nil
			}
			if company != nil && a.Company != nil && *company == *a.Company {
				return a, nil
			}
		}
	}
	return nil, nil
}

func (r *mockAuthorRepository) Update(_ context.Context, author *domain.Author) error {
	r.authors[author.ID] = author
	return nil
}

func (r *mockAuthorRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.authors, id)
	return nil
}

// mockTestimonialRepository implements domain.TestimonialRepository for testing.
type mockTestimonialRepository struct {
	testimonials map[uuid.UUID]*domain.Testimonial
}

func newMockTestimonialRepository() *mockTestimonialRepository {
	return &mockTestimonialRepository{testimonials: make(map[uuid.UUID]*domain.Testimonial)}
}

func (r *mockTestimonialRepository) Create(_ context.Context, testimonial *domain.Testimonial) error {
	if testimonial.ID == uuid.Nil {
		testimonial.ID = uuid.New()
	}
	r.testimonials[testimonial.ID] = testimonial
	return nil
}

func (r *mockTestimonialRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Testimonial, error) {
	return r.testimonials[id], nil
}

func (r *mockTestimonialRepository) GetByProfileID(_ context.Context, profileID uuid.UUID) ([]*domain.Testimonial, error) {
	var result []*domain.Testimonial
	for _, t := range r.testimonials {
		if t.ProfileID == profileID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *mockTestimonialRepository) GetByReferenceLetterID(_ context.Context, referenceLetterID uuid.UUID) ([]*domain.Testimonial, error) {
	var result []*domain.Testimonial
	for _, t := range r.testimonials {
		if t.ReferenceLetterID == referenceLetterID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *mockTestimonialRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.testimonials, id)
	return nil
}

func (r *mockTestimonialRepository) DeleteByReferenceLetterID(_ context.Context, referenceLetterID uuid.UUID) error {
	for id, t := range r.testimonials {
		if t.ReferenceLetterID == referenceLetterID {
			delete(r.testimonials, id)
		}
	}
	return nil
}

// mockSkillValidationRepository implements domain.SkillValidationRepository for testing.
type mockSkillValidationRepository struct {
	validations map[uuid.UUID]*domain.SkillValidation
}

func newMockSkillValidationRepository() *mockSkillValidationRepository {
	return &mockSkillValidationRepository{validations: make(map[uuid.UUID]*domain.SkillValidation)}
}

func (r *mockSkillValidationRepository) Create(_ context.Context, validation *domain.SkillValidation) error {
	if validation.ID == uuid.Nil {
		validation.ID = uuid.New()
	}
	r.validations[validation.ID] = validation
	return nil
}

func (r *mockSkillValidationRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.SkillValidation, error) {
	return r.validations[id], nil
}

func (r *mockSkillValidationRepository) GetByProfileSkillID(_ context.Context, profileSkillID uuid.UUID) ([]*domain.SkillValidation, error) {
	var result []*domain.SkillValidation
	for _, v := range r.validations {
		if v.ProfileSkillID == profileSkillID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockSkillValidationRepository) GetByReferenceLetterID(_ context.Context, referenceLetterID uuid.UUID) ([]*domain.SkillValidation, error) {
	var result []*domain.SkillValidation
	for _, v := range r.validations {
		if v.ReferenceLetterID == referenceLetterID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockSkillValidationRepository) GetByTestimonialID(_ context.Context, testimonialID uuid.UUID) ([]*domain.SkillValidation, error) {
	var result []*domain.SkillValidation
	for _, v := range r.validations {
		if v.TestimonialID != nil && *v.TestimonialID == testimonialID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockSkillValidationRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.validations, id)
	return nil
}

func (r *mockSkillValidationRepository) DeleteByReferenceLetterID(_ context.Context, referenceLetterID uuid.UUID) error {
	for id, v := range r.validations {
		if v.ReferenceLetterID == referenceLetterID {
			delete(r.validations, id)
		}
	}
	return nil
}

func (r *mockSkillValidationRepository) CountByProfileSkillID(_ context.Context, profileSkillID uuid.UUID) (int, error) {
	count := 0
	for _, v := range r.validations {
		if v.ProfileSkillID == profileSkillID {
			count++
		}
	}
	return count, nil
}

// mockLogger implements logger.Logger for testing.
type mockLogger struct{}

func (l *mockLogger) Debug(_ string, _ ...logger.Attr)    {}
func (l *mockLogger) Info(_ string, _ ...logger.Attr)     {}
func (l *mockLogger) Warning(_ string, _ ...logger.Attr)  {}
func (l *mockLogger) Error(_ string, _ ...logger.Attr)    {}
func (l *mockLogger) Critical(_ string, _ ...logger.Attr) {}

type testLetterWorkerMocks struct {
	letterRepo          *mockReferenceLetterRepository
	fileRepo            *mockFileRepository
	profileRepo         *mockProfileRepository
	profileSkillRepo    *mockProfileSkillRepository
	authorRepo          *mockAuthorRepository
	testimonialRepo     *mockTestimonialRepository
	skillValidationRepo *mockSkillValidationRepository
	storage             *mockStorage
	extractor           *mockDocumentExtractor
}

func newTestLetterWorker() (*ReferenceLetterProcessingWorker, *testLetterWorkerMocks) {
	mocks := &testLetterWorkerMocks{
		letterRepo:          newMockReferenceLetterRepository(),
		fileRepo:            newMockFileRepository(),
		profileRepo:         newMockProfileRepository(),
		profileSkillRepo:    newMockProfileSkillRepository(),
		authorRepo:          newMockAuthorRepository(),
		testimonialRepo:     newMockTestimonialRepository(),
		skillValidationRepo: newMockSkillValidationRepository(),
		storage:             newMockStorage(),
		extractor:           &mockDocumentExtractor{},
	}

	worker := NewReferenceLetterProcessingWorker(
		mocks.letterRepo,
		mocks.fileRepo,
		mocks.profileRepo,
		mocks.profileSkillRepo,
		mocks.authorRepo,
		mocks.testimonialRepo,
		mocks.skillValidationRepo,
		mocks.storage,
		mocks.extractor,
		&mockLogger{},
	)

	return worker, mocks
}

func TestReferenceLetterProcessingArgs_Kind(t *testing.T) {
	args := ReferenceLetterProcessingArgs{}
	if args.Kind() != "reference_letter_processing" {
		t.Errorf("Kind() = %q, want %q", args.Kind(), "reference_letter_processing")
	}
}

func TestReferenceLetterProcessingWorker_Work_Success(t *testing.T) {
	worker, mocks := newTestLetterWorker()

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
	mocks.letterRepo.letters[letterID] = letter

	// Create file record
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	// Upload test file
	mocks.storage.data[storageKey] = []byte("fake pdf content")

	// Set up extractor response
	authorTitle := "Engineering Manager"
	authorCompany := "Acme Corp"
	skillContext := "technical skills"
	leadershipContext := "leadership"
	mocks.extractor.extractTextResult = "This is a reference letter for Jane..."
	mocks.extractor.extractLetterData = &domain.ExtractedLetterData{
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
		DiscoveredSkills: []domain.DiscoveredSkill{
			{Skill: "mentoring", Quote: "She mentored junior developers...", Context: &leadershipContext},
			{Skill: "system design", Quote: "Her system design skills...", Context: &skillContext},
		},
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
	updatedLetter := mocks.letterRepo.letters[letterID]
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
	worker, mocks := newTestLetterWorker()

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
	mocks.letterRepo.letters[letterID] = letter

	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	mocks.storage.data[storageKey] = []byte("fake pdf content")

	// Set extractor to fail
	mocks.extractor.extractTextError = fmt.Errorf("OCR failed")

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
	updatedLetter := mocks.letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
	if updatedLetter.ErrorMessage == nil {
		t.Error("ErrorMessage should be set")
	}
}

func TestReferenceLetterProcessingWorker_Work_LetterExtractionFails(t *testing.T) {
	worker, mocks := newTestLetterWorker()

	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	mocks.storage.data[storageKey] = []byte("fake pdf content")

	// Text extraction succeeds but structured extraction fails
	mocks.extractor.extractTextResult = "This is a reference letter..."
	mocks.extractor.extractLetterError = fmt.Errorf("LLM parsing failed")

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

	updatedLetter := mocks.letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
}

func TestReferenceLetterProcessingWorker_Work_StorageDownloadFails(t *testing.T) {
	worker, mocks := newTestLetterWorker()

	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	// Don't upload any file to storage - download will fail
	mocks.extractor.extractTextResult = "unused"

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

	updatedLetter := mocks.letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
}

func TestReferenceLetterProcessingWorker_Work_LetterNotFound(t *testing.T) {
	worker, _ := newTestLetterWorker()

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
	worker, mocks := newTestLetterWorker()

	letterID := uuid.New()
	fileID := uuid.New()
	storageKey := "test/letter.pdf"

	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	// File has the content type
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	mocks.storage.data[storageKey] = []byte("fake pdf content")

	mocks.extractor.extractTextResult = "Letter text"
	mocks.extractor.extractLetterData = &domain.ExtractedLetterData{
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
	updatedLetter := mocks.letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("Status = %q, want %q", updatedLetter.Status, domain.ReferenceLetterStatusCompleted)
	}
}

func TestReferenceLetterProcessingWorker_updateStatus(t *testing.T) {
	worker, mocks := newTestLetterWorker()

	letterID := uuid.New()
	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: uuid.New(),
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	// Test updating to processing
	err := worker.updateStatus(context.Background(), letterID, domain.ReferenceLetterStatusProcessing, nil)
	if err != nil {
		t.Fatalf("updateStatus() returned error: %v", err)
	}

	if mocks.letterRepo.letters[letterID].Status != domain.ReferenceLetterStatusProcessing {
		t.Errorf("Status = %q, want %q", mocks.letterRepo.letters[letterID].Status, domain.ReferenceLetterStatusProcessing)
	}

	// Test updating to failed with error message
	errMsg := "Something went wrong"
	err = worker.updateStatus(context.Background(), letterID, domain.ReferenceLetterStatusFailed, &errMsg)
	if err != nil {
		t.Fatalf("updateStatus() returned error: %v", err)
	}

	if mocks.letterRepo.letters[letterID].Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("Status = %q, want %q", mocks.letterRepo.letters[letterID].Status, domain.ReferenceLetterStatusFailed)
	}
	if mocks.letterRepo.letters[letterID].ErrorMessage == nil || *mocks.letterRepo.letters[letterID].ErrorMessage != errMsg {
		t.Errorf("ErrorMessage = %v, want %q", mocks.letterRepo.letters[letterID].ErrorMessage, errMsg)
	}
}

func TestReferenceLetterProcessingWorker_updateStatus_NotFound(t *testing.T) {
	worker, _ := newTestLetterWorker()

	// Try to update non-existent letter
	err := worker.updateStatus(context.Background(), uuid.New(), domain.ReferenceLetterStatusProcessing, nil)
	if err == nil {
		t.Fatal("expected error when letter not found")
	}
}

func TestReferenceLetterProcessingWorker_PassesProfileSkillsToExtractor(t *testing.T) {
	worker, mocks := newTestLetterWorker()

	// Set up test data
	letterID := uuid.New()
	fileID := uuid.New()
	userID := uuid.New()
	profileID := uuid.New()
	storageKey := "test/letter.pdf"

	// Create profile with existing skills
	profile := &domain.Profile{ID: profileID, UserID: userID}
	mocks.profileRepo.profiles[profileID] = profile

	goSkillID := uuid.New()
	mocks.profileSkillRepo.skills[goSkillID] = &domain.ProfileSkill{
		ID:             goSkillID,
		ProfileID:      profileID,
		Name:           "Go",
		NormalizedName: "go",
		Category:       "TECHNICAL",
	}

	// Create reference letter
	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: userID,
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	// Create file record
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	mocks.storage.data[storageKey] = []byte("fake pdf content")

	// Set up extractor response
	mocks.extractor.extractTextResult = "This is a reference letter..."
	mocks.extractor.extractLetterData = &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "Test Author",
			Relationship: domain.AuthorRelationshipOther,
		},
	}

	// Create and execute job
	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err != nil {
		t.Fatalf("Work() returned error: %v", err)
	}

	// Verify profile skills were passed to extractor
	if mocks.extractor.lastProfileSkillsCtx == nil {
		t.Fatal("Expected profile skills context to be passed to extractor")
	}
	if len(mocks.extractor.lastProfileSkillsCtx) != 1 {
		t.Errorf("Expected 1 profile skill in context, got %d", len(mocks.extractor.lastProfileSkillsCtx))
	}
	if mocks.extractor.lastProfileSkillsCtx[0].Name != "Go" {
		t.Errorf("Expected skill name 'Go', got %q", mocks.extractor.lastProfileSkillsCtx[0].Name)
	}
}

// TestReferenceLetterProcessingWorker_StoresExtractedDataOnly verifies that the job
// extracts data and stores it as JSON, but does NOT create testimonials or validations.
// Record creation is handled by the applyReferenceLetterValidations mutation based on
// user selection, not automatically during processing.
func TestReferenceLetterProcessingWorker_StoresExtractedDataOnly(t *testing.T) {
	worker, mocks := newTestLetterWorker()

	// Set up test data
	letterID := uuid.New()
	fileID := uuid.New()
	userID := uuid.New()
	profileID := uuid.New()
	storageKey := "test/letter.pdf"

	// Create profile with existing skills
	profile := &domain.Profile{ID: profileID, UserID: userID}
	mocks.profileRepo.profiles[profileID] = profile

	goSkillID := uuid.New()
	mocks.profileSkillRepo.skills[goSkillID] = &domain.ProfileSkill{
		ID:             goSkillID,
		ProfileID:      profileID,
		Name:           "Go",
		NormalizedName: "go",
		Category:       "TECHNICAL",
	}

	// Create reference letter
	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: userID,
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	// Create file record
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	mocks.storage.data[storageKey] = []byte("fake pdf content")

	// Set up extractor response with testimonials mentioning skills
	authorTitle := "Engineering Manager"
	authorCompany := "Acme Corp"
	skillContext := "technical skills"
	mocks.extractor.extractTextResult = "This is a reference letter..."
	mocks.extractor.extractLetterData = &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "John Smith",
			Title:        &authorTitle,
			Company:      &authorCompany,
			Relationship: domain.AuthorRelationshipManager,
		},
		Testimonials: []domain.ExtractedTestimonial{
			{Quote: "Their leadership was exceptional.", SkillsMentioned: []string{"leadership"}},
			{Quote: "Their Go expertise helped us greatly.", SkillsMentioned: []string{"Go"}},
		},
		SkillMentions: []domain.ExtractedSkillMention{
			{Skill: "Go", Quote: "Deep Go expertise...", Context: &skillContext},
		},
	}

	// Execute job
	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err != nil {
		t.Fatalf("Work() returned error: %v", err)
	}

	// Verify extracted data is stored in JSON (status becomes completed)
	updatedLetter := mocks.letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("Letter status = %q, want 'completed'", updatedLetter.Status)
	}
	if len(updatedLetter.ExtractedData) == 0 {
		t.Error("Expected ExtractedData to be populated")
	}

	// Verify NO testimonials were created (they're created by apply mutation)
	if len(mocks.testimonialRepo.testimonials) != 0 {
		t.Errorf("Expected 0 testimonials (job should not create records), got %d", len(mocks.testimonialRepo.testimonials))
	}

	// Verify NO authors were created
	if len(mocks.authorRepo.authors) != 0 {
		t.Errorf("Expected 0 authors (job should not create records), got %d", len(mocks.authorRepo.authors))
	}

	// Verify NO skill validations were created
	if len(mocks.skillValidationRepo.validations) != 0 {
		t.Errorf("Expected 0 skill validations (job should not create records), got %d", len(mocks.skillValidationRepo.validations))
	}
}

// TestReferenceLetterProcessingWorker_DoesNotCreateDiscoveredSkills verifies that the job
// extracts discovered skills but does NOT create ProfileSkill records.
// Discovered skill records are created by the applyReferenceLetterValidations mutation.
func TestReferenceLetterProcessingWorker_DoesNotCreateDiscoveredSkills(t *testing.T) {
	worker, mocks := newTestLetterWorker()

	// Set up test data
	letterID := uuid.New()
	fileID := uuid.New()
	userID := uuid.New()
	profileID := uuid.New()
	storageKey := "test/letter.pdf"

	// Create profile with NO existing skills
	profile := &domain.Profile{ID: profileID, UserID: userID}
	mocks.profileRepo.profiles[profileID] = profile

	// Create reference letter
	letter := &domain.ReferenceLetter{
		ID:     letterID,
		UserID: userID,
		FileID: &fileID,
		Status: domain.ReferenceLetterStatusPending,
	}
	mocks.letterRepo.letters[letterID] = letter

	// Create file record
	file := &domain.File{
		ID:          fileID,
		StorageKey:  storageKey,
		ContentType: "application/pdf",
	}
	mocks.fileRepo.files[fileID] = file

	mocks.storage.data[storageKey] = []byte("fake pdf content")

	// Set up extractor response with discovered skills
	technicalContext := "technical programming"
	softContext := "leadership skills"
	mocks.extractor.extractTextResult = "This is a reference letter..."
	mocks.extractor.extractLetterData = &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "Test Author",
			Relationship: domain.AuthorRelationshipOther,
		},
		DiscoveredSkills: []domain.DiscoveredSkill{
			{Skill: "Kubernetes", Quote: "Expert in Kubernetes...", Context: &technicalContext},
			{Skill: "Team Building", Quote: "Great at team building...", Context: &softContext},
		},
	}

	// Execute job
	job := &river.Job[ReferenceLetterProcessingArgs]{
		Args: ReferenceLetterProcessingArgs{
			StorageKey:        storageKey,
			ReferenceLetterID: letterID,
			FileID:            fileID,
			ContentType:       "application/pdf",
		},
	}

	err := worker.Work(context.Background(), job)
	if err != nil {
		t.Fatalf("Work() returned error: %v", err)
	}

	// Verify letter is completed with extracted data
	updatedLetter := mocks.letterRepo.letters[letterID]
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("Letter status = %q, want 'completed'", updatedLetter.Status)
	}
	if len(updatedLetter.ExtractedData) == 0 {
		t.Error("Expected ExtractedData to be populated with discovered skills")
	}

	// Verify NO profile skills were created (job should not create records)
	// Discovered skills are stored in ExtractedData JSON and created later
	// by the applyReferenceLetterValidations mutation based on user selection
	if len(mocks.profileSkillRepo.skills) != 0 {
		t.Errorf("Expected 0 profile skills (job should not create records), got %d", len(mocks.profileSkillRepo.skills))
	}
}

func TestNormalizeSkillName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Go", "go"},
		{"  Go  ", "go"},
		{"KUBERNETES", "kubernetes"},
		{"Team Building", "team building"},
		{"  Mixed CASE  ", "mixed case"},
	}

	for _, tt := range tests {
		result := normalizeSkillName(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeSkillName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestMapAuthorToTestimonialRelationship(t *testing.T) {
	tests := []struct {
		input    domain.AuthorRelationship
		expected domain.TestimonialRelationship
	}{
		{domain.AuthorRelationshipManager, domain.TestimonialRelationshipManager},
		{domain.AuthorRelationshipPeer, domain.TestimonialRelationshipPeer},
		{domain.AuthorRelationshipColleague, domain.TestimonialRelationshipPeer},
		{domain.AuthorRelationshipDirectReport, domain.TestimonialRelationshipDirectReport},
		{domain.AuthorRelationshipClient, domain.TestimonialRelationshipClient},
		{domain.AuthorRelationshipOther, domain.TestimonialRelationshipOther},
	}

	for _, tt := range tests {
		result := mapAuthorToTestimonialRelationship(tt.input)
		if result != tt.expected {
			t.Errorf("mapAuthorToTestimonialRelationship(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func setupLetterTestTracing(t *testing.T) *tracetest.InMemoryExporter {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		otel.SetTracerProvider(prev)
		_ = tp.Shutdown(context.Background()) //nolint:errcheck // best effort cleanup in test
	})
	return exporter
}

func TestExtractLetterData_CreatesParentSpan(t *testing.T) {
	exporter := setupLetterTestTracing(t)

	extractor := &mockDocumentExtractor{
		extractTextResult: "Reference letter text",
		extractLetterData: &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Author Name",
				Relationship: domain.AuthorRelationshipPeer,
			},
			Testimonials:       []domain.ExtractedTestimonial{},
			SkillMentions:      []domain.ExtractedSkillMention{},
			ExperienceMentions: []domain.ExtractedExperienceMention{},
			DiscoveredSkills:   []domain.DiscoveredSkill{},
		},
	}

	worker := &ReferenceLetterProcessingWorker{
		extractor: extractor,
	}

	_, err := worker.extractLetterData(context.Background(), []byte("pdf data"), "application/pdf", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	spans := exporter.GetSpans()
	var found bool
	for _, s := range spans {
		if s.Name == "reference_letter_extraction" {
			found = true
			break
		}
	}
	if !found {
		names := make([]string, len(spans))
		for i, s := range spans {
			names[i] = s.Name
		}
		t.Errorf("expected span named 'reference_letter_extraction', got spans: %v", names)
	}
}
