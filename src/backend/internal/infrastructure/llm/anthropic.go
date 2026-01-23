// Package llm provides LLM provider implementations.
package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/domain"
)

const (
	defaultAnthropicBaseURL = "https://api.anthropic.com"
	defaultAnthropicModel   = "claude-sonnet-4-20250514"
	defaultMaxTokens        = 4096
	anthropicAPIVersion     = "2023-06-01"
)

// AnthropicConfig holds configuration for the Anthropic provider.
type AnthropicConfig struct {
	// APIKey is the Anthropic API key.
	APIKey string

	// BaseURL is the API base URL (optional, for testing).
	BaseURL string

	// DefaultModel is the model to use when not specified per-request.
	// Defaults to claude-sonnet-4-20250514.
	DefaultModel string

	// Timeout is the HTTP client timeout (defaults to 60s).
	Timeout time.Duration
}

// AnthropicProvider implements domain.LLMProvider for Anthropic's Claude API.
type AnthropicProvider struct { //nolint:govet // Field order prioritizes readability
	config AnthropicConfig
	client *http.Client
}

// NewAnthropicProvider creates a new Anthropic provider.
func NewAnthropicProvider(config AnthropicConfig) *AnthropicProvider {
	if config.BaseURL == "" {
		config.BaseURL = defaultAnthropicBaseURL
	}
	if config.DefaultModel == "" {
		config.DefaultModel = defaultAnthropicModel
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	return &AnthropicProvider{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Name returns the provider name.
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Complete sends a request to the Anthropic API.
func (p *AnthropicProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	// Build request body
	body := p.buildRequestBody(req)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, &domain.LLMError{
			Provider: p.Name(),
			Message:  "failed to marshal request",
			Err:      err,
		}
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.config.BaseURL+"/v1/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, &domain.LLMError{
			Provider: p.Name(),
			Message:  "failed to create request",
			Err:      err,
		}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.config.APIKey)
	httpReq.Header.Set("anthropic-version", anthropicAPIVersion)

	// Send request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, &domain.LLMError{
			Provider:  p.Name(),
			Message:   "request failed",
			Retryable: true,
			Err:       err,
		}
	}
	defer resp.Body.Close() //nolint:errcheck // Best effort cleanup

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &domain.LLMError{
			Provider: p.Name(),
			Message:  "failed to read response",
			Err:      err,
		}
	}

	// Handle errors
	if resp.StatusCode != http.StatusOK {
		return nil, p.parseError(resp.StatusCode, respBody)
	}

	// Parse successful response
	return p.parseResponse(respBody)
}

// buildRequestBody constructs the Anthropic API request body.
func (p *AnthropicProvider) buildRequestBody(req domain.LLMRequest) map[string]any {
	// Use request model if specified, otherwise use config default
	model := req.Model
	if model == "" {
		model = p.config.DefaultModel
	}

	body := map[string]any{
		"model":    model,
		"messages": p.convertMessages(req.Messages),
	}

	// Set max_tokens
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}
	body["max_tokens"] = maxTokens

	// Set optional fields
	if req.SystemPrompt != "" {
		body["system"] = req.SystemPrompt
	}

	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}

	return body
}

// convertMessages converts domain messages to Anthropic API format.
func (p *AnthropicProvider) convertMessages(messages []domain.Message) []map[string]any {
	result := make([]map[string]any, 0, len(messages))

	for _, msg := range messages {
		converted := map[string]any{
			"role":    string(msg.Role),
			"content": p.convertContentBlocks(msg.Content),
		}
		result = append(result, converted)
	}

	return result
}

// convertContentBlocks converts domain content blocks to Anthropic API format.
func (p *AnthropicProvider) convertContentBlocks(blocks []domain.ContentBlock) []map[string]any {
	result := make([]map[string]any, 0, len(blocks))

	for _, block := range blocks {
		switch block.Type {
		case domain.ContentTypeText:
			result = append(result, map[string]any{
				"type": "text",
				"text": block.Text,
			})
		case domain.ContentTypeImage:
			result = append(result, map[string]any{
				"type": "image",
				"source": map[string]any{
					"type":       "base64",
					"media_type": string(block.ImageMediaType),
					"data":       base64.StdEncoding.EncodeToString(block.ImageData),
				},
			})
		}
	}

	return result
}

// anthropicErrorResponse represents an error response from the API.
type anthropicErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// parseError converts an API error response to a domain error.
func (p *AnthropicProvider) parseError(statusCode int, body []byte) error {
	var errResp anthropicErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &domain.LLMError{
			Provider: p.Name(),
			Message:  fmt.Sprintf("HTTP %d: %s", statusCode, string(body)),
			Err:      err,
		}
	}

	// Determine if error is retryable
	retryable := statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode >= 500

	return &domain.LLMError{
		Provider:  p.Name(),
		Code:      errResp.Error.Type,
		Message:   errResp.Error.Message,
		Retryable: retryable,
	}
}

// anthropicResponse represents a successful response from the API.
type anthropicResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Role       string `json:"role"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Content    []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// parseResponse converts an API response to a domain response.
func (p *AnthropicProvider) parseResponse(body []byte) (*domain.LLMResponse, error) {
	var resp anthropicResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, &domain.LLMError{
			Provider: p.Name(),
			Message:  "failed to parse response",
			Err:      err,
		}
	}

	// Extract text content
	var content string
	for _, block := range resp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &domain.LLMResponse{
		Content:      content,
		Model:        resp.Model,
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		StopReason:   resp.StopReason,
	}, nil
}

// Verify AnthropicProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*AnthropicProvider)(nil)
