//nolint:errcheck,revive,goconst // Test file - error checks, unused params, and string constants are OK
package llm_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
	"backend/internal/logger"
)

// mockLogger captures log entries for testing.
type mockLogger struct {
	mu      sync.Mutex
	entries []logger.LogEntry
}

func (m *mockLogger) log(severity logger.Severity, message string, attrs ...logger.Attr) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, logger.LogEntry{
		Severity: severity,
		Message:  message,
		Attrs:    attrs,
	})
}

func (m *mockLogger) Debug(message string, attrs ...logger.Attr)    { m.log(logger.Debug, message, attrs...) }
func (m *mockLogger) Info(message string, attrs ...logger.Attr)     { m.log(logger.Info, message, attrs...) }
func (m *mockLogger) Warning(message string, attrs ...logger.Attr)  { m.log(logger.Warning, message, attrs...) }
func (m *mockLogger) Error(message string, attrs ...logger.Attr)    { m.log(logger.Error, message, attrs...) }
func (m *mockLogger) Critical(message string, attrs ...logger.Attr) { m.log(logger.Critical, message, attrs...) }

func (m *mockLogger) getEntries() []logger.LogEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]logger.LogEntry, len(m.entries))
	copy(result, m.entries)
	return result
}

// mockProvider is a simple LLM provider for testing.
type mockProvider struct {
	response *domain.LLMResponse
	err      error
}

func (m *mockProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	return m.response, m.err
}

func (m *mockProvider) Name() string {
	return "mock"
}

func TestLoggingProvider_LogsSuccessfulRequest(t *testing.T) {
	log := &mockLogger{}
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "Hello!",
			Model:        "test-model",
			InputTokens:  10,
			OutputTokens: 5,
			StopReason:   "end_turn",
		},
	}

	provider := llm.NewLoggingProvider(inner, log)

	resp, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hi"),
		},
		MaxTokens: 100,
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if resp.Content != "Hello!" {
		t.Errorf("Content = %q, want %q", resp.Content, "Hello!")
	}

	entries := log.getEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(entries))
	}

	// First entry: request
	if entries[0].Severity != logger.Debug {
		t.Errorf("request log severity = %v, want Debug", entries[0].Severity)
	}
	if entries[0].Message != "LLM request" {
		t.Errorf("request log message = %q, want 'LLM request'", entries[0].Message)
	}

	// Second entry: response
	if entries[1].Severity != logger.Debug {
		t.Errorf("response log severity = %v, want Debug", entries[1].Severity)
	}
	if entries[1].Message != "LLM response" {
		t.Errorf("response log message = %q, want 'LLM response'", entries[1].Message)
	}
}

func TestLoggingProvider_LogsError(t *testing.T) {
	log := &mockLogger{}
	testErr := errors.New("API error")
	inner := &mockProvider{
		err: testErr,
	}

	provider := llm.NewLoggingProvider(inner, log)

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hi"),
		},
	})

	if err == nil {
		t.Fatal("expected error")
	}

	entries := log.getEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(entries))
	}

	// Second entry should be error
	if entries[1].Severity != logger.Error {
		t.Errorf("error log severity = %v, want Error", entries[1].Severity)
	}
	if entries[1].Message != "LLM request failed" {
		t.Errorf("error log message = %q, want 'LLM request failed'", entries[1].Message)
	}
}

func TestLoggingProvider_Name(t *testing.T) {
	log := &mockLogger{}
	inner := &mockProvider{}

	provider := llm.NewLoggingProvider(inner, log)

	if got := provider.Name(); got != "mock" {
		t.Errorf("Name() = %q, want %q", got, "mock")
	}
}
