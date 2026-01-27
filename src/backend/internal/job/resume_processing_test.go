package job

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"backend/internal/domain"
)

// mockProfileRepository implements domain.ProfileRepository for testing.
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

// mockProfileExperienceRepository implements domain.ProfileExperienceRepository for testing.
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

// mockProfileEducationRepository implements domain.ProfileEducationRepository for testing.
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

func stringPtr(s string) *string { return &s }

func newTestWorker() (*ResumeProcessingWorker, *mockProfileRepository, *mockProfileExperienceRepository, *mockProfileEducationRepository) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	worker := &ResumeProcessingWorker{
		profileRepo:    profileRepo,
		profileExpRepo: expRepo,
		profileEduRepo: eduRepo,
	}
	return worker, profileRepo, expRepo, eduRepo
}

func testExtractedData() *domain.ResumeExtractedData {
	return &domain.ResumeExtractedData{
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
	}
}

func TestMaterializeCreatesProfile(t *testing.T) {
	worker, profileRepo, _, _ := newTestWorker()

	err := worker.materializeExtractedData(context.Background(), uuid.New(), uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
	}

	if len(profileRepo.profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profileRepo.profiles))
	}
}

func TestMaterializeCreatesExperience(t *testing.T) {
	worker, _, expRepo, _ := newTestWorker()

	resumeID := uuid.New()
	err := worker.materializeExtractedData(context.Background(), resumeID, uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
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
	worker, _, _, eduRepo := newTestWorker()

	resumeID := uuid.New()
	err := worker.materializeExtractedData(context.Background(), resumeID, uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
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

func TestMaterializeIdempotentReprocessing(t *testing.T) {
	worker, _, expRepo, eduRepo := newTestWorker()

	userID := uuid.New()
	resumeID := uuid.New()

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Old Corp", Title: "Old Role"},
		},
		Education: []domain.Education{
			{Institution: "Old University", Degree: stringPtr("Old Degree")},
		},
	}

	// First materialization
	err := worker.materializeExtractedData(context.Background(), resumeID, userID, data)
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
	}

	err = worker.materializeExtractedData(context.Background(), resumeID, userID, data2)
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
}

func TestMaterializeStoresOriginalData(t *testing.T) {
	worker, _, expRepo, eduRepo := newTestWorker()

	userID := uuid.New()
	resumeID := uuid.New()

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Test Corp", Title: "Engineer", Location: stringPtr("NYC")},
		},
		Education: []domain.Education{
			{Institution: "Test U", Degree: stringPtr("MS"), Field: stringPtr("CS")},
		},
	}

	err := worker.materializeExtractedData(context.Background(), resumeID, userID, data)
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
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
	worker, _, _, eduRepo := newTestWorker()

	data := &domain.ResumeExtractedData{
		Education: []domain.Education{
			{Institution: "No Degree University", Degree: nil},
		},
	}

	err := worker.materializeExtractedData(context.Background(), uuid.New(), uuid.New(), data)
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
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
	worker, _, _, eduRepo := newTestWorker()

	data := &domain.ResumeExtractedData{
		Education: []domain.Education{
			{
				Institution:  "Achievement U",
				Degree:       stringPtr("BS"),
				Achievements: stringPtr("Dean's List, Summa Cum Laude"),
			},
		},
	}

	err := worker.materializeExtractedData(context.Background(), uuid.New(), uuid.New(), data)
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
	}

	for _, edu := range eduRepo.educations {
		if edu.Description == nil || *edu.Description != "Dean's List, Summa Cum Laude" {
			t.Errorf("expected description 'Dean's List, Summa Cum Laude', got %v", edu.Description)
		}
	}
}
