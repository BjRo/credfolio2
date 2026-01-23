//nolint:errcheck,revive // Test file - error checks and unused params are OK in test helpers
package llm_test

import (
	"context"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
)

func TestDocumentExtractor_ExtractText(t *testing.T) {
	// Mock provider that returns extracted text
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "This is the extracted text from the document. It contains important information about the candidate.",
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  500,
			OutputTokens: 50,
			StopReason:   "end_turn",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	// Fake JPEG image data
	imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes

	result, err := extractor.ExtractText(context.Background(), llm.ExtractionRequest{
		Document:  imageData,
		MediaType: domain.ImageMediaTypeJPEG,
	})

	if err != nil {
		t.Fatalf("ExtractText() error = %v", err)
	}

	if result.Text != inner.response.Content {
		t.Errorf("Text = %q, want %q", result.Text, inner.response.Content)
	}
	if result.InputTokens != 500 {
		t.Errorf("InputTokens = %d, want 500", result.InputTokens)
	}
	if result.OutputTokens != 50 {
		t.Errorf("OutputTokens = %d, want 50", result.OutputTokens)
	}
}

func TestDocumentExtractor_ExtractText_PDF(t *testing.T) {
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "PDF content extracted",
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  1000,
			OutputTokens: 100,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	// Fake PDF data
	pdfData := []byte("%PDF-1.4") // PDF magic bytes

	result, err := extractor.ExtractText(context.Background(), llm.ExtractionRequest{
		Document:  pdfData,
		MediaType: domain.ImageMediaTypePDF,
	})

	if err != nil {
		t.Fatalf("ExtractText() error = %v", err)
	}

	if result.Text != "PDF content extracted" {
		t.Errorf("Text = %q, want %q", result.Text, "PDF content extracted")
	}
}

func TestDocumentExtractor_ExtractText_CustomPrompt(t *testing.T) {
	var capturedReq domain.LLMRequest
	inner := &capturingProvider{
		response: &domain.LLMResponse{
			Content: "Extracted with custom prompt",
		},
		captureReq: &capturedReq,
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	customPrompt := "Extract the author name and date from this letter."
	_, err := extractor.ExtractText(context.Background(), llm.ExtractionRequest{
		Document:     []byte{0xFF, 0xD8, 0xFF},
		MediaType:    domain.ImageMediaTypeJPEG,
		CustomPrompt: customPrompt,
	})

	if err != nil {
		t.Fatalf("ExtractText() error = %v", err)
	}

	// Verify custom prompt was used
	if len(capturedReq.Messages) == 0 {
		t.Fatal("expected messages in request")
	}
	msg := capturedReq.Messages[0]
	found := false
	for _, block := range msg.Content {
		if block.Type == domain.ContentTypeText && block.Text == customPrompt {
			found = true
			break
		}
	}
	if !found {
		t.Error("custom prompt not found in request")
	}
}

func TestDocumentExtractor_ExtractText_CustomModel(t *testing.T) {
	var capturedReq domain.LLMRequest
	inner := &capturingProvider{
		response: &domain.LLMResponse{
			Content: "Result",
		},
		captureReq: &capturedReq,
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.ExtractText(context.Background(), llm.ExtractionRequest{
		Document:  []byte{0xFF, 0xD8, 0xFF},
		MediaType: domain.ImageMediaTypeJPEG,
		Model:     "claude-3-haiku-20240307",
	})

	if err != nil {
		t.Fatalf("ExtractText() error = %v", err)
	}

	if capturedReq.Model != "claude-3-haiku-20240307" {
		t.Errorf("Model = %q, want %q", capturedReq.Model, "claude-3-haiku-20240307")
	}
}

func TestDocumentExtractor_ExtractText_Error(t *testing.T) {
	inner := &mockProvider{
		err: &domain.LLMError{
			Provider: "mock",
			Message:  "API error",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.ExtractText(context.Background(), llm.ExtractionRequest{
		Document:  []byte{0xFF, 0xD8, 0xFF},
		MediaType: domain.ImageMediaTypeJPEG,
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

// capturingProvider captures the request for inspection.
type capturingProvider struct {
	response   *domain.LLMResponse
	captureReq *domain.LLMRequest
}

func (p *capturingProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	*p.captureReq = req
	return p.response, nil
}

func (p *capturingProvider) Name() string {
	return "capturing"
}
