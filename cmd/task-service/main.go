package main

import (
	"github.com/distributedmarketplace/internal/task/bootstrap"
	"github.com/distributedmarketplace/internal/task/config"
	configloader "github.com/distributedmarketplace/pkg/config"
)

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config]("task_service")
	bootstrap.StartApplication(cfg)
}
