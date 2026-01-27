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

// ProfileEducationRepository implements domain.ProfileEducationRepository using PostgreSQL.
type ProfileEducationRepository struct {
	db *bun.DB
}

// NewProfileEducationRepository creates a new PostgreSQL profile education repository.
func NewProfileEducationRepository(db *bun.DB) *ProfileEducationRepository {
	return &ProfileEducationRepository{db: db}
}

// Create persists a new profile education entry.
func (r *ProfileEducationRepository) Create(ctx context.Context, education *domain.ProfileEducation) error {
	_, err := r.db.NewInsert().Model(education).Exec(ctx)
	return err
}

// GetByID retrieves a profile education entry by its ID.
func (r *ProfileEducationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProfileEducation, error) {
	education := new(domain.ProfileEducation)
	err := r.db.NewSelect().Model(education).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return education, nil
}

// GetByProfileID retrieves all profile education entries for a profile, ordered by display order.
func (r *ProfileEducationRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.ProfileEducation, error) {
	var educations []*domain.ProfileEducation
	err := r.db.NewSelect().
		Model(&educations).
		Where("profile_id = ?", profileID).
		Order("display_order ASC", "created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return educations, nil
}

// Update persists changes to an existing profile education entry.
func (r *ProfileEducationRepository) Update(ctx context.Context, education *domain.ProfileEducation) error {
	_, err := r.db.NewUpdate().Model(education).WherePK().Exec(ctx)
	return err
}

// Delete removes a profile education entry by its ID.
func (r *ProfileEducationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.ProfileEducation)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// GetNextDisplayOrder returns the next display order value for a profile.
func (r *ProfileEducationRepository) GetNextDisplayOrder(ctx context.Context, profileID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.NewSelect().
		Model((*domain.ProfileEducation)(nil)).
		ColumnExpr("COALESCE(MAX(display_order), -1)").
		Where("profile_id = ?", profileID).
		Scan(ctx, &maxOrder)
	if err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

// Compile-time check that ProfileEducationRepository implements domain.ProfileEducationRepository.
var _ domain.ProfileEducationRepository = (*ProfileEducationRepository)(nil)
