package bootstrap

import (
	app "distributed-analyzer/libs/application"
	"distributed-analyzer/services/worker-manager/internal/config"
	"erro
	"log"
	"time"
)

// StartApplication initializes and starts all application components.
// It sets up the worker-manager service and handles graceful shutdown.
func StartApplication(cfg *config.Config) {
	// Parse shutdown timeout
	shutdownTimeout, err := time.ParseDuration(cfg.ShutdownTimeout)
	if err != nil {
		log.Fatalf("Invalid shutdown timeout: %v", err)
	}

	// Register cleanup handlers
	cleanupHandler := func() error {
		log.Println("Running worker-manager specific cleanup...")
		// Add any worker-manager specific cleanup logic here
		return nil
	}

	// Create and configure the application runner
	runner := app.NewApplicationRunner()

	// Register cleanup handler
	runner.Defer(cleanupHandler)

	// Log the shutdown timeout
	log.Printf("Using shutdown timeout of %s", shutdownTimeout)

	// TODO: Add actual components to the runner once the worker-manager service is implemented

	// Start the application with proper error handling
	if err := runner.Start(); err != nil {
		var appErr *app.AppError
		if errors.As(err, &appErr) {
			switch appErr.Type {
			case app.ErrorTypeStartup:
				log.Fatalf("Failed to start application: %v", err)
			case app.ErrorTypeShutdown:
				log.Fatalf("Error during shutdown: %v", err)
			default:
				log.Fatalf("Application error: %v", err)
			}
		} else {
			log.Fatalf("Failed to start application: %v", err)
		}
	}
}
