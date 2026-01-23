//nolint:errcheck,revive // Test file - error checks and unused params are OK in test helpers
package llm_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
)

// failingProvider fails a specified number of times before succeeding.
type failingProvider struct {
	failCount    int
	currentCount atomic.Int32
	successResp  *domain.LLMResponse
}

func (p *failingProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	count := p.currentCount.Add(1)
	if int(count) <= p.failCount {
		return nil, &domain.LLMError{
			Provider:  "failing",
			Message:   "temporary error",
			Retryable: true,
		}
	}
	return p.successResp, nil
}

func (p *failingProvider) Name() string {
	return "failing"
}

func TestResilientProvider_RetriesAndSucceeds(t *testing.T) {
	inner := &failingProvider{
		failCount: 2,
		successResp: &domain.LLMResponse{
			Content:      "Success after retries",
			Model:        "test-model",
			InputTokens:  10,
			OutputTokens: 5,
		},
	}

	provider := llm.NewResilientProvider(inner, llm.ResilientConfig{
		RetryConfig: llm.RetrierConfig{
			MaxAttempts: 5,
			BaseDelay:   10 * time.Millisecond,
		},
		CircuitBreakerConfig: llm.CircuitBreakerConfig{
			FailureThreshold: 10, // High threshold to not trip
			ResetTimeout:     100 * time.Millisecond,
		},
	})

	resp, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
	})

	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if resp.Content != "Success after retries" {
		t.Errorf("Content = %q, want %q", resp.Content, "Success after retries")
	}
	if inner.currentCount.Load() != 3 {
		t.Errorf("attempt count = %d, want 3", inner.currentCount.Load())
	}
}

func TestResilientProvider_CircuitBreaksAfterFailures(t *testing.T) {
	inner := &failingProvider{
		failCount: 100, // Always fails
	}

	provider := llm.NewResilientProvider(inner, llm.ResilientConfig{
		RetryConfig: llm.RetrierConfig{
			MaxAttempts: 2,
			BaseDelay:   5 * time.Millisecond,
		},
		CircuitBreakerConfig: llm.CircuitBreakerConfig{
			FailureThreshold: 3, // Trip after 3 failures
			ResetTimeout:     1 * time.Second,
		},
	})

	req := domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
	}

	// Make requests until circuit opens (3 failed retries = 6 total attempts)
	var circuitOpenErr error
	for i := 0; i < 10; i++ {
		_, err := provider.Complete(context.Background(), req)
		if errors.Is(err, llm.ErrCircuitOpen) {
			circuitOpenErr = err
			break
		}
	}

	if circuitOpenErr == nil {
		t.Error("expected circuit to open after failures")
	}
}

func TestResilientProvider_DoesNotRetryNonRetryableErrors(t *testing.T) {
	attempts := 0
	inner := &mockProvider{
		err: &domain.LLMError{
			Provider:  "mock",
			Message:   "invalid request",
			Retryable: false,
		},
	}

	// Wrap to count attempts
	countingProvider := &countingProviderWrapper{
		inner:    inner,
		attempts: &attempts,
	}

	provider := llm.NewResilientProvider(countingProvider, llm.ResilientConfig{
		RetryConfig: llm.RetrierConfig{
			MaxAttempts: 5,
			BaseDelay:   10 * time.Millisecond,
		},
	})

	_, err := provider.Complete(context.Background(), domain.LLMRequest{
		Messages: []domain.Message{
			domain.NewTextMessage(domain.RoleUser, "Hello"),
		},
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1 (should not retry non-retryable error)", attempts)
	}
}

type countingProviderWrapper struct {
	inner    domain.LLMProvider
	attempts *int
}

func (c *countingProviderWrapper) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	*c.attempts++
	return c.inner.Complete(ctx, req)
}

func (c *countingProviderWrapper) Name() string {
	return c.inner.Name()
}

func TestResilientProvider_Name(t *testing.T) {
	inner := &mockProvider{}
	provider := llm.NewResilientProvider(inner, llm.ResilientConfig{})

	if got := provider.Name(); got != "mock" {
		t.Errorf("Name() = %q, want %q", got, "mock")
	}
}
