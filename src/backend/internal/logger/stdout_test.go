package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestStdoutLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdoutLogger(WithOutput(&buf))

	log.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected output to contain 'INFO', got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestStdoutLogger_SeverityLevels(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(*StdoutLogger, string, ...Attr)
		expected string
	}{
		{"Debug", (*StdoutLogger).Debug, "DEBUG"},
		{"Info", (*StdoutLogger).Info, "INFO"},
		{"Warning", (*StdoutLogger).Warning, "WARNING"},
		{"Error", (*StdoutLogger).Error, "ERROR"},
		{"Critical", (*StdoutLogger).Critical, "CRITICAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := NewStdoutLogger(WithOutput(&buf))

			tt.logFunc(log, "test")

			if !strings.Contains(buf.String(), tt.expected) {
				t.Errorf("expected output to contain %q, got: %s", tt.expected, buf.String())
			}
		})
	}
}

func TestStdoutLogger_Feature(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdoutLogger(WithOutput(&buf))

	log.Info("test message", Feature("auth"))

	output := buf.String()
	if !strings.Contains(output, "[auth]") {
		t.Errorf("expected output to contain '[auth]', got: %s", output)
	}
}

func TestStdoutLogger_TypedAttrs(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdoutLogger(WithOutput(&buf))

	log.Info("test message", String("user_id", "123"), String("action", "login"))

	output := buf.String()
	if !strings.Contains(output, "user_id") {
		t.Errorf("expected output to contain 'user_id', got: %s", output)
	}
	if !strings.Contains(output, "123") {
		t.Errorf("expected output to contain '123', got: %s", output)
	}
}

func TestStdoutLogger_WithMinLevel(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdoutLogger(WithOutput(&buf), WithMinLevel(Warning))

	log.Debug("debug message")
	log.Info("info message")
	log.Warning("warning message")
	log.Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("debug message should have been filtered")
	}
	if strings.Contains(output, "info message") {
		t.Error("info message should have been filtered")
	}
	if !strings.Contains(output, "warning message") {
		t.Error("warning message should be present")
	}
	if !strings.Contains(output, "error message") {
		t.Error("error message should be present")
	}
}

func TestStdoutLogger_NestedData(t *testing.T) {
	var buf bytes.Buffer
	log := NewStdoutLogger(WithOutput(&buf))

	log.Info("test message",
		Any("user", map[string]any{
			"id":    "123",
			"email": "test@example.com",
		}),
		Any("metadata", map[string]any{
			"ip":      "192.168.1.1",
			"browser": "Chrome",
		}),
	)

	output := buf.String()
	if !strings.Contains(output, "test@example.com") {
		t.Errorf("expected output to contain nested data, got: %s", output)
	}
}

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		severity Severity
		expected string
	}{
		{Debug, "DEBUG"},
		{Info, "INFO"},
		{Warning, "WARNING"},
		{Error, "ERROR"},
		{Critical, "CRITICAL"},
		{Severity(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.severity.String(); got != tt.expected {
				t.Errorf("Severity.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
