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

const testOpenAIModel = "gpt-4o"

func TestOpenAIProvider_Name(t *testing.T) {
	provider := llm.NewOpenAIProvider(llm.OpenAIConfig{
		APIKey: "test-key",
	})

	if got := provider.Name(); got != "openai" {
		t.Errorf("Name() = %q, want %q", got, "openai")
	}
}

func TestOpenAIProvider_Complete_UsesMaxCompletionTokens(t *testing.T) {
	var receivedReq map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedReq)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"model": "gpt-4o",
			"choices": [{
				"index": 0,
				"message": {"role": "assistant", "content": "Hello!"},
				"finish_reason": "stop"
			}],
			"usage": {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
		}`))
	}))
	defer server.Close()

	provider := llm.NewOpenAIProvider(llm.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL + "/v1",
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

	// Must use max_completion_tokens (required by newer models like gpt-5-nano, o1, o3)
	if _, ok := receivedReq["max_completion_tokens"]; !ok {
		t.Error("request should include max_completion_tokens")
	}
	if val, ok := receivedReq["max_completion_tokens"]; ok {
		if val != float64(1024) {
			t.Errorf("max_completion_tokens = %v, want 1024", val)
		}
	}

	// Must NOT use deprecated max_tokens
	if _, ok := receivedReq["max_tokens"]; ok {
		t.Errorf("request should NOT include deprecated max_tokens, got %v", receivedReq["max_tokens"])
	}
}

func TestOpenAIProvider_Complete_TextMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"model": "gpt-4o",
			"choices": [{
				"index": 0,
				"message": {"role": "assistant", "content": "Hello! How can I help?"},
				"finish_reason": "stop"
			}],
			"usage": {"prompt_tokens": 10, "completion_tokens": 8, "total_tokens": 18}
		}`))
	}))
	defer server.Close()

	provider := llm.NewOpenAIProvider(llm.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL + "/v1",
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

	if resp.Content != "Hello! How can I help?" {
		t.Errorf("Content = %q, want %q", resp.Content, "Hello! How can I help?")
	}
	if resp.Model != testOpenAIModel {
		t.Errorf("Model = %q, want %q", resp.Model, testOpenAIModel)
	}
	if resp.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want %d", resp.InputTokens, 10)
	}
	if resp.OutputTokens != 8 {
		t.Errorf("OutputTokens = %d, want %d", resp.OutputTokens, 8)
	}
	if resp.StopReason != "stop" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "stop")
	}
}
