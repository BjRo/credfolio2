package llm

import (
	"context"
	"testing"

	"backend/internal/logger"
)

func TestNewBraintrustTracing_WithEmptyAPIKey(t *testing.T) {
	log := logger.NewStdoutLogger()

	bt, err := NewBraintrustTracing(BraintrustConfig{
		APIKey:  "",
		Project: "test-project",
	}, log)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if bt != nil {
		t.Fatal("expected nil BraintrustTracing when API key is empty")
	}
}

func TestBraintrustTracing_Enabled(t *testing.T) {
	tests := []struct {
		name     string
		bt       *BraintrustTracing
		expected bool
	}{
		{
			name:     "nil tracing returns false",
			bt:       nil,
			expected: false,
		},
		{
			name:     "tracing without client returns false",
			bt:       &BraintrustTracing{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.Enabled(); got != tt.expected {
				t.Errorf("Enabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBraintrustTracing_Project(t *testing.T) {
	tests := []struct {
		name     string
		bt       *BraintrustTracing
		expected string
	}{
		{
			name:     "nil tracing returns empty",
			bt:       nil,
			expected: "",
		},
		{
			name: "returns configured project",
			bt: &BraintrustTracing{
				config: BraintrustConfig{Project: "my-project"},
			},
			expected: "my-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.Project(); got != tt.expected {
				t.Errorf("Project() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBraintrustTracing_Shutdown_NilSafe(t *testing.T) {
	var bt *BraintrustTracing

	err := bt.Shutdown(context.Background())
	if err != nil {
		t.Errorf("Shutdown on nil should not error, got %v", err)
	}
}

//nolint:bodyclose // Test for middleware function, not HTTP response
func TestBraintrustTracing_AnthropicMiddleware_NilSafe(t *testing.T) {
	var bt *BraintrustTracing

	mw := bt.AnthropicMiddleware()
	if mw != nil {
		t.Error("AnthropicMiddleware on nil should return nil")
	}
}

//nolint:bodyclose // Test for middleware function, not HTTP response
func TestBraintrustTracing_OpenAIMiddleware_NilSafe(t *testing.T) {
	var bt *BraintrustTracing

	mw := bt.OpenAIMiddleware()
	if mw != nil {
		t.Error("OpenAIMiddleware on nil should return nil")
	}
}
