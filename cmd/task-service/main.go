package main

import (
	"github.com/distributedmarketplace/internal/gateway/bootstrap"
	"github.com/distributedmarketplace/internal/task/config"
	"github.com/distributedmarketplace/internal/task/service"
	app "github.com/distributedmarketplace/pkg/application"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	runner := app.NewAppRunner()

	taskService := service.NewTaskServiceImpl()

	// --- Kafka
	producer, kafkaComponent := bootstrap.InitKafka(cfg, taskService)
	runner.Register(kafkaComponent)

	// --- gRPC
	_, grpcComponent := bootstrap.InitGrpc(cfg, producer, taskService)
	runner.Register(grpcComponent)

	// --- Run everything
	if err := runner.StartBlocking(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
