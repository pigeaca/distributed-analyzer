package circuitbreaker

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

type Config struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

func DefaultConfig(name string) Config {
	return Config{
		Name:        name,
		MaxRequests: 1,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(c gobreaker.Counts) bool {
			return c.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			fmt.Printf("Circuit breaker %s changed from %s to %s\n", name, from.String(), to.String())
		},
	}
}

func ClientInterceptor(cfg Config) grpc.UnaryClientInterceptor {
	cb := createCircuitBreaker(cfg)

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		_, err := executeWithContext(cb, ctx, func(ctx context.Context) (interface{}, error) {
			return nil, invoker(ctx, method, req, reply, cc, opts...)
		})
		return err
	}
}

func createCircuitBreaker(cfg Config) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: cfg.ReadyToTrip,
	}

	if cfg.OnStateChange != nil {
		settings.OnStateChange = cfg.OnStateChange
	}

	return gobreaker.NewCircuitBreaker(settings)
}

func executeWithContext(cb *gobreaker.CircuitBreaker, ctx context.Context, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	return cb.Execute(func() (interface{}, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return fn(ctx)
	})
}
