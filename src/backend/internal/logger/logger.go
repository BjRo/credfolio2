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

// Attr represents a typed key-value attribute for structured logging.
//
//nolint:govet // Field order optimized for readability over memory alignment
type Attr struct {
	Key   string
	Value any
}

// featureKey is the special key used to identify feature attributes.
const featureKey = "_feature"

// String creates a string attribute.
func String(key, value string) Attr {
	return Attr{Key: key, Value: value}
}

// Int creates an integer attribute.
func Int(key string, value int) Attr {
	return Attr{Key: key, Value: value}
}

// Int64 creates an int64 attribute.
func Int64(key string, value int64) Attr {
	return Attr{Key: key, Value: value}
}

// Float64 creates a float64 attribute.
func Float64(key string, value float64) Attr {
	return Attr{Key: key, Value: value}
}

// Bool creates a boolean attribute.
func Bool(key string, value bool) Attr {
	return Attr{Key: key, Value: value}
}

// Any creates an attribute with any value.
func Any(key string, value any) Attr {
	return Attr{Key: key, Value: value}
}

// Err creates an error attribute with the key "error".
func Err(err error) Attr {
	if err == nil {
		return Attr{Key: "error", Value: nil}
	}
	return Attr{Key: "error", Value: err.Error()}
}

// Feature creates a feature tag attribute for categorizing logs.
func Feature(name string) Attr {
	return Attr{Key: featureKey, Value: name}
}

// LogEntry holds all data for a single log entry.
//
//nolint:govet // Field order optimized for readability over memory alignment
type LogEntry struct {
	Severity  Severity
	Message   string
	Attrs     []Attr
	Timestamp time.Time
}

// Logger is the main interface for logging throughout the application.
type Logger interface {
	// Debug logs a debug-level message.
	Debug(message string, attrs ...Attr)

	// Info logs an info-level message.
	Info(message string, attrs ...Attr)

	// Warning logs a warning-level message.
	Warning(message string, attrs ...Attr)

	// Error logs an error-level message.
	Error(message string, attrs ...Attr)

	// Critical logs a critical-level message.
	Critical(message string, attrs ...Attr)
}
