package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines operations for user persistence.
type UserRepository interface {
	// Create persists a new user.
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by their ID.
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail retrieves a user by their email address.
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update persists changes to an existing user.
	Update(ctx context.Context, user *User) error

	// Delete removes a user by their ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// FileRepository defines operations for file metadata persistence.
type FileRepository interface {
	// Create persists a new file record.
	Create(ctx context.Context, file *File) error

	// GetByID retrieves a file by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*File, error)

	// GetByUserID retrieves all files belonging to a user.
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*File, error)

	// Delete removes a file record by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// ReferenceLetterRepository defines operations for reference letter persistence.
type ReferenceLetterRepository interface {
	// Create persists a new reference letter.
	Create(ctx context.Context, letter *ReferenceLetter) error

	// GetByID retrieves a reference letter by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*ReferenceLetter, error)

	// GetByUserID retrieves all reference letters belonging to a user.
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*ReferenceLetter, error)

	// Update persists changes to an existing reference letter.
	Update(ctx context.Context, letter *ReferenceLetter) error

	// Delete removes a reference letter by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
