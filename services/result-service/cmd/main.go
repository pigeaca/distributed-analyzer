package main

import (
	configloader "distributed-analyzer/libs/config"
	"distributed-analyzer/services/result-service/internal/bootstrap"
	"distributed-analyzer/services/result-service/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("result_service")
	bootstrap.StartApplication(&cfg)
}
