// Package llm provides LLM provider implementations.
package llm

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"backend/internal/domain"
)

const (
	// Claude Sonnet 4.5 supports structured outputs
	defaultAnthropicModel = "claude-sonnet-4-5-20250929"
	defaultMaxTokens      = 4096
	// Beta header for structured outputs feature
	structuredOutputsBeta anthropic.AnthropicBeta = "structured-outputs-2025-11-13"
)

// AnthropicConfig holds configuration for the Anthropic provider.
type AnthropicConfig struct {
	HTTPClient   *http.Client
	APIKey       string
	BaseURL      string
	DefaultModel string
	Timeout      time.Duration
}

// AnthropicProvider implements domain.LLMProvider for Anthropic's Claude API.
type AnthropicProvider struct {
	config AnthropicConfig
	client anthropic.Client
}

// NewAnthropicProvider creates a new Anthropic provider using the official SDK.
func NewAnthropicProvider(config AnthropicConfig) *AnthropicProvider {
	if config.DefaultModel == "" {
		config.DefaultModel = defaultAnthropicModel
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	// Build client options
	opts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
	}

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}

	if config.HTTPClient != nil {
		opts = append(opts, option.WithHTTPClient(config.HTTPClient))
	}

	client := anthropic.NewClient(opts...)

	return &AnthropicProvider{
		client: client,
		config: config,
	}
}

// Name returns the provider name.
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Complete sends a request to the Anthropic API.
func (p *AnthropicProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	// Determine model and max tokens
	model := req.Model
	if model == "" {
		model = p.config.DefaultModel
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	// Choose API based on whether structured output is requested
	if req.OutputSchema != nil {
		return p.completeWithStructuredOutput(ctx, req, model, maxTokens)
	}
	return p.completeStandard(ctx, req, model, maxTokens)
}

// completeStandard handles regular (non-structured) completions.
func (p *AnthropicProvider) completeStandard(
	ctx context.Context,
	req domain.LLMRequest,
	model string,
	maxTokens int,
) (*domain.LLMResponse, error) {
	// Convert domain messages to SDK messages
	messages := p.convertMessages(req.Messages)

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		Messages:  messages,
	}

	// Add system prompt if provided
	if req.SystemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.SystemPrompt},
		}
	}

	// Add temperature if provided
	if req.Temperature > 0 {
		params.Temperature = anthropic.Float(req.Temperature)
	}

	msg, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, p.convertError(err)
	}

	return p.parseResponse(msg)
}

// completeWithStructuredOutput handles completions with JSON schema output.
func (p *AnthropicProvider) completeWithStructuredOutput(
	ctx context.Context,
	req domain.LLMRequest,
	model string,
	maxTokens int,
) (*domain.LLMResponse, error) {
	// Convert domain messages to beta SDK messages
	messages := p.convertBetaMessages(req.Messages)

	params := anthropic.BetaMessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: int64(maxTokens),
		Messages:  messages,
		Betas:     []anthropic.AnthropicBeta{structuredOutputsBeta},
	}

	// Add system prompt if provided
	if req.SystemPrompt != "" {
		params.System = []anthropic.BetaTextBlockParam{
			{Text: req.SystemPrompt},
		}
	}

	// Add temperature if provided
	if req.Temperature > 0 {
		params.Temperature = anthropic.Float(req.Temperature)
	}

	// Add output format with JSON schema using the SDK helper
	params.OutputFormat = anthropic.BetaJSONSchemaOutputFormat(req.OutputSchema)

	msg, err := p.client.Beta.Messages.New(ctx, params)
	if err != nil {
		return nil, p.convertError(err)
	}

	return p.parseBetaResponse(msg)
}

// convertMessages converts domain messages to SDK message params.
func (p *AnthropicProvider) convertMessages(messages []domain.Message) []anthropic.MessageParam {
	result := make([]anthropic.MessageParam, 0, len(messages))

	for _, msg := range messages {
		blocks := p.convertContentBlocks(msg.Content)
		var param anthropic.MessageParam
		switch msg.Role {
		case domain.RoleUser:
			param = anthropic.NewUserMessage(blocks...)
		case domain.RoleAssistant:
			param = anthropic.NewAssistantMessage(blocks...)
		}
		result = append(result, param)
	}

	return result
}

// convertBetaMessages converts domain messages to beta SDK message params.
func (p *AnthropicProvider) convertBetaMessages(messages []domain.Message) []anthropic.BetaMessageParam {
	result := make([]anthropic.BetaMessageParam, 0, len(messages))

	for _, msg := range messages {
		blocks := p.convertBetaContentBlocks(msg.Content)
		var param anthropic.BetaMessageParam
		switch msg.Role {
		case domain.RoleUser:
			param = anthropic.NewBetaUserMessage(blocks...)
		case domain.RoleAssistant:
			// No NewBetaAssistantMessage helper, construct manually
			param = anthropic.BetaMessageParam{
				Role:    anthropic.BetaMessageParamRoleAssistant,
				Content: blocks,
			}
		}
		result = append(result, param)
	}

	return result
}

// convertContentBlocks converts domain content blocks to SDK content block unions.
func (p *AnthropicProvider) convertContentBlocks(blocks []domain.ContentBlock) []anthropic.ContentBlockParamUnion {
	result := make([]anthropic.ContentBlockParamUnion, 0, len(blocks))

	for _, block := range blocks {
		switch block.Type {
		case domain.ContentTypeText:
			result = append(result, anthropic.NewTextBlock(block.Text))
		case domain.ContentTypeImage:
			encoded := base64.StdEncoding.EncodeToString(block.ImageData)
			// PDFs use document blocks, images use image blocks
			if block.ImageMediaType == domain.ImageMediaTypePDF {
				result = append(result, anthropic.NewDocumentBlock(
					anthropic.Base64PDFSourceParam{Data: encoded},
				))
			} else {
				result = append(result, anthropic.NewImageBlockBase64(
					string(block.ImageMediaType),
					encoded,
				))
			}
		}
	}

	return result
}

// convertBetaContentBlocks converts domain content blocks to beta SDK content block unions.
func (p *AnthropicProvider) convertBetaContentBlocks(blocks []domain.ContentBlock) []anthropic.BetaContentBlockParamUnion {
	result := make([]anthropic.BetaContentBlockParamUnion, 0, len(blocks))

	for _, block := range blocks {
		switch block.Type {
		case domain.ContentTypeText:
			result = append(result, anthropic.NewBetaTextBlock(block.Text))
		case domain.ContentTypeImage:
			encoded := base64.StdEncoding.EncodeToString(block.ImageData)
			// PDFs use document blocks, images use image blocks
			if block.ImageMediaType == domain.ImageMediaTypePDF {
				result = append(result, anthropic.NewBetaDocumentBlock(
					anthropic.BetaBase64PDFSourceParam{Data: encoded},
				))
			} else {
				result = append(result, anthropic.NewBetaImageBlock(
					anthropic.BetaBase64ImageSourceParam{
						MediaType: anthropic.BetaBase64ImageSourceMediaType(block.ImageMediaType),
						Data:      encoded,
					},
				))
			}
		}
	}

	return result
}

// parseResponse converts an SDK message to a domain response.
func (p *AnthropicProvider) parseResponse(msg *anthropic.Message) (*domain.LLMResponse, error) {
	// Extract text content from response blocks
	var content string
	for _, block := range msg.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &domain.LLMResponse{
		Content:      content,
		Model:        string(msg.Model),
		InputTokens:  int(msg.Usage.InputTokens),
		OutputTokens: int(msg.Usage.OutputTokens),
		StopReason:   string(msg.StopReason),
	}, nil
}

// parseBetaResponse converts a beta SDK message to a domain response.
func (p *AnthropicProvider) parseBetaResponse(msg *anthropic.BetaMessage) (*domain.LLMResponse, error) {
	// Extract text content from response blocks
	var content string
	for _, block := range msg.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &domain.LLMResponse{
		Content:      content,
		Model:        string(msg.Model),
		InputTokens:  int(msg.Usage.InputTokens),
		OutputTokens: int(msg.Usage.OutputTokens),
		StopReason:   string(msg.StopReason),
	}, nil
}

// convertError converts SDK errors to domain errors.
func (p *AnthropicProvider) convertError(err error) error {
	// Check for API errors from the SDK
	if apiErr, ok := err.(*anthropic.Error); ok {
		retryable := apiErr.StatusCode == http.StatusTooManyRequests ||
			apiErr.StatusCode == http.StatusServiceUnavailable ||
			apiErr.StatusCode >= 500

		return &domain.LLMError{
			Provider:  p.Name(),
			Message:   err.Error(),
			Retryable: retryable,
			Err:       err,
		}
	}

	// Generic error (network, etc.) - assume retryable
	return &domain.LLMError{
		Provider:  p.Name(),
		Message:   err.Error(),
		Retryable: true,
		Err:       err,
	}
}

// Verify AnthropicProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*AnthropicProvider)(nil)
