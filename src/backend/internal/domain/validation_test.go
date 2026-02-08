package domain_test

import (
	"errors"
	"testing"

	"backend/internal/domain"
)

func TestValidationError(t *testing.T) {
	t.Run("formats error message with field and message", func(t *testing.T) {
		err := &domain.ValidationError{
			Field:   "email",
			Message: "invalid format",
			Err:     domain.ErrInvalidCharacter,
		}

		expected := "validation error in email: invalid format"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("unwraps to underlying error", func(t *testing.T) {
		err := &domain.ValidationError{
			Field:   "name",
			Message: "too long",
			Err:     domain.ErrFieldTooLong,
		}

		if !errors.Is(err, domain.ErrFieldTooLong) {
			t.Error("expected errors.Is to match underlying error")
		}
	})
}

func TestValidationConstants(t *testing.T) {
	t.Run("sentinel errors exist", func(t *testing.T) {
		if domain.ErrFieldTooLong == nil {
			t.Error("ErrFieldTooLong should be defined")
		}
		if domain.ErrInvalidCharacter == nil {
			t.Error("ErrInvalidCharacter should be defined")
		}
		if domain.ErrEmptyRequired == nil {
			t.Error("ErrEmptyRequired should be defined")
		}
	})
}
