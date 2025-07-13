package bootstrap

import (
	app "distributed-analyzer/libs/application"
	component "distributed-analyzer/libs/application/http"
	"distributed-analyzer/services/api-gateway/internal/config"
	"distributed-analyzer/services/api-gateway/internal/http"
	"github.com/gin-gonic/gin"
	"log"
)

func StartApplication(cfg *config.Config) {
	httpComponent := initHttpComponent(cfg)
	runner := app.NewApplicationRunner(httpComponent)
	if err := runner.Start(); err != nil {
		log.Println("Error while starting application", err)
	}
}

func initHttpComponent(cfg *config.Config) *component.GinHttpComponent {
	var routes = http.RegisterRoutes(gin.Default(), cfg)
	httpComponent := component.NewGinHttpComponent(&cfg.ServerConfig, routes)
	return httpComponent
}
