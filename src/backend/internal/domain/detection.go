// Package domain contains the core business entities and repository interfaces.
package domain

// DocumentTypeHint describes the detected type of a document.
type DocumentTypeHint string

// Document type hint constants.
const (
	DocumentTypeResume          DocumentTypeHint = "resume"
	DocumentTypeReferenceLetter DocumentTypeHint = "reference_letter"
	DocumentTypeHybrid          DocumentTypeHint = "hybrid"
	DocumentTypeUnknown         DocumentTypeHint = "unknown"
)

// DocumentDetectionResult contains the lightweight classification of a document's content.
// This is used for quick detection before running full extraction.
type DocumentDetectionResult struct { //nolint:govet // Field order prioritizes JSON serialization
	// HasCareerInfo indicates whether the document contains resume/CV content.
	HasCareerInfo bool `json:"hasCareerInfo"`

	// HasTestimonial indicates whether the document contains reference letter / testimonial content.
	HasTestimonial bool `json:"hasTestimonial"`

	// TestimonialAuthor is the detected author name of the testimonial, if any.
	TestimonialAuthor *string `json:"testimonialAuthor,omitempty"`

	// Confidence is the overall detection confidence (0.0 to 1.0).
	Confidence float64 `json:"confidence"`

	// Summary is a brief description of what was found in the document.
	Summary string `json:"summary"`

	// DocumentTypeHint classifies the document as "resume", "reference_letter", "hybrid", or "unknown".
	DocumentTypeHint DocumentTypeHint `json:"documentTypeHint"`
}
