// Package llm provides LLM provider implementations.
package llm

import (
	"context"
	"net/http"

	"github.com/braintrustdata/braintrust-sdk-go"
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

	// Create TracerProvider
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	// Create Braintrust client with API key
	opts := []braintrust.Option{
		braintrust.WithAPIKey(cfg.APIKey),
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
