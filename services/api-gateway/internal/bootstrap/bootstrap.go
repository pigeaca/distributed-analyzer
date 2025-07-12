package bootstrap

import (
	"github.com/gin-gonic/gin"
	app "github.com/pigeaca/DistributedMarketplace/libs/application"
	component "github.com/pigeaca/DistributedMarketplace/libs/application/http"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/http"
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
