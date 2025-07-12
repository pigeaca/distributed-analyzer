package main

import (
	configloader "distributed-analyzer/libs/config"
	"distributed-analyzer/services/scheduler-service/internal/bootstrap"
	"distributed-analyzer/services/scheduler-service/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("scheduler_service")
	bootstrap.StartApplication(&cfg)
}
