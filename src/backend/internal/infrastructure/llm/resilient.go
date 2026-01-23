package llm

import (
	"context"
	"errors"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"

	"backend/internal/domain"
)

// ResilientConfig holds configuration for the resilient provider.
type ResilientConfig struct {
	// Retry settings
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration

	// Circuit breaker settings
	FailureThreshold int
	ResetTimeout     time.Duration

	// Timeout settings
	RequestTimeout time.Duration
}

// ResilientProvider wraps an LLM provider with retry, circuit breaker, and timeout.
type ResilientProvider struct {
	inner    domain.LLMProvider
	executor failsafe.Executor[*domain.LLMResponse]
}

// NewResilientProvider creates a resilient provider with retry, circuit breaker, and timeout.
func NewResilientProvider(inner domain.LLMProvider, config ResilientConfig) *ResilientProvider {
	// Apply defaults if not set
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}
	if config.BaseDelay == 0 {
		config.BaseDelay = 500 * time.Millisecond
	}
	if config.MaxDelay == 0 {
		config.MaxDelay = 30 * time.Second
	}
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 5
	}
	if config.ResetTimeout == 0 {
		config.ResetTimeout = 60 * time.Second
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 120 * time.Second
	}

	// Build retry policy
	retry := retrypolicy.NewBuilder[*domain.LLMResponse]().
		HandleIf(func(_ *domain.LLMResponse, err error) bool {
			return isRetryable(err)
		}).
		WithBackoff(config.BaseDelay, config.MaxDelay).
		WithMaxAttempts(config.MaxAttempts).
		WithJitterFactor(0.1).
		Build()

	// Build circuit breaker
	failureThreshold := config.FailureThreshold
	if failureThreshold < 0 {
		failureThreshold = 0
	}
	cb := circuitbreaker.NewBuilder[*domain.LLMResponse]().
		HandleIf(func(_ *domain.LLMResponse, err error) bool {
			return err != nil
		}).
		WithFailureThreshold(uint(failureThreshold)). //nolint:gosec // Bounds checked above
		WithDelay(config.ResetTimeout).
		Build()

	// Build timeout
	to := timeout.New[*domain.LLMResponse](config.RequestTimeout)

	// Compose policies: timeout -> circuit breaker -> retry -> provider
	// Outer policies are listed first in failsafe.With
	executor := failsafe.With(to, cb, retry)

	return &ResilientProvider{
		inner:    inner,
		executor: executor,
	}
}

// isRetryable checks if an error should be retried.
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if llmErr, ok := err.(*domain.LLMError); ok {
		return llmErr.Retryable
	}
	// Default to retrying unknown errors (network issues, etc.)
	return true
}

// Complete executes the request with retry, circuit breaker, and timeout protection.
func (p *ResilientProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	resp, err := p.executor.WithContext(ctx).GetWithExecution(func(exec failsafe.Execution[*domain.LLMResponse]) (*domain.LLMResponse, error) {
		return p.inner.Complete(exec.Context(), req)
	})

	if err != nil {
		// Convert failsafe errors to domain errors
		return nil, p.convertFailsafeError(err)
	}

	return resp, nil
}

// convertFailsafeError converts failsafe-go errors to domain errors.
func (p *ResilientProvider) convertFailsafeError(err error) error {
	// Check for circuit breaker open
	if errors.Is(err, circuitbreaker.ErrOpen) {
		return &domain.LLMError{
			Provider:  p.inner.Name(),
			Code:      "circuit_open",
			Message:   "circuit breaker is open, requests are being blocked",
			Retryable: true,
			Err:       err,
		}
	}

	// Check for timeout exceeded
	if errors.Is(err, timeout.ErrExceeded) {
		return &domain.LLMError{
			Provider:  p.inner.Name(),
			Code:      "timeout",
			Message:   "request timed out",
			Retryable: true,
			Err:       err,
		}
	}

	// Check for retries exceeded
	if retrypolicy.IsExceededError(err) {
		// Extract the underlying error if available
		if exceeded := retrypolicy.AsExceededError(err); exceeded != nil && exceeded.LastError != nil {
			// Return the last error from the retry chain
			if llmErr, ok := exceeded.LastError.(*domain.LLMError); ok {
				return llmErr
			}
			return &domain.LLMError{
				Provider:  p.inner.Name(),
				Message:   exceeded.LastError.Error(),
				Retryable: false,
				Err:       exceeded.LastError,
			}
		}
		return &domain.LLMError{
			Provider:  p.inner.Name(),
			Code:      "retries_exceeded",
			Message:   "max retries exceeded",
			Retryable: false,
			Err:       err,
		}
	}

	// Return as-is if already a domain error
	if llmErr, ok := err.(*domain.LLMError); ok {
		return llmErr
	}

	// Wrap unknown errors
	return &domain.LLMError{
		Provider:  p.inner.Name(),
		Message:   err.Error(),
		Retryable: false,
		Err:       err,
	}
}

// Name returns the wrapped provider's name.
func (p *ResilientProvider) Name() string {
	return p.inner.Name()
}

// Verify ResilientProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*ResilientProvider)(nil)
