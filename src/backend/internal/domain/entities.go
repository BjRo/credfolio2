// Package domain contains the core business entities and repository interfaces.
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// User represents a user account in the system.
type User struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Email        string    `bun:"email,notnull,unique"`
	PasswordHash string    `bun:"password_hash,notnull"`
	Name         *string   `bun:"name"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

// File represents an uploaded file stored in object storage.
type File struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:files,alias:f"`

	ID          uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID      uuid.UUID `bun:"user_id,notnull,type:uuid"`
	Filename    string    `bun:"filename,notnull"`
	ContentType string    `bun:"content_type,notnull"`
	SizeBytes   int64     `bun:"size_bytes,notnull"`
	StorageKey  string    `bun:"storage_key,notnull,unique"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

// ReferenceLetterStatus represents the processing status of a reference letter.
type ReferenceLetterStatus string

// Reference letter status constants.
const (
	ReferenceLetterStatusPending    ReferenceLetterStatus = "pending"
	ReferenceLetterStatusProcessing ReferenceLetterStatus = "processing"
	ReferenceLetterStatusCompleted  ReferenceLetterStatus = "completed"
	ReferenceLetterStatusFailed     ReferenceLetterStatus = "failed"
)

// ReferenceLetter represents a reference letter with extracted data.
type ReferenceLetter struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	bun.BaseModel `bun:"table:reference_letters,alias:rl"`

	ID            uuid.UUID             `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID             `bun:"user_id,notnull,type:uuid"`
	FileID        *uuid.UUID            `bun:"file_id,type:uuid"`
	Title         *string               `bun:"title"`
	AuthorName    *string               `bun:"author_name"`
	AuthorTitle   *string               `bun:"author_title"`
	Organization  *string               `bun:"organization"`
	DateWritten   *time.Time            `bun:"date_written,type:date"`
	RawText       *string               `bun:"raw_text"`
	ExtractedData json.RawMessage       `bun:"extracted_data,type:jsonb"`
	Status        ReferenceLetterStatus `bun:"status,notnull,default:'pending'"`
	CreatedAt     time.Time             `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt     time.Time             `bun:"updated_at,notnull,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
	File *File `bun:"rel:belongs-to,join:file_id=id"`
}
