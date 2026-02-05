package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

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

// mockProfileSkillRepository implements domain.ProfileSkillRepository for testing.
type mockProfileSkillRepository struct {
	skills           map[uuid.UUID]*domain.ProfileSkill
	normalizedByProfile map[uuid.UUID]map[string]bool // tracks (profile_id, normalized_name) for duplicate detection
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
	// Simulate unique constraint on (profile_id, normalized_name)
	if r.normalizedByProfile[skill.ProfileID] == nil {
		r.normalizedByProfile[skill.ProfileID] = make(map[string]bool)
	}
	if r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] {
		return fmt.Errorf("duplicate key value violates unique constraint \"idx_profile_skills_unique_name\"")
	}
	r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] = true
	r.skills[skill.ID] = skill
	return nil
}

func (r *mockProfileSkillRepository) CreateIgnoreDuplicate(_ context.Context, skill *domain.ProfileSkill) error {
	if skill.ID == uuid.Nil {
		skill.ID = uuid.New()
	}
	// Simulate ON CONFLICT DO NOTHING - silently ignore duplicates
	if r.normalizedByProfile[skill.ProfileID] == nil {
		r.normalizedByProfile[skill.ProfileID] = make(map[string]bool)
	}
	if r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] {
		// Silently ignore duplicate - this is the ON CONFLICT DO NOTHING behavior
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
			// Also remove from normalized tracking
			if r.normalizedByProfile[skill.ProfileID] != nil {
				delete(r.normalizedByProfile[skill.ProfileID], skill.NormalizedName)
			}
			delete(r.skills, id)
		}
	}
	return nil
}

func stringPtr(s string) *string { return &s }

func newTestWorker() (*ResumeProcessingWorker, *mockProfileRepository, *mockProfileExperienceRepository, *mockProfileEducationRepository, *mockProfileSkillRepository) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	skillRepo := newMockProfileSkillRepository()
	worker := &ResumeProcessingWorker{
		profileRepo:      profileRepo,
		profileExpRepo:   expRepo,
		profileEduRepo:   eduRepo,
		profileSkillRepo: skillRepo,
	}
	return worker, profileRepo, expRepo, eduRepo, skillRepo
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
		Skills: []string{"Go", "PostgreSQL", "GraphQL"},
	}
}

func TestMaterializeCreatesProfile(t *testing.T) {
	worker, profileRepo, _, _, _ := newTestWorker()

	err := worker.materializeExtractedData(context.Background(), uuid.New(), uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
	}

	if len(profileRepo.profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profileRepo.profiles))
	}
}

func TestMaterializeCreatesExperience(t *testing.T) {
	worker, _, expRepo, _, _ := newTestWorker()

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
	worker, _, _, eduRepo, _ := newTestWorker()

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
	worker, _, expRepo, eduRepo, skillRepo := newTestWorker()

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
	err := worker.materializeExtractedData(context.Background(), resumeID, userID, data)
	if err != nil {
		t.Fatalf("first materialization failed: %v", err)
	}
	if len(expRepo.experiences) != 1 {
		t.Fatalf("expected 1 experience after first run, got %d", len(expRepo.experiences))
	}
	if len(skillRepo.skills) != 1 {
		t.Fatalf("expected 1 skill after first run, got %d", len(skillRepo.skills))
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
	worker, _, expRepo, eduRepo, _ := newTestWorker()

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
	worker, _, _, eduRepo, _ := newTestWorker()

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
	worker, _, _, eduRepo, _ := newTestWorker()

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

func TestMaterializeCreatesSkills(t *testing.T) {
	worker, _, _, _, skillRepo := newTestWorker()

	resumeID := uuid.New()
	err := worker.materializeExtractedData(context.Background(), resumeID, uuid.New(), testExtractedData())
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
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

// TestDeduplicateSkills tests the skill deduplication helper function.
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
			result := deduplicateSkills(tc.input)
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

// TestMaterializeSkillsWithDuplicatesInExtraction tests that duplicate skills from LLM extraction are deduplicated.
func TestMaterializeSkillsWithDuplicatesInExtraction(t *testing.T) {
	worker, _, _, _, skillRepo := newTestWorker()

	resumeID := uuid.New()
	data := &domain.ResumeExtractedData{
		Skills: []string{"Python", "PYTHON", "python", "Go", "GO", "JavaScript"},
	}

	err := worker.materializeExtractedData(context.Background(), resumeID, uuid.New(), data)
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
	}

	// Should only have 3 unique skills (Python, Go, JavaScript)
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

// TestMaterializeSkillsWithExistingManualSkill tests that skills already existing from manual entry
// don't cause the extraction to fail.
func TestMaterializeSkillsWithExistingManualSkill(t *testing.T) {
	worker, profileRepo, expRepo, eduRepo, skillRepo := newTestWorker()

	userID := uuid.New()
	resumeID := uuid.New()

	// Create a profile first
	profile := &domain.Profile{ID: uuid.New(), UserID: userID}
	profileRepo.profiles[profile.ID] = profile

	// Add an existing manual skill (simulating user added "Python" manually)
	manualSkill := &domain.ProfileSkill{
		ID:             uuid.New(),
		ProfileID:      profile.ID,
		Name:           "Python",
		NormalizedName: "python",
		Category:       "TECHNICAL",
		Source:         domain.ExperienceSourceManual,
		SourceResumeID: nil, // manual skills have no source resume
	}
	// Use Create (not CreateIgnoreDuplicate) to register it in the normalized map
	if err := skillRepo.Create(context.Background(), manualSkill); err != nil {
		t.Fatalf("failed to create manual skill: %v", err)
	}

	// Now try to materialize extracted data that includes "Python" (which already exists)
	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Acme Corp", Title: "Engineer"},
		},
		Education: []domain.Education{
			{Institution: "MIT", Degree: stringPtr("BS")},
		},
		Skills: []string{"Python", "Go", "Rust"}, // Python already exists from manual entry
	}

	// This should NOT fail even though Python already exists
	err := worker.materializeExtractedData(context.Background(), resumeID, userID, data)
	if err != nil {
		t.Fatalf("materializeExtractedData returned error: %v", err)
	}

	// Experience should be created
	if len(expRepo.experiences) != 1 {
		t.Errorf("expected 1 experience, got %d", len(expRepo.experiences))
	}

	// Education should be created
	if len(eduRepo.educations) != 1 {
		t.Errorf("expected 1 education, got %d", len(eduRepo.educations))
	}

	// Skills: manual Python (1) + Go and Rust from extraction (2) = 3 total
	// Python from extraction is silently skipped due to ON CONFLICT DO NOTHING
	if len(skillRepo.skills) != 3 {
		t.Errorf("expected 3 skills (1 manual + 2 new), got %d", len(skillRepo.skills))
	}

	// Verify the manual Python skill is preserved (not overwritten)
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

// mockFailingProfileExperienceRepository always fails on Create for testing partial success.
type mockFailingProfileExperienceRepository struct {
	*mockProfileExperienceRepository
}

func (r *mockFailingProfileExperienceRepository) Create(_ context.Context, _ *domain.ProfileExperience) error {
	return fmt.Errorf("simulated experience create failure")
}

// TestMaterializePartialSuccess_ExperiencesFail tests that education and skills are still saved
// when experiences fail to materialize.
func TestMaterializePartialSuccess_ExperiencesFail(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := &mockFailingProfileExperienceRepository{newMockProfileExperienceRepository()}
	eduRepo := newMockProfileEducationRepository()
	skillRepo := newMockProfileSkillRepository()
	worker := &ResumeProcessingWorker{
		profileRepo:      profileRepo,
		profileExpRepo:   expRepo,
		profileEduRepo:   eduRepo,
		profileSkillRepo: skillRepo,
	}

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Acme Corp", Title: "Engineer"},
		},
		Education: []domain.Education{
			{Institution: "MIT", Degree: stringPtr("BS")},
		},
		Skills: []string{"Go", "Rust"},
	}

	// Should return an error but education and skills should still be saved
	err := worker.materializeExtractedData(context.Background(), uuid.New(), uuid.New(), data)
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
}

// mockFailingProfileSkillRepository always fails on CreateIgnoreDuplicate for testing partial success.
type mockFailingProfileSkillRepository struct {
	*mockProfileSkillRepository
}

func (r *mockFailingProfileSkillRepository) CreateIgnoreDuplicate(_ context.Context, _ *domain.ProfileSkill) error {
	return fmt.Errorf("simulated skill create failure")
}

// TestMaterializePartialSuccess_SkillsFail tests that experiences and education are still saved
// when skills fail to materialize.
func TestMaterializePartialSuccess_SkillsFail(t *testing.T) {
	profileRepo := newMockProfileRepository()
	expRepo := newMockProfileExperienceRepository()
	eduRepo := newMockProfileEducationRepository()
	skillRepo := &mockFailingProfileSkillRepository{newMockProfileSkillRepository()}
	worker := &ResumeProcessingWorker{
		profileRepo:      profileRepo,
		profileExpRepo:   expRepo,
		profileEduRepo:   eduRepo,
		profileSkillRepo: skillRepo,
	}

	data := &domain.ResumeExtractedData{
		Experience: []domain.WorkExperience{
			{Company: "Acme Corp", Title: "Engineer"},
		},
		Education: []domain.Education{
			{Institution: "MIT", Degree: stringPtr("BS")},
		},
		Skills: []string{"Go", "Rust"},
	}

	// Should return an error but experiences and education should still be saved
	err := worker.materializeExtractedData(context.Background(), uuid.New(), uuid.New(), data)
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

// resumeMockExtractor implements domain.DocumentExtractor for span testing.
type resumeMockExtractor struct {
	extractTextResult  string
	extractResumeData  *domain.ResumeExtractedData
	extractTextError   error
	extractResumeError error
}

func (e *resumeMockExtractor) ExtractText(_ context.Context, _ []byte, _ string) (string, error) {
	return e.extractTextResult, e.extractTextError
}

func (e *resumeMockExtractor) ExtractResumeData(_ context.Context, _ string) (*domain.ResumeExtractedData, error) {
	return e.extractResumeData, e.extractResumeError
}

func (e *resumeMockExtractor) ExtractLetterData(_ context.Context, _ string, _ []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	return nil, nil
}

func setupTestTracingForJobs(t *testing.T) *tracetest.InMemoryExporter {
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

func TestExtractResumeData_CreatesParentSpan(t *testing.T) {
	exporter := setupTestTracingForJobs(t)

	extractor := &resumeMockExtractor{
		extractTextResult: "Some resume text",
		extractResumeData: &domain.ResumeExtractedData{
			Name:       "Test User",
			Skills:     []string{},
			Experience: []domain.WorkExperience{},
			Education:  []domain.Education{},
		},
	}

	worker := &ResumeProcessingWorker{
		extractor: extractor,
	}

	_, err := worker.extractResumeData(context.Background(), []byte("pdf data"), "application/pdf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	spans := exporter.GetSpans()
	var found bool
	for _, s := range spans {
		if s.Name == "resume_extraction" {
			found = true
			break
		}
	}
	if !found {
		names := make([]string, len(spans))
		for i, s := range spans {
			names[i] = s.Name
		}
		t.Errorf("expected span named 'resume_extraction', got spans: %v", names)
	}
}

func TestResumeProcessingArgs_InsertOpts_MaxAttempts(t *testing.T) {
	args := ResumeProcessingArgs{}
	opts := args.InsertOpts()
	if opts.MaxAttempts != 2 {
		t.Errorf("MaxAttempts = %d, want 2", opts.MaxAttempts)
	}
}
