package main

import (
	"github.com/distributedmarketplace/internal/scheduler/bootstrap"
	"github.com/distributedmarketplace/internal/scheduler/config"
	service2 "github.com/distributedmarketplace/internal/scheduler/service"
	app "github.com/distributedmarketplace/pkg/application"
	"github.com/distributedmarketplace/pkg/kafka"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func main() {
	cfg := loadConfig()

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	schedulerService, err := service2.NewSchedulerServiceImpl(cfg.GrpcPort, producer)

	if err != nil {
		log.Fatalf("failed to create scheduler service: %v", err)
	}

	kafkaComponent := bootstrap.InitKafka(cfg, schedulerService)
	runner := app.NewApplicationRunner(kafkaComponent)

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
