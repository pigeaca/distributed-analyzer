package main

import (
	bootstrap "github.com/distributedmarketplace/internal/worker"
	"github.com/distributedmarketplace/internal/worker/config"
	configloader "github.com/distributedmarketplace/pkg/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("worker")
	bootstrap.StartApplication(cfg)
}
