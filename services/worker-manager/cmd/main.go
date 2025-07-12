package main

import (
	configloader "distributed-analyzer/libs/config"
	"distributed-analyzer/services/worker-manager/internal/bootstrap"
	"distributed-analyzer/services/worker-manager/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("worker-manager")
	bootstrap.StartApplication(&cfg)
}
