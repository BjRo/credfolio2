package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// AuthorRepository implements domain.AuthorRepository using PostgreSQL.
type AuthorRepository struct {
	db *bun.DB
}

// NewAuthorRepository creates a new PostgreSQL author repository.
func NewAuthorRepository(db *bun.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// Create persists a new author.
func (r *AuthorRepository) Create(ctx context.Context, author *domain.Author) error {
	_, err := r.db.NewInsert().Model(author).Exec(ctx)
	return err
}

// Upsert creates a new author or returns the existing one if already exists.
// Uses ON CONFLICT DO NOTHING + RETURNING to handle race conditions safely.
func (r *AuthorRepository) Upsert(ctx context.Context, author *domain.Author) (*domain.Author, error) {
	// Try insert with ON CONFLICT DO NOTHING and RETURNING to detect if insert happened
	// Note: The unique constraint is on (profile_id, name, COALESCE(company, ''))
	var inserted domain.Author
	err := r.db.NewInsert().
		Model(author).
		On("CONFLICT (profile_id, name, COALESCE(company, '')) DO NOTHING").
		Returning("*").
		Scan(ctx, &inserted)

	if err != nil {
		// If no rows returned (conflict occurred), fetch existing author
		if errors.Is(err, sql.ErrNoRows) {
			existing, findErr := r.FindByNameAndCompany(ctx, author.ProfileID, author.Name, author.Company)
			if findErr != nil {
				return nil, fmt.Errorf("failed to find existing author after conflict: %w", findErr)
			}
			if existing == nil {
				return nil, fmt.Errorf("author not found after conflict (should be impossible)")
			}
			return existing, nil
		}
		return nil, err
	}

	// Insert succeeded, return the inserted author
	return &inserted, nil
}

// GetByID retrieves an author by its ID.
func (r *AuthorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Author, error) {
	author := new(domain.Author)
	err := r.db.NewSelect().Model(author).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return author, nil
}

// GetByProfileID retrieves all authors for a profile.
func (r *AuthorRepository) GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*domain.Author, error) {
	var authors []*domain.Author
	err := r.db.NewSelect().
		Model(&authors).
		Where("profile_id = ?", profileID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return authors, nil
}

// FindByNameAndCompany finds an author by profile, name, and company.
// Returns nil if not found.
func (r *AuthorRepository) FindByNameAndCompany(ctx context.Context, profileID uuid.UUID, name string, company *string) (*domain.Author, error) {
	author := new(domain.Author)
	query := r.db.NewSelect().
		Model(author).
		Where("profile_id = ?", profileID).
		Where("name = ?", name)

	if company == nil {
		query = query.Where("company IS NULL")
	} else {
		query = query.Where("company = ?", *company)
	}

	err := query.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return author, nil
}

// Update persists changes to an existing author.
func (r *AuthorRepository) Update(ctx context.Context, author *domain.Author) error {
	author.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().
		Model(author).
		WherePK().
		Exec(ctx)
	return err
}

// Delete removes an author by its ID.
func (r *AuthorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.Author)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// Compile-time check that AuthorRepository implements domain.AuthorRepository.
var _ domain.AuthorRepository = (*AuthorRepository)(nil)
