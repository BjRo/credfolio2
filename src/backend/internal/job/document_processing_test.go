//nolint:goconst // Test file - string constants are fine inline
package job_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"backend/internal/domain"
	"backend/internal/job"
	"backend/internal/logger"
)

// --- Mock repositories for unified document processing tests ---

type mockResumeRepository struct {
	resumes map[uuid.UUID]*domain.Resume
}

func newMockResumeRepository() *mockResumeRepository {
	return &mockResumeRepository{resumes: make(map[uuid.UUID]*domain.Resume)}
}

func (r *mockResumeRepository) Create(_ context.Context, resume *domain.Resume) error {
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

type mockRefLetterRepository struct {
	letters map[uuid.UUID]*domain.ReferenceLetter
}

func newMockRefLetterRepository() *mockRefLetterRepository {
	return &mockRefLetterRepository{letters: make(map[uuid.UUID]*domain.ReferenceLetter)}
}

func (r *mockRefLetterRepository) Create(_ context.Context, letter *domain.ReferenceLetter) error {
	r.letters[letter.ID] = letter
	return nil
}

func (r *mockRefLetterRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	letter, ok := r.letters[id]
	if !ok {
		return nil, nil
	}
	return letter, nil
}

func (r *mockRefLetterRepository) GetByUserID(_ context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	var result []*domain.ReferenceLetter
	for _, letter := range r.letters {
		if letter.UserID == userID {
			result = append(result, letter)
		}
	}
	return result, nil
}

func (r *mockRefLetterRepository) Update(_ context.Context, letter *domain.ReferenceLetter) error {
	r.letters[letter.ID] = letter
	return nil
}

func (r *mockRefLetterRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.letters, id)
	return nil
}

type mockFileRepository struct {
	files map[uuid.UUID]*domain.File
}

func newMockFileRepository() *mockFileRepository {
	return &mockFileRepository{files: make(map[uuid.UUID]*domain.File)}
}

func (r *mockFileRepository) Create(_ context.Context, file *domain.File) error {
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

func (r *mockFileRepository) Update(_ context.Context, file *domain.File) error {
	r.files[file.ID] = file
	return nil
}

func (r *mockFileRepository) GetByUserID(_ context.Context, _ uuid.UUID) ([]*domain.File, error) {
	return nil, nil
}

func (r *mockFileRepository) GetByUserIDAndContentHash(_ context.Context, _ uuid.UUID, _ string) (*domain.File, error) {
	return nil, nil
}

func (r *mockFileRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.files, id)
	return nil
}

type mockProfileRepository struct {
	profiles map[uuid.UUID]*domain.Profile
}

func newMockProfileRepository() *mockProfileRepository {
	return &mockProfileRepository{profiles: make(map[uuid.UUID]*domain.Profile)}
}

func (r *mockProfileRepository) Create(_ context.Context, profile *domain.Profile) error {
	r.profiles[profile.ID] = profile
	return nil
}

func (r *mockProfileRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Profile, error) {
	p, ok := r.profiles[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *mockProfileRepository) GetByUserID(_ context.Context, userID uuid.UUID) (*domain.Profile, error) {
	for _, p := range r.profiles {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, nil
}

func (r *mockProfileRepository) GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	p, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return p, nil
	}
	profile := &domain.Profile{ID: uuid.New(), UserID: userID}
	r.profiles[profile.ID] = profile
	return profile, nil
}

func (r *mockProfileRepository) Update(_ context.Context, profile *domain.Profile) error {
	r.profiles[profile.ID] = profile
	return nil
}

func (r *mockProfileRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.profiles, id)
	return nil
}

type mockProfileSkillRepository struct {
	skills map[uuid.UUID]*domain.ProfileSkill
}

func newMockProfileSkillRepository() *mockProfileSkillRepository {
	return &mockProfileSkillRepository{skills: make(map[uuid.UUID]*domain.ProfileSkill)}
}

func (r *mockProfileSkillRepository) Create(_ context.Context, skill *domain.ProfileSkill) error {
	r.skills[skill.ID] = skill
	return nil
}

func (r *mockProfileSkillRepository) CreateIgnoreDuplicate(_ context.Context, skill *domain.ProfileSkill) error {
	// Simulate ON CONFLICT DO UPDATE RETURNING * — return existing row's ID on duplicate
	for _, existing := range r.skills {
		if existing.ProfileID == skill.ProfileID && existing.NormalizedName == skill.NormalizedName {
			skill.ID = existing.ID
			return nil
		}
	}
	r.skills[skill.ID] = skill
	return nil
}

func (r *mockProfileSkillRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ProfileSkill, error) {
	skill, ok := r.skills[id]
	if !ok {
		return nil, nil
	}
	return skill, nil
}

func (r *mockProfileSkillRepository) GetByIDs(_ context.Context, ids []uuid.UUID) (map[uuid.UUID]*domain.ProfileSkill, error) {
	result := make(map[uuid.UUID]*domain.ProfileSkill)
	for _, id := range ids {
		if skill, ok := r.skills[id]; ok {
			result[id] = skill
		}
	}
	return result, nil
}

func (r *mockProfileSkillRepository) GetByProfileID(_ context.Context, profileID uuid.UUID) ([]*domain.ProfileSkill, error) {
	var result []*domain.ProfileSkill
	for _, skill := range r.skills {
		if skill.ProfileID == profileID {
			result = append(result, skill)
		}
	}
	return result, nil
}

func (r *mockProfileSkillRepository) Update(_ context.Context, skill *domain.ProfileSkill) error {
	r.skills[skill.ID] = skill
	return nil
}

func (r *mockProfileSkillRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.skills, id)
	return nil
}

func (r *mockProfileSkillRepository) GetNextDisplayOrder(_ context.Context, _ uuid.UUID) (int, error) {
	return len(r.skills), nil
}

func (r *mockProfileSkillRepository) DeleteBySourceResumeID(_ context.Context, sourceResumeID uuid.UUID) error {
	for id, skill := range r.skills {
		if skill.SourceResumeID != nil && *skill.SourceResumeID == sourceResumeID {
			delete(r.skills, id)
		}
	}
	return nil
}

// mockDownloadStorage supports Download with in-memory data.
type mockDownloadStorage struct {
	data  []byte
	dlErr error
}

func newMockDownloadStorage(data []byte) *mockDownloadStorage {
	return &mockDownloadStorage{data: data}
}

func (s *mockDownloadStorage) Upload(_ context.Context, _ string, _ io.Reader, _ int64, _ string) (*domain.StorageObject, error) {
	return &domain.StorageObject{}, nil
}

func (s *mockDownloadStorage) Download(_ context.Context, _ string) (io.ReadCloser, error) {
	if s.dlErr != nil {
		return nil, s.dlErr
	}
	return io.NopCloser(bytes.NewReader(s.data)), nil
}

func (s *mockDownloadStorage) Delete(_ context.Context, _ string) error { return nil }

func (s *mockDownloadStorage) GetPresignedURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", nil
}

func (s *mockDownloadStorage) GetPublicURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", nil
}

func (s *mockDownloadStorage) Exists(_ context.Context, _ string) (bool, error) {
	return true, nil
}

// mockDocExtractor implements domain.DocumentExtractor for testing.
type mockDocExtractor struct {
	extractTextResult string
	extractTextErr    error
	resumeData        *domain.ResumeExtractedData
	resumeErr         error
	letterData        *domain.ExtractedLetterData
	letterErr         error
}

func (e *mockDocExtractor) ExtractText(_ context.Context, _ []byte, _ string) (string, error) {
	return e.extractTextResult, e.extractTextErr
}

func (e *mockDocExtractor) ExtractResumeData(_ context.Context, _ string) (*domain.ResumeExtractedData, error) {
	return e.resumeData, e.resumeErr
}

func (e *mockDocExtractor) ExtractLetterData(_ context.Context, _ string, _ []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	return e.letterData, e.letterErr
}

func (e *mockDocExtractor) DetectDocumentContent(_ context.Context, _ string) (*domain.DocumentDetectionResult, error) {
	return nil, nil
}

// testLogger returns a logger that discards all output (for tests).
func testLogger() logger.Logger {
	return logger.NewStdoutLogger(logger.WithMinLevel(logger.Severity(100))) // level 100 = discard all
}

func uuidPtr(id uuid.UUID) *uuid.UUID { return &id }

// --- Tests ---

func TestDocumentProcessingArgs_Kind(t *testing.T) {
	args := job.DocumentProcessingArgs{}
	if got := args.Kind(); got != "document_processing" {
		t.Errorf("Kind() = %q, want %q", got, "document_processing")
	}
}

func TestDocumentProcessingArgs_InsertOpts_MaxAttempts(t *testing.T) {
	args := job.DocumentProcessingArgs{}
	opts := args.InsertOpts()
	if opts.MaxAttempts != 2 {
		t.Errorf("MaxAttempts = %d, want 2", opts.MaxAttempts)
	}
}

func TestDocumentProcessingWorker_Timeout_TenMinutes(t *testing.T) {
	worker := job.NewDocumentProcessingWorker(
		newMockResumeRepository(), newMockRefLetterRepository(), newMockFileRepository(),
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage(nil), &mockDocExtractor{}, testLogger(),
	)
	timeout := worker.Timeout(nil)
	if timeout != 10*time.Minute {
		t.Errorf("Timeout = %v, want 10m", timeout)
	}
}

func TestDocumentProcessingWorker_ResumeOnly(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                          //nolint:errcheck // test setup

	extractor := &mockDocExtractor{
		extractTextResult: "John Doe, Software Engineer",
		resumeData: &domain.ResumeExtractedData{
			Name:       "John Doe",
			Skills:     []string{"Go", "Python"},
			Experience: []domain.WorkExperience{{Company: "Acme", Title: "Engineer"}},
			Education:  []domain.Education{{Institution: "MIT"}},
		},
	}

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, newMockRefLetterRepository(), fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("pdf data")), extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
			UserID:      userID,
			ResumeID:    uuidPtr(resumeID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	// Verify resume status is completed
	updated, err := resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updated.Status != domain.ResumeStatusCompleted {
		t.Errorf("resume status = %s, want %s", updated.Status, domain.ResumeStatusCompleted)
	}

	// Verify extracted data was saved
	if updated.ExtractedData == nil {
		t.Fatal("expected extracted data to be saved")
	}
	var data domain.ResumeExtractedData
	if err := json.Unmarshal(updated.ExtractedData, &data); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if data.Name != "John Doe" {
		t.Errorf("extracted name = %q, want %q", data.Name, "John Doe")
	}
}

func TestDocumentProcessingWorker_LetterOnly(t *testing.T) {
	ctx := context.Background()
	refLetterRepo := newMockRefLetterRepository()
	fileRepo := newMockFileRepository()

	letterID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = refLetterRepo.Create(ctx, &domain.ReferenceLetter{ID: letterID, UserID: userID, FileID: &fileID, Status: domain.ReferenceLetterStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                                                //nolint:errcheck // test setup

	authorName := "Jane Smith"
	authorTitle := "VP Engineering"
	authorCompany := "BigCo"
	extractor := &mockDocExtractor{
		extractTextResult: "To whom it may concern...",
		letterData: &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:    authorName,
				Title:   &authorTitle,
				Company: &authorCompany,
			},
			Testimonials: []domain.ExtractedTestimonial{{Quote: "Great engineer"}},
		},
	}

	worker := job.NewDocumentProcessingWorker(
		newMockResumeRepository(), refLetterRepo, fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("pdf data")), extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:        "test/letter.pdf",
			FileID:            fileID,
			ContentType:       "application/pdf",
			UserID:            userID,
			ReferenceLetterID: uuidPtr(letterID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	// Verify letter status is completed
	updated, err := refLetterRepo.GetByID(ctx, letterID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updated.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("letter status = %s, want %s", updated.Status, domain.ReferenceLetterStatusCompleted)
	}

	// Verify author fields populated
	if updated.AuthorName == nil || *updated.AuthorName != authorName {
		t.Errorf("author name = %v, want %q", updated.AuthorName, authorName)
	}
	if updated.AuthorTitle == nil || *updated.AuthorTitle != authorTitle {
		t.Errorf("author title = %v, want %q", updated.AuthorTitle, authorTitle)
	}
}

func TestDocumentProcessingWorker_DualExtraction(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	refLetterRepo := newMockRefLetterRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	letterID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending})           //nolint:errcheck // test setup
	_ = refLetterRepo.Create(ctx, &domain.ReferenceLetter{ID: letterID, UserID: userID, FileID: &fileID, Status: domain.ReferenceLetterStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                                                //nolint:errcheck // test setup

	extractor := &mockDocExtractor{
		extractTextResult: "Combined document text",
		resumeData:        &domain.ResumeExtractedData{Name: "Dual User", Skills: []string{"Go"}},
		letterData:        &domain.ExtractedLetterData{Author: domain.ExtractedAuthor{Name: "Referee"}, Testimonials: []domain.ExtractedTestimonial{{Quote: "Outstanding"}}},
	}

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, refLetterRepo, fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("pdf data")), extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:        "test/dual.pdf",
			FileID:            fileID,
			ContentType:       "application/pdf",
			UserID:            userID,
			ResumeID:          uuidPtr(resumeID),
			ReferenceLetterID: uuidPtr(letterID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	// Both should be completed
	updatedResume, err := resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updatedResume.Status != domain.ResumeStatusCompleted {
		t.Errorf("resume status = %s, want %s", updatedResume.Status, domain.ResumeStatusCompleted)
	}
	updatedLetter, err := refLetterRepo.GetByID(ctx, letterID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("letter status = %s, want %s", updatedLetter.Status, domain.ReferenceLetterStatusCompleted)
	}
}

func TestDocumentProcessingWorker_TextExtractionFails_MarksBothFailed(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	refLetterRepo := newMockRefLetterRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	letterID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending})           //nolint:errcheck // test setup
	_ = refLetterRepo.Create(ctx, &domain.ReferenceLetter{ID: letterID, UserID: userID, FileID: &fileID, Status: domain.ReferenceLetterStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                                                //nolint:errcheck // test setup

	extractor := &mockDocExtractor{extractTextErr: errors.New("OCR failure")}

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, refLetterRepo, fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("corrupt data")), extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:        "test/corrupt.pdf",
			FileID:            fileID,
			ContentType:       "application/pdf",
			UserID:            userID,
			ResumeID:          uuidPtr(resumeID),
			ReferenceLetterID: uuidPtr(letterID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err == nil {
		t.Fatal("expected error for text extraction failure")
	}

	// Both should be marked failed
	updatedResume, getErr := resumeRepo.GetByID(ctx, resumeID)
	if getErr != nil {
		t.Fatalf("GetByID failed: %v", getErr)
	}
	if updatedResume.Status != domain.ResumeStatusFailed {
		t.Errorf("resume status = %s, want %s", updatedResume.Status, domain.ResumeStatusFailed)
	}
	updatedLetter, getErr := refLetterRepo.GetByID(ctx, letterID)
	if getErr != nil {
		t.Fatalf("GetByID failed: %v", getErr)
	}
	if updatedLetter.Status != domain.ReferenceLetterStatusFailed {
		t.Errorf("letter status = %s, want %s", updatedLetter.Status, domain.ReferenceLetterStatusFailed)
	}
}

func TestDocumentProcessingWorker_ResumeExtractFails_LetterSucceeds(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	refLetterRepo := newMockRefLetterRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	letterID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending})           //nolint:errcheck // test setup
	_ = refLetterRepo.Create(ctx, &domain.ReferenceLetter{ID: letterID, UserID: userID, FileID: &fileID, Status: domain.ReferenceLetterStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                                                //nolint:errcheck // test setup

	extractor := &mockDocExtractor{
		extractTextResult: "Some text",
		resumeErr:         errors.New("LLM timeout"),
		letterData:        &domain.ExtractedLetterData{Author: domain.ExtractedAuthor{Name: "Referee"}},
	}

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, refLetterRepo, fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("pdf data")), extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:        "test/mixed.pdf",
			FileID:            fileID,
			ContentType:       "application/pdf",
			UserID:            userID,
			ResumeID:          uuidPtr(resumeID),
			ReferenceLetterID: uuidPtr(letterID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err == nil {
		t.Fatal("expected error when resume extraction fails")
	}

	// Resume should be failed
	updatedResume, getErr := resumeRepo.GetByID(ctx, resumeID)
	if getErr != nil {
		t.Fatalf("GetByID failed: %v", getErr)
	}
	if updatedResume.Status != domain.ResumeStatusFailed {
		t.Errorf("resume status = %s, want %s", updatedResume.Status, domain.ResumeStatusFailed)
	}

	// Letter should still succeed
	updatedLetter, getErr := refLetterRepo.GetByID(ctx, letterID)
	if getErr != nil {
		t.Fatalf("GetByID failed: %v", getErr)
	}
	if updatedLetter.Status != domain.ReferenceLetterStatusCompleted {
		t.Errorf("letter status = %s, want %s", updatedLetter.Status, domain.ReferenceLetterStatusCompleted)
	}
}

func TestDocumentProcessingWorker_DownloadFails(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                          //nolint:errcheck // test setup

	storage := newMockDownloadStorage(nil)
	storage.dlErr = fmt.Errorf("network error")

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, newMockRefLetterRepository(), fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		storage, &mockDocExtractor{}, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
			UserID:      userID,
			ResumeID:    uuidPtr(resumeID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err == nil {
		t.Fatal("expected error for download failure")
	}

	// Resume should be marked failed
	updated, getErr := resumeRepo.GetByID(ctx, resumeID)
	if getErr != nil {
		t.Fatalf("GetByID failed: %v", getErr)
	}
	if updated.Status != domain.ResumeStatusFailed {
		t.Errorf("resume status = %s, want %s", updated.Status, domain.ResumeStatusFailed)
	}
}

func TestDocumentProcessingWorker_FallbackToFileContentType(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending}) //nolint:errcheck // test setup
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"})                                        //nolint:errcheck // test setup

	extractor := &mockDocExtractor{
		extractTextResult: "Resume text",
		resumeData:        &domain.ResumeExtractedData{Name: "Test"},
	}

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, newMockRefLetterRepository(), fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("data")), extractor, testLogger(),
	)

	// ContentType is empty — worker should look it up from file record
	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey: "test/key.pdf",
			FileID:     fileID,
			UserID:     userID,
			ResumeID:   uuidPtr(resumeID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	updated, err := resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updated.Status != domain.ResumeStatusCompleted {
		t.Errorf("resume status = %s, want %s", updated.Status, domain.ResumeStatusCompleted)
	}
}

func TestDocumentProcessingWorker_UsesStoredExtractedText(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending}) //nolint:errcheck // test setup

	// File record has extracted_text already set (from detection worker)
	storedText := "John Doe, Software Engineer at Acme Corp"
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf", ExtractedText: &storedText}) //nolint:errcheck // test setup

	// Extractor's ExtractText should NOT be called — use a sentinel error to detect it
	extractor := &mockDocExtractor{
		extractTextErr: errors.New("ExtractText should not be called when stored text exists"),
		resumeData:     &domain.ResumeExtractedData{Name: "John Doe", Skills: []string{"Go"}},
	}

	// Storage download should NOT be called — use a sentinel error to detect it
	storage := newMockDownloadStorage(nil)
	storage.dlErr = errors.New("Download should not be called when stored text exists")

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, newMockRefLetterRepository(), fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		storage, extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
			UserID:      userID,
			ResumeID:    uuidPtr(resumeID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	// Verify resume was processed successfully using stored text
	updated, err := resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updated.Status != domain.ResumeStatusCompleted {
		t.Errorf("resume status = %s, want %s", updated.Status, domain.ResumeStatusCompleted)
	}
}

func TestDocumentProcessingWorker_FallsBackToExtractionWhenNoStoredText(t *testing.T) {
	ctx := context.Background()
	resumeRepo := newMockResumeRepository()
	fileRepo := newMockFileRepository()

	resumeID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	_ = resumeRepo.Create(ctx, &domain.Resume{ID: resumeID, UserID: userID, FileID: fileID, Status: domain.ResumeStatusPending}) //nolint:errcheck // test setup

	// File record has NO extracted_text (legacy flow)
	_ = fileRepo.Create(ctx, &domain.File{ID: fileID, ContentType: "application/pdf"}) //nolint:errcheck // test setup

	extractor := &mockDocExtractor{
		extractTextResult: "Extracted from LLM",
		resumeData:        &domain.ResumeExtractedData{Name: "Fallback User"},
	}

	worker := job.NewDocumentProcessingWorker(
		resumeRepo, newMockRefLetterRepository(), fileRepo,
		newMockProfileRepository(), newMockProfileSkillRepository(),
		newMockDownloadStorage([]byte("pdf data")), extractor, testLogger(),
	)

	riverJob := &river.Job[job.DocumentProcessingArgs]{
		Args: job.DocumentProcessingArgs{
			StorageKey:  "test/key.pdf",
			FileID:      fileID,
			ContentType: "application/pdf",
			UserID:      userID,
			ResumeID:    uuidPtr(resumeID),
		},
	}

	err := worker.Work(ctx, riverJob)
	if err != nil {
		t.Fatalf("Work() error = %v", err)
	}

	updated, err := resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if updated.Status != domain.ResumeStatusCompleted {
		t.Errorf("resume status = %s, want %s", updated.Status, domain.ResumeStatusCompleted)
	}
}
