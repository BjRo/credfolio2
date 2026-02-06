package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
)

// Mock repositories for testing (duplicated from job package - Go test convention)

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

type mockProfileExperienceRepository struct {
	experiences map[uuid.UUID]*domain.ProfileExperience
}

func newMockProfileExperienceRepository() *mockProfileExperienceRepository {
	return &mockProfileExperienceRepository{experiences: make(map[uuid.UUID]*domain.ProfileExperience)}
}

func (r *mockProfileExperienceRepository) Create(_ context.Context, exp *domain.ProfileExperience) error {
	if exp.ID == uuid.Nil {
		exp.ID = uuid.New()
	}
	r.experiences[exp.ID] = exp
	return nil
}

func (r *mockProfileExperienceRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ProfileExperience, error) {
	exp, ok := r.experiences[id]
	if !ok {
		return nil, nil
	}
	return exp, nil
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

func (r *mockProfileExperienceRepository) Update(_ context.Context, exp *domain.ProfileExperience) error {
	r.experiences[exp.ID] = exp
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

type mockProfileEducationRepository struct {
	educations map[uuid.UUID]*domain.ProfileEducation
}

func newMockProfileEducationRepository() *mockProfileEducationRepository {
	return &mockProfileEducationRepository{educations: make(map[uuid.UUID]*domain.ProfileEducation)}
}

func (r *mockProfileEducationRepository) Create(_ context.Context, edu *domain.ProfileEducation) error {
	if edu.ID == uuid.Nil {
		edu.ID = uuid.New()
	}
	r.educations[edu.ID] = edu
	return nil
}

func (r *mockProfileEducationRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ProfileEducation, error) {
	edu, ok := r.educations[id]
	if !ok {
		return nil, nil
	}
	return edu, nil
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

func (r *mockProfileEducationRepository) Update(_ context.Context, edu *domain.ProfileEducation) error {
	r.educations[edu.ID] = edu
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

type mockProfileSkillRepository struct {
	skills              map[uuid.UUID]*domain.ProfileSkill
	normalizedByProfile map[uuid.UUID]map[string]bool
}

func newMockProfileSkillRepository() *mockProfileSkillRepository {
	return &mockProfileSkillRepository{
		skills:              make(map[uuid.UUID]*domain.ProfileSkill),
		normalizedByProfile: make(map[uuid.UUID]map[string]bool),
	}
}

func (r *mockProfileSkillRepository) Create(_ context.Context, skill *domain.ProfileSkill) error {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}
	if r.normalizedByProfile[skill.ProfileID] == nil {
		r.normalizedByProfile[skill.ProfileID] = make(map[string]bool)
	}
	if r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] {
		return fmt.Errorf("duplicate key value violates unique constraint")
	}
	r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] = true
	r.skills[skill.ID] = skill
	return nil
}

func (r *mockProfileSkillRepository) CreateIgnoreDuplicate(_ context.Context, skill *domain.ProfileSkill) error {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}
	if r.normalizedByProfile[skill.ProfileID] == nil {
		r.normalizedByProfile[skill.ProfileID] = make(map[string]bool)
	}
	if r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] {
		return nil
	}
	r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] = true
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
			if r.normalizedByProfile[skill.ProfileID] != nil {
				delete(r.normalizedByProfile[skill.ProfileID], skill.NormalizedName)
			}
			delete(r.skills, id)
		}
	}
	return nil
}

// Test helpers

func stringPtr(s string) *string { return &s }

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
	a, ok := r.authors[id]
	if !ok {
		return nil, nil
	}
	return a, nil
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
	t, ok := r.testimonials[id]
	if !ok {
		return nil, nil
	}
	return t, nil
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

type mockSkillValidationRepository struct {
	validations map[uuid.UUID]*domain.SkillValidation
}

func newMockSkillValidationRepository() *mockSkillValidationRepository {
	return &mockSkillValidationRepository{validations: make(map[uuid.UUID]*domain.SkillValidation)}
}

func (r *mockSkillValidationRepository) Create(_ context.Context, v *domain.SkillValidation) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	r.validations[v.ID] = v
	return nil
}

func (r *mockSkillValidationRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.SkillValidation, error) {
	v, ok := r.validations[id]
	if !ok {
		return nil, nil
	}
	return v, nil
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

func (r *mockSkillValidationRepository) GetByTestimonialID(_ context.Context, testimonialID uuid.UUID) ([]*domain.SkillValidation, error) {
	var result []*domain.SkillValidation
	for _, v := range r.validations {
		if v.TestimonialID != nil && *v.TestimonialID == testimonialID {
			result = append(result, v)
		}
	}
	return result, nil
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

func (r *mockSkillValidationRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.validations, id)
	return nil
}

type mockExpValidationRepository struct {
	validations map[uuid.UUID]*domain.ExperienceValidation
}

func newMockExpValidationRepository() *mockExpValidationRepository {
	return &mockExpValidationRepository{validations: make(map[uuid.UUID]*domain.ExperienceValidation)}
}

func (r *mockExpValidationRepository) Create(_ context.Context, v *domain.ExperienceValidation) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	r.validations[v.ID] = v
	return nil
}

func (r *mockExpValidationRepository) GetByID(_ context.Context, id uuid.UUID) (*domain.ExperienceValidation, error) {
	v, ok := r.validations[id]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (r *mockExpValidationRepository) GetByProfileExperienceID(_ context.Context, profileExpID uuid.UUID) ([]*domain.ExperienceValidation, error) {
	var result []*domain.ExperienceValidation
	for _, v := range r.validations {
		if v.ProfileExperienceID == profileExpID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockExpValidationRepository) GetByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) ([]*domain.ExperienceValidation, error) {
	var result []*domain.ExperienceValidation
	for _, v := range r.validations {
		if v.ReferenceLetterID == refLetterID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *mockExpValidationRepository) DeleteByReferenceLetterID(_ context.Context, refLetterID uuid.UUID) error {
	for id, v := range r.validations {
		if v.ReferenceLetterID == refLetterID {
			delete(r.validations, id)
		}
	}
	return nil
}

func (r *mockExpValidationRepository) CountByProfileExperienceID(_ context.Context, profileExpID uuid.UUID) (int, error) {
	count := 0
	for _, v := range r.validations {
		if v.ProfileExperienceID == profileExpID {
			count++
		}
	}
	return count, nil
}

func (r *mockExpValidationRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.validations, id)
	return nil
}

func newTestService() (*MaterializationService, *mockProfileRepository, *mockProfileExperienceRepository, *mockProfileEducationRepository, *mockProfileSkillRepository) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	skillRepo := newMockProfileSkillRepository()
	authorRepo := newMockAuthorRepository()
	testimonialRepo := newMockTestimonialRepository()
	svc := NewMaterializationService(profileRepo, expRepo, eduRepo, skillRepo, authorRepo, testimonialRepo, newMockSkillValidationRepository(), newMockExpValidationRepository())
	return svc, profileRepo, expRepo, eduRepo, skillRepo
}

func newTestServiceWithRefs() (*MaterializationService, *mockProfileRepository, *mockAuthorRepository, *mockTestimonialRepository) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	skillRepo := newMockProfileSkillRepository()
	authorRepo := newMockAuthorRepository()
	testimonialRepo := newMockTestimonialRepository()
	svc := NewMaterializationService(profileRepo, expRepo, eduRepo, skillRepo, authorRepo, testimonialRepo, newMockSkillValidationRepository(), newMockExpValidationRepository())
	return svc, profileRepo, authorRepo, testimonialRepo
}

func testExtractedData() *domain.ResumeExtractedData {
	return &domain.ResumeExtractedData{
		Name:     "Jane Doe",
		Email:    stringPtr("jane@example.com"),
		Phone:    stringPtr("+1 555-0100"),
		Location: stringPtr("San Francisco, CA"),
		Summary:  stringPtr("Experienced engineer"),
		Experience: []domain.WorkExperience{
			{
				Company:     "Acme Corp",
				Title:       "Software Engineer",
				Location:    stringPtr("San Francisco"),
				StartDate:   stringPtr("2020-01"),
				EndDate:     stringPtr("2023-06"),
				IsCurrent:   false,
				Description: stringPtr("Built things"),
			},
		},
		Education: []domain.Education{
			{
				Institution: "MIT",
				Degree:      stringPtr("Bachelor of Science"),
				Field:       stringPtr("Computer Science"),
				StartDate:   stringPtr("2016-09"),
				EndDate:     stringPtr("2020-05"),
				GPA:         stringPtr("3.9"),
			},
		},
		Skills: []string{"Go", "PostgreSQL", "GraphQL"},
	}
}

// Tests

func TestMaterializeCreatesProfile(t *testing.T) {
	svc, profileRepo, _, _, _ := newTestService()

	result, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(profileRepo.profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profileRepo.profiles))
	}

	if result.Experiences != 1 {
		t.Errorf("expected 1 experience count, got %d", result.Experiences)
	}
	if result.Educations != 1 {
		t.Errorf("expected 1 education count, got %d", result.Educations)
	}
	if result.Skills != 3 {
		t.Errorf("expected 3 skills count, got %d", result.Skills)
	}
}

func TestMaterializePopulatesProfileHeader(t *testing.T) {
	svc, profileRepo, _, _, _ := newTestService()

	_, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	for _, profile := range profileRepo.profiles {
		if profile.Name == nil || *profile.Name != "Jane Doe" {
			t.Errorf("expected name 'Jane Doe', got %v", profile.Name)
		}
		if profile.Email == nil || *profile.Email != "jane@example.com" {
			t.Errorf("expected email 'jane@example.com', got %v", profile.Email)
		}
		if profile.Phone == nil || *profile.Phone != "+1 555-0100" {
			t.Errorf("expected phone '+1 555-0100', got %v", profile.Phone)
		}
		if profile.Location == nil || *profile.Location != "San Francisco, CA" {
			t.Errorf("expected location 'San Francisco, CA', got %v", profile.Location)
		}
		if profile.Summary == nil || *profile.Summary != "Experienced engineer" {
			t.Errorf("expected summary 'Experienced engineer', got %v", profile.Summary)
		}
	}
}

func TestMaterializePreservesExistingProfileHeader(t *testing.T) {
	svc, profileRepo, _, _, _ := newTestService()

	// Pre-create a profile with existing header data (simulates user-edited profile)
	userID := uuid.New()
	existingName := "Existing Name"
	existingEmail := "existing@example.com"
	profile := &domain.Profile{
		ID:     uuid.New(),
		UserID: userID,
		Name:   &existingName,
		Email:  &existingEmail,
	}
	profileRepo.profiles[profile.ID] = profile

	_, err := svc.MaterializeResumeData(context.Background(), uuid.New(), userID, testExtractedData())
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	updated := profileRepo.profiles[profile.ID]
	// Existing fields should NOT be overwritten
	if *updated.Name != "Existing Name" {
		t.Errorf("expected name to be preserved as 'Existing Name', got %q", *updated.Name)
	}
	if *updated.Email != "existing@example.com" {
		t.Errorf("expected email to be preserved as 'existing@example.com', got %q", *updated.Email)
	}
	// Nil fields should be populated from extraction
	if updated.Phone == nil || *updated.Phone != "+1 555-0100" {
		t.Errorf("expected phone to be filled from extraction, got %v", updated.Phone)
	}
	if updated.Location == nil || *updated.Location != "San Francisco, CA" {
		t.Errorf("expected location to be filled from extraction, got %v", updated.Location)
	}
}

func TestMaterializeCreatesExperience(t *testing.T) {
	svc, _, expRepo, _, _ := newTestService()

	resumeID := uuid.New()
	_, err := svc.MaterializeResumeData(context.Background(), resumeID, uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(expRepo.experiences) != 1 {
		t.Fatalf("expected 1 experience, got %d", len(expRepo.experiences))
	}
	for _, exp := range expRepo.experiences {
		if exp.Company != "Acme Corp" {
			t.Errorf("expected company 'Acme Corp', got %q", exp.Company)
		}
		if exp.Title != "Software Engineer" {
			t.Errorf("expected title 'Software Engineer', got %q", exp.Title)
		}
		if exp.Source != domain.ExperienceSourceResumeExtracted {
			t.Errorf("expected source 'resume_extracted', got %q", exp.Source)
		}
		if exp.SourceResumeID == nil || *exp.SourceResumeID != resumeID {
			t.Error("expected source_resume_id to match resume ID")
		}
		if exp.OriginalData == nil {
			t.Error("expected original_data to be set")
		}
	}
}

func TestMaterializeCreatesEducation(t *testing.T) {
	svc, _, _, eduRepo, _ := newTestService()

	resumeID := uuid.New()
	_, err := svc.MaterializeResumeData(context.Background(), resumeID, uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(eduRepo.educations) != 1 {
		t.Fatalf("expected 1 education, got %d", len(eduRepo.educations))
	}
	for _, edu := range eduRepo.educations {
		if edu.Institution != "MIT" {
			t.Errorf("expected institution 'MIT', got %q", edu.Institution)
		}
		if edu.Degree != "Bachelor of Science" {
			t.Errorf("expected degree 'Bachelor of Science', got %q", edu.Degree)
		}
		if edu.Source != domain.ExperienceSourceResumeExtracted {
			t.Errorf("expected source 'resume_extracted', got %q", edu.Source)
		}
		if edu.SourceResumeID == nil || *edu.SourceResumeID != resumeID {
			t.Error("expected source_resume_id to match resume ID")
		}
		if edu.OriginalData == nil {
			t.Error("expected original_data to be set")
		}
	}
}

func TestMaterializeCreatesSkills(t *testing.T) {
	svc, _, _, _, skillRepo := newTestService()

	resumeID := uuid.New()
	_, err := svc.MaterializeResumeData(context.Background(), resumeID, uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(skillRepo.skills) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(skillRepo.skills))
	}

	names := make(map[string]bool)
	for _, skill := range skillRepo.skills {
		names[skill.Name] = true
		if skill.Category != "TECHNICAL" {
			t.Errorf("expected category 'TECHNICAL', got %q", skill.Category)
		}
		if skill.Source != domain.ExperienceSourceResumeExtracted {
			t.Errorf("expected source 'resume_extracted', got %q", skill.Source)
		}
		if skill.SourceResumeID == nil || *skill.SourceResumeID != resumeID {
			t.Error("expected source_resume_id to match resume ID")
		}
		if skill.NormalizedName != strings.ToLower(skill.Name) {
			t.Errorf("expected normalized name %q, got %q", strings.ToLower(skill.Name), skill.NormalizedName)
		}
	}

	for _, expected := range []string{"Go", "PostgreSQL", "GraphQL"} {
		if !names[expected] {
			t.Errorf("expected skill %q to be created", expected)
		}
	}
}

func TestMaterializeIdempotentReprocessing(t *testing.T) {
	svc, _, expRepo, eduRepo, skillRepo := newTestService()

	userID := uuid.New()
	resumeID := uuid.New()

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Old Corp", Title: "Old Role"},
		},
		Education: []domain.Education{
			{Institution: "Old University", Degree: stringPtr("Old Degree")},
		},
		Skills: []string{"Old Skill"},
	}

	// First materialization
	_, err := svc.MaterializeResumeData(context.Background(), resumeID, userID, data)
	if err != nil {
		t.Fatalf("first materialization failed: %v", err)
	}
	if len(expRepo.experiences) != 1 {
		t.Fatalf("expected 1 experience after first run, got %d", len(expRepo.experiences))
	}

	// Second materialization with different data
	data2 := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "New Corp", Title: "New Role"},
			{Company: "Another Corp", Title: "Another Role"},
		},
		Education: []domain.Education{
			{Institution: "New University", Degree: stringPtr("New Degree")},
		},
		Skills: []string{"New Skill A", "New Skill B"},
	}

	_, err = svc.MaterializeResumeData(context.Background(), resumeID, userID, data2)
	if err != nil {
		t.Fatalf("second materialization failed: %v", err)
	}

	if len(expRepo.experiences) != 2 {
		t.Fatalf("expected 2 experiences after re-processing, got %d", len(expRepo.experiences))
	}
	for _, exp := range expRepo.experiences {
		if exp.Company == "Old Corp" {
			t.Error("old experience should have been deleted during re-processing")
		}
	}

	if len(eduRepo.educations) != 1 {
		t.Fatalf("expected 1 education after re-processing, got %d", len(eduRepo.educations))
	}
	for _, edu := range eduRepo.educations {
		if edu.Institution == "Old University" {
			t.Error("old education should have been deleted during re-processing")
		}
	}

	if len(skillRepo.skills) != 2 {
		t.Fatalf("expected 2 skills after re-processing, got %d", len(skillRepo.skills))
	}
	for _, skill := range skillRepo.skills {
		if skill.Name == "Old Skill" {
			t.Error("old skill should have been deleted during re-processing")
		}
	}
}

func TestMaterializeStoresOriginalData(t *testing.T) {
	svc, _, expRepo, eduRepo, _ := newTestService()

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Test Corp", Title: "Engineer", Location: stringPtr("NYC")},
		},
		Education: []domain.Education{
			{Institution: "Test U", Degree: stringPtr("MS"), Field: stringPtr("CS")},
		},
	}

	_, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), data)
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	for _, exp := range expRepo.experiences {
		var original domain.WorkExperience
		if unmarshalErr := json.Unmarshal(exp.OriginalData, &original); unmarshalErr != nil {
			t.Fatalf("failed to unmarshal experience original data: %v", unmarshalErr)
		}
		if original.Company != "Test Corp" {
			t.Errorf("expected original company 'Test Corp', got %q", original.Company)
		}
		if original.Location == nil || *original.Location != "NYC" {
			t.Error("expected original location 'NYC'")
		}
	}

	for _, edu := range eduRepo.educations {
		var original domain.Education
		if unmarshalErr := json.Unmarshal(edu.OriginalData, &original); unmarshalErr != nil {
			t.Fatalf("failed to unmarshal education original data: %v", unmarshalErr)
		}
		if original.Institution != "Test U" {
			t.Errorf("expected original institution 'Test U', got %q", original.Institution)
		}
	}
}

func TestMaterializeDefaultsDegreeWhenNil(t *testing.T) {
	svc, _, _, eduRepo, _ := newTestService()

	data := &domain.ResumeExtractedData{
		Education: []domain.Education{
			{Institution: "No Degree University", Degree: nil},
		},
	}

	_, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), data)
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(eduRepo.educations) != 1 {
		t.Fatalf("expected 1 education, got %d", len(eduRepo.educations))
	}
	for _, edu := range eduRepo.educations {
		if edu.Degree != "Degree" {
			t.Errorf("expected default degree 'Degree', got %q", edu.Degree)
		}
	}
}

func TestMaterializeMapsAchievementsToDescription(t *testing.T) {
	svc, _, _, eduRepo, _ := newTestService()

	data := &domain.ResumeExtractedData{
		Education: []domain.Education{
			{
				Institution:  "Achievement U",
				Degree:       stringPtr("BS"),
				Achievements: stringPtr("Dean's List, Summa Cum Laude"),
			},
		},
	}

	_, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), data)
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	for _, edu := range eduRepo.educations {
		if edu.Description == nil || *edu.Description != "Dean's List, Summa Cum Laude" {
			t.Errorf("expected description 'Dean's List, Summa Cum Laude', got %v", edu.Description)
		}
	}
}

func TestMaterializeSkillsWithDuplicatesInExtraction(t *testing.T) {
	svc, _, _, _, skillRepo := newTestService()

	resumeID := uuid.New()
	data := &domain.ResumeExtractedData{
		Skills: []string{"Python", "PYTHON", "python", "Go", "GO", "JavaScript"},
	}

	_, err := svc.MaterializeResumeData(context.Background(), resumeID, uuid.New(), data)
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(skillRepo.skills) != 3 {
		t.Fatalf("expected 3 unique skills, got %d", len(skillRepo.skills))
	}

	normalizedNames := make(map[string]bool)
	for _, skill := range skillRepo.skills {
		normalizedNames[skill.NormalizedName] = true
	}

	for _, expected := range []string{"python", "go", "javascript"} {
		if !normalizedNames[expected] {
			t.Errorf("expected normalized skill %q to be created", expected)
		}
	}
}

func TestMaterializeSkillsWithExistingManualSkill(t *testing.T) {
	svc, profileRepo, expRepo, eduRepo, skillRepo := newTestService()

	userID := uuid.New()
	resumeID := uuid.New()

	// Create a profile first
	profile := &domain.Profile{ID: uuid.New(), UserID: userID}
	profileRepo.profiles[profile.ID] = profile

	// Add an existing manual skill
	manualSkill := &domain.ProfileSkill{
		ID:             uuid.New(),
		ProfileID:      profile.ID,
		Name:           "Python",
		NormalizedName: "python",
		Category:       "TECHNICAL",
		Source:         domain.ExperienceSourceManual,
	}
	if err := skillRepo.Create(context.Background(), manualSkill); err != nil {
		t.Fatalf("failed to create manual skill: %v", err)
	}

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Acme Corp", Title: "Engineer"},
		},
		Education: []domain.Education{
			{Institution: "MIT", Degree: stringPtr("BS")},
		},
		Skills: []string{"Python", "Go", "Rust"},
	}

	_, err := svc.MaterializeResumeData(context.Background(), resumeID, userID, data)
	if err != nil {
		t.Fatalf("MaterializeResumeData returned error: %v", err)
	}

	if len(expRepo.experiences) != 1 {
		t.Errorf("expected 1 experience, got %d", len(expRepo.experiences))
	}
	if len(eduRepo.educations) != 1 {
		t.Errorf("expected 1 education, got %d", len(eduRepo.educations))
	}
	// manual Python (1) + Go and Rust from extraction (2) = 3 total
	if len(skillRepo.skills) != 3 {
		t.Errorf("expected 3 skills (1 manual + 2 new), got %d", len(skillRepo.skills))
	}

	var manualPythonFound bool
	for _, skill := range skillRepo.skills {
		if skill.NormalizedName == "python" && skill.Source == domain.ExperienceSourceManual {
			manualPythonFound = true
			break
		}
	}
	if !manualPythonFound {
		t.Error("expected manual Python skill to be preserved")
	}
}

// Mock failing repositories for partial success testing

type mockFailingProfileExperienceRepository struct {
	*mockProfileExperienceRepository
}

func (r *mockFailingProfileExperienceRepository) Create(_ context.Context, _ *domain.ProfileExperience) error {
	return fmt.Errorf("simulated experience create failure")
}

func TestMaterializePartialSuccess_ExperiencesFail(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := &mockFailingProfileExperienceRepository{newMockProfileExperienceRepository()}
	eduRepo := newMockProfileEducationRepository()
	skillRepo := newMockProfileSkillRepository()
	svc := NewMaterializationService(profileRepo, expRepo, eduRepo, skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExpValidationRepository())

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Acme Corp", Title: "Engineer"},
		},
		Education: []domain.Education{
			{Institution: "MIT", Degree: stringPtr("BS")},
		},
		Skills: []string{"Go", "Rust"},
	}

	result, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), data)
	if err == nil {
		t.Fatal("expected error when experiences fail")
	}

	// Education should be created despite experience failure
	if len(eduRepo.educations) != 1 {
		t.Errorf("expected 1 education despite experience failure, got %d", len(eduRepo.educations))
	}
	// Skills should be created despite experience failure
	if len(skillRepo.skills) != 2 {
		t.Errorf("expected 2 skills despite experience failure, got %d", len(skillRepo.skills))
	}
	// Result should still have partial counts
	if result.Educations != 1 {
		t.Errorf("expected 1 education in result, got %d", result.Educations)
	}
	if result.Skills != 2 {
		t.Errorf("expected 2 skills in result, got %d", result.Skills)
	}
}

type mockFailingProfileSkillRepository struct {
	*mockProfileSkillRepository
}

func (r *mockFailingProfileSkillRepository) CreateIgnoreDuplicate(_ context.Context, _ *domain.ProfileSkill) error {
	return fmt.Errorf("simulated skill create failure")
}

func TestMaterializePartialSuccess_SkillsFail(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	skillRepo := &mockFailingProfileSkillRepository{newMockProfileSkillRepository()}
	svc := NewMaterializationService(profileRepo, expRepo, eduRepo, skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), newMockSkillValidationRepository(), newMockExpValidationRepository())

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Acme Corp", Title: "Engineer"},
		},
		Education: []domain.Education{
			{Institution: "MIT", Degree: stringPtr("BS")},
		},
		Skills: []string{"Go", "Rust"},
	}

	_, err := svc.MaterializeResumeData(context.Background(), uuid.New(), uuid.New(), data)
	if err == nil {
		t.Fatal("expected error when skills fail")
	}

	// Experiences should be created despite skill failure
	if len(expRepo.experiences) != 1 {
		t.Errorf("expected 1 experience despite skill failure, got %d", len(expRepo.experiences))
	}
	// Education should be created despite skill failure
	if len(eduRepo.educations) != 1 {
		t.Errorf("expected 1 education despite skill failure, got %d", len(eduRepo.educations))
	}
}

func TestDeduplicateSkills(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "mixed case duplicates",
			input:    []string{"Python", "PYTHON", "python"},
			expected: []string{"Python"},
		},
		{
			name:     "no duplicates",
			input:    []string{"Go", "Rust", "Python"},
			expected: []string{"Go", "Rust", "Python"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "whitespace handling",
			input:    []string{"  Go  ", "go", "GO"},
			expected: []string{"Go"},
		},
		{
			name:     "empty strings filtered",
			input:    []string{"Go", "", "  ", "Rust"},
			expected: []string{"Go", "Rust"},
		},
		{
			name:     "preserves first occurrence case",
			input:    []string{"JavaScript", "javascript", "JAVASCRIPT", "TypeScript"},
			expected: []string{"JavaScript", "TypeScript"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := DeduplicateSkills(tc.input)
			if len(result) != len(tc.expected) {
				t.Fatalf("expected %d skills, got %d: %v", len(tc.expected), len(result), result)
			}
			for i, expected := range tc.expected {
				if result[i] != expected {
					t.Errorf("at index %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

// Reference letter materialization tests

func testExtractedLetterData() *domain.ExtractedLetterData {
	return &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "Jane Smith",
			Title:        stringPtr("CTO"),
			Company:      stringPtr("TechCo"),
			Relationship: domain.AuthorRelationshipManager,
		},
		Testimonials: []domain.ExtractedTestimonial{
			{
				Quote:           "John is an exceptional engineer.",
				SkillsMentioned: []string{"Go", "TypeScript"},
			},
			{
				Quote:           "He excels at problem-solving.",
				SkillsMentioned: nil,
			},
		},
	}
}

func TestMaterializeReferenceLetterCreatesTestimonials(t *testing.T) {
	svc, _, _, testimonialRepo := newTestServiceWithRefs()

	refLetterID := uuid.New()
	result, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, uuid.New(), testExtractedLetterData())
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	if result.Testimonials != 2 {
		t.Errorf("expected 2 testimonials, got %d", result.Testimonials)
	}
	if len(testimonialRepo.testimonials) != 2 {
		t.Fatalf("expected 2 testimonials in repo, got %d", len(testimonialRepo.testimonials))
	}

	var quotes []string
	for _, tm := range testimonialRepo.testimonials {
		quotes = append(quotes, tm.Quote)
	}
	for _, expected := range []string{"John is an exceptional engineer.", "He excels at problem-solving."} {
		found := false
		for _, q := range quotes {
			if q == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected quote %q to be created", expected)
		}
	}
}

func TestMaterializeReferenceLetterCreatesAuthor(t *testing.T) {
	svc, _, authorRepo, testimonialRepo := newTestServiceWithRefs()

	refLetterID := uuid.New()
	_, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, uuid.New(), testExtractedLetterData())
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	if len(authorRepo.authors) != 1 {
		t.Fatalf("expected 1 author, got %d", len(authorRepo.authors))
	}

	var author *domain.Author
	for _, a := range authorRepo.authors {
		author = a
	}
	if author.Name != "Jane Smith" {
		t.Errorf("expected author name 'Jane Smith', got %q", author.Name)
	}
	if author.Title == nil || *author.Title != "CTO" {
		t.Errorf("expected author title 'CTO', got %v", author.Title)
	}
	if author.Company == nil || *author.Company != "TechCo" {
		t.Errorf("expected author company 'TechCo', got %v", author.Company)
	}

	// Verify testimonials are linked to author
	for _, testimonial := range testimonialRepo.testimonials {
		if testimonial.AuthorID == nil || *testimonial.AuthorID != author.ID {
			t.Error("expected testimonial to be linked to author")
		}
	}
}

func TestMaterializeReferenceLetterReusesExistingAuthor(t *testing.T) {
	svc, profileRepo, authorRepo, _ := newTestServiceWithRefs()

	userID := uuid.New()
	refLetterID := uuid.New()

	// Pre-create a profile and author
	profile := &domain.Profile{ID: uuid.New(), UserID: userID}
	profileRepo.profiles[profile.ID] = profile

	existingAuthor := &domain.Author{
		ID:        uuid.New(),
		ProfileID: profile.ID,
		Name:      "Jane Smith",
		Company:   stringPtr("TechCo"),
	}
	authorRepo.authors[existingAuthor.ID] = existingAuthor

	_, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, userID, testExtractedLetterData())
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	// Should still have exactly 1 author (reused existing)
	if len(authorRepo.authors) != 1 {
		t.Errorf("expected 1 author (reused), got %d", len(authorRepo.authors))
	}
}

func TestMaterializeReferenceLetterSetsRelationship(t *testing.T) {
	svc, _, _, testimonialRepo := newTestServiceWithRefs()

	refLetterID := uuid.New()
	_, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, uuid.New(), testExtractedLetterData())
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	for _, testimonial := range testimonialRepo.testimonials {
		if testimonial.Relationship != domain.TestimonialRelationshipManager {
			t.Errorf("expected relationship 'manager', got %q", testimonial.Relationship)
		}
	}
}

func TestMaterializeReferenceLetterSetsAuthorFields(t *testing.T) {
	svc, _, _, testimonialRepo := newTestServiceWithRefs()

	refLetterID := uuid.New()
	_, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, uuid.New(), testExtractedLetterData())
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	for _, testimonial := range testimonialRepo.testimonials {
		if testimonial.AuthorName == nil || *testimonial.AuthorName != "Jane Smith" {
			t.Errorf("expected author name 'Jane Smith', got %v", testimonial.AuthorName)
		}
		if testimonial.AuthorTitle == nil || *testimonial.AuthorTitle != "CTO" {
			t.Errorf("expected author title 'CTO', got %v", testimonial.AuthorTitle)
		}
		if testimonial.AuthorCompany == nil || *testimonial.AuthorCompany != "TechCo" {
			t.Errorf("expected author company 'TechCo', got %v", testimonial.AuthorCompany)
		}
	}
}

func TestMaterializeReferenceLetterSetsSkillsMentioned(t *testing.T) {
	svc, _, _, testimonialRepo := newTestServiceWithRefs()

	refLetterID := uuid.New()
	_, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, uuid.New(), testExtractedLetterData())
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	var foundWithSkills bool
	for _, testimonial := range testimonialRepo.testimonials {
		if testimonial.Quote == "John is an exceptional engineer." {
			if len(testimonial.SkillsMentioned) != 2 {
				t.Errorf("expected 2 skills mentioned, got %d", len(testimonial.SkillsMentioned))
			}
			foundWithSkills = true
		}
	}
	if !foundWithSkills {
		t.Error("expected to find testimonial with skills mentioned")
	}
}

func TestMaterializeReferenceLetterIdempotent(t *testing.T) {
	svc, _, _, testimonialRepo := newTestServiceWithRefs()

	userID := uuid.New()
	refLetterID := uuid.New()

	// First materialization
	_, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, userID, testExtractedLetterData())
	if err != nil {
		t.Fatalf("first materialization failed: %v", err)
	}
	if len(testimonialRepo.testimonials) != 2 {
		t.Fatalf("expected 2 testimonials after first run, got %d", len(testimonialRepo.testimonials))
	}

	// Second materialization with different data
	data2 := &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "Jane Smith",
			Company:      stringPtr("TechCo"),
			Relationship: domain.AuthorRelationshipManager,
		},
		Testimonials: []domain.ExtractedTestimonial{
			{Quote: "New testimonial quote."},
		},
	}

	result, err := svc.MaterializeReferenceLetterData(context.Background(), refLetterID, userID, data2)
	if err != nil {
		t.Fatalf("second materialization failed: %v", err)
	}

	if result.Testimonials != 1 {
		t.Errorf("expected 1 testimonial after re-processing, got %d", result.Testimonials)
	}
	if len(testimonialRepo.testimonials) != 1 {
		t.Errorf("expected 1 testimonial in repo after re-processing, got %d", len(testimonialRepo.testimonials))
	}

	for _, testimonial := range testimonialRepo.testimonials {
		if testimonial.Quote != "New testimonial quote." {
			t.Errorf("expected new quote, got %q", testimonial.Quote)
		}
	}
}

func TestMaterializeReferenceLetterWithNoTestimonials(t *testing.T) {
	svc, _, authorRepo, testimonialRepo := newTestServiceWithRefs()

	data := &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         "Jane Smith",
			Relationship: domain.AuthorRelationshipOther,
		},
		Testimonials: []domain.ExtractedTestimonial{},
	}

	result, err := svc.MaterializeReferenceLetterData(context.Background(), uuid.New(), uuid.New(), data)
	if err != nil {
		t.Fatalf("MaterializeReferenceLetterData returned error: %v", err)
	}

	if result.Testimonials != 0 {
		t.Errorf("expected 0 testimonials, got %d", result.Testimonials)
	}
	if len(testimonialRepo.testimonials) != 0 {
		t.Errorf("expected 0 testimonials in repo, got %d", len(testimonialRepo.testimonials))
	}
	// Should not create an author when there are no testimonials
	if len(authorRepo.authors) != 0 {
		t.Errorf("expected 0 authors when no testimonials, got %d", len(authorRepo.authors))
	}
}

func TestCrossReferenceValidationsMatchesSkills(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	skillRepo := newMockProfileSkillRepository()
	skillValRepo := newMockSkillValidationRepository()
	expValRepo := newMockExpValidationRepository()
	svc := NewMaterializationService(profileRepo, expRepo, eduRepo, skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), skillValRepo, expValRepo)

	profileID := uuid.New()

	// Pre-populate profile skills (as if materialized from resume)
	goSkill := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Go", NormalizedName: "go"}
	pgSkill := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "PostgreSQL", NormalizedName: "postgresql"}
	skillRepo.skills[goSkill.ID] = goSkill
	skillRepo.skills[pgSkill.ID] = pgSkill

	// Pre-populate profile experience
	exp := &domain.ProfileExperience{ID: uuid.New(), ProfileID: profileID, Company: "Acme Corp", Title: "Engineer"}
	expRepo.experiences[exp.ID] = exp

	refLetterID := uuid.New()
	letterData := &domain.ExtractedLetterData{
		SkillMentions: []domain.ExtractedSkillMention{
			{Skill: "Go", Quote: "Expert Go programmer"},
			{Skill: "PostgreSQL", Quote: "Deep PostgreSQL knowledge"},
			{Skill: "Kubernetes", Quote: "K8s experience"}, // no match
		},
		ExperienceMentions: []domain.ExtractedExperienceMention{
			{Company: "Acme Corp", Role: "Engineer", Quote: "Great work at Acme"},
			{Company: "Other Inc", Role: "Intern", Quote: "Summer internship"},  // no match
		},
	}

	result, err := svc.CrossReferenceValidations(context.Background(), profileID, refLetterID, letterData)
	if err != nil {
		t.Fatalf("CrossReferenceValidations returned error: %v", err)
	}

	if result.SkillValidations != 2 {
		t.Errorf("expected 2 skill validations, got %d", result.SkillValidations)
	}
	if result.ExperienceValidations != 1 {
		t.Errorf("expected 1 experience validation, got %d", result.ExperienceValidations)
	}
	if len(skillValRepo.validations) != 2 {
		t.Errorf("expected 2 skill validation records, got %d", len(skillValRepo.validations))
	}
	if len(expValRepo.validations) != 1 {
		t.Errorf("expected 1 experience validation record, got %d", len(expValRepo.validations))
	}

	// Verify the skill validations reference the correct skills and quotes
	for _, sv := range skillValRepo.validations {
		if sv.ReferenceLetterID != refLetterID {
			t.Errorf("expected reference letter ID %s, got %s", refLetterID, sv.ReferenceLetterID)
		}
		if sv.ProfileSkillID != goSkill.ID && sv.ProfileSkillID != pgSkill.ID {
			t.Errorf("unexpected skill ID %s", sv.ProfileSkillID)
		}
	}
}

func TestCrossReferenceValidationsCaseInsensitive(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	skillRepo := newMockProfileSkillRepository()
	skillValRepo := newMockSkillValidationRepository()
	expValRepo := newMockExpValidationRepository()
	svc := NewMaterializationService(profileRepo, expRepo, newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), skillValRepo, expValRepo)

	profileID := uuid.New()

	// Profile skill with lowercase normalized name
	skill := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "AWS", NormalizedName: "aws"}
	skillRepo.skills[skill.ID] = skill

	// Experience with mixed case
	exp := &domain.ProfileExperience{ID: uuid.New(), ProfileID: profileID, Company: "AL Talent, Inc. DBA Wellfound", Title: "SEM"}
	expRepo.experiences[exp.ID] = exp

	letterData := &domain.ExtractedLetterData{
		SkillMentions: []domain.ExtractedSkillMention{
			{Skill: "AWS", Quote: "Used AWS extensively"}, // uppercase mention vs lowercase normalized
		},
		ExperienceMentions: []domain.ExtractedExperienceMention{
			{Company: "al talent, inc. dba wellfound", Role: "SEM", Quote: "Led teams at Wellfound"}, // lowercase mention
		},
	}

	result, err := svc.CrossReferenceValidations(context.Background(), profileID, uuid.New(), letterData)
	if err != nil {
		t.Fatalf("CrossReferenceValidations returned error: %v", err)
	}

	if result.SkillValidations != 1 {
		t.Errorf("expected 1 skill validation (case-insensitive match), got %d", result.SkillValidations)
	}
	if result.ExperienceValidations != 1 {
		t.Errorf("expected 1 experience validation (case-insensitive match), got %d", result.ExperienceValidations)
	}
}

func TestCrossReferenceValidationsMatchesDiscoveredSkills(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	skillRepo := newMockProfileSkillRepository()
	skillValRepo := newMockSkillValidationRepository()
	expValRepo := newMockExpValidationRepository()
	svc := NewMaterializationService(profileRepo, expRepo, newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), skillValRepo, expValRepo)

	profileID := uuid.New()

	// Pre-populate profile skills
	mentoring := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Mentoring", NormalizedName: "mentoring"}
	techDoc := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Technical Documentation", NormalizedName: "technical documentation"}
	skillRepo.skills[mentoring.ID] = mentoring
	skillRepo.skills[techDoc.ID] = techDoc

	refLetterID := uuid.New()
	// SkillMentions is empty; skills are in DiscoveredSkills (as the LLM often produces)
	letterData := &domain.ExtractedLetterData{
		SkillMentions: []domain.ExtractedSkillMention{},
		DiscoveredSkills: []domain.DiscoveredSkill{
			{Skill: "mentoring", Quote: "excellent mentor"},
			{Skill: "Technical Documentation", Quote: "introduced Diátaxis framework"},
			{Skill: "cloud cost optimization", Quote: "reduced AWS costs"}, // no match
		},
	}

	result, err := svc.CrossReferenceValidations(context.Background(), profileID, refLetterID, letterData)
	if err != nil {
		t.Fatalf("CrossReferenceValidations returned error: %v", err)
	}

	if result.SkillValidations != 2 {
		t.Errorf("expected 2 skill validations from discovered skills, got %d", result.SkillValidations)
	}
	if len(skillValRepo.validations) != 2 {
		t.Errorf("expected 2 skill validation records, got %d", len(skillValRepo.validations))
	}
}

func TestCrossReferenceValidationsDeduplicatesAcrossSources(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	skillRepo := newMockProfileSkillRepository()
	skillValRepo := newMockSkillValidationRepository()
	expValRepo := newMockExpValidationRepository()
	svc := NewMaterializationService(profileRepo, expRepo, newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), skillValRepo, expValRepo)

	profileID := uuid.New()

	goSkill := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Go", NormalizedName: "go"}
	skillRepo.skills[goSkill.ID] = goSkill

	refLetterID := uuid.New()
	// Same skill appears in both SkillMentions and DiscoveredSkills
	letterData := &domain.ExtractedLetterData{
		SkillMentions: []domain.ExtractedSkillMention{
			{Skill: "Go", Quote: "Expert Go developer"},
		},
		DiscoveredSkills: []domain.DiscoveredSkill{
			{Skill: "Go", Quote: "Strong Go skills"},
		},
	}

	result, err := svc.CrossReferenceValidations(context.Background(), profileID, refLetterID, letterData)
	if err != nil {
		t.Fatalf("CrossReferenceValidations returned error: %v", err)
	}

	// Should only create 1 validation, not 2 (deduplication)
	if result.SkillValidations != 1 {
		t.Errorf("expected 1 skill validation (deduplicated), got %d", result.SkillValidations)
	}
	if len(skillValRepo.validations) != 1 {
		t.Errorf("expected 1 skill validation record, got %d", len(skillValRepo.validations))
	}
}

func TestCrossReferenceValidationsSubstringMatching(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	skillRepo := newMockProfileSkillRepository()
	skillValRepo := newMockSkillValidationRepository()
	expValRepo := newMockExpValidationRepository()
	svc := NewMaterializationService(profileRepo, expRepo, newMockProfileEducationRepository(), skillRepo, newMockAuthorRepository(), newMockTestimonialRepository(), skillValRepo, expValRepo)

	profileID := uuid.New()

	// Profile skills with concise names (from resume extraction)
	techDoc := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Technical Documentation", NormalizedName: "technical documentation"}
	incResp := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Incident Response", NormalizedName: "incident response"}
	seo := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "SEO", NormalizedName: "seo"}
	teamLead := &domain.ProfileSkill{ID: uuid.New(), ProfileID: profileID, Name: "Team Leadership", NormalizedName: "team leadership"}
	skillRepo.skills[techDoc.ID] = techDoc
	skillRepo.skills[incResp.ID] = incResp
	skillRepo.skills[seo.ID] = seo
	skillRepo.skills[teamLead.ID] = teamLead

	refLetterID := uuid.New()
	// Reference letter uses more verbose/descriptive skill names
	letterData := &domain.ExtractedLetterData{
		DiscoveredSkills: []domain.DiscoveredSkill{
			{Skill: "technical documentation (Diátaxis framework)", Quote: "introduced the Diátaxis framework"},
			{Skill: "incident response program design", Quote: "designed the incident response program"},
			{Skill: "search engine optimization (SEO)", Quote: "doubled organic traffic"},
			{Skill: "interim team leadership", Quote: "served as interim leader"},
			{Skill: "cross-functional collaboration", Quote: "contributed across multiple dimensions"}, // no match
		},
	}

	result, err := svc.CrossReferenceValidations(context.Background(), profileID, refLetterID, letterData)
	if err != nil {
		t.Fatalf("CrossReferenceValidations returned error: %v", err)
	}

	if result.SkillValidations != 4 {
		t.Errorf("expected 4 skill validations (substring matches), got %d", result.SkillValidations)
	}

	// Verify the correct profile skills were matched
	matchedSkillIDs := make(map[uuid.UUID]bool)
	for _, sv := range skillValRepo.validations {
		matchedSkillIDs[sv.ProfileSkillID] = true
	}
	for _, expected := range []*domain.ProfileSkill{techDoc, incResp, seo, teamLead} {
		if !matchedSkillIDs[expected.ID] {
			t.Errorf("expected skill %q to be matched, but it wasn't", expected.Name)
		}
	}
}
