package llm

import (
	"context"
	"errors"
	"testing"

	btlogger "github.com/braintrustdata/braintrust-sdk-go/logger"

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

func TestOtelErrorHandler_Handle(t *testing.T) {
	log := &capturingLogger{}
	handler := &otelErrorHandler{log: log}

	handler.Handle(errors.New("OTLP export failed: connection refused"))

	if len(log.errors) != 1 {
		t.Fatalf("expected 1 error log, got %d", len(log.errors))
	}
	if log.errors[0] != "OpenTelemetry error (possible Braintrust export failure)" {
		t.Errorf("unexpected error message: %s", log.errors[0])
	}
}

func TestBraintrustLogAdapter_ImplementsInterface(t *testing.T) {
	t.Helper()
	log := &capturingLogger{}

	// Verify it satisfies the Braintrust Logger interface and calls work
	var adapter btlogger.Logger = &braintrustLogAdapter{log: log}
	adapter.Debug("compile check")
	if len(log.debugs) != 1 {
		t.Error("expected adapter to forward Debug call")
	}
}

func TestBraintrustLogAdapter_ForwardsMessages(t *testing.T) {
	log := &capturingLogger{}
	adapter := &braintrustLogAdapter{log: log}

	adapter.Debug("test debug", "key", "value")
	adapter.Info("test info", "key", "value")
	adapter.Warn("test warn", "key", "value")
	adapter.Error("test error", "key", "value")

	if len(log.debugs) != 1 {
		t.Errorf("expected 1 debug log, got %d", len(log.debugs))
	}
	if len(log.infos) != 1 {
		t.Errorf("expected 1 info log, got %d", len(log.infos))
	}
	if len(log.warnings) != 1 {
		t.Errorf("expected 1 warning log, got %d", len(log.warnings))
	}
	if len(log.errors) != 1 {
		t.Errorf("expected 1 error log, got %d", len(log.errors))
	}
}

// capturingLogger captures log messages for test assertions.
type capturingLogger struct {
	debugs   []string
	infos    []string
	warnings []string
	errors   []string
}

func (l *capturingLogger) Debug(msg string, _ ...logger.Attr)    { l.debugs = append(l.debugs, msg) }
func (l *capturingLogger) Info(msg string, _ ...logger.Attr)     { l.infos = append(l.infos, msg) }
func (l *capturingLogger) Warning(msg string, _ ...logger.Attr)  { l.warnings = append(l.warnings, msg) }
func (l *capturingLogger) Error(msg string, _ ...logger.Attr)    { l.errors = append(l.errors, msg) }
func (l *capturingLogger) Critical(_ string, _ ...logger.Attr) {}
