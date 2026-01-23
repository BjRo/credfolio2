package llm

import (
	"context"
	"time"

	"backend/internal/domain"
)

// ResilientConfig holds configuration for the resilient provider.
type ResilientConfig struct {
	RetryConfig          RetrierConfig
	CircuitBreakerConfig CircuitBreakerConfig
}

// ResilientProvider wraps an LLM provider with retry and circuit breaker.
type ResilientProvider struct {
	inner          domain.LLMProvider
	retrier        *Retrier
	circuitBreaker *CircuitBreaker
}

// NewResilientProvider creates a resilient provider with retry and circuit breaker.
func NewResilientProvider(inner domain.LLMProvider, config ResilientConfig) *ResilientProvider {
	// Apply defaults if not set
	if config.RetryConfig.MaxAttempts == 0 {
		config.RetryConfig.MaxAttempts = 3
	}
	if config.RetryConfig.BaseDelay == 0 {
		config.RetryConfig.BaseDelay = 500 * time.Millisecond
	}
	if config.RetryConfig.MaxDelay == 0 {
		config.RetryConfig.MaxDelay = 30 * time.Second
	}
	if config.CircuitBreakerConfig.FailureThreshold == 0 {
		config.CircuitBreakerConfig.FailureThreshold = 5
	}
	if config.CircuitBreakerConfig.ResetTimeout == 0 {
		config.CircuitBreakerConfig.ResetTimeout = 60 * time.Second
	}

	// Configure retry to only retry retryable errors
	config.RetryConfig.ShouldRetry = isRetryable

	return &ResilientProvider{
		inner:          inner,
		retrier:        NewRetrier(config.RetryConfig),
		circuitBreaker: NewCircuitBreaker(config.CircuitBreakerConfig),
	}
}

// isRetryable checks if an error should be retried.
func isRetryable(err error) bool {
	if llmErr, ok := err.(*domain.LLMError); ok {
		return llmErr.Retryable
	}
	// Default to retrying unknown errors (network issues, etc.)
	return true
}

// Complete executes the request with retry and circuit breaker protection.
func (p *ResilientProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	var resp *domain.LLMResponse
	var innerErr error

	err := p.circuitBreaker.Execute(ctx, func(ctx context.Context) error {
		return p.retrier.Execute(ctx, func(ctx context.Context) error {
			var err error
			resp, err = p.inner.Complete(ctx, req)
			innerErr = err
			return err
		})
	})

	if err != nil {
		// If circuit is open, return that error
		if err == ErrCircuitOpen {
			return nil, &domain.LLMError{
				Provider:  p.inner.Name(),
				Code:      "circuit_open",
				Message:   "circuit breaker is open, requests are being blocked",
				Retryable: true,
				Err:       err,
			}
		}
		return nil, innerErr
	}

	return resp, nil
}

// Name returns the wrapped provider's name.
func (p *ResilientProvider) Name() string {
	return p.inner.Name()
}

// Verify ResilientProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*ResilientProvider)(nil)
