// Package application provides core functionality for application lifecycle management.
package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ErrorType represents the category of an application error
type ErrorType int

const (
	// ErrorTypeStartup indicates an error occurred during component startup
	ErrorTypeStartup ErrorType = iota
	// ErrorTypeRuntime indicates an error occurred during normal operation
	ErrorTypeRuntime
	// ErrorTypeShutdown indicates an error occurred during shutdown
	ErrorTypeShutdown
)

// AppError is a custom error type that provides context about where and why an error occurred
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Details map[string]interface{}
}

func (e *AppError) Error() string {
	var sb strings.Builder
	sb.WriteString(e.Message)

	if e.Cause != nil {
		sb.WriteString(": ")
		sb.WriteString(e.Cause.Error())
	}

	return sb.String()
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewStartupError creates a new error of type ErrorTypeStartup
func NewStartupError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeStartup,
		Message: message,
		Cause:   cause,
	}
}

// NewRuntimeError creates a new error of type ErrorTypeRuntime
func NewRuntimeError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeRuntime,
		Message: message,
		Cause:   cause,
	}
}

// NewShutdownError creates a new error of type ErrorTypeShutdown
func NewShutdownError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeShutdown,
		Message: message,
		Cause:   cause,
	}
}

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
		r.RegisterComponent(cn)
	}
	return r
}

// RegisterComponent adds a component to the Runner.
// Returns the Runner for method chaining.
func (r *Runner) RegisterComponent(c Component) *Runner {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.components = append(r.components, c)
	return r
}

// Start initializes and starts all registered components.
// It blocks until a shutdown signal is received, then gracefully stops all components.
// Returns an error if any component fails to start or if there are errors during shutdown.
func (r *Runner) Start() error {
	// Create a cancellable context for the application
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel

	// Variable to store any startup error
	var startupErr error

	// Ensure components are stopped even if we encounter an error during startup
	defer func() {
		// Create a timeout context for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Stop all components and collect any shutdown errors
		if shutdownErr := r.StopAll(shutdownCtx); shutdownErr != nil {
			// If we already have a startup error, log the shutdown error but keep the startup error as the return value
			if startupErr != nil {
				log.Printf("Additional errors during shutdown: %v", shutdownErr)
			} else {
				// If there was no startup error, return the shutdown error
				startupErr = shutdownErr
			}
		}
	}()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal when all components have started
	started := make(chan struct{})

	// Start components in a goroutine so we can handle signals during startup
	go func() {
		// Start all components
		for _, c := range r.components {
			log.Printf("Starting %s...", c.Name())
			if err := c.Start(ctx); err != nil {
				startupErr = NewStartupError(fmt.Sprintf("Failed to start %s", c.Name()), err)
				cancel() // Cancel the context to signal other components to stop
				return
			}
		}

		log.Println("All components started. Waiting for shutdown signal...")
		close(started)
	}()

	// Wait for either all components to start or a shutdown signal
	select {
	case <-started:
		// All components started successfully, now wait for the shutdown signal
		<-sigCh
		log.Println("Shutdown signal received. Stopping components...")
		return startupErr // Will be nil if everything went well
	case <-ctx.Done():
		// Context was canceled during startup, likely due to a startup error
		if startupErr != nil {
			return startupErr
		}
		return NewStartupError("Startup canceled", ctx.Err())
	case sig := <-sigCh:
		// Shutdown signal received during startup
		log.Printf("Shutdown signal (%v) received during startup. Stopping components...", sig)
		return NewRuntimeError("Shutdown signal received during startup", nil)
	}
}

// StopAll stops all registered components and executes deferred functions.
// It uses the provided context for timeout/cancellation during the shutdown process.
// Returns an error that aggregates all errors encountered during shutdown.
func (r *Runner) StopAll(ctx context.Context) error {
	var shutdownErrors []string

	// Stop all components
	for _, c := range r.components {
		log.Printf("Stopping %s...", c.Name())
		if err := c.Stop(ctx); err != nil {
			errMsg := fmt.Sprintf("Error stopping %s: %v", c.Name(), err)
			log.Print(errMsg)
			shutdownErrors = append(shutdownErrors, errMsg)
		}
	}

	// Execute deferred functions in reverse order
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := len(r.defers) - 1; i >= 0; i-- {
		if ctx.Err() != nil {
			errMsg := fmt.Sprintf("Context cancelled, skipping remaining defer functions: %v", ctx.Err())
			log.Print(errMsg)
			shutdownErrors = append(shutdownErrors, errMsg)
			break
		}

		if err := r.defers[i](); err != nil {
			errMsg := fmt.Sprintf("Error running defer function: %v", err)
			log.Print(errMsg)
			shutdownErrors = append(shutdownErrors, errMsg)
		}
	}

	// Cancel the application context if it hasn't been canceled already
	if r.cancel != nil {
		r.cancel()
	}

	// If we encountered any errors, return them as a single error
	if len(shutdownErrors) > 0 {
		return NewShutdownError(
			fmt.Sprintf("Encountered %d errors during shutdown", len(shutdownErrors)),
			fmt.Errorf(strings.Join(shutdownErrors, "; ")),
		)
	}

	return nil
}
