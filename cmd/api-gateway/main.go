package main

import (
	"github.com/distributedmarketplace/internal/gateway/config"
	"github.com/distributedmarketplace/internal/gateway/http"
	app "github.com/distributedmarketplace/pkg/application"
	app2 "github.com/distributedmarketplace/pkg/application/http"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func main() {
	var cfg = loadConfig()

	var routes = http.RegisterRoutes(gin.Default())

	httpComponent := app2.NewGinHttpComponent(cfg.BillingServiceGrpcAddr, routes)
	runner := app.NewApplicationRunner(httpComponent)

	if err := runner.StartBlocking(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func loadConfig() config.Config {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}
	return cfg
}
