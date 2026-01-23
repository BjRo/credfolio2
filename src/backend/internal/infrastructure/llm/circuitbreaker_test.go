//nolint:errcheck,revive // Test file - error checks and unused params are OK in test helpers
package llm_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"backend/internal/infrastructure/llm"
)

func TestCircuitBreaker_ClosedState(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     100 * time.Millisecond,
	})

	// Successful calls should pass through
	called := false
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if !called {
		t.Error("function was not called")
	}
}

func TestCircuitBreaker_OpensAfterFailures(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     100 * time.Millisecond,
	})

	testErr := errors.New("test error")

	// Fail 3 times to trip the circuit
	for i := 0; i < 3; i++ {
		_ = cb.Execute(context.Background(), func(ctx context.Context) error {
			return testErr
		})
	}

	// Next call should fail immediately with circuit open error
	called := false
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})

	if err == nil {
		t.Error("expected error when circuit is open")
	}
	if called {
		t.Error("function should not be called when circuit is open")
	}
	if !errors.Is(err, llm.ErrCircuitOpen) {
		t.Errorf("error = %v, want ErrCircuitOpen", err)
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
	})

	testErr := errors.New("test error")

	// Trip the circuit
	for i := 0; i < 2; i++ {
		_ = cb.Execute(context.Background(), func(ctx context.Context) error {
			return testErr
		})
	}

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// Should allow one call through (half-open state)
	called := false
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if !called {
		t.Error("function should be called in half-open state")
	}
}

func TestCircuitBreaker_ClosesAfterSuccessInHalfOpen(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
	})

	testErr := errors.New("test error")

	// Trip the circuit
	for i := 0; i < 2; i++ {
		_ = cb.Execute(context.Background(), func(ctx context.Context) error {
			return testErr
		})
	}

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// Succeed in half-open state
	_ = cb.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})

	// Should now be closed and allow calls
	callCount := 0
	for i := 0; i < 5; i++ {
		err := cb.Execute(context.Background(), func(ctx context.Context) error {
			callCount++
			return nil
		})
		if err != nil {
			t.Errorf("call %d error = %v, want nil", i, err)
		}
	}

	if callCount != 5 {
		t.Errorf("callCount = %d, want 5", callCount)
	}
}

func TestCircuitBreaker_ReopensAfterFailureInHalfOpen(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     50 * time.Millisecond,
	})

	testErr := errors.New("test error")

	// Trip the circuit
	for i := 0; i < 2; i++ {
		_ = cb.Execute(context.Background(), func(ctx context.Context) error {
			return testErr
		})
	}

	// Wait for reset timeout
	time.Sleep(60 * time.Millisecond)

	// Fail in half-open state
	_ = cb.Execute(context.Background(), func(ctx context.Context) error {
		return testErr
	})

	// Should be open again
	called := false
	err := cb.Execute(context.Background(), func(ctx context.Context) error {
		called = true
		return nil
	})

	if !errors.Is(err, llm.ErrCircuitOpen) {
		t.Errorf("error = %v, want ErrCircuitOpen", err)
	}
	if called {
		t.Error("function should not be called when circuit is open")
	}
}

func TestCircuitBreaker_ContextCancellation(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 3,
		ResetTimeout:     100 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := cb.Execute(ctx, func(ctx context.Context) error {
		return nil
	})

	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestCircuitBreaker_ConcurrentAccess(t *testing.T) {
	cb := llm.NewCircuitBreaker(llm.CircuitBreakerConfig{
		FailureThreshold: 100, // High threshold to avoid tripping
		ResetTimeout:     100 * time.Millisecond,
	})

	var successCount atomic.Int32

	// Run many concurrent requests
	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func() {
			err := cb.Execute(context.Background(), func(ctx context.Context) error {
				successCount.Add(1)
				return nil
			})
			if err != nil {
				t.Errorf("Execute() error = %v", err)
			}
			done <- struct{}{}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}

	if successCount.Load() != 100 {
		t.Errorf("successCount = %d, want 100", successCount.Load())
	}
}
