package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

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
Extract all information present; use null for missing optional fields.
For dates, use the format found in the document (e.g., "Jan 2020", "2020", "January 2020").
Set isCurrent to true if the job end date is "Present" or similar.
Skills should be a flat array of individual skills.
Confidence should reflect how clear and complete the resume text was (0.0 to 1.0).

Resume text:
`

// resumeOutputSchema defines the JSON schema for structured resume extraction.
// This schema is used with Anthropic's structured output feature to guarantee valid JSON.
var resumeOutputSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"name": map[string]any{
			"type":        "string",
			"description": "Full name of the candidate",
		},
		"email": map[string]any{
			"type":        "string",
			"description": "Email address if found",
		},
		"phone": map[string]any{
			"type":        "string",
			"description": "Phone number if found",
		},
		"location": map[string]any{
			"type":        "string",
			"description": "City, State/Country if found",
		},
		"summary": map[string]any{
			"type":        "string",
			"description": "Professional summary or objective if found",
		},
		"experience": map[string]any{
			"type":        "array",
			"description": "Work experience entries",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"company": map[string]any{
						"type":        "string",
						"description": "Company name",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "Job title",
					},
					"location": map[string]any{
						"type":        "string",
						"description": "Job location if found",
					},
					"startDate": map[string]any{
						"type":        "string",
						"description": "Start date (e.g., 'Jan 2020')",
					},
					"endDate": map[string]any{
						"type":        "string",
						"description": "End date or 'Present'",
					},
					"isCurrent": map[string]any{
						"type":        "boolean",
						"description": "True if this is the current job",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Job description or responsibilities",
					},
				},
				"required": []string{"company", "title"},
			},
		},
		"education": map[string]any{
			"type":        "array",
			"description": "Education entries",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"institution": map[string]any{
						"type":        "string",
						"description": "School/University name",
					},
					"degree": map[string]any{
						"type":        "string",
						"description": "Degree type (e.g., 'Bachelor of Science')",
					},
					"field": map[string]any{
						"type":        "string",
						"description": "Field of study",
					},
					"startDate": map[string]any{
						"type":        "string",
						"description": "Start date if found",
					},
					"endDate": map[string]any{
						"type":        "string",
						"description": "End/graduation date",
					},
					"gpa": map[string]any{
						"type":        "string",
						"description": "GPA if mentioned",
					},
					"achievements": map[string]any{
						"type":        "string",
						"description": "Notable achievements or honors",
					},
				},
				"required": []string{"institution"},
			},
		},
		"skills": map[string]any{
			"type":        "array",
			"description": "List of skills",
			"items": map[string]any{
				"type": "string",
			},
		},
		"confidence": map[string]any{
			"type":        "number",
			"description": "Confidence in extraction accuracy (0.0 to 1.0)",
		},
	},
	"required": []string{"name", "experience", "education", "skills", "confidence"},
}

// stripMarkdownCodeBlock removes markdown code block delimiters from LLM responses.
// LLMs often wrap JSON responses in ```json ... ``` blocks.
func stripMarkdownCodeBlock(content string) string {
	content = strings.TrimSpace(content)

	// Check for ```json or ``` at the start
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}

	// Remove ``` at the end if present
	content = strings.TrimSuffix(content, "```")

	return strings.TrimSpace(content)
}

// trailingCommaRegex matches trailing commas before ] or }
var trailingCommaRegex = regexp.MustCompile(`,\s*([\]}])`)

// fixTrailingCommas removes trailing commas from JSON which Go's parser doesn't accept.
func fixTrailingCommas(content string) string {
	return trailingCommaRegex.ReplaceAllString(content, "$1")
}

// ExtractResumeData implements domain.DocumentExtractor interface.
// It extracts structured resume data from text using LLM with structured output.
func (e *DocumentExtractor) ExtractResumeData(ctx context.Context, text string) (*domain.ResumeExtractedData, error) {
	llmReq := domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, resumeExtractionPrompt+text),
		},
		Model:        e.config.DefaultModel,
		MaxTokens:    e.config.MaxTokens,
		OutputSchema: resumeOutputSchema,
	}

	resp, err := e.provider.Complete(ctx, llmReq)
	if err != nil {
		return nil, fmt.Errorf("LLM extraction failed: %w", err)
	}

	// Clean up the JSON response
	jsonContent := stripMarkdownCodeBlock(resp.Content)
	jsonContent = fixTrailingCommas(jsonContent)

	// Parse JSON response
	var data domain.ResumeExtractedData
	if err := json.Unmarshal([]byte(jsonContent), &data); err != nil {
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
