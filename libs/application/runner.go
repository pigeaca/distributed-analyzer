// Package application provides core functionality for application lifecycle management.
package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Runner manages the lifecycle of application components.
// It handles starting, stopping, and graceful shutdown of registered components.
type Runner struct {
	components []Component
	defers     []func() error
	cancel     context.CancelFunc
	mu         sync.Mutex // Protects concurrent access to components and defers
}

// Defer registers a function to be executed during shutdown.
// Functions are executed in reverse order (last-in-first-out).
func (r *Runner) Defer(fn func() error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defers = append(r.defers, fn)
}

// NewApplicationRunner creates a new Runner with the provided components.
func NewApplicationRunner(components ...Component) *Runner {
	r := &Runner{}
	for _, cn := range components {
		r.Register(cn)
	}
	return r
}

// Register adds a component to the Runner.
// Returns the Runner for method chaining.
func (r *Runner) Register(c Component) *Runner {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.components = append(r.components, c)
	return r
}

// Start initializes and starts all registered components.
// It blocks until a shutdown signal is received, then gracefully stops all components.
func (r *Runner) Start() error {
	// Create a cancellable context for the application
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	// Ensure components are stopped even if we encounter an error during startup
	defer func() {
		// Create a timeout context for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		r.StopAll(shutdownCtx)
	}()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start all components
	for _, c := range r.components {
		log.Printf("Starting %s...", c.Name())
		if err := c.Start(ctx); err != nil {
			return fmt.Errorf("failed to start %s: %w", c.Name(), err)
		}
	}

	log.Println("All components started. Waiting for shutdown signal...")

	// Wait for a shutdown signal
	<-sigCh
	log.Println("Shutdown signal received. Stopping components...")
	return nil
}

// StopAll stops all registered components and executes deferred functions.
// It uses the provided context for timeout/cancellation during the shutdown process.
func (r *Runner) StopAll(ctx context.Context) {
	// Stop all components
	for _, c := range r.components {
		log.Printf("Stopping %s...", c.Name())
		if err := c.Stop(ctx); err != nil {
			log.Printf("Error stopping %s: %v", c.Name(), err)
		}
	}

	// Execute deferred functions in reverse order
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := len(r.defers) - 1; i >= 0; i-- {
		if ctx.Err() != nil {
			log.Printf("Context cancelled, skipping remaining defer functions")
			break
		}

		if err := r.defers[i](); err != nil {
			log.Printf("Error running defer function: %v", err)
		}
	}

	// Cancel the application context if it hasn't been canceled already
	if r.cancel != nil {
		r.cancel()
	}
}
