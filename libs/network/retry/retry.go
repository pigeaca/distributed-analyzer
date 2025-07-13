// Package retry provides retry functionality with exponential backoff for service communication.
package retry

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Config holds configuration for the retry mechanism.
type Config struct {
	// InitialInterval is the initial interval between retries.
	InitialInterval time.Duration

	// MaxInterval is the maximum interval between retries.
	MaxInterval time.Duration

	// MaxElapsedTime is the maximum elapsed time for retries.
	MaxElapsedTime time.Duration

	// Multiplier is the multiplier for the exponential backoff.
	Multiplier float64

	// RandomizationFactor is the randomization factor for the exponential backoff.
	RandomizationFactor float64

	// MaxRetries is the maximum number of retries.
	// If 0, there is no limit.
	MaxRetries int

	// RetryableErrors is a function that determines if an error is retryable.
	// If nil, all errors are considered retryable.
	RetryableErrors func(error) bool
}

func ClientInterceptor(cfg Config) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return Retry(ctx, cfg, func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}

// DefaultConfig returns a default configuration for the retry mechanism.
func DefaultConfig() Config {
	return Config{
		InitialInterval:     100 * time.Millisecond,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      1 * time.Minute,
		Multiplier:          1.5,
		RandomizationFactor: 0.5,
		MaxRetries:          5,
		RetryableErrors:     nil,
	}
}

// Retry retries the given function with exponential backoff.
func Retry(ctx context.Context, config Config, operation func() error) error {
	// Create a new exponential backoff
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = config.InitialInterval
	expBackoff.MaxInterval = config.MaxInterval
	expBackoff.MaxElapsedTime = config.MaxElapsedTime
	expBackoff.Multiplier = config.Multiplier
	expBackoff.RandomizationFactor = config.RandomizationFactor

	// Reset the backoff
	expBackoff.Reset()

	// Create a backoff with context
	b := backoff.WithContext(expBackoff, ctx)

	// Create a retry counter
	var retryCount int

	// Create a retry function
	retryFunc := func() error {
		// Check if we've reached the maximum number of retries
		if config.MaxRetries > 0 && retryCount >= config.MaxRetries {
			return backoff.Permanent(errors.New("maximum number of retries reached"))
		}

		// Increment the retry counter
		retryCount++

		// Execute the operation
		err := operation()
		if err == nil {
			return nil
		}

		// Check if the error is retryable
		if config.RetryableErrors != nil && !config.RetryableErrors(err) {
			return backoff.Permanent(err)
		}

		// Return the error to trigger a retry
		return err
	}

	// Retry the operation with backoff
	return backoff.Retry(retryFunc, b)
}
