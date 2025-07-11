package main

import (
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
	bootstrap "github.com/pigeaca/DistributedMarketplace/services/worker/internal"
	"github.com/pigeaca/DistributedMarketplace/services/worker/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("worker")
	bootstrap.StartApplication(&cfg)
}
