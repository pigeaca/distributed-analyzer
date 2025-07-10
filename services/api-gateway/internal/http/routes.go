package http

import (
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/http/handlers"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	RegisterTaskRoutes(api)
	return r
}

func RegisterTaskRoutes(rg *gin.RouterGroup) {
	handler := handlers.NewTaskHandler(service.NewTaskServiceImpl())
	handler.Register(rg.Group("/task"))
}
