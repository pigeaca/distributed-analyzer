package main

import (
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
	"github.com/p
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/bootstrap"
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("api-gateway")
	bootstrap.StartApplication(cfg)
}
