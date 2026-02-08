package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// ResumeRepository implements domain.ResumeRepository using PostgreSQL.
type ResumeRepository struct {
	db bun.IDB
}

// NewResumeRepository creates a new PostgreSQL resume repository.
func NewResumeRepository(db bun.IDB) *ResumeRepository {
	return &ResumeRepository{db: db}
}

// Create persists a new resume.
func (r *ResumeRepository) Create(ctx context.Context, resume *domain.Resume) error {
	_, err := r.db.NewInsert().Model(resume).Exec(ctx)
	return err
}

// GetByID retrieves a resume by its ID.
func (r *ResumeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Resume, error) {
	resume := new(domain.Resume)
	err := r.db.NewSelect().Model(resume).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return resume, nil
}

// GetByUserID retrieves all resumes belonging to a user.
func (r *ResumeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Resume, error) {
	var resumes []*domain.Resume
	err := r.db.NewSelect().Model(&resumes).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return resumes, nil
}

// Update persists changes to an existing resume.
func (r *ResumeRepository) Update(ctx context.Context, resume *domain.Resume) error {
	_, err := r.db.NewUpdate().Model(resume).WherePK().Exec(ctx)
	return err
}

// Delete removes a resume by its ID.
func (r *ResumeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.Resume)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// Compile-time check that ResumeRepository implements domain.ResumeRepository.
var _ domain.ResumeRepository = (*ResumeRepository)(nil)
