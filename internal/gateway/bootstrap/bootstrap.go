package bootstrap

import (
	"github.com/distributedmarketplace/internal/gateway/config"
	"github.com/distributedmarketplace/internal/gateway/http"
	app "github.com/distributedmarketplace/pkg/application"
	component "github.com/distributedmarketplace/pkg/application/http"
	"github.com/gin-gonic/gin"
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
