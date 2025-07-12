// Package grpc provides application components for gRPC server integration.
package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"distributed-analyzer/libs/common/config"
	"google.golang.org/grpc"
)

// Component wraps a gRPC server as an application component.
type Component struct {
	server       *grpc.Server
	listener     net.Listener
	startTimeout time.Duration // Timeout for checking server startup errors
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
	return &Component{
		server:       server,
		listener:     listener,
		startTimeout: 500 * time.Millisecond, // Default timeout for checking server startup
	}
}

// Start begins serving gRPC requests.
// It starts the server in a separate goroutine to avoid blocking.
// It respects the provided context for cancellation and uses the configured
// startTimeout to check for immediate startup errors.
func (g *Component) Start(ctx context.Context) error {
	if g.server == nil {
		return fmt.Errorf("failed to start gRPC server: server is nil")
	}

	if g.listener == nil {
		return fmt.Errorf("failed to start gRPC server: listener is nil")
	}

	// Create a context that will be canceled when the server starts or when an error occurs
	serverCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is canceled to prevent goroutine leaks

	// Use a channel to capture server errors
	errCh := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on %s", g.listener.Addr())
		if err := g.server.Serve(g.listener); err != nil {
			select {
			case errCh <- fmt.Errorf("gRPC server failed: %w", err):
				// Error sent successfully
			case <-serverCtx.Done():
				// Context already canceled, server is shutting down
				log.Printf("gRPC server shutting down: %v", err)
			}
		}
		cancel() // Signal that the server has stopped
	}()

	// Check for immediate errors or context cancellation
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// Parent context was canceled
		return fmt.Errorf("gRPC server start canceled: %w", ctx.Err())
	case <-time.After(g.startTimeout):
		// No immediate errors detected within the timeout period
		log.Printf("gRPC server started successfully on %s", g.listener.Addr())
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
