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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	otelTrace "go.opentelemetry.io/otel/trace"

	"backend/internal/domain"
	"backend/internal/logger"
)

const tracerName = "credfolio"

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

// Reference letter extraction prompts
//
//go:embed prompts/reference_letter_extraction_system.txt
var letterExtractionSystemPrompt string

//go:embed prompts/reference_letter_extraction_user.txt
var letterExtractionUserTemplate string

// Document detection prompts
//
//go:embed prompts/document_detection_system.txt
var detectionSystemPrompt string

//go:embed prompts/document_detection_user.txt
var detectionUserTemplate string

// Compiled templates for user prompts with placeholder substitution
var resumeUserTemplate = template.Must(template.New("resume_user").Parse(resumeExtractionUserTemplate))
var letterUserTemplate = template.Must(template.New("letter_user").Parse(letterExtractionUserTemplate))
var detectionUserTmpl = template.Must(template.New("detection_user").Parse(detectionUserTemplate))

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

	// ReferenceExtractionChain specifies the provider chain for reference letter data extraction.
	// If nil or empty, falls back to ResumeExtractionChain, then the default provider.
	ReferenceExtractionChain ProviderChain

	// DetectionChain specifies the provider chain for lightweight document content detection.
	// If nil or empty, falls back to ResumeExtractionChain, then the default provider.
	DetectionChain ProviderChain

	// Logger for logging chain fallback events. If nil, fallbacks are silent.
	Logger logger.Logger
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
		if e.config.Logger != nil {
			e.config.Logger.Warning("Provider chain creation failed, falling back to default provider",
				logger.Feature("llm"),
				logger.String("chain_provider", chain.Primary().Provider),
				logger.String("chain_model", chain.Primary().Model),
				logger.Err(err),
			)
		}
	}
	return e.defaultProvider
}

// ExtractTextWithRequest extracts text from a document image or PDF using a detailed request.
// For PDFs, it first attempts local (Go-native) text extraction which is nearly instant.
// If the local extraction produces usable text, the LLM vision call is skipped entirely.
// For scanned/image-based PDFs or non-PDF documents, the LLM vision path is used as before.
func (e *DocumentExtractor) ExtractTextWithRequest(ctx context.Context, req ExtractionRequest) (*ExtractionResult, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "pdf_text_extraction",
		otelTrace.WithAttributes(
			attribute.String("content_type", string(req.MediaType)),
		),
	)
	defer span.End()

	// For PDFs without a custom prompt, try fast local extraction first.
	if req.MediaType == domain.ImageMediaTypePDF && req.CustomPrompt == "" {
		text, localErr := extractTextFromPDF(req.Document)
		if localErr == nil && isUsableText(text) {
			span.SetAttributes(attribute.String("extraction_method", "local"))
			return &ExtractionResult{Text: text}, nil
		}
		// Local extraction failed or produced unusable text — fall through to LLM.
		span.SetAttributes(attribute.String("local_extraction_skipped_reason", localFallbackReason(localErr, text)))
	}

	span.SetAttributes(attribute.String("extraction_method", "llm"))

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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &ExtractionResult{
		Text:         resp.Content,
		InputTokens:  resp.InputTokens,
		OutputTokens: resp.OutputTokens,
	}, nil
}

// localFallbackReason returns a human-readable reason why local PDF extraction was skipped.
func localFallbackReason(err error, text string) string {
	if err != nil {
		return "parse_error"
	}
	if strings.TrimSpace(text) == "" {
		return "empty_text"
	}
	return "low_quality_text"
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
// This schema is used with OpenAI/Anthropic structured output features to guarantee valid JSON.
// Note: OpenAI requires "additionalProperties": false at every object level.
var resumeOutputSchema = map[string]any{
	"type":                 "object",
	"additionalProperties": false,
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
			"description": "Professional summary or objective. If none is explicitly present, synthesize a brief 2-3 sentence summary from the candidate's experience and skills.",
		},
		"experience": map[string]any{
			"type":        "array",
			"description": "Work experience entries",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
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
				"required": []string{"company", "title", "location", "startDate", "endDate", "isCurrent", "description"},
			},
		},
		"education": map[string]any{
			"type":        "array",
			"description": "Education entries",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
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
				"required": []string{"institution", "degree", "field", "startDate", "endDate", "gpa", "achievements"},
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
	"required": []string{"name", "email", "phone", "location", "summary", "experience", "education", "skills", "confidence"},
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
	ctx, span := otel.Tracer(tracerName).Start(ctx, "resume_data_extraction",
		otelTrace.WithAttributes(
			attribute.Int("text_length", len(text)),
		),
	)
	defer span.End()

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
		MaxTokens:    e.config.MaxTokens,
		OutputSchema: resumeOutputSchema,
	}

	resp, err := provider.Complete(ctx, llmReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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

// letterOutputSchema defines the JSON schema for structured reference letter extraction.
// This schema is used with OpenAI/Anthropic structured output features to guarantee valid JSON.
var letterOutputSchema = map[string]any{
	"type":                 "object",
	"additionalProperties": false,
	"properties": map[string]any{
		"author": map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"description":          "Information about the letter author",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "Full name of the letter author",
				},
				"title": map[string]any{
					"type":        "string",
					"description": "Job title/position of the author if mentioned",
				},
				"company": map[string]any{
					"type":        "string",
					"description": "Organization where the author works if mentioned",
				},
				"relationship": map[string]any{
					"type":        "string",
					"description": "Relationship to candidate: manager, peer, direct_report, client, mentor, professor, colleague, or other",
					"enum":        []string{"manager", "peer", "direct_report", "client", "mentor", "professor", "colleague", "other"},
				},
			},
			"required": []string{"name", "title", "company", "relationship"},
		},
		"testimonials": map[string]any{
			"type":        "array",
			"description": "2-4 meaningful quotes suitable for display on profile",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"quote": map[string]any{
						"type":        "string",
						"description": "Complete, impactful statement about the candidate",
					},
					"skillsMentioned": map[string]any{
						"type":        "array",
						"description": "Skills referenced in this quote",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
				"required": []string{"quote", "skillsMentioned"},
			},
		},
		"skillMentions": map[string]any{
			"type":        "array",
			"description": "Specific mentions of technical or professional skills",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"skill": map[string]any{
						"type":        "string",
						"description": "Skill name (normalized, e.g., 'Golang' → 'Go')",
					},
					"quote": map[string]any{
						"type":        "string",
						"description": "Exact sentence(s) mentioning this skill",
					},
					"context": map[string]any{
						"type":        "string",
						"description": "Brief category like 'technical skills', 'leadership', 'communication'",
					},
				},
				"required": []string{"skill", "quote", "context"},
			},
		},
		"experienceMentions": map[string]any{
			"type":        "array",
			"description": "References to specific roles or companies the candidate held",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"company": map[string]any{
						"type":        "string",
						"description": "Company/organization name mentioned",
					},
					"role": map[string]any{
						"type":        "string",
						"description": "Job title or role mentioned",
					},
					"quote": map[string]any{
						"type":        "string",
						"description": "Sentence(s) discussing this role/company",
					},
				},
				"required": []string{"company", "role", "quote"},
			},
		},
		"discoveredSkills": map[string]any{
			"type":        "array",
			"description": "Skills mentioned that are NOT in the candidate's existing profile skills",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"skill": map[string]any{
						"type":        "string",
						"description": "Skill name (normalized, e.g., 'Golang' → 'Go')",
					},
					"quote": map[string]any{
						"type":        "string",
						"description": "Exact sentence(s) mentioning this skill",
					},
					"context": map[string]any{
						"type":        "string",
						"description": "Brief description of how the skill was demonstrated",
					},
					"category": map[string]any{
						"type":        "string",
						"description": "Skill category: TECHNICAL (tools, languages, methodologies), SOFT (interpersonal, leadership), or DOMAIN (industry knowledge)",
						"enum":        []string{"TECHNICAL", "SOFT", "DOMAIN"},
					},
				},
				"required": []string{"skill", "quote", "context", "category"},
			},
		},
	},
	"required": []string{"author", "testimonials", "skillMentions", "experienceMentions", "discoveredSkills"},
}

// LetterTemplateData holds the data for rendering the letter extraction user prompt.
type LetterTemplateData struct {
	Text          string
	ProfileSkills []domain.ProfileSkillContext
}

// ExtractLetterData implements domain.DocumentExtractor interface.
// It extracts structured credibility data from reference letter text using LLM with structured output.
// The profileSkills parameter provides existing skills context for the LLM to distinguish
// between mentions of existing skills (for validation) and newly discovered skills.
//
//nolint:gocyclo // Complex extraction logic with multiple validation paths
func (e *DocumentExtractor) ExtractLetterData(ctx context.Context, text string, profileSkills []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "letter_data_extraction",
		otelTrace.WithAttributes(
			attribute.Int("text_length", len(text)),
		),
	)
	defer span.End()

	// Get the appropriate provider for reference letter extraction
	chain := e.config.ReferenceExtractionChain
	if len(chain) == 0 {
		chain = e.config.ResumeExtractionChain // Fall back to resume chain for backwards compatibility
	}
	provider := e.getProviderForChain(chain)

	// Render the user prompt template with the letter text and profile skills context
	var userPromptBuf bytes.Buffer
	if err := letterUserTemplate.Execute(&userPromptBuf, LetterTemplateData{Text: text, ProfileSkills: profileSkills}); err != nil {
		return nil, fmt.Errorf("failed to render user prompt template: %w", err)
	}

	llmReq := domain.LLMRequest{
		SystemPrompt: letterExtractionSystemPrompt,
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, userPromptBuf.String()),
		},
		MaxTokens:    e.config.MaxTokens,
		OutputSchema: letterOutputSchema,
	}

	resp, err := provider.Complete(ctx, llmReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("LLM extraction failed: %w", err)
	}

	// Clean up the JSON response
	jsonContent := stripMarkdownCodeBlock(resp.Content)
	jsonContent = fixTrailingCommas(jsonContent)

	// Parse JSON response into raw structure first
	var rawData struct {
		Author struct {
			Name         string `json:"name"`
			Title        string `json:"title"`
			Company      string `json:"company"`
			Relationship string `json:"relationship"`
		} `json:"author"`
		Testimonials []struct {
			Quote           string   `json:"quote"`
			SkillsMentioned []string `json:"skillsMentioned"`
		} `json:"testimonials"`
		SkillMentions []struct {
			Skill   string `json:"skill"`
			Quote   string `json:"quote"`
			Context string `json:"context"`
		} `json:"skillMentions"`
		ExperienceMentions []struct {
			Company string `json:"company"`
			Role    string `json:"role"`
			Quote   string `json:"quote"`
		} `json:"experienceMentions"`
		DiscoveredSkills []struct {
			Skill    string `json:"skill"`
			Quote    string `json:"quote"`
			Context  string `json:"context"`
			Category string `json:"category"`
		} `json:"discoveredSkills"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse extraction response: %w", err)
	}

	// Convert to domain types with proper pointer handling
	data := &domain.ExtractedLetterData{
		Author: domain.ExtractedAuthor{
			Name:         rawData.Author.Name,
			Relationship: domain.AuthorRelationship(rawData.Author.Relationship),
		},
		Testimonials:       make([]domain.ExtractedTestimonial, 0, len(rawData.Testimonials)),
		SkillMentions:      make([]domain.ExtractedSkillMention, 0, len(rawData.SkillMentions)),
		ExperienceMentions: make([]domain.ExtractedExperienceMention, 0, len(rawData.ExperienceMentions)),
		DiscoveredSkills:   make([]domain.DiscoveredSkill, 0, len(rawData.DiscoveredSkills)),
	}

	// Handle optional author fields
	if rawData.Author.Title != "" {
		data.Author.Title = &rawData.Author.Title
	}
	if rawData.Author.Company != "" {
		data.Author.Company = &rawData.Author.Company
	}

	// Convert testimonials
	for _, t := range rawData.Testimonials {
		data.Testimonials = append(data.Testimonials, domain.ExtractedTestimonial{
			Quote:           t.Quote,
			SkillsMentioned: t.SkillsMentioned,
		})
	}

	// Convert skill mentions
	for _, s := range rawData.SkillMentions {
		mention := domain.ExtractedSkillMention{
			Skill: s.Skill,
			Quote: s.Quote,
		}
		if s.Context != "" {
			mention.Context = &s.Context
		}
		data.SkillMentions = append(data.SkillMentions, mention)
	}

	// Convert experience mentions
	for _, e := range rawData.ExperienceMentions {
		data.ExperienceMentions = append(data.ExperienceMentions, domain.ExtractedExperienceMention{
			Company: e.Company,
			Role:    e.Role,
			Quote:   e.Quote,
		})
	}

	// Convert discovered skills
	for _, ds := range rawData.DiscoveredSkills {
		skill := domain.DiscoveredSkill{
			Skill:    ds.Skill,
			Quote:    ds.Quote,
			Category: mapDiscoveredSkillCategory(ds.Category),
		}
		if ds.Context != "" {
			skill.Context = &ds.Context
		}
		data.DiscoveredSkills = append(data.DiscoveredSkills, skill)
	}

	// Ensure slices are initialized (not nil)
	if data.Testimonials == nil {
		data.Testimonials = []domain.ExtractedTestimonial{}
	}
	if data.SkillMentions == nil {
		data.SkillMentions = []domain.ExtractedSkillMention{}
	}
	if data.ExperienceMentions == nil {
		data.ExperienceMentions = []domain.ExtractedExperienceMention{}
	}
	if data.DiscoveredSkills == nil {
		data.DiscoveredSkills = []domain.DiscoveredSkill{}
	}

	// Set model version from the actual LLM response so callers don't need to hardcode it
	data.Metadata.ModelVersion = resp.Model

	return data, nil
}

// detectionOutputSchema defines the JSON schema for structured document detection.
var detectionOutputSchema = map[string]any{
	"type":                 "object",
	"additionalProperties": false,
	"properties": map[string]any{
		"hasCareerInfo": map[string]any{
			"type":        "boolean",
			"description": "Whether the document contains resume/CV career content (work experience, education, skills)",
		},
		"hasTestimonial": map[string]any{
			"type":        "boolean",
			"description": "Whether the document contains recommendation/reference letter content",
		},
		"testimonialAuthor": map[string]any{
			"type":        "string",
			"description": "Name of the person who wrote the recommendation, or empty string if none",
		},
		"confidence": map[string]any{
			"type":        "number",
			"description": "Confidence in the classification (0.0 to 1.0)",
		},
		"summary": map[string]any{
			"type":        "string",
			"description": "Brief one-sentence summary of the document content",
		},
		"documentTypeHint": map[string]any{
			"type":        "string",
			"description": "Document type: resume, reference_letter, hybrid, or unknown",
			"enum":        []string{"resume", "reference_letter", "hybrid", "unknown"},
		},
	},
	"required": []string{"hasCareerInfo", "hasTestimonial", "testimonialAuthor", "confidence", "summary", "documentTypeHint"},
}

// DetectionTemplateData holds the data for rendering the detection user prompt.
type DetectionTemplateData struct {
	Text string
}

// DetectDocumentContent performs lightweight classification of a document's content.
// It quickly identifies whether the document contains career information, testimonials, or both,
// without running full extraction. This is significantly faster and cheaper than full extraction.
func (e *DocumentExtractor) DetectDocumentContent(ctx context.Context, text string) (*domain.DocumentDetectionResult, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "document_content_detection",
		otelTrace.WithAttributes(
			attribute.Int("text_length", len(text)),
		),
	)
	defer span.End()

	// Get the appropriate provider for detection
	chain := e.config.DetectionChain
	if len(chain) == 0 {
		chain = e.config.ResumeExtractionChain // Fall back to resume chain
	}
	provider := e.getProviderForChain(chain)

	// Render the user prompt template with the document text
	var userPromptBuf bytes.Buffer
	if err := detectionUserTmpl.Execute(&userPromptBuf, DetectionTemplateData{Text: text}); err != nil {
		return nil, fmt.Errorf("failed to render detection user prompt template: %w", err)
	}

	llmReq := domain.LLMRequest{
		SystemPrompt: detectionSystemPrompt,
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, userPromptBuf.String()),
		},
		MaxTokens:    1024, // Detection is lightweight, doesn't need many tokens
		OutputSchema: detectionOutputSchema,
	}

	resp, err := provider.Complete(ctx, llmReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("LLM detection failed: %w", err)
	}

	// Clean up the JSON response
	jsonContent := stripMarkdownCodeBlock(resp.Content)
	jsonContent = fixTrailingCommas(jsonContent)

	// Parse JSON response
	var rawData struct { //nolint:govet // Field order matches JSON output for readability
		HasCareerInfo     bool    `json:"hasCareerInfo"`
		HasTestimonial    bool    `json:"hasTestimonial"`
		TestimonialAuthor string  `json:"testimonialAuthor"`
		Confidence        float64 `json:"confidence"`
		Summary           string  `json:"summary"`
		DocumentTypeHint  string  `json:"documentTypeHint"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse detection response: %w", err)
	}

	result := &domain.DocumentDetectionResult{
		HasCareerInfo:    rawData.HasCareerInfo,
		HasTestimonial:   rawData.HasTestimonial,
		Confidence:       rawData.Confidence,
		Summary:          rawData.Summary,
		DocumentTypeHint: domain.DocumentTypeHint(rawData.DocumentTypeHint),
	}

	// Handle optional testimonial author
	if rawData.TestimonialAuthor != "" {
		result.TestimonialAuthor = &rawData.TestimonialAuthor
	}

	return result, nil
}

// mapDiscoveredSkillCategory maps a raw category string from LLM output to domain.SkillCategory.
// Defaults to SOFT if the category is not recognized.
func mapDiscoveredSkillCategory(raw string) domain.SkillCategory {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "TECHNICAL":
		return domain.SkillCategoryTechnical
	case "SOFT":
		return domain.SkillCategorySoft
	case "DOMAIN":
		return domain.SkillCategoryDomain
	default:
		return domain.SkillCategorySoft
	}
}

// Verify DocumentExtractor implements domain.DocumentExtractor interface.
var _ domain.DocumentExtractor = (*DocumentExtractor)(nil)
