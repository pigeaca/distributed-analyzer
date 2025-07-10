package main

import (
	"github.com/distributedmarketplace/internal/result/bootstrap"
	"github.com/distributedmarketplace/internal/result/config"
	configloader "github.com/distributedmarketplace/pkg/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("result_service")
	bootstrap.StartApplication(cfg)
}
