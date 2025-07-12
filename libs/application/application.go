// Package application provides core functionality for application lifecycle management.
package application

import (
	"context"
)

// Component represents an application component that can be started and stopped.
// Components are managed by the Runner, which handles their lifecycle.
type Component interface {
	// Start initializes and starts the component.
	// It should be non-blocking and return quickly.
	// Long-running operations should be started in a separate goroutine.
	// The provided context can be used for cancellation.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the component and releases any resources.
	// It should respect the provided context's deadline or cancellation.
	Stop(ctx context.Context) error

	// Name returns the component's name, used for logging and identification.
	Name() string
}
