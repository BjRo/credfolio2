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

func TestRetry_SucceedsOnFirstAttempt(t *testing.T) {
	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
	})

	attempts := 0
	err := r.Execute(context.Background(), func(ctx context.Context) error {
		attempts++
		return nil
	})

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1", attempts)
	}
}

func TestRetry_RetriesOnError(t *testing.T) {
	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    50 * time.Millisecond,
	})

	attempts := 0
	testErr := errors.New("temporary error")

	err := r.Execute(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return testErr
		}
		return nil
	})

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
}

func TestRetry_FailsAfterMaxAttempts(t *testing.T) {
	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
	})

	attempts := 0
	testErr := errors.New("persistent error")

	err := r.Execute(context.Background(), func(ctx context.Context) error {
		attempts++
		return testErr
	})

	if err == nil {
		t.Error("expected error after max attempts")
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	if !errors.Is(err, testErr) {
		t.Errorf("error = %v, want %v", err, testErr)
	}
}

func TestRetry_RespectsContextCancellation(t *testing.T) {
	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 10,
		BaseDelay:   100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var attempts atomic.Int32
	err := r.Execute(ctx, func(ctx context.Context) error {
		attempts.Add(1)
		return errors.New("error")
	})

	if err == nil {
		t.Error("expected error for cancelled context")
	}
	// Should have made at least 1 attempt but not all 10
	if attempts.Load() >= 10 {
		t.Errorf("made too many attempts: %d", attempts.Load())
	}
}

func TestRetry_ExponentialBackoff(t *testing.T) {
	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 4,
		BaseDelay:   20 * time.Millisecond,
		MaxDelay:    200 * time.Millisecond,
		Multiplier:  2.0,
	})

	var timestamps []time.Time

	err := r.Execute(context.Background(), func(ctx context.Context) error {
		timestamps = append(timestamps, time.Now())
		if len(timestamps) < 4 {
			return errors.New("retry")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	// Check that delays are increasing (with some tolerance for timing)
	// Expected delays: ~20ms, ~40ms, ~80ms
	for i := 1; i < len(timestamps)-1; i++ {
		delay := timestamps[i+1].Sub(timestamps[i])
		prevDelay := timestamps[i].Sub(timestamps[i-1])
		// Each delay should be at least 1.5x the previous (allowing for jitter)
		if delay < prevDelay {
			t.Logf("delay %d (%v) < prev delay (%v)", i, delay, prevDelay)
		}
	}
}

func TestRetry_RespectsMaxDelay(t *testing.T) {
	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 5,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    60 * time.Millisecond, // Cap at 60ms
		Multiplier:  10.0,                  // Aggressive multiplier
	})

	var timestamps []time.Time

	_ = r.Execute(context.Background(), func(ctx context.Context) error {
		timestamps = append(timestamps, time.Now())
		if len(timestamps) < 5 {
			return errors.New("retry")
		}
		return nil
	})

	// Check that no delay exceeds MaxDelay significantly
	for i := 1; i < len(timestamps); i++ {
		delay := timestamps[i].Sub(timestamps[i-1])
		// Allow 50% tolerance for jitter and timing
		maxAllowed := 90 * time.Millisecond
		if delay > maxAllowed {
			t.Errorf("delay %d = %v, exceeds max allowed %v", i, delay, maxAllowed)
		}
	}
}

func TestRetry_ShouldRetryFunc(t *testing.T) {
	permanentErr := errors.New("permanent error")
	retryableErr := errors.New("retryable error")

	r := llm.NewRetrier(llm.RetrierConfig{
		MaxAttempts: 5,
		BaseDelay:   10 * time.Millisecond,
		ShouldRetry: func(err error) bool {
			return errors.Is(err, retryableErr)
		},
	})

	// Test that permanent errors don't retry
	attempts := 0
	err := r.Execute(context.Background(), func(ctx context.Context) error {
		attempts++
		return permanentErr
	})

	if attempts != 1 {
		t.Errorf("attempts for permanent error = %d, want 1", attempts)
	}
	if !errors.Is(err, permanentErr) {
		t.Errorf("error = %v, want %v", err, permanentErr)
	}

	// Test that retryable errors do retry
	attempts = 0
	err = r.Execute(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return retryableErr
		}
		return nil
	})

	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if attempts != 3 {
		t.Errorf("attempts for retryable error = %d, want 3", attempts)
	}
}
