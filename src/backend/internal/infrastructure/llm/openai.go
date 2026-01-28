// Package llm provides LLM provider implementations.
package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"backend/internal/domain"
)

const (
	// GPT-5-nano supports vision and structured outputs
	defaultOpenAIModel   = "gpt-5-nano-2025-08-07"
	defaultOpenAITimeout = 60 * time.Second
)

// OpenAIConfig holds configuration for the OpenAI provider.
type OpenAIConfig struct {
	HTTPClient   *http.Client
	APIKey       string
	BaseURL      string
	DefaultModel string
	Timeout      time.Duration
}

// OpenAIProvider implements domain.LLMProvider for OpenAI's API.
type OpenAIProvider struct {
	config OpenAIConfig
	client openai.Client
}

// NewOpenAIProvider creates a new OpenAI provider.
func NewOpenAIProvider(config OpenAIConfig) *OpenAIProvider {
	if config.DefaultModel == "" {
		config.DefaultModel = defaultOpenAIModel
	}
	if config.Timeout == 0 {
		config.Timeout = defaultOpenAITimeout
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

	client := openai.NewClient(opts...)

	return &OpenAIProvider{
		client: client,
		config: config,
	}
}

// Name returns the provider name.
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Complete sends a request to the OpenAI API.
func (p *OpenAIProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	// Determine model and max tokens
	model := req.Model
	if model == "" {
		model = p.config.DefaultModel
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	// Convert domain messages to OpenAI messages
	messages := p.convertMessages(req)

	params := openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(model),
		MaxTokens: openai.Int(int64(maxTokens)),
		Messages:  messages,
	}

	// Add temperature if provided
	if req.Temperature > 0 {
		params.Temperature = openai.Float(req.Temperature)
	}

	// Add structured output if requested
	if req.OutputSchema != nil {
		schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
			Name:   "extraction_result",
			Schema: req.OutputSchema,
			Strict: openai.Bool(true),
		}
		params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		}
	}

	msg, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, p.convertError(err)
	}

	return p.parseResponse(msg)
}

// convertMessages converts domain messages to OpenAI message params.
func (p *OpenAIProvider) convertMessages(req domain.LLMRequest) []openai.ChatCompletionMessageParamUnion {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages)+1)

	// Add system message if provided
	if req.SystemPrompt != "" {
		result = append(result, openai.SystemMessage(req.SystemPrompt))
	}

	for _, msg := range req.Messages {
		switch msg.Role {
		case domain.RoleUser:
			result = append(result, p.convertUserMessage(msg))
		case domain.RoleAssistant:
			result = append(result, p.convertAssistantMessage(msg))
		}
	}

	return result
}

// convertUserMessage converts a domain user message to OpenAI format.
func (p *OpenAIProvider) convertUserMessage(msg domain.Message) openai.ChatCompletionMessageParamUnion {
	// Check if there are any image blocks
	hasImages := false
	for _, block := range msg.Content {
		if block.Type == domain.ContentTypeImage {
			hasImages = true
			break
		}
	}

	// If no images, use simple text message
	if !hasImages {
		var text string
		for _, block := range msg.Content {
			if block.Type == domain.ContentTypeText {
				text += block.Text
			}
		}
		return openai.UserMessage(text)
	}

	// Build content parts for multimodal message
	parts := make([]openai.ChatCompletionContentPartUnionParam, 0, len(msg.Content))
	for _, block := range msg.Content {
		switch block.Type {
		case domain.ContentTypeText:
			parts = append(parts, openai.TextContentPart(block.Text))
		case domain.ContentTypeImage:
			parts = append(parts, p.convertImageBlock(block))
		}
	}

	return openai.UserMessage(parts)
}

// convertAssistantMessage converts a domain assistant message to OpenAI format.
func (p *OpenAIProvider) convertAssistantMessage(msg domain.Message) openai.ChatCompletionMessageParamUnion {
	var text string
	for _, block := range msg.Content {
		if block.Type == domain.ContentTypeText {
			text += block.Text
		}
	}
	return openai.AssistantMessage(text)
}

// convertImageBlock converts a domain image content block to OpenAI format.
func (p *OpenAIProvider) convertImageBlock(block domain.ContentBlock) openai.ChatCompletionContentPartUnionParam {
	encoded := base64.StdEncoding.EncodeToString(block.ImageData)

	// For PDFs, use file content part
	if block.ImageMediaType == domain.ImageMediaTypePDF {
		return openai.FileContentPart(openai.ChatCompletionContentPartFileFileParam{
			FileData: openai.String(encoded),
			Filename: openai.String("document.pdf"),
		})
	}

	// For images, use image URL with data URI
	mediaType := string(block.ImageMediaType)
	dataURI := fmt.Sprintf("data:%s;base64,%s", mediaType, encoded)

	return openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
		URL:    dataURI,
		Detail: "auto",
	})
}

// parseResponse converts an OpenAI response to a domain response.
func (p *OpenAIProvider) parseResponse(msg *openai.ChatCompletion) (*domain.LLMResponse, error) {
	if len(msg.Choices) == 0 {
		return nil, &domain.LLMError{
			Provider:  p.Name(),
			Message:   "no choices in response",
			Retryable: false,
		}
	}

	choice := msg.Choices[0]
	content := choice.Message.Content

	var stopReason string
	if choice.FinishReason != "" {
		stopReason = string(choice.FinishReason)
	}

	return &domain.LLMResponse{
		Content:      content,
		Model:        msg.Model,
		InputTokens:  int(msg.Usage.PromptTokens),
		OutputTokens: int(msg.Usage.CompletionTokens),
		StopReason:   stopReason,
	}, nil
}

// convertError converts OpenAI errors to domain errors.
func (p *OpenAIProvider) convertError(err error) error {
	// Check for API errors from the SDK
	if apiErr, ok := err.(*openai.Error); ok {
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

// Verify OpenAIProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*OpenAIProvider)(nil)
