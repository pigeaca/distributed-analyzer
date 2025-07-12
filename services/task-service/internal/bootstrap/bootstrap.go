// Package bootstrap provides functionality to initialize and start the application components.
package bootstrap

import (
	"context"
	"fmt"
	"github.com/pigeaca/DistributedMarketplace/libs/application"
	grpcApp "github.com/pigeaca/DistributedMarketplace/libs/application/grpc"
	kafkaApp "github.com/pigeaca/DistributedMarketplace/libs/application/kafka"
	"github.com/pigeaca/DistributedMarketplace/libs/kafka"
	pb "github.com/pigeaca/DistributedMarketplace/libs/proto/task"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/grpc"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/handler"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/producer"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/service"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

// StartApplication initializes and starts all application components.
// It sets up the task service, Kafka components, and gRPC server.
func StartApplication(cfg *config.Config) {
	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize service and components
	taskService := service.NewTaskServiceImpl()
	kafkaConsumerComponent, kafkaProducerComponent := initKafka(cfg, taskService)
	grpcComponent := initGrpc(cfg, kafkaProducerComponent.Producer(), taskService)

	// Start the application
	runner := application.NewApplicationRunner(grpcComponent, kafkaConsumerComponent, kafkaProducerComponent)
	if err := runner.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

// validateConfig checks if the configuration is valid.
func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	if len(cfg.Kafka.Brokers) == 0 {
		return fmt.Errorf("no Kafka brokers specified")
	}

	if cfg.Kafka.GroupID == "" {
		return fmt.Errorf("kafka group ID is empty")
	}

	if cfg.ServerConfig.GrpcPort == "" {
		return fmt.Errorf("gRPC port is empty")
	}

	return nil
}

// initKafka initializes both Kafka consumer and producer components.
func initKafka(cfg *config.Config, taskService service.TaskService) (*kafkaApp.ConsumerComponent, *kafkaApp.ProducerComponent) {
	consumerComponent := initKafkaConsumerComponent(cfg, taskService)
	producerComponent := initKafkaProducerComponent(cfg)
	return consumerComponent, producerComponent
}

// initKafkaConsumerComponent creates and configures a Kafka consumer component.
// It sets up the message handler and subscribes to the required topics.
func initKafkaConsumerComponent(cfg *config.Config, taskService service.TaskService) *kafkaApp.ConsumerComponent {
	if taskService == nil {
		log.Fatalf("Task service is nil")
	}

	taskHandler := handler.NewTaskMessageHandler(taskService)
	topics := []string{"task-status-changed", "task-completed", "task-failed"}

	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)

	return kafkaApp.NewKafkaComponent(consumer)
}

// initKafkaProducerComponent creates and configures a Kafka producer component.
func initKafkaProducerComponent(cfg *config.Config) *kafkaApp.ProducerComponent {
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers)
	return kafkaApp.NewKafkaProducerComponent(kafkaProducer)
}

// initGrpc initializes the gRPC component with the configured server.
func initGrpc(cfg *config.Config, kafkaProducer *kafka.Producer, service service.TaskService) *grpcApp.Component {
	if kafkaProducer == nil {
		log.Fatalf("Kafka producer is nil")
	}

	if service == nil {
		log.Fatalf("Task service is nil")
	}

	grpcServer := registerGrpcServer(kafkaProducer, service)

	return grpcApp.NewGrpcComponent(grpcServer, &cfg.ServerConfig)
}

// registerGrpcServer creates a new gRPC server and registers the task service.
// It also enables server reflection for debugging purposes.
func registerGrpcServer(kafkaProducer *kafka.Producer, service service.TaskService) *stdgrpc.Server {
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

	// Create task producer and server
	taskProducer := producer.NewTaskProducer(kafkaProducer)
	taskGrpcServer := grpc.NewTaskServer(service, taskProducer)

	// Register services
	pb.RegisterTaskServiceServer(grpcServer, taskGrpcServer)
	reflection.Register(grpcServer)

	return grpcServer
}
