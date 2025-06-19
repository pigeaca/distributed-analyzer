package bootstrap

import (
	"github.com/distributedmarketplace/internal/task/config"
	"github.com/distributedmarketplace/internal/task/kafka/handler"
	"github.com/distributedmarketplace/internal/task/service"
	app "github.com/distributedmarketplace/pkg/application/kafka"
	"github.com/distributedmarketplace/pkg/kafka"
)

func InitKafka(cfg config.Config, taskService service.TaskService) (*kafka.Producer, *app.KafkaComponent) {
	producer := kafka.NewTaskProducer(cfg.KafkaBrokers)
	taskHandler := handler.NewTaskMessageHandler(taskService)
	topics := []string{"task-status-changed", "task-completed", "task-failed"}
	consumer := kafka.NewConsumer(topics, cfg.KafkaBrokers, cfg.KafkaGroupID, taskHandler)
	return producer, app.NewKafkaComponent(consumer)
}
