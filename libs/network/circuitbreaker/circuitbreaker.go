// Package circuitbreaker provides a circuit breaker implementation for service communication.
package circuitbreaker

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker is a wrapper around gobreaker.CircuitBreaker with additional functionality.
type CircuitBreaker struct {
	name   string
	cb     *gobreaker.CircuitBreaker
	config Config
}

// Config holds configuration for the circuit breaker.
type Config struct {
	// Name is the name of the circuit breaker.
	Name string

	// MaxRequests is the maximum number of requests allowed to pass through
	// when the CircuitBreaker is half-open.
	// If MaxRequests is 0, the default value is 1.
	MaxRequests uint32

	// Interval is the cyclic period of the closed state for the CircuitBreaker to
	// clear the internal counts. If Interval is 0, the CircuitBreaker doesn't clear
	// internal counts during the closed state.
	Interval time.Duration

	// Timeout is the period of the open state, after which the state of the
	// CircuitBreaker becomes half-open.
	// If Timeout is 0, the default value is 60 seconds.
	Timeout time.Duration

	// ReadyToTrip is called with a copy of Counts whenever a request fails in the
	// closed state. If ReadyToTrip returns true, the CircuitBreaker will be placed
	// into the open state.
	// If ReadyToTrip is nil, the default value is used, which returns true when the
	// number of consecutive failures is more than 5.
	ReadyToTrip func(counts gobreaker.Counts) bool

	// OnStateChange is called whenever the state of the CircuitBreaker changes.
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// DefaultConfig returns a default configuration for the circuit breaker.
func DefaultConfig(name string) Config {
	return Config{
		Name:        name,
		MaxRequests: 1,
		Interval:    0,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker %s changed from %s to %s\n", name, from, to)
		},
	}
}

func (cb *CircuitBreaker) CircuitBreakerInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		_, err := cb.ExecuteContext(ctx, func(ctx context.Context) (interface{}, error) {
			return nil, invoker(ctx, method, req, reply, cc, opts...)
		})
		return err
	}
}

// New creates a new CircuitBreaker with the given configuration.
func New(config Config) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: config.ReadyToTrip,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if config.OnStateChange != nil {
				config.OnStateChange(name, from, to)
			}
		},
	}

	return &CircuitBreaker{
		name:   config.Name,
		cb:     gobreaker.NewCircuitBreaker(settings),
		config: config,
	}
}

// Execute executes the given function if the circuit breaker is closed or half-open.
// If the circuit breaker is open, it returns an error immediately.
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return cb.cb.Execute(fn)
}

// ExecuteContext executes the given function with a context if the circuit breaker is closed or half-open.
// If the circuit breaker is open, it returns an error immediately.
// If the context is canceled, it returns the context error.
func (cb *CircuitBreaker) ExecuteContext(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	return cb.cb.Execute(func() (interface{}, error) {
		// Check if the context is already canceled
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		return fn(ctx)
	})
}

// Name returns the name of the circuit breaker.
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.cb.State()
}

// Counts return the current counts of the circuit breaker.
func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	return cb.cb.Counts()
}
