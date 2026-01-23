//nolint:errcheck,revive // Test file - error checks and unused params are OK in test helpers
package llm_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
)

func TestAnthropicProvider_Name(t *testing.T) {
	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey: "test-key",
	})

	if got := provider.Name(); got != "anthropic" {
		t.Errorf("Name() = %q, want %q", got, "anthropic")
	}
}

func TestAnthropicProvider_Complete_TextMessage(t *testing.T) {
	// Create a mock server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("x-api-key") != "test-api-key" {
			t.Errorf("x-api-key header = %q, want %q", r.Header.Get("x-api-key"), "test-api-key")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("anthropic-version header = %q, want %q", r.Header.Get("anthropic-version"), "2023-06-01")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type header = %q, want %q", r.Header.Get("Content-Type"), "application/json")
		}

		// Parse and verify request body
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		// Verify model
		if req["model"] != "claude-sonnet-4-20250514" {
			t.Errorf("model = %v, want claude-sonnet-4-20250514", req["model"])
		}

		// Verify max_tokens
		if req["max_tokens"] != float64(1024) {
			t.Errorf("max_tokens = %v, want 1024", req["max_tokens"])
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Hello! How can I help you today?"}],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {
				"input_tokens": 10,
				"output_tokens": 15
			}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
	})

	resp, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens: 1024,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if resp.Content != "Hello! How can I help you today?" {
		t.Errorf("Content = %q, want %q", resp.Content, "Hello! How can I help you today?")
	}
	if resp.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Model = %q, want %q", resp.Model, "claude-sonnet-4-20250514")
	}
	if resp.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want %d", resp.InputTokens, 10)
	}
	if resp.OutputTokens != 15 {
		t.Errorf("OutputTokens = %d, want %d", resp.OutputTokens, 15)
	}
	if resp.StopReason != "end_turn" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "end_turn")
	}
}

func TestAnthropicProvider_Complete_WithSystemPrompt(t *testing.T) {
	var receivedReq map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedReq)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Response"}],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 20, "output_tokens": 5}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		SystemPrompt: "You are a helpful assistant.",
		MaxTokens:    1024,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	// Verify system prompt was sent (SDK sends it as an array of text blocks)
	systemBlocks, ok := receivedReq["system"].([]any)
	if !ok || len(systemBlocks) == 0 {
		t.Fatalf("expected system to be an array of text blocks, got %v", receivedReq["system"])
	}
	firstBlock, ok := systemBlocks[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first system block to be a map, got %T", systemBlocks[0])
	}
	if firstBlock["text"] != "You are a helpful assistant." {
		t.Errorf("system[0].text = %v, want %q", firstBlock["text"], "You are a helpful assistant.")
	}
}

func TestAnthropicProvider_Complete_WithImage(t *testing.T) {
	var receivedReq map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedReq)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "I see an image"}],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 100, "output_tokens": 10}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	imageData := []byte{0xFF, 0xD8, 0xFF} // Fake JPEG header
	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewImageMessage(domain.RoleUser, domain.ImageMediaTypeJPEG, imageData, "What's in this image?"),
		},
		MaxTokens: 1024,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	// Verify image was sent correctly
	messages, ok := receivedReq["messages"].([]any)
	if !ok || len(messages) == 0 {
		t.Fatal("expected messages in request")
	}
	msg := messages[0].(map[string]any)
	content := msg["content"].([]any)
	if len(content) != 2 {
		t.Fatalf("expected 2 content blocks (image + text), got %d", len(content))
	}

	// First block should be the image
	imgBlock := content[0].(map[string]any)
	if imgBlock["type"] != "image" {
		t.Errorf("first block type = %v, want image", imgBlock["type"])
	}
	source := imgBlock["source"].(map[string]any)
	if source["type"] != "base64" {
		t.Errorf("source type = %v, want base64", source["type"])
	}
	if source["media_type"] != "image/jpeg" {
		t.Errorf("media_type = %v, want image/jpeg", source["media_type"])
	}
}

func TestAnthropicProvider_Complete_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{
			"type": "error",
			"error": {
				"type": "invalid_request_error",
				"message": "Invalid API key"
			}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "invalid-key",
		BaseURL: server.URL,
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens: 1024,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify it's an LLMError
	llmErr, ok := err.(*domain.LLMError)
	if !ok {
		t.Fatalf("expected *domain.LLMError, got %T", err)
	}

	if llmErr.Provider != "anthropic" {
		t.Errorf("Provider = %q, want %q", llmErr.Provider, "anthropic")
	}
	// Note: Code is not extracted from SDK errors as they don't expose it in a structured way
	// The error message should contain details from the API response
	if llmErr.Message == "" {
		t.Error("expected error message to be non-empty")
	}
}

func TestAnthropicProvider_Complete_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{
			"type": "error",
			"error": {
				"type": "rate_limit_error",
				"message": "Rate limit exceeded"
			}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens: 1024,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	llmErr, ok := err.(*domain.LLMError)
	if !ok {
		t.Fatalf("expected *domain.LLMError, got %T", err)
	}

	if !llmErr.Retryable {
		t.Error("rate limit error should be retryable")
	}
}

func TestAnthropicProvider_Complete_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{
			"type": "error",
			"error": {
				"type": "api_error",
				"message": "Internal server error"
			}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens: 1024,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	llmErr, ok := err.(*domain.LLMError)
	if !ok {
		t.Fatalf("expected *domain.LLMError, got %T", err)
	}

	if !llmErr.Retryable {
		t.Error("server error should be retryable")
	}
}

func TestAnthropicProvider_Complete_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response that will be cancelled
		<-r.Context().Done()
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := provider.Complete(ctx, domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens: 1024,
	})

	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestAnthropicProvider_Complete_DefaultModel(t *testing.T) {
	var receivedReq map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedReq)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Response"}],
			"model": "claude-3-haiku-20240307",
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 10, "output_tokens": 5}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:       "test-key",
		BaseURL:      server.URL,
		DefaultModel: "claude-3-haiku-20240307",
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens: 1024,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if receivedReq["model"] != "claude-3-haiku-20240307" {
		t.Errorf("model = %v, want claude-3-haiku-20240307", receivedReq["model"])
	}
}

func TestAnthropicProvider_Complete_PerRequestModel(t *testing.T) {
	var receivedReq map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedReq)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Response"}],
			"model": "claude-3-opus-20240229",
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 10, "output_tokens": 5}
		}`))
	}))
	defer server.Close()

	// Provider configured with haiku as default
	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:       "test-key",
		BaseURL:      server.URL,
		DefaultModel: "claude-3-haiku-20240307",
	})

	// But request specifies opus
	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		Model:     "claude-3-opus-20240229",
		MaxTokens: 1024,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	// Per-request model should override default
	if receivedReq["model"] != "claude-3-opus-20240229" {
		t.Errorf("model = %v, want claude-3-opus-20240229", receivedReq["model"])
	}
}

func TestAnthropicProvider_Complete_Temperature(t *testing.T) {
	var receivedReq map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedReq)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"content": [{"type": "text", "text": "Response"}],
			"model": "claude-sonnet-4-20250514",
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 10, "output_tokens": 5}
		}`))
	}))
	defer server.Close()

	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
		MaxTokens:   1024,
		Temperature: 0.7,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if receivedReq["temperature"] != 0.7 {
		t.Errorf("temperature = %v, want 0.7", receivedReq["temperature"])
	}
}
