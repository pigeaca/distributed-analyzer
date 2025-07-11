package main

import (
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
	"github.com/pigeaca/DistributedMarketplace/services/worker-manager/internal/bootstrap"
	"github.com/pigeaca/DistributedMarketplace/services/worker-manager/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("api-gateway")
	bootstrap.StartApplication(&cfg)
}
