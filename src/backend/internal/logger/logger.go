// Package logger provides a logging abstraction for the application.
// Currently logs to stdout, but designed for future extensibility
// to forward logs to services like Datadog, Sentry, etc.
package logger

import "time"

// Severity represents the importance level of a log entry.
type Severity int

// Severity levels for log entries.
const (
	Debug Severity = iota
	Info
	Warning
	Error
	Critical
)

// String returns the string representation of a severity level.
func (s Severity) String() string {
	switch s {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Critical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry holds all data for a single log entry.
//
//nolint:govet // Field order optimized for readability over memory alignment
type LogEntry struct {
	Severity  Severity
	Message   string
	Feature   string         // optional: categorizes the log by feature area
	Data      map[string]any // optional: structured context data
	Timestamp time.Time
}

// LogOption is a functional option for configuring log entries.
type LogOption func(*LogEntry)

// WithFeature adds a feature tag to the log entry.
func WithFeature(feature string) LogOption {
	return func(e *LogEntry) {
		e.Feature = feature
	}
}

// WithData adds structured data to the log entry.
// The data should be JSON-compatible (strings, numbers, booleans, maps, slices).
func WithData(data map[string]any) LogOption {
	return func(e *LogEntry) {
		e.Data = data
	}
}

// Logger is the main interface for logging throughout the application.
type Logger interface {
	// Log emits a log entry with the given severity and message.
	Log(severity Severity, message string, opts ...LogOption)

	// Debug logs a debug-level message.
	Debug(message string, opts ...LogOption)

	// Info logs an info-level message.
	Info(message string, opts ...LogOption)

	// Warning logs a warning-level message.
	Warning(message string, opts ...LogOption)

	// Error logs an error-level message.
	Error(message string, opts ...LogOption)

	// Critical logs a critical-level message.
	Critical(message string, opts ...LogOption)
}
