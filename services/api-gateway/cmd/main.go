// @title Distributed Marketplace API
// @version 1.0
// @description API Gateway for the Distributed Marketplace system.
// @host localhost:8081
// @BasePath /api
package main

import (
	configloader "distributed-analyzer/libs/config"
	"distributed-analyzer/services/api-gateway/internal/bootstrap"
	"distributed-analyzer/services/api-gateway/internal/config"
)

const ApiGatewayApplicationName = "api-gateway"

func main() {
	var cfg = configloader.LoadApplicationConfig[config.Config](ApiGatewayApplicationName)
	bootstrap.StartApplication(&cfg)
}
