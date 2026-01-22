package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBgRed  = "\033[41m"
	colorWhite  = "\033[97m"
	colorCyan   = "\033[36m"
	colorDim    = "\033[2m"
)

// StdoutLogger implements Logger by writing to stdout with colored output.
type StdoutLogger struct {
	out      io.Writer
	mu       sync.Mutex
	minLevel Severity
}

// StdoutOption is a functional option for configuring StdoutLogger.
type StdoutOption func(*StdoutLogger)

// WithMinLevel sets the minimum severity level to log.
// Messages below this level will be ignored.
func WithMinLevel(level Severity) StdoutOption {
	return func(l *StdoutLogger) {
		l.minLevel = level
	}
}

// WithOutput sets a custom output writer (useful for testing).
func WithOutput(w io.Writer) StdoutOption {
	return func(l *StdoutLogger) {
		l.out = w
	}
}

// NewStdoutLogger creates a new logger that writes to stdout.
func NewStdoutLogger(opts ...StdoutOption) *StdoutLogger {
	l := &StdoutLogger{
		out:      os.Stdout,
		minLevel: Debug, // log everything by default
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// log emits a log entry with the given severity and message.
func (l *StdoutLogger) log(severity Severity, message string, attrs []Attr) {
	if severity < l.minLevel {
		return
	}

	entry := &LogEntry{
		Severity:  severity,
		Message:   message,
		Attrs:     attrs,
		Timestamp: time.Now(),
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprintln(l.out, l.format(entry)) //nolint:errcheck // Best effort logging
}

// Debug logs a debug-level message.
func (l *StdoutLogger) Debug(message string, attrs ...Attr) {
	l.log(Debug, message, attrs)
}

// Info logs an info-level message.
func (l *StdoutLogger) Info(message string, attrs ...Attr) {
	l.log(Info, message, attrs)
}

// Warning logs a warning-level message.
func (l *StdoutLogger) Warning(message string, attrs ...Attr) {
	l.log(Warning, message, attrs)
}

// Error logs an error-level message.
func (l *StdoutLogger) Error(message string, attrs ...Attr) {
	l.log(Error, message, attrs)
}

// Critical logs a critical-level message.
func (l *StdoutLogger) Critical(message string, attrs ...Attr) {
	l.log(Critical, message, attrs)
}

// format creates a colored, human-readable log line.
func (l *StdoutLogger) format(entry *LogEntry) string {
	var b strings.Builder

	// Timestamp in dim gray
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
	b.WriteString(colorDim)
	b.WriteString(timestamp)
	b.WriteString(colorReset)
	b.WriteString(" ")

	// Severity with color
	b.WriteString(l.severityColor(entry.Severity))
	b.WriteString(fmt.Sprintf("[%-8s]", entry.Severity.String()))
	b.WriteString(colorReset)

	// Extract feature and data attrs
	var feature string
	data := make(map[string]any)
	for _, attr := range entry.Attrs {
		if attr.Key == featureKey {
			if s, ok := attr.Value.(string); ok {
				feature = s
			}
		} else {
			data[attr.Key] = attr.Value
		}
	}

	// Feature tag in cyan (if present)
	if feature != "" {
		b.WriteString(" ")
		b.WriteString(colorCyan)
		b.WriteString("[")
		b.WriteString(feature)
		b.WriteString("]")
		b.WriteString(colorReset)
	}

	// Message
	b.WriteString(" ")
	b.WriteString(entry.Message)

	// Data as JSON (if present)
	if len(data) > 0 {
		b.WriteString(" ")
		b.WriteString(colorDim)
		dataJSON, err := json.Marshal(data)
		if err != nil {
			b.WriteString("{\"_error\": \"failed to marshal data\"}")
		} else {
			b.Write(dataJSON)
		}
		b.WriteString(colorReset)
	}

	return b.String()
}

// severityColor returns the ANSI color code for a severity level.
func (l *StdoutLogger) severityColor(s Severity) string {
	switch s {
	case Debug:
		return colorGray
	case Info:
		return colorGreen
	case Warning:
		return colorYellow
	case Error:
		return colorRed
	case Critical:
		return colorBgRed + colorWhite
	default:
		return colorReset
	}
}
