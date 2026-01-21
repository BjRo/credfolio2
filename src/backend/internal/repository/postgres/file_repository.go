package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"backend/internal/domain"
)

// FileRepository implements domain.FileRepository using PostgreSQL.
type FileRepository struct {
	db *bun.DB
}

// NewFileRepository creates a new PostgreSQL file repository.
func NewFileRepository(db *bun.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create persists a new file record.
func (r *FileRepository) Create(ctx context.Context, file *domain.File) error {
	_, err := r.db.NewInsert().Model(file).Exec(ctx)
	return err
}

// GetByID retrieves a file by its ID.
func (r *FileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	file := new(domain.File)
	err := r.db.NewSelect().Model(file).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetByUserID retrieves all files belonging to a user.
func (r *FileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.File, error) {
	var files []*domain.File
	err := r.db.NewSelect().Model(&files).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// Delete removes a file record by its ID.
func (r *FileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*domain.File)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// Compile-time check that FileRepository implements domain.FileRepository.
var _ domain.FileRepository = (*FileRepository)(nil)
