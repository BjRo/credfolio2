package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// ReferenceLetterRepository implements domain.ReferenceLetterRepository using PostgreSQL.
type ReferenceLetterRepository struct {
	db bun.IDB
}

// NewReferenceLetterRepository creates a new PostgreSQL reference letter repository.
func NewReferenceLetterRepository(db bun.IDB) *ReferenceLetterRepository {
	return &ReferenceLetterRepository{db: db}
}

// Create persists a new reference letter.
func (r *ReferenceLetterRepository) Create(ctx context.Context, letter *domain.ReferenceLetter) error {
	_, err := r.db.NewInsert().Model(letter).Exec(ctx)
	return err
}

// GetByID retrieves a reference letter by its ID.
func (r *ReferenceLetterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ReferenceLetter, error) {
	letter := new(domain.ReferenceLetter)
	err := r.db.NewSelect().Model(letter).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return letter, nil
}

// GetByUserID retrieves all reference letters belonging to a user.
func (r *ReferenceLetterRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ReferenceLetter, error) {
	var letters []*domain.ReferenceLetter
	err := r.db.NewSelect().Model(&letters).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return letters, nil
}

// Update persists changes to an existing reference letter.
func (r *ReferenceLetterRepository) Update(ctx context.Context, letter *domain.ReferenceLetter) error {
	_, err := r.db.NewUpdate().Model(letter).WherePK().Exec(ctx)
	return err
}

// Delete removes a reference letter by its ID.
func (r *ReferenceLetterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.ReferenceLetter)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// Compile-time check that ReferenceLetterRepository implements domain.ReferenceLetterRepository.
var _ domain.ReferenceLetterRepository = (*ReferenceLetterRepository)(nil)
