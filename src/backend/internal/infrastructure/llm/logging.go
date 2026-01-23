package llm

import (
	"context"
	"time"

	"backend/internal/domain"
	"backend/internal/logger"
)

// LoggingProvider wraps an LLM provider to add request/response logging.
type LoggingProvider struct {
	inner  domain.LLMProvider
	logger logger.Logger
}

// NewLoggingProvider creates a new logging provider wrapper.
func NewLoggingProvider(inner domain.LLMProvider, log logger.Logger) *LoggingProvider {
	return &LoggingProvider{
		inner:  inner,
		logger: log,
	}
}

// Complete logs the request and response, then delegates to the inner provider.
func (p *LoggingProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	start := time.Now()

	// Log request
	p.logger.Debug("LLM request",
		logger.Feature("llm"),
		logger.String("provider", p.inner.Name()),
		logger.String("model", req.Model),
		logger.Int("message_count", len(req.Messages)),
		logger.Int("max_tokens", req.MaxTokens),
		logger.Bool("has_system_prompt", req.SystemPrompt != ""),
	)

	// Execute request
	resp, err := p.inner.Complete(ctx, req)
	duration := time.Since(start)

	if err != nil {
		// Log error
		p.logger.Error("LLM request failed",
			logger.Feature("llm"),
			logger.String("provider", p.inner.Name()),
			logger.String("model", req.Model),
			logger.Int64("duration_ms", duration.Milliseconds()),
			logger.Err(err),
		)
		return nil, err
	}

	// Log successful response
	p.logger.Debug("LLM response",
		logger.Feature("llm"),
		logger.String("provider", p.inner.Name()),
		logger.String("model", resp.Model),
		logger.Int("input_tokens", resp.InputTokens),
		logger.Int("output_tokens", resp.OutputTokens),
		logger.String("stop_reason", resp.StopReason),
		logger.Int64("duration_ms", duration.Milliseconds()),
	)

	return resp, nil
}

// Name returns the wrapped provider's name.
func (p *LoggingProvider) Name() string {
	return p.inner.Name()
}

// Verify LoggingProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*LoggingProvider)(nil)
