// Package ratelimit provides rate limiting functionality for HTTP and gRPC services.
package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"sync"
)

// GRPCLimiter is a rate limiter for gRPC requests.
type GRPCLimiter struct {
	limiter *rate.Limiter

	// mu protects the limiter map.
	mu sync.RWMutex

	// config is the configuration for the rate limiter.
	config GRPCConfig
}

// GRPCConfig holds configuration for the gRPC rate limiter.
type GRPCConfig struct {
	// Rate is the maximum number of requests per second.
	Rate float64

	// Burst is the maximum number of requests that can be made in a burst.
	Burst int

	// ExcludedMethods is a list of methods that are excluded from rate limiting.
	ExcludedMethods []string
}

// DefaultGRPCConfig returns a default configuration for the gRPC rate limiter.
func DefaultGRPCConfig() GRPCConfig {
	return GRPCConfig{
		Rate:  100,
		Burst: 200,
	}
}

func (l *GRPCLimiter) ClientRateLimiterInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if err := l.limiter.Wait(ctx); err != nil {
			return grpc.Errorf(grpc.Code(err), "rate limit exceeded: %v", err)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewGRPCLimiter creates a new gRPC rate limiter with the given configuration.
func NewGRPCLimiter(config GRPCConfig) *GRPCLimiter {
	limiter := &GRPCLimiter{
		limiter: rate.NewLimiter(rate.Limit(config.Rate), config.Burst),
		config:  config,
	}

	return limiter
}

// there isExcluded checks if the method is excluded from rate limiting.
func (l *GRPCLimiter) isExcluded(method string) bool {
	for _, excludedMethod := range l.config.ExcludedMethods {
		if method == excludedMethod {
			return true
		}
	}
	return false
}
