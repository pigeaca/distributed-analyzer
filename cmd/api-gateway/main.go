package main

import (
	"github.com/distributedmarketplace/internal/gateway/bootstrap"
	"github.com/distributedmarketplace/internal/gateway/config"
	configloader "github.com/distributedmarketplace/pkg/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("api_gateway")
	bootstrap.StartApplication(cfg)
}
