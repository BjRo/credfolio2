//nolint:errcheck,revive // Test file - error checks and unused params are OK in test helpers
package llm_test

import (
	"context"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"

	"go.opentelemetry.io/otel/codes"
)

func TestDocumentExtractor_DetectDocumentContent_Resume(t *testing.T) {
	jsonResponse := `{
		"hasCareerInfo": true,
		"hasTestimonial": false,
		"testimonialAuthor": "",
		"confidence": 0.95,
		"summary": "A professional resume listing work experience, education, and technical skills.",
		"documentTypeHint": "resume"
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      jsonResponse,
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  500,
			OutputTokens: 100,
			StopReason:   "end_turn",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.DetectDocumentContent(context.Background(), "Resume text with experience and skills...")

	if err != nil {
		t.Fatalf("DetectDocumentContent() error = %v", err)
	}

	if !result.HasCareerInfo {
		t.Error("HasCareerInfo = false, want true")
	}
	if result.HasTestimonial {
		t.Error("HasTestimonial = true, want false")
	}
	if result.TestimonialAuthor != nil {
		t.Errorf("TestimonialAuthor = %v, want nil", result.TestimonialAuthor)
	}
	if result.Confidence != 0.95 {
		t.Errorf("Confidence = %v, want 0.95", result.Confidence)
	}
	if result.Summary != "A professional resume listing work experience, education, and technical skills." {
		t.Errorf("Summary = %q, unexpected", result.Summary)
	}
	if result.DocumentTypeHint != domain.DocumentTypeResume {
		t.Errorf("DocumentTypeHint = %q, want %q", result.DocumentTypeHint, domain.DocumentTypeResume)
	}
}

func TestDocumentExtractor_DetectDocumentContent_ReferenceLetter(t *testing.T) {
	jsonResponse := `{
		"hasCareerInfo": false,
		"hasTestimonial": true,
		"testimonialAuthor": "Jane Smith",
		"confidence": 0.92,
		"summary": "A reference letter from Jane Smith recommending the candidate for their leadership and technical skills.",
		"documentTypeHint": "reference_letter"
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      jsonResponse,
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  600,
			OutputTokens: 120,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.DetectDocumentContent(context.Background(), "To Whom It May Concern, I am writing to recommend...")

	if err != nil {
		t.Fatalf("DetectDocumentContent() error = %v", err)
	}

	if result.HasCareerInfo {
		t.Error("HasCareerInfo = true, want false")
	}
	if !result.HasTestimonial {
		t.Error("HasTestimonial = false, want true")
	}
	if result.TestimonialAuthor == nil || *result.TestimonialAuthor != "Jane Smith" {
		t.Errorf("TestimonialAuthor = %v, want %q", result.TestimonialAuthor, "Jane Smith")
	}
	if result.DocumentTypeHint != domain.DocumentTypeReferenceLetter {
		t.Errorf("DocumentTypeHint = %q, want %q", result.DocumentTypeHint, domain.DocumentTypeReferenceLetter)
	}
}

func TestDocumentExtractor_DetectDocumentContent_Hybrid(t *testing.T) {
	jsonResponse := `{
		"hasCareerInfo": true,
		"hasTestimonial": true,
		"testimonialAuthor": "John Manager",
		"confidence": 0.85,
		"summary": "A reference letter that also contains detailed career history of the candidate.",
		"documentTypeHint": "hybrid"
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: jsonResponse,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.DetectDocumentContent(context.Background(), "Letter with career details...")

	if err != nil {
		t.Fatalf("DetectDocumentContent() error = %v", err)
	}

	if !result.HasCareerInfo {
		t.Error("HasCareerInfo = false, want true")
	}
	if !result.HasTestimonial {
		t.Error("HasTestimonial = false, want true")
	}
	if result.TestimonialAuthor == nil || *result.TestimonialAuthor != "John Manager" {
		t.Errorf("TestimonialAuthor = %v, want %q", result.TestimonialAuthor, "John Manager")
	}
	if result.DocumentTypeHint != domain.DocumentTypeHybrid {
		t.Errorf("DocumentTypeHint = %q, want %q", result.DocumentTypeHint, domain.DocumentTypeHybrid)
	}
}

func TestDocumentExtractor_DetectDocumentContent_Unknown(t *testing.T) {
	jsonResponse := `{
		"hasCareerInfo": false,
		"hasTestimonial": false,
		"testimonialAuthor": "",
		"confidence": 0.3,
		"summary": "The document appears to be a grocery shopping list.",
		"documentTypeHint": "unknown"
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: jsonResponse,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.DetectDocumentContent(context.Background(), "Milk, eggs, bread...")

	if err != nil {
		t.Fatalf("DetectDocumentContent() error = %v", err)
	}

	if result.HasCareerInfo {
		t.Error("HasCareerInfo = true, want false")
	}
	if result.HasTestimonial {
		t.Error("HasTestimonial = true, want false")
	}
	if result.TestimonialAuthor != nil {
		t.Errorf("TestimonialAuthor = %v, want nil", result.TestimonialAuthor)
	}
	if result.Confidence != 0.3 {
		t.Errorf("Confidence = %v, want 0.3", result.Confidence)
	}
	if result.DocumentTypeHint != domain.DocumentTypeUnknown {
		t.Errorf("DocumentTypeHint = %q, want %q", result.DocumentTypeHint, domain.DocumentTypeUnknown)
	}
}

func TestDocumentExtractor_DetectDocumentContent_MarkdownCodeBlock(t *testing.T) {
	jsonResponse := "```json\n" + `{
		"hasCareerInfo": true,
		"hasTestimonial": false,
		"testimonialAuthor": "",
		"confidence": 0.9,
		"summary": "A resume.",
		"documentTypeHint": "resume"
	}` + "\n```"

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: jsonResponse,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.DetectDocumentContent(context.Background(), "Resume text")

	if err != nil {
		t.Fatalf("DetectDocumentContent() error = %v", err)
	}

	if !result.HasCareerInfo {
		t.Error("HasCareerInfo = false, want true")
	}
	if result.DocumentTypeHint != domain.DocumentTypeResume {
		t.Errorf("DocumentTypeHint = %q, want %q", result.DocumentTypeHint, domain.DocumentTypeResume)
	}
}

func TestDocumentExtractor_DetectDocumentContent_Error(t *testing.T) {
	inner := &mockProvider{
		err: &domain.LLMError{
			Provider: "mock",
			Message:  "API error",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.DetectDocumentContent(context.Background(), "Some text")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDocumentExtractor_DetectDocumentContent_InvalidJSON(t *testing.T) {
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: "not valid json at all",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.DetectDocumentContent(context.Background(), "Some text")

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDocumentExtractor_DetectDocumentContent_CreatesSpan(t *testing.T) {
	exporter := setupTestTracing(t)

	jsonResponse := `{
		"hasCareerInfo": true,
		"hasTestimonial": false,
		"testimonialAuthor": "",
		"confidence": 0.9,
		"summary": "A resume.",
		"documentTypeHint": "resume"
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: jsonResponse,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.DetectDocumentContent(context.Background(), "Resume text here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	spans := exporter.GetSpans()
	var found bool
	for _, s := range spans {
		if s.Name == "document_content_detection" {
			found = true
			// Check for text_length attribute
			attrFound := false
			for _, attr := range s.Attributes {
				if string(attr.Key) == "text_length" && attr.Value.AsInt64() == int64(len("Resume text here")) {
					attrFound = true
				}
			}
			if !attrFound {
				t.Error("expected text_length attribute on span")
			}
			break
		}
	}
	if !found {
		t.Errorf("expected span named 'document_content_detection', got spans: %v", spanNames(spans))
	}
}

func TestDocumentExtractor_DetectDocumentContent_SpanRecordsErrorOnFailure(t *testing.T) {
	exporter := setupTestTracing(t)

	inner := &mockProvider{
		err: &domain.LLMError{
			Provider: "mock",
			Message:  "API error",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.DetectDocumentContent(context.Background(), "Some text")
	if err == nil {
		t.Fatal("expected error")
	}

	spans := exporter.GetSpans()
	var found bool
	for _, s := range spans {
		if s.Name == "document_content_detection" {
			found = true
			if s.Status.Code != codes.Error {
				t.Errorf("expected span status Error, got %v", s.Status.Code)
			}
			break
		}
	}
	if !found {
		t.Error("expected span even on error")
	}
}

func TestDocumentExtractor_DetectDocumentContent_UsesDetectionUserPrompt(t *testing.T) {
	var capturedReq domain.LLMRequest
	inner := &capturingProvider{
		response: &domain.LLMResponse{
			Content: `{
				"hasCareerInfo": true,
				"hasTestimonial": false,
				"testimonialAuthor": "",
				"confidence": 0.9,
				"summary": "A resume.",
				"documentTypeHint": "resume"
			}`,
		},
		captureReq: &capturedReq,
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.DetectDocumentContent(context.Background(), "My resume text goes here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify a system prompt was set
	if capturedReq.SystemPrompt == "" {
		t.Error("expected system prompt to be set")
	}

	// Verify user message contains the document text
	if len(capturedReq.Messages) == 0 {
		t.Fatal("expected at least one message")
	}
	msg := capturedReq.Messages[0]
	found := false
	for _, block := range msg.Content {
		if block.Type == domain.ContentTypeText && block.Text != "" {
			// The user prompt should contain the original text
			if contains(block.Text, "My resume text goes here") {
				found = true
			}
		}
	}
	if !found {
		t.Error("document text not found in user message")
	}

	// Verify structured output schema was set
	if capturedReq.OutputSchema == nil {
		t.Error("expected OutputSchema to be set for structured output")
	}
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
