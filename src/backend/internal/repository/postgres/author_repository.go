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
	db bun.IDB
}

// NewAuthorRepository creates a new PostgreSQL author repository.
// Accepts bun.IDB to support both regular DB operations and transactions.
func NewAuthorRepository(db bun.IDB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// Create persists a new author.
func (r *AuthorRepository) Create(ctx context.Context, author *domain.Author) error {
	_, err := r.db.NewInsert().Model(author).Exec(ctx)
	return err
}

// Upsert creates a new author or returns the existing one if already exists.
// Uses ON CONFLICT DO UPDATE to safely handle concurrent inserts without a race condition.
// Returns the inserted author on success, or the existing author if a conflict occurred.
func (r *AuthorRepository) Upsert(ctx context.Context, author *domain.Author) (*domain.Author, error) {
	// Use ON CONFLICT DO UPDATE with a no-op update to return the existing row on conflict
	// This avoids the race condition of DO NOTHING (which returns no rows) + separate SELECT
	// The unique index is on (profile_id, name, COALESCE(company, ''))
	// We qualify with table name to reference existing value: "SET updated_at = a.updated_at"
	var result domain.Author
	err := r.db.NewInsert().
		Model(author).
		On("CONFLICT (profile_id, name, COALESCE(company, '')) DO UPDATE SET updated_at = a.updated_at").
		Returning("*").
		Scan(ctx, &result)

	if err != nil {
		return nil, fmt.Errorf("upsert author (profile=%s, name=%s): %w", author.ProfileID, author.Name, err)
	}

	return &result, nil
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

// GetByIDs retrieves multiple authors by their IDs in a single query.
func (r *AuthorRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*domain.Author, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]*domain.Author{}, nil
	}

	var authors []*domain.Author
	err := r.db.NewSelect().
		Model(&authors).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// Build map for efficient lookup
	result := make(map[uuid.UUID]*domain.Author, len(authors))
	for _, author := range authors {
		result[author.ID] = author
	}

	return result, nil
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
