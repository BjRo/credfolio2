package domain

import (
	"errors"
	"fmt"
)

// Validation errors - sentinel errors for different validation failures.
var (
	ErrFieldTooLong     = errors.New("field exceeds maximum length")
	ErrInvalidCharacter = errors.New("field contains invalid characters")
	ErrEmptyRequired    = errors.New("required field is empty")
)

// ValidationError wraps field-specific validation failures.
type ValidationError struct { //nolint:govet // Field order prioritizes readability
	Field   string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// Unwrap returns the underlying error.
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ExtractedDataValidator validates and sanitizes extracted LLM data.
type ExtractedDataValidator interface {
	ValidateResumeData(data *ResumeExtractedData) error
	ValidateLetterData(data *ExtractedLetterData) error
}
