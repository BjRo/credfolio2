// Package domain contains the core business entities and repository interfaces.
package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// ResumeStatus represents the processing status of a resume.
type ResumeStatus string

// Resume status constants.
const (
	ResumeStatusPending    ResumeStatus = "pending"
	ResumeStatusProcessing ResumeStatus = "processing"
	ResumeStatusCompleted  ResumeStatus = "completed"
	ResumeStatusFailed     ResumeStatus = "failed"
)

// Resume represents an uploaded resume with extracted profile data.
type Resume struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:resumes,alias:r"`

	ID            uuid.UUID       `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID       `bun:"user_id,notnull,type:uuid"`
	FileID        uuid.UUID       `bun:"file_id,notnull,type:uuid"`
	Status        ResumeStatus    `bun:"status,notnull,default:'pending'"`
	ExtractedData json.RawMessage `bun:"extracted_data,type:jsonb"`
	ErrorMessage  *string         `bun:"error_message"`
	CreatedAt     time.Time       `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time       `bun:"updated_at,notnull,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
	File *File `bun:"rel:belongs-to,join:file_id=id"`
}

// WorkExperience represents a single work experience entry from a resume.
type WorkExperience struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Company     string  `json:"company"`
	Title       string  `json:"title"`
	Location    *string `json:"location,omitempty"`
	StartDate   *string `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
	IsCurrent   bool    `json:"isCurrent"`
	Description *string `json:"description,omitempty"`
}

// Education represents a single education entry from a resume.
type Education struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Institution  string  `json:"institution"`
	Degree       *string `json:"degree,omitempty"`
	Field        *string `json:"field,omitempty"`
	StartDate    *string `json:"startDate,omitempty"`
	EndDate      *string `json:"endDate,omitempty"`
	GPA          *string `json:"gpa,omitempty"`
	Achievements *string `json:"achievements,omitempty"`
}

// ResumeExtractedData is the complete extracted data from a resume.
type ResumeExtractedData struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Name        string           `json:"name"`
	Email       *string          `json:"email,omitempty"`
	Phone       *string          `json:"phone,omitempty"`
	Location    *string          `json:"location,omitempty"`
	Summary     *string          `json:"summary,omitempty"`
	Experience  []WorkExperience `json:"experience"`
	Education   []Education      `json:"education"`
	Skills      []string         `json:"skills"`
	ExtractedAt time.Time        `json:"extractedAt"`
	Confidence  float64          `json:"confidence"`
}

// ResumeRepository defines operations for resume persistence.
type ResumeRepository interface {
	// Create persists a new resume.
	Create(ctx context.Context, resume *Resume) error

	// GetByID retrieves a resume by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Resume, error)

	// GetByUserID retrieves all resumes belonging to a user.
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Resume, error)

	// Update persists changes to an existing resume.
	Update(ctx context.Context, resume *Resume) error

	// Delete removes a resume by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// ResumeProcessingRequest contains the data needed to enqueue a resume processing job.
type ResumeProcessingRequest struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	StorageKey  string
	ResumeID    uuid.UUID
	FileID      uuid.UUID
	ContentType string
}
