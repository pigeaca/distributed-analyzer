package client

import (
	"distributed-analyzer/libs/network/circuitbreaker"
	"distributed-analyzer/libs/network/ratelimit"
	"distributed-analyzer/libs/network/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config aggregates configuration for all resilience mechanisms.
type Config struct {
	CBConfig     circuitbreaker.Config
	RetryConfig  retry.Config
	RateLimitCfg ratelimit.GRPCConfig
}

// NewGrpcResilientClient constructs a new client with circuit breaker, rate limiter, and retry.
func NewGrpcResilientClient(cfg *Config, grpcTarget string) (*grpc.ClientConn, error) {
	if cfg == nil {
		cfg = &Config{
			CBConfig:     circuitbreaker.DefaultConfig(""),
			RetryConfig:  retry.DefaultConfig(),
			RateLimitCfg: ratelimit.DefaultGRPCConfig(),
		}
	}

	limiterInterceptor := ratelimit.NewGRPCLimiter(cfg.RateLimitCfg).ClientRateLimiterInterceptor()
	circuitBreakerInterceptor := circuitbreaker.New(cfg.CBConfig).CircuitBreakerInterceptor()
	retryInterceptor := retry.RetryInterceptor(cfg.RetryConfig)

	dialOption := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	dialOption = append(dialOption, grpc.WithChainUnaryInterceptor(
		limiterInterceptor,
		circuitBreakerInterceptor,
		retryInterceptor,
	))

	conn, err := grpc.NewClient(grpcTarget, dialOption...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
