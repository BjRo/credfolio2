package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// ExperienceValidationRepository implements domain.ExperienceValidationRepository using PostgreSQL.
type ExperienceValidationRepository struct {
	db bun.IDB
}

// NewExperienceValidationRepository creates a new PostgreSQL experience validation repository.
func NewExperienceValidationRepository(db bun.IDB) *ExperienceValidationRepository {
	return &ExperienceValidationRepository{db: db}
}

// Create persists a new experience validation.
func (r *ExperienceValidationRepository) Create(ctx context.Context, validation *domain.ExperienceValidation) error {
	_, err := r.db.NewInsert().Model(validation).Exec(ctx)
	return err
}

// GetByID retrieves an experience validation by its ID.
func (r *ExperienceValidationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ExperienceValidation, error) {
	validation := new(domain.ExperienceValidation)
	err := r.db.NewSelect().Model(validation).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return validation, nil
}

// GetByProfileExperienceID retrieves all validations for a specific experience.
func (r *ExperienceValidationRepository) GetByProfileExperienceID(ctx context.Context, profileExperienceID uuid.UUID) ([]*domain.ExperienceValidation, error) {
	var validations []*domain.ExperienceValidation
	err := r.db.NewSelect().
		Model(&validations).
		Where("profile_experience_id = ?", profileExperienceID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return validations, nil
}

// GetByReferenceLetterID retrieves all experience validations from a reference letter.
func (r *ExperienceValidationRepository) GetByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) ([]*domain.ExperienceValidation, error) {
	var validations []*domain.ExperienceValidation
	err := r.db.NewSelect().
		Model(&validations).
		Where("reference_letter_id = ?", referenceLetterID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return validations, nil
}

// Delete removes an experience validation by its ID.
func (r *ExperienceValidationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.ExperienceValidation)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// DeleteByReferenceLetterID removes all experience validations from a reference letter.
func (r *ExperienceValidationRepository) DeleteByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.ExperienceValidation)(nil)).
		Where("reference_letter_id = ?", referenceLetterID).
		Exec(ctx)
	return err
}

// CountByProfileExperienceID returns the number of validations for an experience.
func (r *ExperienceValidationRepository) CountByProfileExperienceID(ctx context.Context, profileExperienceID uuid.UUID) (int, error) {
	count, err := r.db.NewSelect().
		Model((*domain.ExperienceValidation)(nil)).
		Where("profile_experience_id = ?", profileExperienceID).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// BatchCountByProfileExperienceIDs returns validation counts for multiple experiences in one query.
func (r *ExperienceValidationRepository) BatchCountByProfileExperienceIDs(ctx context.Context, profileExperienceIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(profileExperienceIDs) == 0 {
		return map[uuid.UUID]int{}, nil
	}

	type Result struct {
		ProfileExperienceID uuid.UUID `bun:"profile_experience_id"`
		Count               int       `bun:"count"`
	}

	var results []Result
	err := r.db.NewSelect().
		Model((*domain.ExperienceValidation)(nil)).
		Column("profile_experience_id").
		ColumnExpr("COUNT(*) as count").
		Where("profile_experience_id IN (?)", bun.In(profileExperienceIDs)).
		Group("profile_experience_id").
		Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	// Initialize all IDs to 0 so callers can distinguish "no validations" from "ID not queried"
	counts := make(map[uuid.UUID]int, len(profileExperienceIDs))
	for _, id := range profileExperienceIDs {
		counts[id] = 0
	}

	// Update with actual counts from query
	for _, r := range results {
		counts[r.ProfileExperienceID] = r.Count
	}

	return counts, nil
}

// Compile-time check that ExperienceValidationRepository implements domain.ExperienceValidationRepository.
var _ domain.ExperienceValidationRepository = (*ExperienceValidationRepository)(nil)
