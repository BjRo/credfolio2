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

// ProfileRepository implements domain.ProfileRepository using PostgreSQL.
type ProfileRepository struct {
	db bun.IDB
}

// NewProfileRepository creates a new PostgreSQL profile repository.
func NewProfileRepository(db bun.IDB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create persists a new profile.
func (r *ProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	_, err := r.db.NewInsert().Model(profile).Exec(ctx)
	return err
}

// GetByID retrieves a profile by its ID.
func (r *ProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Profile, error) {
	profile := new(domain.Profile)
	err := r.db.NewSelect().Model(profile).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// GetByUserID retrieves a profile by user ID.
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	profile := new(domain.Profile)
	err := r.db.NewSelect().Model(profile).Where("user_id = ?", userID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// GetOrCreateByUserID retrieves a profile by user ID, creating one if it doesn't exist.
func (r *ProfileRepository) GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	profile, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile != nil {
		return profile, nil
	}

	// Create a new profile
	profile = &domain.Profile{
		ID:     uuid.New(),
		UserID: userID,
	}
	if err := r.Create(ctx, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// Update persists changes to an existing profile.
func (r *ProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	_, err := r.db.NewUpdate().Model(profile).WherePK().Exec(ctx)
	return err
}

// Delete removes a profile by its ID.
func (r *ProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.Profile)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// Compile-time check that ProfileRepository implements domain.ProfileRepository.
var _ domain.ProfileRepository = (*ProfileRepository)(nil)
