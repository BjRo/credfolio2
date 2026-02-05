package domain

import (
	"time"

	"github.com/google/uuid"
)

// DocumentFeedbackType represents the kind of feedback being reported.
type DocumentFeedbackType string

// Feedback type constants.
const (
	DocumentFeedbackDetectionCorrection DocumentFeedbackType = "detection_correction"
	DocumentFeedbackExtractionQuality   DocumentFeedbackType = "extraction_quality"
)

// DocumentFeedback captures user feedback about document detection or extraction quality.
// For MVP this is logged only (no dedicated table); a future iteration may persist these.
type DocumentFeedback struct { //nolint:govet // Field ordering prioritizes readability
	UserID       uuid.UUID            `json:"userId"`
	FileID       uuid.UUID            `json:"fileId"`
	FeedbackType DocumentFeedbackType `json:"feedbackType"`
	Message      string               `json:"message"`
	CreatedAt    time.Time            `json:"createdAt"`
}
