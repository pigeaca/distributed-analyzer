package bootstrap

import (
	"github.com/distributedmarketplace/internal/scheduler/config"
	handler2 "github.com/distributedmarketplace/internal/scheduler/kafka/handler"
	"github.com/distributedmarketplace/internal/scheduler/service"
	app "github.com/distributedmarketplace/pkg/application/kafka"
	"github.com/distributedmarketplace/pkg/kafka"
)

func InitKafka(cfg config.Config, schedulerService service.SchedulerService) *app.KafkaComponent {
	taskHandler := handler2.NewSchedulerHandler(schedulerService)
	topics := []string{"task-created"}
	consumer := kafka.NewConsumer(topics, cfg.KafkaBrokers, cfg.KafkaGroupID, taskHandler)
	return app.NewKafkaComponent(consumer)
}
