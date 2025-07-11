package main

import (
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
	"github.com/pigeaca/DistributedMarketplace/services/scheduler-service/internal/bootstrap"
	"github.com/pigeaca/DistributedMarketplace/services/scheduler-service/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("scheduler_service")
	bootstrap.StartApplication(&cfg)
}
