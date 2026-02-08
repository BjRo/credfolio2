package job

import (
	"context"
	"fmt"

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

// mockProfileSkillRepository implements domain.ProfileSkillRepository for testing.
type mockProfileSkillRepository struct {
	skills              map[uuid.UUID]*domain.ProfileSkill
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
	// Simulate ON CONFLICT DO UPDATE RETURNING * — return existing row's ID on duplicate
	if r.normalizedByProfile[skill.ProfileID] == nil {
		r.normalizedByProfile[skill.ProfileID] = make(map[string]bool)
	}
	if r.normalizedByProfile[skill.ProfileID][skill.NormalizedName] {
		for _, existing := range r.skills {
			if existing.ProfileID == skill.ProfileID && existing.NormalizedName == skill.NormalizedName {
				skill.ID = existing.ID
				break
			}
		}
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
			// Also remove from normalized tracking
			if r.normalizedByProfile[skill.ProfileID] != nil {
				delete(r.normalizedByProfile[skill.ProfileID], skill.NormalizedName)
			}
			delete(r.skills, id)
		}
	}
	return nil
}
