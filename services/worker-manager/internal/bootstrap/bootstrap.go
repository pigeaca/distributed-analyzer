package bootstrap

import (
	"context"
	app "distributed-analyzer/libs/application"
	grpcApp "distributed-analyzer/libs/application/grpc"
	"distributed-analyzer/libs/proto/worker"
	"distributed-analyzer/services/worker-manager/internal/config"
	workerGrpc "distributed-analyzer/services/worker-manager/internal/grpc"
	"distributed-analyzer/services/worker-manager/internal/service"
	"errors"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// Initialize worker manager service
	workerManager, err := service.NewWorkerManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create worker manager: %v", err)
	}

	// Initialize gRPC server
	grpcComponent := initGrpc(cfg, workerManager)

	// Register cleanup handlers
	cleanupHandler := func() error {
		log.Println("Running worker-manager specific cleanup...")
		// Add any worker-manager-specific cleanup logic here
		return nil
	}

	// Create and configure the application runner
	runner := app.NewApplicationRunner(grpcComponent)

	// Register cleanup handler
	runner.Defer(cleanupHandler)

	// Log the shutdown timeout
	log.Printf("Using shutdown timeout of %s", shutdownTimeout)

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

// initGrpc initializes the gRPC component with the configured server.
func initGrpc(cfg *config.Config, workerManager *service.WorkerManager) *grpcApp.Component {
	grpcServer := registerGrpcServer(workerManager)
	return grpcApp.NewGrpcComponent(grpcServer, &cfg.ServerConfig)
}

// registerGrpcServer creates a new gRPC server and registers the worker manager service.
// It also enables server reflection for debugging purposes.
func registerGrpcServer(workerManager *service.WorkerManager) *stdgrpc.Server {
	// Create a server with appropriate options
	grpcServer := stdgrpc.NewServer(
		stdgrpc.ChainUnaryInterceptor(
			// Add a logging interceptor
			func(ctx context.Context, req interface{}, info *stdgrpc.UnaryServerInfo, handler stdgrpc.UnaryHandler) (interface{}, error) {
				log.Printf("gRPC request: %s", info.FullMethod)
				resp, err := handler(ctx, req)
				if err != nil {
					log.Printf("gRPC error: %v", err)
				}
				return resp, err
			},
		),
	)

	// Create a worker manager server
	workerManagerServer := workerGrpc.NewWorkerManagerServer(workerManager)

	// Register services
	worker.RegisterWorkerManagerServiceServer(grpcServer, workerManagerServer)
	reflection.Register(grpcServer)

	return grpcServer
}
