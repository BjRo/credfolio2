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

	// GetByUserIDAndContentHash retrieves a file by user ID and content hash.
	// Returns nil if no matching file exists.
	GetByUserIDAndContentHash(ctx context.Context, userID uuid.UUID, contentHash string) (*File, error)

	// Update persists changes to an existing file record.
	Update(ctx context.Context, file *File) error

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

// AuthorRepository defines operations for author persistence.
type AuthorRepository interface {
	// Create persists a new author.
	Create(ctx context.Context, author *Author) error

	// Upsert creates a new author or returns the existing one if a duplicate exists.
	// This handles concurrent creation attempts safely using the database unique constraint.
	Upsert(ctx context.Context, author *Author) (*Author, error)

	// GetByID retrieves an author by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Author, error)

	// GetByProfileID retrieves all authors for a profile.
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*Author, error)

	// FindByNameAndCompany finds an author by profile, name, and company.
	// Returns nil if not found.
	FindByNameAndCompany(ctx context.Context, profileID uuid.UUID, name string, company *string) (*Author, error)

	// Update persists changes to an existing author.
	Update(ctx context.Context, author *Author) error

	// Delete removes an author by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// TestimonialRepository defines operations for testimonial persistence.
type TestimonialRepository interface {
	// Create persists a new testimonial.
	Create(ctx context.Context, testimonial *Testimonial) error

	// GetByID retrieves a testimonial by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Testimonial, error)

	// GetByProfileID retrieves all testimonials for a profile.
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*Testimonial, error)

	// GetByReferenceLetterID retrieves all testimonials from a reference letter.
	GetByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) ([]*Testimonial, error)

	// Delete removes a testimonial by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByReferenceLetterID removes all testimonials from a reference letter.
	DeleteByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) error
}

// SkillValidationRepository defines operations for skill validation persistence.
type SkillValidationRepository interface {
	// Create persists a new skill validation.
	Create(ctx context.Context, validation *SkillValidation) error

	// GetByID retrieves a skill validation by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*SkillValidation, error)

	// GetByProfileSkillID retrieves all validations for a specific skill.
	GetByProfileSkillID(ctx context.Context, profileSkillID uuid.UUID) ([]*SkillValidation, error)

	// GetByReferenceLetterID retrieves all skill validations from a reference letter.
	GetByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) ([]*SkillValidation, error)

	// GetByTestimonialID retrieves all skill validations for a specific testimonial.
	GetByTestimonialID(ctx context.Context, testimonialID uuid.UUID) ([]*SkillValidation, error)

	// Delete removes a skill validation by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByReferenceLetterID removes all skill validations from a reference letter.
	DeleteByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) error

	// CountByProfileSkillID returns the number of validations for a skill.
	CountByProfileSkillID(ctx context.Context, profileSkillID uuid.UUID) (int, error)
}

// ExperienceValidationRepository defines operations for experience validation persistence.
type ExperienceValidationRepository interface {
	// Create persists a new experience validation.
	Create(ctx context.Context, validation *ExperienceValidation) error

	// GetByID retrieves an experience validation by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*ExperienceValidation, error)

	// GetByProfileExperienceID retrieves all validations for a specific experience.
	GetByProfileExperienceID(ctx context.Context, profileExperienceID uuid.UUID) ([]*ExperienceValidation, error)

	// GetByReferenceLetterID retrieves all experience validations from a reference letter.
	GetByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) ([]*ExperienceValidation, error)

	// Delete removes an experience validation by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByReferenceLetterID removes all experience validations from a reference letter.
	DeleteByReferenceLetterID(ctx context.Context, referenceLetterID uuid.UUID) error

	// CountByProfileExperienceID returns the number of validations for an experience.
	CountByProfileExperienceID(ctx context.Context, profileExperienceID uuid.UUID) (int, error)
}
