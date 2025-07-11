package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	RegisterTaskRoutes(api)
	return r
}

func RegisterTaskRoutes(rg *gin.RouterGroup) {
	//handler := handlers.NewTaskHandler(service.NewTaskServiceImpl())
	//handler.Register(rg.Group("/task"))
}
