package bootstrap

import (
	"github.com/gin-gonic/gin"
	app "github.com/pigeaca/DistributedMarketplace/libs/application"
	component "github.com/pigeaca/DistributedMarketplace/libs/application/http"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/http"
)

func StartApplication(cfg config.Config) error {
	httpComponent := initHttp(cfg)
	runner := app.NewApplicationRunner(httpComponent)
	return runner.StartBlocking()
}

func initHttp(cfg config.Config) *component.GinHttpComponent {
	var routes = http.RegisterRoutes(gin.Default())
	addr := ":" + cfg.Port
	httpComponent := component.NewGinHttpComponent(addr, routes)
	return httpComponent
}
