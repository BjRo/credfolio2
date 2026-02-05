// Package llm provides LLM provider implementations.
package llm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/braintrustdata/braintrust-sdk-go"
	btlogger "github.com/braintrustdata/braintrust-sdk-go/logger"
	traceanthropic "github.com/braintrustdata/braintrust-sdk-go/trace/contrib/anthropic"
	traceopenai "github.com/braintrustdata/braintrust-sdk-go/trace/contrib/openai"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"backend/internal/logger"
)

// BraintrustConfig holds configuration for Braintrust tracing.
type BraintrustConfig struct {
	// APIKey is the Braintrust API key. Required for tracing to be enabled.
	APIKey string
	// Project is the Braintrust project name for organizing traces.
	Project string
}

// BraintrustTracing manages Braintrust tracing infrastructure.
type BraintrustTracing struct {
	client         *braintrust.Client
	tracerProvider *trace.TracerProvider
	log            logger.Logger
	config         BraintrustConfig
}

// NewBraintrustTracing creates a new Braintrust tracing instance.
// Returns nil if APIKey is empty (tracing disabled).
func NewBraintrustTracing(cfg BraintrustConfig, log logger.Logger) (*BraintrustTracing, error) {
	if cfg.APIKey == "" {
		return nil, nil
	}

	// Set up OpenTelemetry error handler BEFORE creating TracerProvider.
	// Without this, OTLP export failures (network errors, auth failures) are
	// silently swallowed by the batch span processor — the most common reason
	// for "no data in Braintrust."
	otel.SetErrorHandler(&otelErrorHandler{log: log})

	// Create TracerProvider
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	// Create Braintrust client with API key and SDK logger for visibility
	// into login, export, and span processing activity.
	opts := []braintrust.Option{
		braintrust.WithAPIKey(cfg.APIKey),
		braintrust.WithLogger(&braintrustLogAdapter{log: log}),
	}
	if cfg.Project != "" {
		opts = append(opts, braintrust.WithProject(cfg.Project))
	}

	client, err := braintrust.New(tp, opts...)
	if err != nil {
		// Clean up tracer provider on error (best effort)
		if shutdownErr := tp.Shutdown(context.Background()); shutdownErr != nil {
			// Log both errors - the original error is still returned
			_ = shutdownErr // best effort cleanup
		}
		return nil, err
	}

	return &BraintrustTracing{
		client:         client,
		tracerProvider: tp,
		config:         cfg,
		log:            log,
	}, nil
}

// AnthropicMiddleware returns the tracing middleware for Anthropic clients.
// Returns nil if tracing is not configured.
//
//nolint:bodyclose // This returns middleware that passes through responses, not a response to close
func (bt *BraintrustTracing) AnthropicMiddleware() func(*http.Request, traceanthropic.NextMiddleware) (*http.Response, error) {
	if bt == nil {
		return nil
	}
	return traceanthropic.NewMiddleware(
		traceanthropic.WithTracerProvider(bt.tracerProvider),
	)
}

// OpenAIMiddleware returns the tracing middleware for OpenAI clients.
// Returns nil if tracing is not configured.
//
//nolint:bodyclose // This returns middleware that passes through responses, not a response to close
func (bt *BraintrustTracing) OpenAIMiddleware() func(*http.Request, traceopenai.NextMiddleware) (*http.Response, error) {
	if bt == nil {
		return nil
	}
	return traceopenai.NewMiddleware(
		traceopenai.WithTracerProvider(bt.tracerProvider),
	)
}

// Shutdown gracefully shuts down the Braintrust tracing infrastructure.
func (bt *BraintrustTracing) Shutdown(ctx context.Context) error {
	if bt == nil || bt.tracerProvider == nil {
		return nil
	}
	return bt.tracerProvider.Shutdown(ctx)
}

// Enabled returns true if Braintrust tracing is configured and active.
func (bt *BraintrustTracing) Enabled() bool {
	return bt != nil && bt.client != nil
}

// Project returns the configured project name.
func (bt *BraintrustTracing) Project() string {
	if bt == nil {
		return ""
	}
	return bt.config.Project
}

// otelErrorHandler routes OpenTelemetry internal errors (e.g. OTLP export
// failures) to our application logger so they are visible in server logs.
type otelErrorHandler struct {
	log logger.Logger
}

func (h *otelErrorHandler) Handle(err error) {
	h.log.Error("OpenTelemetry error (possible Braintrust export failure)",
		logger.Feature("braintrust"),
		logger.Err(err),
	)
}

// braintrustLogAdapter adapts our logger.Logger to the Braintrust SDK's
// logger.Logger interface so SDK-internal messages (login, tracing setup,
// span processing) are visible in server logs.
type braintrustLogAdapter struct {
	log logger.Logger
}

var _ btlogger.Logger = (*braintrustLogAdapter)(nil)

func (a *braintrustLogAdapter) Debug(msg string, args ...any) {
	a.log.Debug(fmt.Sprintf("[braintrust] %s %v", msg, args), logger.Feature("braintrust"))
}

func (a *braintrustLogAdapter) Info(msg string, args ...any) {
	a.log.Info(fmt.Sprintf("[braintrust] %s %v", msg, args), logger.Feature("braintrust"))
}

func (a *braintrustLogAdapter) Warn(msg string, args ...any) {
	a.log.Warning(fmt.Sprintf("[braintrust] %s %v", msg, args), logger.Feature("braintrust"))
}

func (a *braintrustLogAdapter) Error(msg string, args ...any) {
	a.log.Error(fmt.Sprintf("[braintrust] %s %v", msg, args), logger.Feature("braintrust"))
}
