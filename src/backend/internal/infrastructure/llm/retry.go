package llm

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// RetrierConfig holds configuration for the retrier.
type RetrierConfig struct { //nolint:govet // Field order prioritizes readability
	// MaxAttempts is the maximum number of attempts (including the first).
	MaxAttempts int

	// BaseDelay is the initial delay between retries.
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration

	// Multiplier is the exponential backoff multiplier.
	Multiplier float64

	// Jitter adds randomness to delays (0.0-1.0).
	Jitter float64

	// ShouldRetry determines if an error should be retried.
	// If nil, all errors are retried.
	ShouldRetry func(error) bool
}

// Retrier implements retry with exponential backoff.
type Retrier struct {
	config RetrierConfig
}

// NewRetrier creates a new retrier.
func NewRetrier(config RetrierConfig) *Retrier {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	if config.BaseDelay <= 0 {
		config.BaseDelay = 100 * time.Millisecond
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = 30 * time.Second
	}
	if config.Multiplier <= 0 {
		config.Multiplier = 2.0
	}
	if config.Jitter < 0 || config.Jitter > 1 {
		config.Jitter = 0.1
	}

	return &Retrier{
		config: config,
	}
}

// Execute runs the function with retries.
func (r *Retrier) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error

	for attempt := 0; attempt < r.config.MaxAttempts; attempt++ {
		// Check context before each attempt
		if ctx.Err() != nil {
			if lastErr != nil {
				return lastErr
			}
			return ctx.Err()
		}

		// Execute the function
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry
		if r.config.ShouldRetry != nil && !r.config.ShouldRetry(err) {
			return err
		}

		// Don't wait after the last attempt
		if attempt < r.config.MaxAttempts-1 {
			delay := r.calculateDelay(attempt)
			if err := r.sleep(ctx, delay); err != nil {
				return lastErr
			}
		}
	}

	return lastErr
}

// calculateDelay computes the delay for a given attempt.
func (r *Retrier) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * multiplier^attempt
	delay := float64(r.config.BaseDelay) * math.Pow(r.config.Multiplier, float64(attempt))

	// Apply jitter
	if r.config.Jitter > 0 {
		jitterRange := delay * r.config.Jitter
		jitter := (rand.Float64() * 2 * jitterRange) - jitterRange //nolint:gosec // G404: crypto random not needed for jitter
		delay += jitter
	}

	// Cap at max delay
	if delay > float64(r.config.MaxDelay) {
		delay = float64(r.config.MaxDelay)
	}

	return time.Duration(delay)
}

// sleep waits for the given duration or until context is cancelled.
func (r *Retrier) sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
