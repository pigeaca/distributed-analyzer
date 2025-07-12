package main

import (
	configloader "distributed-analyzer/libs/config"
	bootstrap "distributed-analyzer/services/worker/internal"
	"distributed-analyzer/services/worker/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("worker")
	bootstrap.StartApplication(&cfg)
}
