package main

import (
	configloader "distributed-analyzer/libs/config"
	"distributed-analyzer/services/task-service/internal/bootstrap"
	"distributed-analyzer/services/task-service/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("task-service")
	bootstrap.StartApplication(&cfg)
}
