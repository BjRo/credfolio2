// Package domain contains the core business entities and repository interfaces.
package domain

import (
	"context"
)

// Role represents the role of a message sender in a conversation.
type Role string

// Role constants for LLM messages.
const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// ContentType represents the type of content in a message.
type ContentType string

// Content type constants.
const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
)

// ImageMediaType represents supported image formats.
type ImageMediaType string

// Image media type constants.
const (
	ImageMediaTypeJPEG ImageMediaType = "image/jpeg"
	ImageMediaTypePNG  ImageMediaType = "image/png"
	ImageMediaTypeGIF  ImageMediaType = "image/gif"
	ImageMediaTypeWebP ImageMediaType = "image/webp"
	ImageMediaTypePDF  ImageMediaType = "application/pdf"
)

// ContentBlock represents a single block of content within a message.
// A message can contain multiple content blocks (e.g., text + images).
type ContentBlock struct {
	Type ContentType `json:"type"`

	// Text content (when Type is ContentTypeText)
	Text string `json:"text,omitempty"`

	// Image content (when Type is ContentTypeImage)
	ImageMediaType ImageMediaType `json:"imageMediaType,omitempty"`
	ImageData      []byte         `json:"-"` // Base64-decoded image data
}

// Message represents a single message in an LLM conversation.
type Message struct {
	Role    Role           `json:"role"`
	Content []ContentBlock `json:"content"`
}

// NewTextMessage creates a simple text-only message.
func NewTextMessage(role Role, text string) Message {
	return Message{
		Role: role,
		Content: []ContentBlock{
			{Type: ContentTypeText, Text: text},
		},
	}
}

// NewImageMessage creates a message with an image and optional text.
func NewImageMessage(role Role, mediaType ImageMediaType, imageData []byte, text string) Message {
	blocks := []ContentBlock{
		{Type: ContentTypeImage, ImageMediaType: mediaType, ImageData: imageData},
	}
	if text != "" {
		blocks = append(blocks, ContentBlock{Type: ContentTypeText, Text: text})
	}
	return Message{
		Role:    role,
		Content: blocks,
	}
}

// LLMResponse represents the response from an LLM provider.
type LLMResponse struct { //nolint:govet // Field order prioritizes readability
	// Content is the text response from the model.
	Content string `json:"content"`

	// Model is the identifier of the model that generated this response.
	Model string `json:"model"`

	// InputTokens is the number of tokens in the input.
	InputTokens int `json:"inputTokens"`

	// OutputTokens is the number of tokens in the output.
	OutputTokens int `json:"outputTokens"`

	// StopReason indicates why the model stopped generating.
	StopReason string `json:"stopReason,omitempty"`
}

// LLMRequest represents a request to an LLM provider.
type LLMRequest struct { //nolint:govet // Field order prioritizes readability
	// Messages is the conversation history to send to the model.
	Messages []Message `json:"messages"`

	// SystemPrompt is an optional system prompt to guide the model's behavior.
	SystemPrompt string `json:"systemPrompt,omitempty"`

	// Model specifies which model to use. If empty, uses provider default.
	Model string `json:"model,omitempty"`

	// MaxTokens limits the response length. If 0, uses provider default.
	MaxTokens int `json:"maxTokens,omitempty"`

	// Temperature controls randomness (0.0-1.0). If 0, uses provider default.
	Temperature float64 `json:"temperature,omitempty"`

	// OutputSchema defines a JSON schema for structured output responses.
	// When set, the response Content will be valid JSON matching this schema.
	// Uses the provider's native structured output feature (e.g., Anthropic's
	// constrained decoding). If nil, the response is unstructured text.
	OutputSchema map[string]any `json:"outputSchema,omitempty"`
}

// LLMProvider defines the interface for LLM service providers.
// Implementations handle provider-specific details like authentication,
// request formatting, and response parsing.
type LLMProvider interface {
	// Complete sends a request to the LLM and returns the response.
	// It handles retries internally based on the provider's configuration.
	Complete(ctx context.Context, req LLMRequest) (*LLMResponse, error)

	// Name returns the provider's identifier (e.g., "anthropic", "openai").
	Name() string
}

// LLMError represents an error from an LLM provider with additional context.
type LLMError struct { //nolint:govet // Field order prioritizes readability
	// Provider is the name of the provider that returned the error.
	Provider string `json:"provider"`

	// Code is the error code from the provider (if available).
	Code string `json:"code,omitempty"`

	// Message is the human-readable error message.
	Message string `json:"message"`

	// Retryable indicates whether the request can be retried.
	Retryable bool `json:"retryable"`

	// Err is the underlying error.
	Err error `json:"-"`
}

// Error implements the error interface.
func (e *LLMError) Error() string {
	if e.Code != "" {
		return e.Provider + ": " + e.Code + ": " + e.Message
	}
	return e.Provider + ": " + e.Message
}

// Unwrap returns the underlying error.
func (e *LLMError) Unwrap() error {
	return e.Err
}

// DocumentExtractor defines the interface for extracting data from documents using LLM.
type DocumentExtractor interface {
	// ExtractText extracts raw text from a document (image or PDF).
	ExtractText(ctx context.Context, document []byte, contentType string) (string, error)

	// ExtractResumeData extracts structured resume data from text using LLM.
	ExtractResumeData(ctx context.Context, text string) (*ResumeExtractedData, error)

	// ExtractLetterData extracts structured reference letter data from text using LLM.
	ExtractLetterData(ctx context.Context, text string) (*ExtractedLetterData, error)
}
