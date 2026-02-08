package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// SkillValidationRepository implements domain.SkillValidationRepository using PostgreSQL.
type SkillValidationRepository struct {
	db bun.IDB
}

// NewSkillValidationRepository creates a new PostgreSQL skill validation repository.
func NewSkillValidationRepository(db bun.IDB) *SkillValidationRepository {
	return &SkillValidationRepository{db: db}
}

// Create persists a new skill validation.
func (r *SkillValidationRepository) Create(ctx context.Context, validation *domain.SkillValidation) error {
	_, err := r.db.NewInsert().Model(validation).Exec(ctx)
	return err
}

// GetByID retrieves a skill validation by its ID.
func (r *SkillValidationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.SkillValidation, error) {
	validation := new(domain.SkillValidation)
	err := r.db.NewSelect().Model(validation).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return validation, nil
}

// GetByProfileSkillID retrieves all validations for a specific skill.
func (r *SkillValidationRepository) GetByProfileSkillID(ctx context.Context, profileSkillID uuid.UUID) ([]*domain.SkillValidation, error) {
	var validations []*domain.SkillValidation
	err := r.db.NewSelect().
		Model(&validations).
		Where("profile_skill_id = ?", profileSkillID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return validations, nil
}

// GetByReferenceLetterID retrieves all skill validations from a reference letter.
func (r *SkillValidationRepository) GetByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) ([]*domain.SkillValidation, error) {
	var validations []*domain.SkillValidation
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

// GetByTestimonialID retrieves all skill validations for a specific testimonial.
func (r *SkillValidationRepository) GetByTestimonialID(ctx context.Context, testimonialID uuid.UUID) ([]*domain.SkillValidation, error) {
	var validations []*domain.SkillValidation
	err := r.db.NewSelect().
		Model(&validations).
		Where("testimonial_id = ?", testimonialID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return validations, nil
}

// Delete removes a skill validation by its ID.
func (r *SkillValidationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.SkillValidation)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// DeleteByReferenceLetterID removes all skill validations from a reference letter.
func (r *SkillValidationRepository) DeleteByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.SkillValidation)(nil)).
		Where("reference_letter_id = ?", referenceLetterID).
		Exec(ctx)
	return err
}

// CountByProfileSkillID returns the number of validations for a skill.
func (r *SkillValidationRepository) CountByProfileSkillID(ctx context.Context, profileSkillID uuid.UUID) (int, error) {
	count, err := r.db.NewSelect().
		Model((*domain.SkillValidation)(nil)).
		Where("profile_skill_id = ?", profileSkillID).
		Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Compile-time check that SkillValidationRepository implements domain.SkillValidationRepository.
var _ domain.SkillValidationRepository = (*SkillValidationRepository)(nil)
