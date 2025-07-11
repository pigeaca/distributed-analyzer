package bootstrap

import (
	"github.com/pigeaca/DistributedMarketplace/libs/application"
	app "github.com/pigeaca/DistributedMarketplace/libs/application/kafka"
	"github.com/pigeaca/DistributedMarketplace/libs/kafka"
	"github.com/pigeaca/DistributedMarketplace/services/scheduler-service/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/scheduler-service/internal/kafka/handler"
	"github.com/pigeaca/DistributedMarketplace/services/scheduler-service/internal/service"
	"log"
)

func StartApplication(cfg *config.Config) {
	producer := kafka.NewProducer(cfg.Kafka.Brokers)
	schedulerService, err := service.NewSchedulerServiceImpl(cfg.GrpcPort, producer)
	if err != nil {
		log.Fatalf("failed to create scheduler service: %v", err)
	}
	kafkaComponent := initKafka(cfg, schedulerService)
	runner := application.NewApplicationRunner(kafkaComponent)
	if err := runner.Start(); err != nil {
		log.Println("Error while starting application", err)
	}
}

func initKafka(cfg *config.Config, schedulerService service.SchedulerService) *app.KafkaComponent {
	taskHandler := handler.NewSchedulerHandler(schedulerService)
	topics := []string{"task-created"}
	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)
	return app.NewKafkaComponent(consumer)
}
