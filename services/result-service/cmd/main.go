package main

import (
	configloader "github.com/pigeaca/DistributedMarketplace/libs/config"
	"github.com/pigeaca/DistributedMarketplace/services/result-service/internal/bootstrap"
	"github.com/pigeaca/DistributedMarketplace/services/result-service/internal/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("result_service")
	bootstrap.StartApplication(&cfg)
}
