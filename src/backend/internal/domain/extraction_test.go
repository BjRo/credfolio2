package domain_test

import (
	"testing"
	"time"

	"backend/internal/domain"
)

// TestPromptVersionConstants verifies prompt version constants exist and follow semantic versioning.
func TestPromptVersionConstants(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"ResumeExtractionPromptVersion", domain.ResumeExtractionPromptVersion},
		{"LetterExtractionPromptVersion", domain.LetterExtractionPromptVersion},
		{"DocumentDetectionPromptVersion", domain.DocumentDetectionPromptVersion},
		{"DocumentExtractionPromptVersion", domain.DocumentExtractionPromptVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.version == "" {
				t.Errorf("%s is empty, want semantic version like 'v1.0.0'", tt.name)
			}
			// Verify it starts with 'v' and has at least major.minor.patch format
			if len(tt.version) < 5 || tt.version[0] != 'v' {
				t.Errorf("%s = %q, want format 'vX.Y.Z'", tt.name, tt.version)
			}
		})
	}
}

// TestExtractionMetadata_EnhancedFields verifies ExtractionMetadata contains new observability fields.
func TestExtractionMetadata_EnhancedFields(t *testing.T) {
	now := time.Now()
	metadata := domain.ExtractionMetadata{
		ExtractedAt:   now,
		ModelVersion:  "claude-sonnet-4-5-20250929",
		PromptVersion: "v1.1.0",
		InputTokens:   1000,
		OutputTokens:  500,
		DurationMs:    1234,
	}

	// Verify all fields are accessible
	if metadata.ExtractedAt != now {
		t.Errorf("ExtractedAt = %v, want %v", metadata.ExtractedAt, now)
	}
	if metadata.ModelVersion != "claude-sonnet-4-5-20250929" {
		t.Errorf("ModelVersion = %q, want %q", metadata.ModelVersion, "claude-sonnet-4-5-20250929")
	}
	if metadata.PromptVersion != "v1.1.0" {
		t.Errorf("PromptVersion = %q, want %q", metadata.PromptVersion, "v1.1.0")
	}
	if metadata.InputTokens != 1000 {
		t.Errorf("InputTokens = %d, want %d", metadata.InputTokens, 1000)
	}
	if metadata.OutputTokens != 500 {
		t.Errorf("OutputTokens = %d, want %d", metadata.OutputTokens, 500)
	}
	if metadata.DurationMs != 1234 {
		t.Errorf("DurationMs = %d, want %d", metadata.DurationMs, 1234)
	}
}
