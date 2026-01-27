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

// ProfileExperienceRepository implements domain.ProfileExperienceRepository using PostgreSQL.
type ProfileExperienceRepository struct {
	db *bun.DB
}

// NewProfileExperienceRepository creates a new PostgreSQL profile experience repository.
func NewProfileExperienceRepository(db *bun.DB) *ProfileExperienceRepository {
	return &ProfileExperienceRepository{db: db}
}

// Create persists a new profile experience.
func (r *ProfileExperienceRepository) Create(ctx context.Context, experience *domain.ProfileExperience) error {
	_, err := r.db.NewInsert().Model(experience).Exec(ctx)
	return err
}

// GetByID retrieves a profile experience by its ID.
func (r *ProfileExperienceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProfileExperience, error) {
	experience := new(domain.ProfileExperience)
	err := r.db.NewSelect().Model(experience).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return experience, nil
}

// GetByProfileID retrieves all profile experiences for a profile, ordered by display order.
func (r *ProfileExperienceRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.ProfileExperience, error) {
	var experiences []*domain.ProfileExperience
	err := r.db.NewSelect().
		Model(&experiences).
		Where("profile_id = ?", profileID).
		Order("display_order ASC", "created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return experiences, nil
}

// Update persists changes to an existing profile experience.
func (r *ProfileExperienceRepository) Update(ctx context.Context, experience *domain.ProfileExperience) error {
	_, err := r.db.NewUpdate().Model(experience).WherePK().Exec(ctx)
	return err
}

// Delete removes a profile experience by its ID.
func (r *ProfileExperienceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.ProfileExperience)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// GetNextDisplayOrder returns the next display order value for a profile.
func (r *ProfileExperienceRepository) GetNextDisplayOrder(ctx context.Context, profileID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.NewSelect().
		Model((*domain.ProfileExperience)(nil)).
		ColumnExpr("COALESCE(MAX(display_order), -1)").
		Where("profile_id = ?", profileID).
		Scan(ctx, &maxOrder)
	if err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

// DeleteBySourceResumeID removes all experiences extracted from a specific resume.
func (r *ProfileExperienceRepository) DeleteBySourceResumeID(ctx context.Context, sourceResumeID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.ProfileExperience)(nil)).
		Where("source_resume_id = ?", sourceResumeID).
		Exec(ctx)
	return err
}

// Compile-time check that ProfileExperienceRepository implements domain.ProfileExperienceRepository.
var _ domain.ProfileExperienceRepository = (*ProfileExperienceRepository)(nil)
