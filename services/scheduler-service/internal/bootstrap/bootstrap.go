package bootstrap

import (
	"distributed-analyzer/libs/application"
	app "distributed-analyzer/libs/application/kafka"
	"distributed-analyzer/libs/kafka"
	"distributed-analyzer/services/scheduler-service/internal/config"
	"distributed-analyzer/services/scheduler-service/internal/kafka/handler"
	"distributed-analyzer/services/scheduler-service/internal/service"
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

	// Create and configure the application runner
	runner := application.NewApplicationRunner(kafkaComponent, kafkaProducerComponent)

	// Log the shutdown timeout
	log.Printf("Using shutdown timeout of %s", shutdownTimeout)

	runner.DefaultStart()
}

func initKafka(cfg *config.Config, schedulerService service.SchedulerService, producer *kafka.Producer) (*app.ConsumerComponent, *app.ProducerComponent) {
	taskHandler := handler.NewSchedulerHandler(schedulerService)
	topics := []string{"task-created"}
	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)
	return app.NewKafkaComponent(consumer), app.NewKafkaProducerComponent(producer)
}
