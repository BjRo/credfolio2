package llm

import (
	"context"
	"encoding/json"
	"fmt"

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

// ExtractTextWithRequest extracts text from a document image or PDF using a detailed request.
func (e *DocumentExtractor) ExtractTextWithRequest(ctx context.Context, req ExtractionRequest) (*ExtractionResult, error) {
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

// ExtractText implements domain.DocumentExtractor interface.
// It extracts raw text from a document using LLM vision capabilities.
func (e *DocumentExtractor) ExtractText(ctx context.Context, document []byte, contentType string) (string, error) {
	// Map content type to media type
	mediaType, err := contentTypeToMediaType(contentType)
	if err != nil {
		return "", err
	}

	result, err := e.ExtractTextWithRequest(ctx, ExtractionRequest{
		Document:  document,
		MediaType: mediaType,
	})
	if err != nil {
		return "", err
	}

	return result.Text, nil
}

// contentTypeToMediaType converts a content type string to domain.ImageMediaType.
func contentTypeToMediaType(contentType string) (domain.ImageMediaType, error) {
	switch contentType {
	case "application/pdf":
		return domain.ImageMediaTypePDF, nil
	case "image/jpeg":
		return domain.ImageMediaTypeJPEG, nil
	case "image/png":
		return domain.ImageMediaTypePNG, nil
	case "image/gif":
		return domain.ImageMediaTypeGIF, nil
	case "image/webp":
		return domain.ImageMediaTypeWebP, nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}

const resumeExtractionPrompt = `Extract structured profile data from the following resume text.

Return a JSON object with the following structure:
{
  "name": "Full name of the candidate",
  "email": "Email address (if found)",
  "phone": "Phone number (if found)",
  "location": "City, State/Country (if found)",
  "summary": "Professional summary or objective (if found)",
  "experience": [
    {
      "company": "Company name",
      "title": "Job title",
      "location": "Job location (if found)",
      "startDate": "Start date (e.g., 'Jan 2020')",
      "endDate": "End date or 'Present'",
      "isCurrent": true/false,
      "description": "Job description or responsibilities"
    }
  ],
  "education": [
    {
      "institution": "School/University name",
      "degree": "Degree type (e.g., 'Bachelor of Science')",
      "field": "Field of study",
      "startDate": "Start date (if found)",
      "endDate": "End/graduation date",
      "gpa": "GPA (if mentioned)",
      "achievements": "Notable achievements or honors"
    }
  ],
  "skills": ["Skill 1", "Skill 2", "..."],
  "confidence": 0.0 to 1.0 (your confidence in the accuracy of extraction)
}

Rules:
- Extract all information present; use null for missing fields
- For dates, use the format found in the document (e.g., "Jan 2020", "2020", "January 2020")
- Set isCurrent to true if the job end date is "Present" or similar
- Skills should be a flat array of individual skills
- Confidence should reflect how clear and complete the resume text was

Resume text:
`

// ExtractResumeData implements domain.DocumentExtractor interface.
// It extracts structured resume data from text using LLM.
func (e *DocumentExtractor) ExtractResumeData(ctx context.Context, text string) (*domain.ResumeExtractedData, error) {
	llmReq := domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, resumeExtractionPrompt+text),
		},
		Model:     e.config.DefaultModel,
		MaxTokens: e.config.MaxTokens,
	}

	resp, err := e.provider.Complete(ctx, llmReq)
	if err != nil {
		return nil, fmt.Errorf("LLM extraction failed: %w", err)
	}

	// Parse JSON response
	var data domain.ResumeExtractedData
	if err := json.Unmarshal([]byte(resp.Content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse extraction response: %w", err)
	}

	// Ensure slices are initialized
	if data.Experience == nil {
		data.Experience = []domain.WorkExperience{}
	}
	if data.Education == nil {
		data.Education = []domain.Education{}
	}
	if data.Skills == nil {
		data.Skills = []string{}
	}

	return &data, nil
}

// Verify DocumentExtractor implements domain.DocumentExtractor interface.
var _ domain.DocumentExtractor = (*DocumentExtractor)(nil)
