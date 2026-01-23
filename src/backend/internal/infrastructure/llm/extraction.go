package llm

import (
	"context"

	"backend/internal/domain"
)

const defaultExtractionPrompt = `Please extract all text from this document image.
Return the complete text content exactly as it appears in the document.
Preserve the original formatting, paragraphs, and structure as much as possible.
If the document contains handwritten text, transcribe it as accurately as possible.
Do not add any commentary or explanation - just return the extracted text.`

// DocumentExtractorConfig holds configuration for the document extractor.
type DocumentExtractorConfig struct {
	// DefaultModel is the model to use for extraction if not specified per-request.
	// If empty, the provider's default model is used.
	DefaultModel string

	// MaxTokens for extraction responses. Defaults to 8192.
	MaxTokens int
}

// ExtractionRequest represents a request to extract text from a document.
type ExtractionRequest struct { //nolint:govet // Field order prioritizes readability
	// Document is the raw document data (image or PDF bytes).
	Document []byte

	// MediaType indicates the document format.
	MediaType domain.ImageMediaType

	// CustomPrompt overrides the default extraction prompt.
	CustomPrompt string

	// Model specifies which model to use. If empty, uses config default.
	Model string
}

// ExtractionResult contains the extracted text and metadata.
type ExtractionResult struct {
	// Text is the extracted text content.
	Text string

	// InputTokens used for the extraction.
	InputTokens int

	// OutputTokens generated.
	OutputTokens int
}

// DocumentExtractor extracts text from documents using LLM vision capabilities.
type DocumentExtractor struct {
	provider domain.LLMProvider
	config   DocumentExtractorConfig
}

// NewDocumentExtractor creates a new document extractor.
func NewDocumentExtractor(provider domain.LLMProvider, config DocumentExtractorConfig) *DocumentExtractor {
	if config.MaxTokens == 0 {
		config.MaxTokens = 8192
	}
	return &DocumentExtractor{
		provider: provider,
		config:   config,
	}
}

// ExtractText extracts text from a document image or PDF.
func (e *DocumentExtractor) ExtractText(ctx context.Context, req ExtractionRequest) (*ExtractionResult, error) {
	// Determine prompt
	prompt := req.CustomPrompt
	if prompt == "" {
		prompt = defaultExtractionPrompt
	}

	// Determine model
	model := req.Model
	if model == "" {
		model = e.config.DefaultModel
	}

	// Build LLM request with image
	llmReq := domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewImageMessage(domain.RoleUser, req.MediaType, req.Document, prompt),
		},
		Model:     model,
		MaxTokens: e.config.MaxTokens,
	}

	// Execute extraction
	resp, err := e.provider.Complete(ctx, llmReq)
	if err != nil {
		return nil, err
	}

	return &ExtractionResult{
		Text:         resp.Content,
		InputTokens:  resp.InputTokens,
		OutputTokens: resp.OutputTokens,
	}, nil
}
