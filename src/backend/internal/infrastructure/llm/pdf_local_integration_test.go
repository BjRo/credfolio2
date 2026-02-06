//nolint:errcheck,revive,goconst // Test file - error checks, unused params, and string constants are OK
package llm_test

import (
	"context"
	"os"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestExtractTextWithRequest_PDF_UsesLocalExtraction(t *testing.T) {
	// Load a real text-based PDF
	data, err := os.ReadFile("../../../../../fixtures/CV_TEMPLATE_0004.pdf")
	if err != nil {
		t.Skipf("Fixture PDF not available: %v", err)
	}

	exporter := setupTestTracingHelper(t)

	// The mock provider should NOT be called for a text-based PDF
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "LLM should not be called",
			InputTokens:  9999,
			OutputTokens: 9999,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  data,
		MediaType: domain.ImageMediaTypePDF,
	})
	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	// Should have extracted text locally (not from LLM)
	if result.Text == "LLM should not be called" {
		t.Fatal("Expected local extraction, but LLM was called")
	}
	if result.Text == "" {
		t.Fatal("Expected non-empty text from local extraction")
	}

	// Token counts should be zero (no LLM call)
	if result.InputTokens != 0 {
		t.Errorf("InputTokens = %d, want 0 (no LLM call)", result.InputTokens)
	}
	if result.OutputTokens != 0 {
		t.Errorf("OutputTokens = %d, want 0 (no LLM call)", result.OutputTokens)
	}

	// Check span attributes
	spans := exporter.GetSpans()
	for _, s := range spans {
		if s.Name == "pdf_text_extraction" {
			assertSpanAttribute(t, s, "extraction_method", "local")
			return
		}
	}
	t.Error("expected pdf_text_extraction span")
}

func TestExtractTextWithRequest_PDF_FallsBackToLLM_ForScannedPDF(t *testing.T) {
	exporter := setupTestTracingHelper(t)

	// Simulate a "scanned" PDF by providing data that the Go parser can't extract text from.
	// A minimal valid PDF structure with no text content.
	fakePDF := []byte("%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n2 0 obj<</Type/Pages/Kids[]/Count 0>>endobj\nxref\n0 3\n0000000000 65535 f \n0000000009 00000 n \n0000000052 00000 n \ntrailer<</Root 1 0 R/Size 3>>\nstartxref\n101\n%%EOF")

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "LLM extracted text from scanned PDF",
			InputTokens:  500,
			OutputTokens: 50,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  fakePDF,
		MediaType: domain.ImageMediaTypePDF,
	})
	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	// Should have fallen back to LLM
	if result.Text != "LLM extracted text from scanned PDF" {
		t.Errorf("Text = %q, want LLM response", result.Text)
	}

	// Check span attributes
	spans := exporter.GetSpans()
	for _, s := range spans {
		if s.Name == "pdf_text_extraction" {
			assertSpanAttribute(t, s, "extraction_method", "llm")
			return
		}
	}
	t.Error("expected pdf_text_extraction span")
}

func TestExtractTextWithRequest_PDF_FallsBackToLLM_WithCustomPrompt(t *testing.T) {
	// Even for a text-based PDF, a custom prompt should force LLM extraction
	data, err := os.ReadFile("../../../../../fixtures/CV_TEMPLATE_0004.pdf")
	if err != nil {
		t.Skipf("Fixture PDF not available: %v", err)
	}

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "Custom prompt result from LLM",
			InputTokens:  500,
			OutputTokens: 50,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:     data,
		MediaType:    domain.ImageMediaTypePDF,
		CustomPrompt: "Extract only the candidate's name and email.",
	})
	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	// Should use LLM because of custom prompt
	if result.Text != "Custom prompt result from LLM" {
		t.Errorf("Text = %q, want LLM response", result.Text)
	}
}

func TestExtractTextWithRequest_JPEG_AlwaysUsesLLM(t *testing.T) {
	// Images should always go through LLM
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "Text from image via LLM",
			InputTokens:  500,
			OutputTokens: 50,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  []byte{0xFF, 0xD8, 0xFF, 0xE0},
		MediaType: domain.ImageMediaTypeJPEG,
	})
	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	if result.Text != "Text from image via LLM" {
		t.Errorf("Text = %q, want LLM response", result.Text)
	}
}

// setupTestTracingHelper sets up in-memory tracing for tests.
func setupTestTracingHelper(t *testing.T) *tracetest.InMemoryExporter {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		otel.SetTracerProvider(prev)
		_ = tp.Shutdown(context.Background())
	})
	return exporter
}

// assertSpanAttribute checks that a span has an attribute with the expected value.
func assertSpanAttribute(t *testing.T, span tracetest.SpanStub, key, wantValue string) {
	t.Helper()
	for _, attr := range span.Attributes {
		if string(attr.Key) == key {
			if attr.Value.AsString() != wantValue {
				t.Errorf("span attribute %q = %q, want %q", key, attr.Value.AsString(), wantValue)
			}
			return
		}
	}
	t.Errorf("span missing attribute %q", key)
}
