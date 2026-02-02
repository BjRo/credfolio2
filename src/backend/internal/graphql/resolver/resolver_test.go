package resolver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
	"backend/internal/graphql/model"
	"backend/internal/graphql/resolver"
	"backend/internal/infrastructure/storage"
	"backend/internal/logger"
)

const testUserName = "Test User"

// Error message constants for test assertions.
const (
	errMsgInvalidUserIDFormat = "invalid user ID format"
	errMsgUserNotFound        = "user not found"
)

// stringPtr returns a pointer to a string (test helper).
func stringPtr(s string) *string {
	return &s
}

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

// mockProfileRepository is a mock implementation of domain.ProfileRepository.
type mockProfileRepository struct {
	profiles map[uuid.UUID]*domain.Profile
}

func newMockProfileRepository() *mockProfileRepository {
	return &mockProfileRepository{profiles: make(map[uuid.UUID]*domain.Profile)}
}

func (r *mockProfileRepository) Create(_ context.Context, profile *domain.Profile) error {
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	r.profiles[profile.ID] = profile
	return nil
}

func (r *mockProfileRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.Profile, error) {
	profile, ok := r.profiles[id]
	if !ok {
		return nil, nil
	}
	return profile, nil
}

func (r *mockProfileRepository) GetByUserID(_ context.Context, userID uuid.UUID) (*domain.Profile, error) {
	for _, profile := range r.profiles {
		if profile.UserID == userID {
			return profile, nil
		}
	}
	return nil, nil
}

func (r *mockProfileRepository) GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	profile, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile != nil {
		return profile, nil
	}
	profile = &domain.Profile{
		ID:     uuid.New(),
		UserID: userID,
	}
	if err := r.Create(ctx, profile); err != nil {
		return nil, err
	}
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

// mockProfileExperienceRepository is a mock implementation of domain.ProfileExperienceRepository.
type mockProfileExperienceRepository struct {
	experiences map[uuid.UUID]*domain.ProfileExperience
}

func newMockProfileExperienceRepository() *mockProfileExperienceRepository {
	return &mockProfileExperienceRepository{experiences: make(map[uuid.UUID]*domain.ProfileExperience)}
}

func (r *mockProfileExperienceRepository) Create(_ context.Context, experience *domain.ProfileExperience) error {
	if experience.ID == uuid.Nil {
		experience.ID = uuid.New()
	}
	r.experiences[experience.ID] = experience
	return nil
}

func (r *mockProfileExperienceRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ProfileExperience, error) {
	experience, ok := r.experiences[id]
	if !ok {
		return nil, nil
	}
	return experience, nil
}

func (r *mockProfileExperienceRepository) GetByProfileID(_ context.Context, profileID uuid.UUID) ([]*domain.ProfileExperience, error) {
	var result []*domain.ProfileExperience
	for _, exp := range r.experiences {
		if exp.ProfileID == profileID {
			result = append(result, exp)
		}
	}
	return result, nil
}

func (r *mockProfileExperienceRepository) Update(_ context.Context, experience *domain.ProfileExperience) error {
	r.experiences[experience.ID] = experience
	return nil
}

func (r *mockProfileExperienceRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.experiences, id)
	return nil
}

func (r *mockProfileExperienceRepository) GetNextDisplayOrder(_ context.Context, _ uuid.UUID) (int, error) {
	return len(r.experiences), nil
}

func (r *mockProfileExperienceRepository) DeleteBySourceResumeID(_ context.Context, sourceResumeID uuid.UUID) error {
	for id, exp := range r.experiences {
		if exp.SourceResumeID != nil && *exp.SourceResumeID == sourceResumeID {
			delete(r.experiences, id)
		}
	}
	return nil
}

// mockProfileEducationRepository is a mock implementation of domain.ProfileEducationRepository.
type mockProfileEducationRepository struct {
	educations map[uuid.UUID]*domain.ProfileEducation
}

func newMockProfileEducationRepository() *mockProfileEducationRepository {
	return &mockProfileEducationRepository{educations: make(map[uuid.UUID]*domain.ProfileEducation)}
}

func (r *mockProfileEducationRepository) Create(_ context.Context, education *domain.ProfileEducation) error {
	if education.ID == uuid.Nil {
		education.ID = uuid.New()
	}
	r.educations[education.ID] = education
	return nil
}

func (r *mockProfileEducationRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ProfileEducation, error) {
	education, ok := r.educations[id]
	if !ok {
		return nil, nil
	}
	return education, nil
}

func (r *mockProfileEducationRepository) GetByProfileID(_ context.Context, profileID uuid.UUID) ([]*domain.ProfileEducation, error) {
	var result []*domain.ProfileEducation
	for _, edu := range r.educations {
		if edu.ProfileID == profileID {
			result = append(result, edu)
		}
	}
	return result, nil
}

func (r *mockProfileEducationRepository) Update(_ context.Context, education *domain.ProfileEducation) error {
	r.educations[education.ID] = education
	return nil
}

func (r *mockProfileEducationRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.educations, id)
	return nil
}

func (r *mockProfileEducationRepository) GetNextDisplayOrder(_ context.Context, _ uuid.UUID) (int, error) {
	return len(r.educations), nil
}

func (r *mockProfileEducationRepository) DeleteBySourceResumeID(_ context.Context, sourceResumeID uuid.UUID) error {
	for id, edu := range r.educations {
		if edu.SourceResumeID != nil && *edu.SourceResumeID == sourceResumeID {
			delete(r.educations, id)
		}
	}
	return nil
}

// mockProfileSkillRepository is a mock implementation of domain.ProfileSkillRepository.
type mockProfileSkillRepository struct {
	skills map[uuid.UUID]*domain.ProfileSkill
}

func newMockProfileSkillRepository() *mockProfileSkillRepository {
	return &mockProfileSkillRepository{skills: make(map[uuid.UUID]*domain.ProfileSkill)}
}

func (r *mockProfileSkillRepository) Create(_ context.Context, skill *domain.ProfileSkill) error {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}
	r.skills[skill.ID] = skill
	return nil
}

func (r *mockProfileSkillRepository) CreateIgnoreDuplicate(_ context.Context, skill *domain.ProfileSkill) error {
	// Same as Create for test purposes - just silently succeeds
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
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

func (r *mockProfileSkillRepository) GetByProfileID(_ context.Context, profileID uuid.UUID) ([]*domain.ProfileSkill, error) {
	var result []*domain.ProfileSkill
	for _, s := range r.skills {
		if s.ProfileID == profileID {
			result = append(result, s)
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
	for id, s := range r.skills {
		if s.SourceResumeID != nil && *s.SourceResumeID == sourceResumeID {
			delete(r.skills, id)
		}
	}
	return nil
}

// mockAuthorRepository is a mock implementation of domain.AuthorRepository.
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
	author, ok := r.authors[id]
	if !ok {
		return nil, nil
	}
	return author, nil
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

// mockTestimonialRepository is a mock implementation of domain.TestimonialRepository.
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
	testimonial, ok := r.testimonials[id]
	if !ok {
		return nil, nil
	}
	return testimonial, nil
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

func (r *mockTestimonialRepository) GetByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) ([]*domain.Testimonial, error) {
	var result []*domain.Testimonial
	for _, t := range r.testimonials {
		if t.ReferenceLetterID == refLetterID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *mockTestimonialRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.testimonials, id)
	return nil
}

func (r *mockTestimonialRepository) DeleteByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) error {
	for id, t := range r.testimonials {
		if t.ReferenceLetterID == refLetterID {
			delete(r.testimonials, id)
		}
	}
	return nil
}

// mockSkillValidationRepository is a mock implementation of domain.SkillValidationRepository.
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
	validation, ok := r.validations[id]
	if !ok {
		return nil, nil
	}
	return validation, nil
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

func (r *mockSkillValidationRepository) GetByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) ([]*domain.SkillValidation, error) {
	var result []*domain.SkillValidation
	for _, v := range r.validations {
		if v.ReferenceLetterID == refLetterID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockSkillValidationRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.validations, id)
	return nil
}

func (r *mockSkillValidationRepository) DeleteByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) error {
	for id, v := range r.validations {
		if v.ReferenceLetterID == refLetterID {
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

// mockExperienceValidationRepository is a mock implementation of domain.ExperienceValidationRepository.
type mockExperienceValidationRepository struct {
	validations map[uuid.UUID]*domain.ExperienceValidation
}

func newMockExperienceValidationRepository() *mockExperienceValidationRepository {
	return &mockExperienceValidationRepository{validations: make(map[uuid.UUID]*domain.ExperienceValidation)}
}

func (r *mockExperienceValidationRepository) Create(_ context.Context, validation *domain.ExperienceValidation) error {
	if validation.ID == uuid.Nil {
		validation.ID = uuid.New()
	}
	r.validations[validation.ID] = validation
	return nil
}

func (r *mockExperienceValidationRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ExperienceValidation, error) {
	validation, ok := r.validations[id]
	if !ok {
		return nil, nil
	}
	return validation, nil
}

func (r *mockExperienceValidationRepository) GetByProfileExperienceID(_ context.Context, profileExpID uuid.UUID) ([]*domain.ExperienceValidation, error) {
	var result []*domain.ExperienceValidation
	for _, v := range r.validations {
		if v.ProfileExperienceID == profileExpID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockExperienceValidationRepository) GetByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) ([]*domain.ExperienceValidation, error) {
	var result []*domain.ExperienceValidation
	for _, v := range r.validations {
		if v.ReferenceLetterID == refLetterID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockExperienceValidationRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.validations, id)
	return nil
}

func (r *mockExperienceValidationRepository) DeleteByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) error {
	for id, v := range r.validations {
		if v.ReferenceLetterID == refLetterID {
			delete(r.validations, id)
		}
	}
	return nil
}

func (r *mockExperienceValidationRepository) CountByProfileExperienceID(_ context.Context, profileExpID uuid.UUID) (int, error) {
	count := 0
	for _, v := range r.validations {
		if v.ProfileExperienceID == profileExpID {
			count++
		}
	}
	return count, nil
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
	name := testUserName
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), newMockProfileRepository(), newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
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
		errorR := resolver.NewResolver(&errorUserRepository{}, fileRepo, refLetterRepo, newMockResumeRepository(), newMockProfileRepository(), newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
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

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), newMockProfileRepository(), newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
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

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), newMockProfileRepository(), newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
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

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), newMockProfileRepository(), newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
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

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), newMockProfileRepository(), newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
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

func TestCreateEducation(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	eduRepo := newMockProfileEducationRepository()

	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "edu-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), eduRepo, newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	t.Run("creates education entry", func(t *testing.T) {
		result, err := mutation.CreateEducation(ctx, user.ID.String(), model.CreateEducationInput{
			Institution: "MIT",
			Degree:      "Bachelor of Science",
			Field:       stringPtr("Computer Science"),
			IsCurrent:   false,
		})
		if err != nil {
			t.Fatalf("CreateEducation failed: %v", err)
		}

		eduResult, ok := result.(*model.EducationResult)
		if !ok {
			t.Fatalf("expected EducationResult, got %T", result)
		}

		if eduResult.Education.Institution != "MIT" {
			t.Errorf("Institution mismatch: got %s, want MIT", eduResult.Education.Institution)
		}
		if eduResult.Education.Degree != "Bachelor of Science" {
			t.Errorf("Degree mismatch: got %s, want Bachelor of Science", eduResult.Education.Degree)
		}
		if eduResult.Education.Field == nil || *eduResult.Education.Field != "Computer Science" {
			t.Errorf("Field mismatch: got %v, want Computer Science", eduResult.Education.Field)
		}
	})

	t.Run("returns validation error for invalid user ID", func(t *testing.T) {
		result, err := mutation.CreateEducation(ctx, "invalid-uuid", model.CreateEducationInput{
			Institution: "MIT",
			Degree:      "BS",
			IsCurrent:   false,
		})
		if err != nil {
			t.Fatalf("CreateEducation failed: %v", err)
		}

		validationErr, ok := result.(*model.EducationValidationError)
		if !ok {
			t.Fatalf("expected EducationValidationError, got %T", result)
		}
		if validationErr.Message != errMsgInvalidUserIDFormat {
			t.Errorf("unexpected error message: %s", validationErr.Message)
		}
	})

	t.Run("returns validation error for non-existent user", func(t *testing.T) {
		result, err := mutation.CreateEducation(ctx, uuid.New().String(), model.CreateEducationInput{
			Institution: "MIT",
			Degree:      "BS",
			IsCurrent:   false,
		})
		if err != nil {
			t.Fatalf("CreateEducation failed: %v", err)
		}

		validationErr, ok := result.(*model.EducationValidationError)
		if !ok {
			t.Fatalf("expected EducationValidationError, got %T", result)
		}
		if validationErr.Message != errMsgUserNotFound {
			t.Errorf("unexpected error message: %s", validationErr.Message)
		}
	})
}

func TestUpdateEducation(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	eduRepo := newMockProfileEducationRepository()

	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "edu-update@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), eduRepo, newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	// Create an education entry first
	createResult, err := mutation.CreateEducation(ctx, user.ID.String(), model.CreateEducationInput{
		Institution: "MIT",
		Degree:      "Bachelor of Science",
		IsCurrent:   false,
	})
	if err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}
	eduResult, ok := createResult.(*model.EducationResult)
	if !ok {
		t.Fatalf("expected EducationResult, got %T", createResult)
	}
	eduID := eduResult.Education.ID

	t.Run("updates education entry", func(t *testing.T) {
		newInstitution := "Stanford University"
		result, err := mutation.UpdateEducation(ctx, eduID, model.UpdateEducationInput{
			Institution: &newInstitution,
		})
		if err != nil {
			t.Fatalf("UpdateEducation failed: %v", err)
		}

		updated, ok := result.(*model.EducationResult)
		if !ok {
			t.Fatalf("expected EducationResult, got %T", result)
		}
		if updated.Education.Institution != "Stanford University" {
			t.Errorf("Institution mismatch: got %s, want Stanford University", updated.Education.Institution)
		}
		// Degree should remain unchanged
		if updated.Education.Degree != "Bachelor of Science" {
			t.Errorf("Degree should not change: got %s, want Bachelor of Science", updated.Education.Degree)
		}
	})

	t.Run("returns validation error for non-existent education", func(t *testing.T) {
		newInstitution := "Harvard"
		result, err := mutation.UpdateEducation(ctx, uuid.New().String(), model.UpdateEducationInput{
			Institution: &newInstitution,
		})
		if err != nil {
			t.Fatalf("UpdateEducation failed: %v", err)
		}

		validationErr, ok := result.(*model.EducationValidationError)
		if !ok {
			t.Fatalf("expected EducationValidationError, got %T", result)
		}
		if validationErr.Message != "education not found" {
			t.Errorf("unexpected error message: %s", validationErr.Message)
		}
	})

	t.Run("returns validation error for invalid ID", func(t *testing.T) {
		newInstitution := "Harvard"
		result, err := mutation.UpdateEducation(ctx, "invalid-uuid", model.UpdateEducationInput{
			Institution: &newInstitution,
		})
		if err != nil {
			t.Fatalf("UpdateEducation failed: %v", err)
		}

		validationErr, ok := result.(*model.EducationValidationError)
		if !ok {
			t.Fatalf("expected EducationValidationError, got %T", result)
		}
		if validationErr.Message != "invalid education ID format" {
			t.Errorf("unexpected error message: %s", validationErr.Message)
		}
	})
}

func TestDeleteEducation(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	eduRepo := newMockProfileEducationRepository()

	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "edu-delete@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), eduRepo, newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	// Create an education entry first
	createResult, err := mutation.CreateEducation(ctx, user.ID.String(), model.CreateEducationInput{
		Institution: "MIT",
		Degree:      "Bachelor of Science",
		IsCurrent:   false,
	})
	if err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}
	eduResult, ok := createResult.(*model.EducationResult)
	if !ok {
		t.Fatalf("expected EducationResult, got %T", createResult)
	}
	eduID := eduResult.Education.ID

	t.Run("deletes education entry", func(t *testing.T) {
		result, err := mutation.DeleteEducation(ctx, eduID)
		if err != nil {
			t.Fatalf("DeleteEducation failed: %v", err)
		}

		if !result.Success {
			t.Error("expected success to be true")
		}
		if result.DeletedID != eduID {
			t.Errorf("DeletedID mismatch: got %s, want %s", result.DeletedID, eduID)
		}
	})

	t.Run("returns false for non-existent education", func(t *testing.T) {
		result, err := mutation.DeleteEducation(ctx, uuid.New().String())
		if err != nil {
			t.Fatalf("DeleteEducation failed: %v", err)
		}

		if result.Success {
			t.Error("expected success to be false for non-existent education")
		}
	})

	t.Run("returns false for invalid ID", func(t *testing.T) {
		result, err := mutation.DeleteEducation(ctx, "invalid-uuid")
		if err != nil {
			t.Fatalf("DeleteEducation failed: %v", err)
		}

		if result.Success {
			t.Error("expected success to be false for invalid ID")
		}
	})
}

func TestProfileEducationQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	eduRepo := newMockProfileEducationRepository()

	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "edu-query@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), eduRepo, newMockProfileSkillRepository(), newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()
	query := r.Query()

	// Create an education entry
	createResult, err := mutation.CreateEducation(ctx, user.ID.String(), model.CreateEducationInput{
		Institution: "MIT",
		Degree:      "PhD",
		Field:       stringPtr("Physics"),
		IsCurrent:   false,
		Gpa:         stringPtr("4.0"),
	})
	if err != nil {
		t.Fatalf("CreateEducation failed: %v", err)
	}
	eduResult, ok := createResult.(*model.EducationResult)
	if !ok {
		t.Fatalf("expected EducationResult, got %T", createResult)
	}
	eduID := eduResult.Education.ID

	t.Run("returns education when found", func(t *testing.T) {
		result, err := query.ProfileEducation(ctx, eduID)
		if err != nil {
			t.Fatalf("ProfileEducation query failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected education, got nil")
		}

		if result.Institution != "MIT" {
			t.Errorf("Institution mismatch: got %s, want MIT", result.Institution)
		}
		if result.Degree != "PhD" {
			t.Errorf("Degree mismatch: got %s, want PhD", result.Degree)
		}
		if result.Field == nil || *result.Field != "Physics" {
			t.Errorf("Field mismatch: got %v, want Physics", result.Field)
		}
		if result.Gpa == nil || *result.Gpa != "4.0" {
			t.Errorf("GPA mismatch: got %v, want 4.0", result.Gpa)
		}
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		result, err := query.ProfileEducation(ctx, uuid.New().String())
		if err != nil {
			t.Fatalf("ProfileEducation query failed: %v", err)
		}

		if result != nil {
			t.Error("expected nil for non-existent education")
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.ProfileEducation(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestCreateSkill(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	skillRepo := newMockProfileSkillRepository()

	ctx := context.Background()

	// Create a test user
	name := testUserName
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	t.Run("creates skill successfully", func(t *testing.T) {
		input := model.CreateSkillInput{
			Name:     "Go",
			Category: "technical",
		}

		result, err := mutation.CreateSkill(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("CreateSkill failed: %v", err)
		}

		skillResult, ok := result.(*model.SkillResult)
		if !ok {
			t.Fatalf("expected SkillResult, got %T", result)
		}

		if skillResult.Skill.Name != "Go" {
			t.Errorf("Name mismatch: got %s, want Go", skillResult.Skill.Name)
		}
		if skillResult.Skill.NormalizedName != "go" {
			t.Errorf("NormalizedName mismatch: got %s, want go", skillResult.Skill.NormalizedName)
		}
	})

	t.Run("returns validation error for empty name", func(t *testing.T) {
		input := model.CreateSkillInput{
			Name:     "  ",
			Category: "technical",
		}

		result, err := mutation.CreateSkill(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("CreateSkill failed: %v", err)
		}

		_, ok := result.(*model.SkillValidationError)
		if !ok {
			t.Fatalf("expected SkillValidationError, got %T", result)
		}
	})

	t.Run("returns validation error for invalid user ID", func(t *testing.T) {
		input := model.CreateSkillInput{
			Name:     "Go",
			Category: "technical",
		}

		result, err := mutation.CreateSkill(ctx, "invalid-uuid", input)
		if err != nil {
			t.Fatalf("CreateSkill failed: %v", err)
		}

		valErr, ok := result.(*model.SkillValidationError)
		if !ok {
			t.Fatalf("expected SkillValidationError, got %T", result)
		}
		if valErr.Message != errMsgInvalidUserIDFormat {
			t.Errorf("unexpected message: %s", valErr.Message)
		}
	})

	t.Run("returns validation error for non-existent user", func(t *testing.T) {
		input := model.CreateSkillInput{
			Name:     "Go",
			Category: "technical",
		}

		result, err := mutation.CreateSkill(ctx, uuid.New().String(), input)
		if err != nil {
			t.Fatalf("CreateSkill failed: %v", err)
		}

		valErr, ok := result.(*model.SkillValidationError)
		if !ok {
			t.Fatalf("expected SkillValidationError, got %T", result)
		}
		if valErr.Message != errMsgUserNotFound {
			t.Errorf("unexpected message: %s", valErr.Message)
		}
	})
}

func TestUpdateSkill(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	skillRepo := newMockProfileSkillRepository()

	ctx := context.Background()

	// Create a test user
	name := testUserName
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	// Create a skill first
	input := model.CreateSkillInput{
		Name:     "JavaScript",
		Category: "technical",
	}
	createResult, err := mutation.CreateSkill(ctx, user.ID.String(), input)
	if err != nil {
		t.Fatalf("setup: CreateSkill failed: %v", err)
	}
	skillResult, ok := createResult.(*model.SkillResult)
	if !ok {
		t.Fatalf("setup: expected SkillResult, got %T", createResult)
	}
	skillID := skillResult.Skill.ID

	t.Run("updates skill name", func(t *testing.T) {
		newName := "TypeScript"
		updateInput := model.UpdateSkillInput{
			Name: &newName,
		}

		result, err := mutation.UpdateSkill(ctx, skillID, updateInput)
		if err != nil {
			t.Fatalf("UpdateSkill failed: %v", err)
		}

		updated, ok := result.(*model.SkillResult)
		if !ok {
			t.Fatalf("expected SkillResult, got %T", result)
		}

		if updated.Skill.Name != "TypeScript" {
			t.Errorf("Name mismatch: got %s, want TypeScript", updated.Skill.Name)
		}
		if updated.Skill.NormalizedName != "typescript" {
			t.Errorf("NormalizedName mismatch: got %s, want typescript", updated.Skill.NormalizedName)
		}
	})

	t.Run("returns validation error for non-existent skill", func(t *testing.T) {
		newName := "Rust"
		updateInput := model.UpdateSkillInput{
			Name: &newName,
		}

		result, err := mutation.UpdateSkill(ctx, uuid.New().String(), updateInput)
		if err != nil {
			t.Fatalf("UpdateSkill failed: %v", err)
		}

		_, ok := result.(*model.SkillValidationError)
		if !ok {
			t.Fatalf("expected SkillValidationError, got %T", result)
		}
	})
}

func TestDeleteSkill(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	skillRepo := newMockProfileSkillRepository()

	ctx := context.Background()

	// Create a test user
	name := testUserName
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	// Create a skill first
	input := model.CreateSkillInput{
		Name:     "Python",
		Category: "technical",
	}
	createResult, err := mutation.CreateSkill(ctx, user.ID.String(), input)
	if err != nil {
		t.Fatalf("setup: CreateSkill failed: %v", err)
	}
	skillResult, ok := createResult.(*model.SkillResult)
	if !ok {
		t.Fatalf("setup: expected SkillResult, got %T", createResult)
	}
	skillID := skillResult.Skill.ID

	t.Run("deletes skill successfully", func(t *testing.T) {
		result, err := mutation.DeleteSkill(ctx, skillID)
		if err != nil {
			t.Fatalf("DeleteSkill failed: %v", err)
		}

		if !result.Success {
			t.Error("expected success to be true")
		}
		if result.DeletedID != skillID {
			t.Errorf("DeletedID mismatch: got %s, want %s", result.DeletedID, skillID)
		}
	})

	t.Run("returns success false for non-existent skill", func(t *testing.T) {
		result, err := mutation.DeleteSkill(ctx, uuid.New().String())
		if err != nil {
			t.Fatalf("DeleteSkill failed: %v", err)
		}

		if result.Success {
			t.Error("expected success to be false for non-existent skill")
		}
	})

	t.Run("returns success false for invalid ID", func(t *testing.T) {
		result, err := mutation.DeleteSkill(ctx, "invalid-uuid")
		if err != nil {
			t.Fatalf("DeleteSkill failed: %v", err)
		}

		if result.Success {
			t.Error("expected success to be false for invalid UUID")
		}
	})
}

func TestProfileSkillQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	skillRepo := newMockProfileSkillRepository()

	ctx := context.Background()

	// Create a test user
	name := testUserName
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	r := resolver.NewResolver(userRepo, newMockFileRepository(), newMockReferenceLetterRepository(), newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()
	query := r.Query()

	// Create a skill
	input := model.CreateSkillInput{
		Name:     "React",
		Category: "technical",
	}
	createResult, err := mutation.CreateSkill(ctx, user.ID.String(), input)
	if err != nil {
		t.Fatalf("setup: CreateSkill failed: %v", err)
	}
	skillResult, ok := createResult.(*model.SkillResult)
	if !ok {
		t.Fatalf("setup: expected SkillResult, got %T", createResult)
	}
	skillID := skillResult.Skill.ID

	t.Run("returns skill when found", func(t *testing.T) {
		result, err := query.ProfileSkill(ctx, skillID)
		if err != nil {
			t.Fatalf("ProfileSkill query failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected skill, got nil")
		}

		if result.Name != "React" {
			t.Errorf("Name mismatch: got %s, want React", result.Name)
		}
		if result.NormalizedName != "react" {
			t.Errorf("NormalizedName mismatch: got %s, want react", result.NormalizedName)
		}
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		result, err := query.ProfileSkill(ctx, uuid.New().String())
		if err != nil {
			t.Fatalf("ProfileSkill query failed: %v", err)
		}

		if result != nil {
			t.Error("expected nil for non-existent skill")
		}
	})

	t.Run("returns error for invalid ID", func(t *testing.T) {
		_, err := query.ProfileSkill(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestTestimonialsQuery(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	testimonialRepo := newMockTestimonialRepository()
	refLetterRepo := newMockReferenceLetterRepository()

	ctx := context.Background()

	// Create a test user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "testimonials-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	// Create a profile for the user
	profile := &domain.Profile{
		ID:     uuid.New(),
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("setup: failed to create profile: %v", err)
	}

	// Create a reference letter
	refLetter := &domain.ReferenceLetter{
		ID:     uuid.New(),
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	mustCreateReferenceLetter(refLetterRepo, refLetter)

	// Create testimonials
	testimonial1 := &domain.Testimonial{
		ID:                uuid.New(),
		ProfileID:         profile.ID,
		ReferenceLetterID: refLetter.ID,
		Quote:             "Great team player with excellent leadership skills.",
		AuthorName:        stringPtr("John Manager"),
		AuthorTitle:       stringPtr("Engineering Manager"),
		AuthorCompany:     stringPtr("Acme Corp"),
		Relationship:      domain.TestimonialRelationshipManager,
	}
	if err := testimonialRepo.Create(ctx, testimonial1); err != nil {
		t.Fatalf("setup: failed to create testimonial: %v", err)
	}

	testimonial2 := &domain.Testimonial{
		ID:                uuid.New(),
		ProfileID:         profile.ID,
		ReferenceLetterID: refLetter.ID,
		Quote:             "A brilliant collaborator who consistently delivers high-quality work.",
		AuthorName:        stringPtr("Sarah Peer"),
		AuthorTitle:       stringPtr("Senior Engineer"),
		AuthorCompany:     stringPtr("Acme Corp"),
		Relationship:      domain.TestimonialRelationshipPeer,
	}
	if err := testimonialRepo.Create(ctx, testimonial2); err != nil {
		t.Fatalf("setup: failed to create testimonial: %v", err)
	}

	r := resolver.NewResolver(userRepo, newMockFileRepository(), refLetterRepo, newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), newMockProfileEducationRepository(), newMockProfileSkillRepository(), newMockAuthorRepository(), testimonialRepo, newMockSkillValidationRepository(), newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns testimonials for profile", func(t *testing.T) {
		results, err := query.Testimonials(ctx, profile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("expected 2 testimonials, got %d", len(results))
		}

		// Verify the testimonials have expected data
		quotes := make(map[string]bool)
		for _, testimonial := range results {
			quotes[testimonial.Quote] = true
			if testimonial.AuthorName == "" {
				t.Error("expected author name to be set")
			}
		}

		if !quotes["Great team player with excellent leadership skills."] {
			t.Error("expected first testimonial quote")
		}
		if !quotes["A brilliant collaborator who consistently delivers high-quality work."] {
			t.Error("expected second testimonial quote")
		}
	})

	t.Run("returns empty slice for profile with no testimonials", func(t *testing.T) {
		// Create a new profile without testimonials
		emptyProfile := &domain.Profile{
			ID:     uuid.New(),
			UserID: user.ID,
		}
		if err := profileRepo.Create(ctx, emptyProfile); err != nil {
			t.Fatalf("setup: failed to create empty profile: %v", err)
		}

		results, err := query.Testimonials(ctx, emptyProfile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		if results == nil {
			t.Error("expected empty slice, got nil")
		}

		if len(results) != 0 {
			t.Errorf("expected 0 testimonials, got %d", len(results))
		}
	})

	t.Run("returns error for invalid profile ID", func(t *testing.T) {
		_, err := query.Testimonials(ctx, "invalid-uuid")
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})

	t.Run("returns empty validatedSkills when no skill validations exist", func(t *testing.T) {
		results, err := query.Testimonials(ctx, profile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		for _, testimonial := range results {
			if testimonial.ValidatedSkills == nil {
				t.Error("expected ValidatedSkills to be non-nil empty slice")
			}
			if len(testimonial.ValidatedSkills) != 0 {
				t.Errorf("expected 0 validated skills, got %d", len(testimonial.ValidatedSkills))
			}
		}
	})
}

func TestTestimonialsWithValidatedSkills(t *testing.T) {
	userRepo := newMockUserRepository()
	profileRepo := newMockProfileRepository()
	testimonialRepo := newMockTestimonialRepository()
	refLetterRepo := newMockReferenceLetterRepository()
	profileSkillRepo := newMockProfileSkillRepository()
	skillValidationRepo := newMockSkillValidationRepository()

	ctx := context.Background()

	// Create a test user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "testimonials-skills-test@example.com",
		PasswordHash: "hashed",
	}
	mustCreateUser(userRepo, user)

	// Create a profile for the user
	profile := &domain.Profile{
		ID:     uuid.New(),
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("setup: failed to create profile: %v", err)
	}

	// Create profile skills
	skill1 := &domain.ProfileSkill{
		ID:             uuid.New(),
		ProfileID:      profile.ID,
		Name:           "Leadership",
		NormalizedName: "leadership",
		Category:       "soft",
		Source:         domain.ExperienceSourceManual,
	}
	if err := profileSkillRepo.Create(ctx, skill1); err != nil {
		t.Fatalf("setup: failed to create skill: %v", err)
	}

	skill2 := &domain.ProfileSkill{
		ID:             uuid.New(),
		ProfileID:      profile.ID,
		Name:           "Go Programming",
		NormalizedName: "go programming",
		Category:       "technical",
		Source:         domain.ExperienceSourceManual,
	}
	if err := profileSkillRepo.Create(ctx, skill2); err != nil {
		t.Fatalf("setup: failed to create skill: %v", err)
	}

	// Create a reference letter
	refLetter := &domain.ReferenceLetter{
		ID:     uuid.New(),
		UserID: user.ID,
		Status: domain.ReferenceLetterStatusCompleted,
	}
	mustCreateReferenceLetter(refLetterRepo, refLetter)

	// Create skill validations for this reference letter
	sv1 := &domain.SkillValidation{
		ID:                uuid.New(),
		ProfileSkillID:    skill1.ID,
		ReferenceLetterID: refLetter.ID,
		QuoteSnippet:      stringPtr("Excellent leadership skills"),
	}
	if err := skillValidationRepo.Create(ctx, sv1); err != nil {
		t.Fatalf("setup: failed to create skill validation: %v", err)
	}

	sv2 := &domain.SkillValidation{
		ID:                uuid.New(),
		ProfileSkillID:    skill2.ID,
		ReferenceLetterID: refLetter.ID,
		QuoteSnippet:      stringPtr("Expert Go programmer"),
	}
	if err := skillValidationRepo.Create(ctx, sv2); err != nil {
		t.Fatalf("setup: failed to create skill validation: %v", err)
	}

	// Create a testimonial with skills mentioned
	testimonial := &domain.Testimonial{
		ID:                uuid.New(),
		ProfileID:         profile.ID,
		ReferenceLetterID: refLetter.ID,
		Quote:             "A great leader and skilled Go developer.",
		AuthorName:        stringPtr("Jane Manager"),
		AuthorTitle:       stringPtr("Engineering Director"),
		AuthorCompany:     stringPtr("Tech Corp"),
		Relationship:      domain.TestimonialRelationshipManager,
		SkillsMentioned:   []string{"Leadership", "Go Programming"},
	}
	if err := testimonialRepo.Create(ctx, testimonial); err != nil {
		t.Fatalf("setup: failed to create testimonial: %v", err)
	}

	r := resolver.NewResolver(userRepo, newMockFileRepository(), refLetterRepo, newMockResumeRepository(), profileRepo, newMockProfileExperienceRepository(), newMockProfileEducationRepository(), profileSkillRepo, newMockAuthorRepository(), testimonialRepo, skillValidationRepo, newMockExperienceValidationRepository(), storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	query := r.Query()

	t.Run("returns validatedSkills for testimonial", func(t *testing.T) {
		results, err := query.Testimonials(ctx, profile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		if len(results) != 1 {
			t.Fatalf("expected 1 testimonial, got %d", len(results))
		}

		testimonialResult := results[0]
		if testimonialResult.ValidatedSkills == nil {
			t.Fatal("expected ValidatedSkills to be non-nil")
		}

		if len(testimonialResult.ValidatedSkills) != 2 {
			t.Fatalf("expected 2 validated skills, got %d", len(testimonialResult.ValidatedSkills))
		}

		// Verify the skills are correct
		skillNames := make(map[string]bool)
		for _, skill := range testimonialResult.ValidatedSkills {
			skillNames[skill.Name] = true
		}

		if !skillNames["Leadership"] {
			t.Error("expected Leadership skill in validated skills")
		}
		if !skillNames["Go Programming"] {
			t.Error("expected Go Programming skill in validated skills")
		}
	})

	t.Run("different testimonials from different reference letters have different validated skills", func(t *testing.T) {
		// Create another reference letter with a different skill validation
		refLetter2 := &domain.ReferenceLetter{
			ID:     uuid.New(),
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		mustCreateReferenceLetter(refLetterRepo, refLetter2)

		// Only validate skill1 for this reference letter
		sv3 := &domain.SkillValidation{
			ID:                uuid.New(),
			ProfileSkillID:    skill1.ID,
			ReferenceLetterID: refLetter2.ID,
			QuoteSnippet:      stringPtr("Strong leadership"),
		}
		if err := skillValidationRepo.Create(ctx, sv3); err != nil {
			t.Fatalf("setup: failed to create skill validation: %v", err)
		}

		// Create testimonial from second reference letter with skills mentioned
		testimonial2 := &domain.Testimonial{
			ID:                uuid.New(),
			ProfileID:         profile.ID,
			ReferenceLetterID: refLetter2.ID,
			Quote:             "A natural leader.",
			AuthorName:        stringPtr("Bob Peer"),
			AuthorTitle:       stringPtr("Staff Engineer"),
			AuthorCompany:     stringPtr("Tech Corp"),
			Relationship:      domain.TestimonialRelationshipPeer,
			SkillsMentioned:   []string{"Leadership"},
		}
		if err := testimonialRepo.Create(ctx, testimonial2); err != nil {
			t.Fatalf("setup: failed to create testimonial: %v", err)
		}

		results, err := query.Testimonials(ctx, profile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("expected 2 testimonials, got %d", len(results))
		}

		// Find the testimonial from the second reference letter
		var testimonialFromRefLetter2 *model.Testimonial
		for _, t := range results {
			if t.Quote == "A natural leader." {
				testimonialFromRefLetter2 = t
				break
			}
		}

		if testimonialFromRefLetter2 == nil {
			t.Fatal("could not find testimonial from second reference letter")
		}

		// Should only have 1 validated skill (Leadership)
		if len(testimonialFromRefLetter2.ValidatedSkills) != 1 {
			t.Fatalf("expected 1 validated skill for second testimonial, got %d", len(testimonialFromRefLetter2.ValidatedSkills))
		}

		if testimonialFromRefLetter2.ValidatedSkills[0].Name != "Leadership" {
			t.Errorf("expected Leadership skill, got %s", testimonialFromRefLetter2.ValidatedSkills[0].Name)
		}
	})

	t.Run("filters validatedSkills to only skills mentioned in testimonial", func(t *testing.T) {
		// Create another reference letter that validates both skills
		refLetter3 := &domain.ReferenceLetter{
			ID:     uuid.New(),
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		mustCreateReferenceLetter(refLetterRepo, refLetter3)

		// Create skill validations for both skills
		sv4 := &domain.SkillValidation{
			ID:                uuid.New(),
			ProfileSkillID:    skill1.ID,
			ReferenceLetterID: refLetter3.ID,
			QuoteSnippet:      stringPtr("Great leadership"),
		}
		if err := skillValidationRepo.Create(ctx, sv4); err != nil {
			t.Fatalf("setup: failed to create skill validation: %v", err)
		}

		sv5 := &domain.SkillValidation{
			ID:                uuid.New(),
			ProfileSkillID:    skill2.ID,
			ReferenceLetterID: refLetter3.ID,
			QuoteSnippet:      stringPtr("Expert in Go"),
		}
		if err := skillValidationRepo.Create(ctx, sv5); err != nil {
			t.Fatalf("setup: failed to create skill validation: %v", err)
		}

		// Create a testimonial that only mentions one skill (Leadership)
		// Even though the reference letter validates both skills
		testimonial3 := &domain.Testimonial{
			ID:                uuid.New(),
			ProfileID:         profile.ID,
			ReferenceLetterID: refLetter3.ID,
			Quote:             "Exceptional leadership abilities.",
			AuthorName:        stringPtr("Alice Director"),
			AuthorTitle:       stringPtr("VP Engineering"),
			AuthorCompany:     stringPtr("Tech Corp"),
			Relationship:      domain.TestimonialRelationshipManager,
			SkillsMentioned:   []string{"Leadership"}, // Only Leadership mentioned, not Go Programming
		}
		if err := testimonialRepo.Create(ctx, testimonial3); err != nil {
			t.Fatalf("setup: failed to create testimonial: %v", err)
		}

		results, err := query.Testimonials(ctx, profile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		// Find the testimonial we just created
		var filteredTestimonial *model.Testimonial
		for _, t := range results {
			if t.Quote == "Exceptional leadership abilities." {
				filteredTestimonial = t
				break
			}
		}

		if filteredTestimonial == nil {
			t.Fatal("could not find newly created testimonial")
		}

		// Should only have 1 validated skill (Leadership), not Go Programming
		// even though the reference letter validates both
		if len(filteredTestimonial.ValidatedSkills) != 1 {
			t.Fatalf("expected 1 validated skill (filtered by skillsMentioned), got %d", len(filteredTestimonial.ValidatedSkills))
		}

		if filteredTestimonial.ValidatedSkills[0].Name != "Leadership" {
			t.Errorf("expected Leadership skill, got %s", filteredTestimonial.ValidatedSkills[0].Name)
		}
	})

	t.Run("returns empty validatedSkills when skillsMentioned is empty", func(t *testing.T) {
		// Create a reference letter
		refLetter4 := &domain.ReferenceLetter{
			ID:     uuid.New(),
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusCompleted,
		}
		mustCreateReferenceLetter(refLetterRepo, refLetter4)

		// Create a skill validation
		sv6 := &domain.SkillValidation{
			ID:                uuid.New(),
			ProfileSkillID:    skill1.ID,
			ReferenceLetterID: refLetter4.ID,
			QuoteSnippet:      stringPtr("Demonstrates leadership"),
		}
		if err := skillValidationRepo.Create(ctx, sv6); err != nil {
			t.Fatalf("setup: failed to create skill validation: %v", err)
		}

		// Create a testimonial WITHOUT skillsMentioned (empty/nil)
		testimonial4 := &domain.Testimonial{
			ID:                uuid.New(),
			ProfileID:         profile.ID,
			ReferenceLetterID: refLetter4.ID,
			Quote:             "A dependable team member.",
			AuthorName:        stringPtr("Charlie CEO"),
			AuthorTitle:       stringPtr("CEO"),
			AuthorCompany:     stringPtr("Tech Corp"),
			Relationship:      domain.TestimonialRelationshipOther,
			SkillsMentioned:   nil, // No skills mentioned
		}
		if err := testimonialRepo.Create(ctx, testimonial4); err != nil {
			t.Fatalf("setup: failed to create testimonial: %v", err)
		}

		results, err := query.Testimonials(ctx, profile.ID.String())
		if err != nil {
			t.Fatalf("Testimonials query failed: %v", err)
		}

		// Find the testimonial with no skills mentioned
		var emptySkillsTestimonial *model.Testimonial
		for _, t := range results {
			if t.Quote == "A dependable team member." {
				emptySkillsTestimonial = t
				break
			}
		}

		if emptySkillsTestimonial == nil {
			t.Fatal("could not find testimonial with no skills mentioned")
		}

		// Should have 0 validated skills because skillsMentioned is empty
		if len(emptySkillsTestimonial.ValidatedSkills) != 0 {
			t.Fatalf("expected 0 validated skills when skillsMentioned is empty, got %d", len(emptySkillsTestimonial.ValidatedSkills))
		}
	})
}

func TestApplyReferenceLetterValidations(t *testing.T) {
	userRepo := newMockUserRepository()
	fileRepo := newMockFileRepository()
	refLetterRepo := newMockReferenceLetterRepository()
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	skillRepo := newMockProfileSkillRepository()
	testimonialRepo := newMockTestimonialRepository()
	skillValidationRepo := newMockSkillValidationRepository()
	expValidationRepo := newMockExperienceValidationRepository()

	ctx := context.Background()

	// Create a test user
	name := testUserName
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         &name,
	}
	mustCreateUser(userRepo, user)

	// Create a profile for the user
	profile := &domain.Profile{
		ID:     uuid.New(),
		UserID: user.ID,
	}
	if err := profileRepo.Create(ctx, profile); err != nil {
		t.Fatalf("setup: failed to create profile: %v", err)
	}

	// Create a skill for the profile
	skill := &domain.ProfileSkill{
		ID:             uuid.New(),
		ProfileID:      profile.ID,
		Name:           "Go",
		NormalizedName: "go",
		Category:       "TECHNICAL",
		DisplayOrder:   0,
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(ctx, skill); err != nil {
		t.Fatalf("setup: failed to create skill: %v", err)
	}

	// Create an experience for the profile
	experience := &domain.ProfileExperience{
		ID:           uuid.New(),
		ProfileID:    profile.ID,
		Company:      "Acme Inc",
		Title:        "Software Engineer",
		DisplayOrder: 0,
		Source:       domain.ExperienceSourceManual,
	}
	if err := expRepo.Create(ctx, experience); err != nil {
		t.Fatalf("setup: failed to create experience: %v", err)
	}

	// Create extracted data for reference letter
	kubernetesContext := "infrastructure"
	extractedData := domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "John Manager",
			Relationship: domain.AuthorRelationshipManager,
		},
		Testimonials: []domain.ExtractedTestimonial{
			{Quote: "Great team player", SkillsMentioned: []string{"teamwork"}},
		},
		SkillMentions: []domain.ExtractedSkillMention{
			{Skill: "Go", Quote: "Expert in Go programming"},
		},
		ExperienceMentions: []domain.ExtractedExperienceMention{
			{Company: "Acme Inc", Role: "Software Engineer", Quote: "Led the team at Acme Inc"},
		},
		DiscoveredSkills: []domain.DiscoveredSkill{
			{Skill: "Kubernetes", Quote: "Deployed applications on Kubernetes", Context: &kubernetesContext},
		},
	}
	extractedDataJSON, err := json.Marshal(extractedData)
	if err != nil {
		t.Fatalf("setup: failed to marshal extracted data: %v", err)
	}

	// Create a completed reference letter
	refLetter := &domain.ReferenceLetter{
		ID:            uuid.New(),
		UserID:        user.ID,
		Status:        domain.ReferenceLetterStatusCompleted,
		ExtractedData: extractedDataJSON,
	}
	mustCreateReferenceLetter(refLetterRepo, refLetter)

	r := resolver.NewResolver(userRepo, fileRepo, refLetterRepo, newMockResumeRepository(), profileRepo, expRepo, newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), testimonialRepo, skillValidationRepo, expValidationRepo, storage.NewMockStorage(), newMockJobEnqueuer(), testLogger())
	mutation := r.Mutation()

	t.Run("applies skill validations successfully", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID: refLetter.ID.String(),
			SkillValidations: []*model.SkillValidationInput{
				{
					ProfileSkillID: skill.ID.String(),
					QuoteSnippet:   "Expert in Go programming",
				},
			},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}

		successResult, ok := result.(*model.ApplyValidationsResult)
		if !ok {
			t.Fatalf("expected ApplyValidationsResult, got %T", result)
		}

		if successResult.AppliedCount.SkillValidations != 1 {
			t.Errorf("expected 1 skill validation, got %d", successResult.AppliedCount.SkillValidations)
		}
	})

	t.Run("applies experience validations successfully", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID: refLetter.ID.String(),
			SkillValidations:  []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{
				{
					ProfileExperienceID: experience.ID.String(),
					QuoteSnippet:        "Led the team at Acme Inc",
				},
			},
			Testimonials: []*model.TestimonialInput{},
			NewSkills:    []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}

		successResult, ok := result.(*model.ApplyValidationsResult)
		if !ok {
			t.Fatalf("expected ApplyValidationsResult, got %T", result)
		}

		if successResult.AppliedCount.ExperienceValidations != 1 {
			t.Errorf("expected 1 experience validation, got %d", successResult.AppliedCount.ExperienceValidations)
		}
	})

	t.Run("creates testimonials successfully", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID:     refLetter.ID.String(),
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials: []*model.TestimonialInput{
				{
					Quote:           "Great team player",
					SkillsMentioned: []string{"teamwork"},
				},
			},
			NewSkills: []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}

		successResult, ok := result.(*model.ApplyValidationsResult)
		if !ok {
			t.Fatalf("expected ApplyValidationsResult, got %T", result)
		}

		if successResult.AppliedCount.Testimonials != 1 {
			t.Errorf("expected 1 testimonial, got %d", successResult.AppliedCount.Testimonials)
		}
	})

	t.Run("creates new skills successfully", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID:     refLetter.ID.String(),
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills: []*model.NewSkillInput{
				{
					Name:         "Kubernetes",
					Category:     domain.SkillCategoryTechnical,
					QuoteContext: stringPtr("Deployed services to Kubernetes"),
				},
			},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}

		successResult, ok := result.(*model.ApplyValidationsResult)
		if !ok {
			t.Fatalf("expected ApplyValidationsResult, got %T", result)
		}

		if successResult.AppliedCount.NewSkills != 1 {
			t.Errorf("expected 1 new skill, got %d", successResult.AppliedCount.NewSkills)
		}
	})

	t.Run("returns error for invalid user ID", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID:     refLetter.ID.String(),
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, "invalid-uuid", input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		errorResult, ok := result.(*model.ApplyValidationsError)
		if !ok {
			t.Fatalf("expected ApplyValidationsError, got %T", result)
		}

		if errorResult.Message != errMsgInvalidUserIDFormat {
			t.Errorf("expected '%s', got %s", errMsgInvalidUserIDFormat, errorResult.Message)
		}
	})

	t.Run("returns error for non-existent user", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID:     refLetter.ID.String(),
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, uuid.New().String(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		errorResult, ok := result.(*model.ApplyValidationsError)
		if !ok {
			t.Fatalf("expected ApplyValidationsError, got %T", result)
		}

		if errorResult.Message != errMsgUserNotFound {
			t.Errorf("expected '%s', got %s", errMsgUserNotFound, errorResult.Message)
		}
	})

	t.Run("returns error for non-existent reference letter", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID:     uuid.New().String(),
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		errorResult, ok := result.(*model.ApplyValidationsError)
		if !ok {
			t.Fatalf("expected ApplyValidationsError, got %T", result)
		}

		if errorResult.Message != "reference letter not found" {
			t.Errorf("expected 'reference letter not found', got %s", errorResult.Message)
		}
	})

	t.Run("returns error for reference letter belonging to different user", func(t *testing.T) {
		// Create another user
		otherName := "Other User"
		otherUser := &domain.User{
			ID:           uuid.New(),
			Email:        "other@example.com",
			PasswordHash: "hashed",
			Name:         &otherName,
		}
		mustCreateUser(userRepo, otherUser)

		input := model.ApplyValidationsInput{
			ReferenceLetterID:     refLetter.ID.String(), // belongs to first user
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, otherUser.ID.String(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		errorResult, ok := result.(*model.ApplyValidationsError)
		if !ok {
			t.Fatalf("expected ApplyValidationsError, got %T", result)
		}

		if errorResult.Message != "reference letter does not belong to user" {
			t.Errorf("expected 'reference letter does not belong to user', got %s", errorResult.Message)
		}
	})

	t.Run("returns error for pending reference letter", func(t *testing.T) {
		// Create a pending reference letter
		pendingLetter := &domain.ReferenceLetter{
			ID:     uuid.New(),
			UserID: user.ID,
			Status: domain.ReferenceLetterStatusPending,
		}
		mustCreateReferenceLetter(refLetterRepo, pendingLetter)

		input := model.ApplyValidationsInput{
			ReferenceLetterID:     pendingLetter.ID.String(),
			SkillValidations:      []*model.SkillValidationInput{},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		errorResult, ok := result.(*model.ApplyValidationsError)
		if !ok {
			t.Fatalf("expected ApplyValidationsError, got %T", result)
		}

		if errorResult.Message != "reference letter extraction not completed" {
			t.Errorf("expected 'reference letter extraction not completed', got %s", errorResult.Message)
		}
	})

	t.Run("skips non-existent skills gracefully", func(t *testing.T) {
		input := model.ApplyValidationsInput{
			ReferenceLetterID: refLetter.ID.String(),
			SkillValidations: []*model.SkillValidationInput{
				{
					ProfileSkillID: uuid.New().String(), // non-existent skill
					QuoteSnippet:   "Some quote",
				},
			},
			ExperienceValidations: []*model.ExperienceValidationInput{},
			Testimonials:          []*model.TestimonialInput{},
			NewSkills:             []*model.NewSkillInput{},
		}

		result, err := mutation.ApplyReferenceLetterValidations(ctx, user.ID.String(), input)
		if err != nil {
			t.Fatalf("mutation failed: %v", err)
		}

		successResult, ok := result.(*model.ApplyValidationsResult)
		if !ok {
			t.Fatalf("expected ApplyValidationsResult, got %T", result)
		}

		// Should succeed but with 0 skill validations since skill doesn't exist
		if successResult.AppliedCount.SkillValidations != 0 {
			t.Errorf("expected 0 skill validations for non-existent skill, got %d", successResult.AppliedCount.SkillValidations)
		}
	})
}

func TestFileResolver_URL(t *testing.T) {
	ctx := context.Background()

	// Create mock storage with a file
	mockStorage := storage.NewMockStorage()
	storageKey := "test-files/doc.pdf"
	_, err := mockStorage.Upload(ctx, storageKey, bytes.NewReader([]byte("test content")), 12, "application/pdf")
	if err != nil {
		t.Fatalf("failed to upload test file: %v", err)
	}

	// Create resolver with mock storage
	r := resolver.NewResolver(
		newMockUserRepository(),
		newMockFileRepository(),
		newMockReferenceLetterRepository(),
		newMockResumeRepository(),
		newMockProfileRepository(),
		newMockProfileExperienceRepository(),
		newMockProfileEducationRepository(),
		newMockProfileSkillRepository(),
		newMockAuthorRepository(),
		newMockTestimonialRepository(),
		newMockSkillValidationRepository(),
		newMockExperienceValidationRepository(),
		mockStorage,
		newMockJobEnqueuer(),
		testLogger(),
	)

	fileResolver := r.File()

	t.Run("returns presigned URL for existing file", func(t *testing.T) {
		file := &model.File{
			ID:         uuid.New().String(),
			StorageKey: storageKey,
		}

		url, err := fileResolver.URL(ctx, file)
		if err != nil {
			t.Fatalf("URL resolver failed: %v", err)
		}

		// MockStorage returns URL in format: http://mock-storage/{key}?expiry={duration}
		expectedPrefix := "http://mock-storage/" + storageKey
		if !strings.HasPrefix(url, expectedPrefix) {
			t.Errorf("expected URL to start with %s, got %s", expectedPrefix, url)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		file := &model.File{
			ID:         uuid.New().String(),
			StorageKey: "non-existent-key",
		}

		_, err := fileResolver.URL(ctx, file)
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})
}
