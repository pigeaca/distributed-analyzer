package bootstrap

import (
	"distributed-analyzer/libs/application"
	app "distributed-analyzer/libs/application/kafka"
	"distributed-analyzer/libs/kafka"
	"distributed-analyzer/services/scheduler-service/internal/config"
	"distributed-analyzer/services/scheduler-service/internal/kafka/handler"
	"distributed-analyzer/services/scheduler-service/internal/service"
	"errors"
	"log"
	"time"
)

// StartApplication initializes and starts all application components.
// It sets up the scheduler service, Kafka components, and handles graceful shutdown.
func StartApplication(cfg *config.Config) {
	// Parse shutdown timeout
	shutdownTimeout, err := time.ParseDuration(cfg.ShutdownTimeout)
	if err != nil {
		log.Fatalf("Invalid shutdown timeout: %v", err)
	}

	// Initialize components
	producer := kafka.NewProducer(cfg.Kafka.Brokers)
	schedulerService, err := service.NewSchedulerServiceImpl(cfg.GrpcPort, producer)
	if err != nil {
		log.Fatalf("Failed to create scheduler service: %v", err)
	}

	kafkaComponent, kafkaProducerComponent := initKafka(cfg, schedulerService, producer)

	// Register cleanup handlers
	cleanupHandler := func() error {
		log.Println("Running scheduler-service specific cleanup...")
		// Add any scheduler-service-specific cleanup logic here
		return nil
	}

	// Create and configure the application runner
	runner := application.NewApplicationRunner(kafkaComponent, kafkaProducerComponent)

	// Register cleanup handler
	runner.Defer(cleanupHandler)

	// Log the shutdown timeout
	log.Printf("Using shutdown timeout of %s", shutdownTimeout)

	// Start the application with proper error handling
	if err := runner.Start(); err != nil {
		var appErr *application.AppError
		if errors.As(err, &appErr) {
			switch appErr.Type {
			case application.ErrorTypeStartup:
				log.Fatalf("Failed to start application: %v", err)
			case application.ErrorTypeShutdown:
				log.Fatalf("Error during shutdown: %v", err)
			default:
				log.Fatalf("Application error: %v", err)
			}
		} else {
			log.Fatalf("Failed to start application: %v", err)
		}
	}
}

func initKafka(cfg *config.Config, schedulerService service.SchedulerService, producer *kafka.Producer) (*app.ConsumerComponent, *app.ProducerComponent) {
	taskHandler := handler.NewSchedulerHandler(schedulerService)
	topics := []string{"task-created"}
	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)
	return app.NewKafkaComponent(consumer), app.NewKafkaProducerComponent(producer)
}
