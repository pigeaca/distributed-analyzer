package main

import (
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/bootstrap"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("task-service")
	bootstrap.StartApplication(&cfg)
}
