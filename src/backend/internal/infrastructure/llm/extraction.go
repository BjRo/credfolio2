package llm

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"backend/internal/domain"
)

// Embedded prompts from external files for easier maintenance and review.
// Prompts are split into system (instructions) and user (content) for better
// token caching and clearer separation of concerns.

// Document extraction prompts
//
//go:embed prompts/document_extraction_system.txt
var documentExtractionSystemPrompt string

//go:embed prompts/document_extraction_user.txt
var documentExtractionUserPrompt string

// Resume extraction prompts
//
//go:embed prompts/resume_extraction_system.txt
var resumeExtractionSystemPrompt string

//go:embed prompts/resume_extraction_user.txt
var resumeExtractionUserTemplate string

// Compiled templates for user prompts with placeholder substitution
var resumeUserTemplate = template.Must(template.New("resume_user").Parse(resumeExtractionUserTemplate))

// DocumentExtractorConfig holds configuration for the document extractor.
type DocumentExtractorConfig struct { //nolint:govet // Field order prioritizes readability
	// DefaultModel is the model to use for extraction if not specified per-request.
	// If empty, the provider's default model is used.
	// Deprecated: Use DocumentExtractionChain and ResumeExtractionChain instead.
	DefaultModel string

	// MaxTokens for extraction responses. Defaults to 8192.
	MaxTokens int

	// ProviderRegistry holds all available providers for chain-based access.
	// If nil, falls back to the single provider passed to NewDocumentExtractor.
	ProviderRegistry *ProviderRegistry

	// DocumentExtractionChain specifies the provider chain for document text extraction.
	// If nil or empty, uses the default provider.
	DocumentExtractionChain ProviderChain

	// ResumeExtractionChain specifies the provider chain for resume data extraction.
	// If nil or empty, uses the default provider.
	ResumeExtractionChain ProviderChain
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
	defaultProvider domain.LLMProvider
	config          DocumentExtractorConfig
}

// NewDocumentExtractor creates a new document extractor.
func NewDocumentExtractor(provider domain.LLMProvider, config DocumentExtractorConfig) *DocumentExtractor {
	if config.MaxTokens == 0 {
		config.MaxTokens = 8192
	}
	return &DocumentExtractor{
		defaultProvider: provider,
		config:          config,
	}
}

// getProviderForChain returns the appropriate provider for a given chain.
// If the chain is configured and registry is available, returns a chained provider.
// Otherwise, falls back to the default provider.
func (e *DocumentExtractor) getProviderForChain(chain ProviderChain) domain.LLMProvider {
	if len(chain) > 0 && e.config.ProviderRegistry != nil {
		chained, err := NewChainedProvider(e.config.ProviderRegistry, chain)
		if err == nil {
			return chained
		}
		// Fall back to default on error
	}
	return e.defaultProvider
}

// ExtractTextWithRequest extracts text from a document image or PDF using a detailed request.
func (e *DocumentExtractor) ExtractTextWithRequest(ctx context.Context, req ExtractionRequest) (*ExtractionResult, error) {
	// Get the appropriate provider for document extraction
	provider := e.getProviderForChain(e.config.DocumentExtractionChain)

	// Determine prompts - use system/user split for better token caching
	systemPrompt := documentExtractionSystemPrompt
	userPrompt := req.CustomPrompt
	if userPrompt == "" {
		userPrompt = documentExtractionUserPrompt
	}

	// Determine model
	model := req.Model
	if model == "" {
		model = e.config.DefaultModel
	}

	// Build LLM request with system prompt and image in user message
	llmReq := domain.LLMRequest{
		SystemPrompt: systemPrompt,
		Messages: []domain.Message{
			domain.NewImageMessage(domain.RoleUser, req.MediaType, req.Document, userPrompt),
		},
		Model:     model,
		MaxTokens: e.config.MaxTokens,
	}

	// Execute extraction
	resp, err := provider.Complete(ctx, llmReq)
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
						"description": "Start date in ISO format YYYY-MM-DD. Year is REQUIRED. Use 01 for unknown day/month. Return null if year cannot be determined.",
					},
					"endDate": map[string]any{
						"type":        "string",
						"description": "End date in ISO format YYYY-MM-DD, or null if isCurrent is true or year cannot be determined.",
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
						"description": "Start date in ISO format YYYY-MM-DD. Year is REQUIRED. Use 01 for unknown day/month. Return null if year cannot be determined.",
					},
					"endDate": map[string]any{
						"type":        "string",
						"description": "Graduation/end date in ISO format YYYY-MM-DD. Year is REQUIRED. Use 01 for unknown day/month. Return null if year cannot be determined.",
					},
					"gpa": map[string]any{
						"type":        "string",
						"description": "GPA as numeric string like 3.8 or 3.8/4.0. Only if explicitly stated. Not a date.",
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

// ResumeTemplateData holds the data for rendering the resume extraction user prompt.
type ResumeTemplateData struct {
	Text string
}

// ExtractResumeData implements domain.DocumentExtractor interface.
// It extracts structured resume data from text using LLM with structured output.
func (e *DocumentExtractor) ExtractResumeData(ctx context.Context, text string) (*domain.ResumeExtractedData, error) {
	// Get the appropriate provider for resume extraction
	provider := e.getProviderForChain(e.config.ResumeExtractionChain)

	// Render the user prompt template with the resume text
	var userPromptBuf bytes.Buffer
	if err := resumeUserTemplate.Execute(&userPromptBuf, ResumeTemplateData{Text: text}); err != nil {
		return nil, fmt.Errorf("failed to render user prompt template: %w", err)
	}

	llmReq := domain.LLMRequest{
		SystemPrompt: resumeExtractionSystemPrompt,
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, userPromptBuf.String()),
		},
		Model:        e.config.DefaultModel,
		MaxTokens:    e.config.MaxTokens,
		OutputSchema: resumeOutputSchema,
	}

	resp, err := provider.Complete(ctx, llmReq)
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

	// Normalize extracted data to fix OCR/PDF text extraction artifacts
	normalizeResumeData(&data)

	return &data, nil
}

// normalizeResumeData cleans up extracted data to fix common extraction artifacts
// such as spurious spaces in words from PDF text extraction.
func normalizeResumeData(data *domain.ResumeExtractedData) {
	// Normalize top-level fields
	data.Name = NormalizeSpacedText(data.Name)
	data.Email = normalizeOptionalText(data.Email)
	data.Location = normalizeOptionalText(data.Location)

	// Normalize education entries
	for i := range data.Education {
		normalizeEducation(&data.Education[i])
	}

	// Normalize experience entries
	for i := range data.Experience {
		normalizeExperience(&data.Experience[i])
	}

	// Normalize skills
	for i := range data.Skills {
		data.Skills[i] = NormalizeSpacedText(data.Skills[i])
	}
}

// normalizeOptionalText normalizes an optional string pointer
func normalizeOptionalText(s *string) *string {
	if s == nil {
		return nil
	}
	normalized := NormalizeSpacedText(*s)
	return &normalized
}

// normalizeOptionalDate normalizes an optional date string pointer
func normalizeOptionalDate(s *string) *string {
	if s == nil {
		return nil
	}
	normalized := NormalizeDate(*s)
	if normalized == "" {
		return nil
	}
	return &normalized
}

// normalizeEducation normalizes a single education entry
func normalizeEducation(edu *domain.Education) {
	edu.Institution = NormalizeSpacedText(edu.Institution)
	edu.Degree = normalizeOptionalText(edu.Degree)
	edu.Field = normalizeOptionalText(edu.Field)
	edu.StartDate = normalizeOptionalDate(edu.StartDate)
	edu.EndDate = normalizeOptionalDate(edu.EndDate)
}

// normalizeExperience normalizes a single work experience entry
func normalizeExperience(exp *domain.WorkExperience) {
	exp.Company = NormalizeSpacedText(exp.Company)
	exp.Title = NormalizeSpacedText(exp.Title)
	exp.Location = normalizeOptionalText(exp.Location)
	exp.StartDate = normalizeOptionalDate(exp.StartDate)
	exp.EndDate = normalizeOptionalDate(exp.EndDate)
}

// Verify DocumentExtractor implements domain.DocumentExtractor interface.
var _ domain.DocumentExtractor = (*DocumentExtractor)(nil)
