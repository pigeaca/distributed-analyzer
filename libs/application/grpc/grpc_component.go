// Package grpc provides application components for gRPC server integration.
package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pigeaca/DistributedMarketplace/libs/common/config"
	"google.golang.org/grpc"
)

// Component wraps a gRPC server as an application component.
type Component struct {
	server   *grpc.Server
	listener net.Listener
}

// NewGrpcComponent creates a new GrpcComponent with the provided gRPC server and configuration.
// It sets up a TCP listener on the configured port.
func NewGrpcComponent(server *grpc.Server, cfg *config.ServerConfig) *Component {
	if server == nil {
		log.Fatal("grpc server cannot be nil")
	}

	if cfg == nil {
		log.Fatal("server config cannot be nil")
	}

	if cfg.GrpcPort == "" {
		log.Fatal("grpc port cannot be empty")
	}

	grpcAddr := fmt.Sprintf(":%s", cfg.GrpcPort)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}

	log.Printf("gRPC server will listen on %s", grpcAddr)
	return &Component{server: server, listener: listener}
}

// Start begins serving gRPC requests.
// It starts the server in a separate goroutine to avoid blocking.
func (g *Component) Start(ctx context.Context) error {
	if g.server == nil {
		return fmt.Errorf("grpc server is nil")
	}

	if g.listener == nil {
		return fmt.Errorf("grpc listener is nil")
	}

	// Use a channel to capture server errors
	errCh := make(chan error, 1)

	go func() {
		log.Printf("Starting gRPC server on %s", g.listener.Addr())
		if err := g.server.Serve(g.listener); err != nil {
			errCh <- fmt.Errorf("grpc server error: %w", err)
		}
	}()

	// Check for immediate errors
	select {
	case err := <-errCh:
		return err
	case <-time.After(100 * time.Millisecond):
		// Server started successfully
		return nil
	}
}

// Stop gracefully stops the gRPC server.
// It respects the provided context for cancellation.
func (g *Component) Stop(ctx context.Context) error {
	if g.server == nil {
		return nil // Nothing to stop
	}

	// Create a channel to signal when GracefulStop completes
	stopped := make(chan struct{})

	go func() {
		g.server.GracefulStop()
		close(stopped)
	}()

	// Wait for either GracefulStop to complete or context to be canceled
	select {
	case <-stopped:
		return nil
	case <-ctx.Done():
		// Context canceled, force stop
		g.server.Stop()
		return ctx.Err()
	}
}

// Name returns the component name for logging and identification.
func (g *Component) Name() string {
	return "gRPC"
}
