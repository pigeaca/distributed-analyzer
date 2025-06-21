package main

import (
	"github.com/distributedmarketplace/internal/task/bootstrap"
	"github.com/distributedmarketplace/internal/task/config"
	"github.com/distributedmarketplace/internal/task/service"
	app "github.com/distributedmarketplace/pkg/application"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func main() {
	cfg := loadConfig()

	taskService := service.NewTaskServiceImpl()

	producer, kafkaComponent := bootstrap.InitKafka(cfg, taskService)
	_, grpcComponent := bootstrap.InitGrpc(cfg, producer, taskService)

	runner := app.NewApplicationRunner(grpcComponent, kafkaComponent)

	if err := runner.StartBlocking(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func loadConfig() config.Config {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}
	return cfg
}
