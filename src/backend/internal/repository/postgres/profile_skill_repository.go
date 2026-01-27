// Package postgres provides PostgreSQL implementations of domain repositories using Bun ORM.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// ProfileSkillRepository implements domain.ProfileSkillRepository using PostgreSQL.
type ProfileSkillRepository struct {
	db *bun.DB
}

// NewProfileSkillRepository creates a new PostgreSQL profile skill repository.
func NewProfileSkillRepository(db *bun.DB) *ProfileSkillRepository {
	return &ProfileSkillRepository{db: db}
}

// Create persists a new profile skill.
func (r *ProfileSkillRepository) Create(ctx context.Context, skill *domain.ProfileSkill) error {
	_, err := r.db.NewInsert().Model(skill).Exec(ctx)
	return err
}

// GetByID retrieves a profile skill by its ID.
func (r *ProfileSkillRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProfileSkill, error) {
	skill := new(domain.ProfileSkill)
	err := r.db.NewSelect().Model(skill).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return skill, nil
}

// GetByProfileID retrieves all profile skills for a profile, ordered by display order.
func (r *ProfileSkillRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.ProfileSkill, error) {
	var skills []*domain.ProfileSkill
	err := r.db.NewSelect().
		Model(&skills).
		Where("profile_id = ?", profileID).
		Order("display_order ASC", "created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

// Update persists changes to an existing profile skill.
func (r *ProfileSkillRepository) Update(ctx context.Context, skill *domain.ProfileSkill) error {
	_, err := r.db.NewUpdate().Model(skill).WherePK().Exec(ctx)
	return err
}

// Delete removes a profile skill by its ID.
func (r *ProfileSkillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.ProfileSkill)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// GetNextDisplayOrder returns the next display order value for a profile.
func (r *ProfileSkillRepository) GetNextDisplayOrder(ctx context.Context, profileID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.NewSelect().
		Model((*domain.ProfileSkill)(nil)).
		ColumnExpr("COALESCE(MAX(display_order), -1)").
		Where("profile_id = ?", profileID).
		Scan(ctx, &maxOrder)
	if err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

// DeleteBySourceResumeID removes all skills extracted from a specific resume.
func (r *ProfileSkillRepository) DeleteBySourceResumeID(ctx context.Context, sourceResumeID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.ProfileSkill)(nil)).
		Where("source_resume_id = ?", sourceResumeID).
		Exec(ctx)
	return err
}

// Compile-time check that ProfileSkillRepository implements domain.ProfileSkillRepository.
var _ domain.ProfileSkillRepository = (*ProfileSkillRepository)(nil)
