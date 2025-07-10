package bootstrap

import (
	"github.com/distributedmarketplace/internal/scheduler/config"
	"github.com/distributedmarketplace/internal/scheduler/kafka/handler"
	"github.com/distributedmarketplace/internal/scheduler/service"
	"github.com/distributedmarketplace/pkg/application"
	app "github.com/distributedmarketplace/pkg/application/kafka"
	"github.com/distributedmarketplace/pkg/kafka"
	"log"
)

func StartApplication(cfg config.Config) error {
	producer := kafka.NewProducer(cfg.Kafka.Brokers)
	schedulerService, err := service.NewSchedulerServiceImpl(cfg.GrpcPort, producer)
	if err != nil {
		log.Fatalf("failed to create scheduler service: %v", err)
	}
	kafkaComponent := initKafka(cfg, schedulerService)
	runner := application.NewApplicationRunner(kafkaComponent)
	return runner.StartBlocking()
}

func initKafka(cfg config.Config, schedulerService service.SchedulerService) *app.KafkaComponent {
	taskHandler := handler.NewSchedulerHandler(schedulerService)
	topics := []string{"task-created"}
	consumer := kafka.NewConsumer(topics, cfg.Kafka.Brokers, cfg.Kafka.GroupID, taskHandler)
	return app.NewKafkaComponent(consumer)
}
