package main

import (
	"github.com/distributedmarketplace/internal/scheduler/bootstrap"
	"github.com/distributedmarketplace/internal/scheduler/config"
	configloader "github.com/distributedmarketplace/pkg/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("scheduler_service")
	bootstrap.StartApplication(cfg)
}
