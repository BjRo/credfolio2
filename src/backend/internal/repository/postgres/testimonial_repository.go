package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// TestimonialRepository implements domain.TestimonialRepository using PostgreSQL.
type TestimonialRepository struct {
	db bun.IDB
}

// NewTestimonialRepository creates a new PostgreSQL testimonial repository.
func NewTestimonialRepository(db bun.IDB) *TestimonialRepository {
	return &TestimonialRepository{db: db}
}

// Create persists a new testimonial.
func (r *TestimonialRepository) Create(ctx context.Context, testimonial *domain.Testimonial) error {
	_, err := r.db.NewInsert().Model(testimonial).Exec(ctx)
	return err
}

// GetByID retrieves a testimonial by its ID.
func (r *TestimonialRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Testimonial, error) {
	testimonial := new(domain.Testimonial)
	err := r.db.NewSelect().Model(testimonial).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return testimonial, nil
}

// GetByProfileID retrieves all testimonials for a profile.
func (r *TestimonialRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.Testimonial, error) {
	var testimonials []*domain.Testimonial
	err := r.db.NewSelect().
		Model(&testimonials).
		Where("profile_id = ?", profileID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return testimonials, nil
}

// GetByReferenceLetterID retrieves all testimonials from a reference letter.
func (r *TestimonialRepository) GetByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) ([]*domain.Testimonial, error) {
	var testimonials []*domain.Testimonial
	err := r.db.NewSelect().
		Model(&testimonials).
		Where("reference_letter_id = ?", referenceLetterID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return testimonials, nil
}

// Delete removes a testimonial by its ID.
func (r *TestimonialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.Testimonial)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// DeleteByReferenceLetterID removes all testimonials from a reference letter.
func (r *TestimonialRepository) DeleteByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.Testimonial)(nil)).
		Where("reference_letter_id = ?", referenceLetterID).
		Exec(ctx)
	return err
}

// Compile-time check that TestimonialRepository implements domain.TestimonialRepository.
var _ domain.TestimonialRepository = (*TestimonialRepository)(nil)
