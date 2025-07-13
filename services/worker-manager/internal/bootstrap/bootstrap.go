package bootstrap

import (
	app "distributed-analyzer/libs/application"
	grpcApp "distributed-analyzer/libs/application/grpc"
	"distributed-analyzer/libs/network/logging"
	"distributed-analyzer/libs/proto/worker"
	"distributed-analyzer/services/worker-manager/internal/config"
	workerGrpc "distributed-analyzer/services/worker-manager/internal/grpc"
	"distributed-analyzer/services/worker-manager/internal/service"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

// StartApplication initializes and starts all application components.
// It sets up the worker-manager service and handles graceful shutdown.
func StartApplication(cfg *config.Config) {
	// Initialize worker manager service
	workerManager, err := service.NewWorkerManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create worker manager: %v", err)
	}

	// Initialize gRPC server
	grpcComponent := initGrpc(cfg, workerManager)

	// Create and configure the application runner
	runner := app.NewApplicationRunner(grpcComponent)

	runner.DefaultStart()
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
	grpcServer := stdgrpc.NewServer(stdgrpc.ChainUnaryInterceptor(logging.ServerInterceptor()))

	// Create a worker manager server
	workerManagerServer := workerGrpc.NewWorkerManagerServer(workerManager)

	// Register services
	worker.RegisterWorkerManagerServiceServer(grpcServer, workerManagerServer)
	reflection.Register(grpcServer)

	return grpcServer
}
