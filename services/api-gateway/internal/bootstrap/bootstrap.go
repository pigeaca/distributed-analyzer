package bootstrap

import (
	app "distributed-analyzer/libs/application"
	component "distributed-analyzer/libs/application/http"
	"distributed-analyzer/services/api-gateway/internal/config"
	"distributed-analyzer/services/api-gateway/internal/http"
	"gith
	"log"
)

func StartApplication(cfg *config.Config) {
	httpComponent := initHttp(cfg)
	runner := app.NewApplicationRunner(httpComponent)
	if err := runner.Start(); err != nil {
		log.Println("Error while starting application", err)
	}
}

func initHttp(cfg *config.Config) *component.GinHttpComponent {
	var routes = http.RegisterRoutes(gin.Default(), cfg)
	httpComponent := component.NewGinHttpComponent(&cfg.ServerConfig, routes)
	return httpComponent
}
