// Package domain contains the core business entities and repository interfaces.
package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

// ExperienceSource represents where an experience entry came from.
type ExperienceSource string

// Experience source constants.
const (
	ExperienceSourceManual          ExperienceSource = "manual"
	ExperienceSourceResumeExtracted ExperienceSource = "resume_extracted"
)

// Profile represents a user's profile containing manually editable data.
type Profile struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:profiles,alias:p"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID    uuid.UUID `bun:"user_id,notnull,type:uuid"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`

	// Header fields (user-editable, override resume extraction if set)
	Name     *string `bun:"name"`
	Email    *string `bun:"email"`
	Phone    *string `bun:"phone"`
	Location *string `bun:"location"`
	Summary  *string `bun:"summary"`

	// Profile photo
	ProfilePhotoFileID *uuid.UUID `bun:"profile_photo_file_id,type:uuid"`

	// Relations
	ProfilePhotoFile *File `bun:"rel:belongs-to,join:profile_photo_file_id=id"`
	User        *User                `bun:"rel:belongs-to,join:user_id=id"`
	Experiences []*ProfileExperience `bun:"rel:has-many,join:id=profile_id"`
	Educations  []*ProfileEducation  `bun:"rel:has-many,join:id=profile_id"`
	Skills      []*ProfileSkill      `bun:"rel:has-many,join:id=profile_id"`
}

// ProfileExperience represents a work experience entry in a user's profile.
type ProfileExperience struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:profile_experiences,alias:pe"`

	ID             uuid.UUID        `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	ProfileID      uuid.UUID        `bun:"profile_id,notnull,type:uuid"`
	Company        string           `bun:"company,notnull"`
	Title          string           `bun:"title,notnull"`
	Location       *string          `bun:"location"`
	StartDate      *string          `bun:"start_date"`
	EndDate        *string          `bun:"end_date"`
	IsCurrent      bool             `bun:"is_current,notnull,default:false"`
	Description    *string          `bun:"description"`
	Highlights     pq.StringArray   `bun:"highlights,type:text[],array"`
	DisplayOrder   int              `bun:"display_order,notnull,default:0"`
	Source         ExperienceSource `bun:"source,notnull,default:'manual'"`
	SourceResumeID *uuid.UUID       `bun:"source_resume_id,type:uuid"`
	OriginalData   json.RawMessage  `bun:"original_data,type:jsonb"`
	CreatedAt      time.Time        `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time        `bun:"updated_at,notnull,default:current_timestamp"`

	// Relations
	Profile      *Profile `bun:"rel:belongs-to,join:profile_id=id"`
	SourceResume *Resume  `bun:"rel:belongs-to,join:source_resume_id=id"`
}

// ProfileEducation represents an education entry in a user's profile.
type ProfileEducation struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:profile_education,alias:ped"`

	ID             uuid.UUID        `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	ProfileID      uuid.UUID        `bun:"profile_id,notnull,type:uuid"`
	Institution    string           `bun:"institution,notnull"`
	Degree         string           `bun:"degree,notnull"`
	Field          *string          `bun:"field"`
	StartDate      *string          `bun:"start_date"`
	EndDate        *string          `bun:"end_date"`
	IsCurrent      bool             `bun:"is_current,notnull,default:false"`
	Description    *string          `bun:"description"`
	GPA            *string          `bun:"gpa"`
	DisplayOrder   int              `bun:"display_order,notnull,default:0"`
	Source         ExperienceSource `bun:"source,notnull,default:'manual'"`
	SourceResumeID *uuid.UUID       `bun:"source_resume_id,type:uuid"`
	OriginalData   json.RawMessage  `bun:"original_data,type:jsonb"`
	CreatedAt      time.Time        `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time        `bun:"updated_at,notnull,default:current_timestamp"`

	// Relations
	Profile      *Profile `bun:"rel:belongs-to,join:profile_id=id"`
	SourceResume *Resume  `bun:"rel:belongs-to,join:source_resume_id=id"`
}

// ProfileRepository defines operations for profile persistence.
type ProfileRepository interface {
	// Create persists a new profile.
	Create(ctx context.Context, profile *Profile) error

	// GetByID retrieves a profile by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Profile, error)

	// GetByUserID retrieves a profile by user ID.
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error)

	// GetOrCreateByUserID retrieves a profile by user ID, creating one if it doesn't exist.
	GetOrCreateByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error)

	// Update persists changes to an existing profile.
	Update(ctx context.Context, profile *Profile) error

	// Delete removes a profile by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProfileExperienceRepository defines operations for profile experience persistence.
type ProfileExperienceRepository interface {
	// Create persists a new profile experience.
	Create(ctx context.Context, experience *ProfileExperience) error

	// GetByID retrieves a profile experience by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*ProfileExperience, error)

	// GetByProfileID retrieves all profile experiences for a profile, ordered by display order.
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*ProfileExperience, error)

	// Update persists changes to an existing profile experience.
	Update(ctx context.Context, experience *ProfileExperience) error

	// Delete removes a profile experience by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetNextDisplayOrder returns the next display order value for a profile.
	GetNextDisplayOrder(ctx context.Context, profileID uuid.UUID) (int, error)

	// DeleteBySourceResumeID removes all experiences extracted from a specific resume.
	DeleteBySourceResumeID(ctx context.Context, sourceResumeID uuid.UUID) error
}

// ProfileEducationRepository defines operations for profile education persistence.
type ProfileEducationRepository interface {
	// Create persists a new profile education entry.
	Create(ctx context.Context, education *ProfileEducation) error

	// GetByID retrieves a profile education entry by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*ProfileEducation, error)

	// GetByProfileID retrieves all profile education entries for a profile, ordered by display order.
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*ProfileEducation, error)

	// Update persists changes to an existing profile education entry.
	Update(ctx context.Context, education *ProfileEducation) error

	// Delete removes a profile education entry by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetNextDisplayOrder returns the next display order value for a profile.
	GetNextDisplayOrder(ctx context.Context, profileID uuid.UUID) (int, error)

	// DeleteBySourceResumeID removes all education entries extracted from a specific resume.
	DeleteBySourceResumeID(ctx context.Context, sourceResumeID uuid.UUID) error
}

// ProfileSkill represents a skill entry in a user's profile.
type ProfileSkill struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:profile_skills,alias:ps"`

	ID             uuid.UUID        `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	ProfileID      uuid.UUID        `bun:"profile_id,notnull,type:uuid"`
	Name           string           `bun:"name,notnull"`
	NormalizedName string           `bun:"normalized_name,notnull"`
	Category       string           `bun:"category,notnull,default:'TECHNICAL'"`
	DisplayOrder   int              `bun:"display_order,notnull,default:0"`
	Source         ExperienceSource `bun:"source,notnull,default:'manual'"`
	SourceResumeID *uuid.UUID       `bun:"source_resume_id,type:uuid"`
	OriginalData   json.RawMessage  `bun:"original_data,type:jsonb"`
	CreatedAt      time.Time        `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time        `bun:"updated_at,notnull,default:current_timestamp"`

	// Relations
	Profile      *Profile `bun:"rel:belongs-to,join:profile_id=id"`
	SourceResume *Resume  `bun:"rel:belongs-to,join:source_resume_id=id"`
}

// ProfileSkillRepository defines operations for profile skill persistence.
type ProfileSkillRepository interface {
	// Create persists a new profile skill.
	Create(ctx context.Context, skill *ProfileSkill) error

	// CreateIgnoreDuplicate persists a new profile skill, silently ignoring duplicates.
	// This uses ON CONFLICT DO NOTHING to handle unique constraint violations on (profile_id, normalized_name).
	CreateIgnoreDuplicate(ctx context.Context, skill *ProfileSkill) error

	// GetByID retrieves a profile skill by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*ProfileSkill, error)

	// GetByProfileID retrieves all profile skills for a profile, ordered by display order.
	GetByProfileID(ctx context.Context, profileID uuid.UUID) ([]*ProfileSkill, error)

	// Update persists changes to an existing profile skill.
	Update(ctx context.Context, skill *ProfileSkill) error

	// Delete removes a profile skill by its ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetNextDisplayOrder returns the next display order value for a profile.
	GetNextDisplayOrder(ctx context.Context, profileID uuid.UUID) (int, error)

	// DeleteBySourceResumeID removes all skills extracted from a specific resume.
	DeleteBySourceResumeID(ctx context.Context, sourceResumeID uuid.UUID) error
}
