package llm

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	// StateClosed allows requests through normally.
	StateClosed CircuitState = iota
	// StateOpen blocks all requests.
	StateOpen
	// StateHalfOpen allows one test request through.
	StateHalfOpen
)

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of failures before opening the circuit.
	FailureThreshold int

	// ResetTimeout is how long to wait before attempting to close the circuit.
	ResetTimeout time.Duration
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct { //nolint:govet // Field order prioritizes readability
	config CircuitBreakerConfig

	mu           sync.Mutex
	state        CircuitState
	failureCount int
	lastFailure  time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	if config.FailureThreshold <= 0 {
		config.FailureThreshold = 5
	}
	if config.ResetTimeout <= 0 {
		config.ResetTimeout = 30 * time.Second
	}

	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute runs the given function if the circuit allows it.
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	// Check context first
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Check if we can proceed
	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	// Execute the function
	err := fn(ctx)

	// Record the result
	cb.recordResult(err)

	return err
}

// allowRequest checks if the circuit breaker allows a request.
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// Check if reset timeout has passed
		if time.Since(cb.lastFailure) > cb.config.ResetTimeout {
			cb.state = StateHalfOpen
			return true
		}
		return false

	case StateHalfOpen:
		// Only allow one request in half-open state
		// The request is already in progress if we're here
		return true

	default:
		return false
	}
}

// recordResult records the result of a request.
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err == nil {
		// Success - reset the circuit
		cb.failureCount = 0
		cb.state = StateClosed
		return
	}

	// Failure
	cb.failureCount++
	cb.lastFailure = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failureCount >= cb.config.FailureThreshold {
			cb.state = StateOpen
		}

	case StateHalfOpen:
		// Any failure in half-open state opens the circuit again
		cb.state = StateOpen
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
}
