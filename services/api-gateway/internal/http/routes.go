package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/config"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/http/handlers"
	"github.com/pigeaca/DistributedMarketplace/services/api-gateway/internal/service/grpc"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) *gin.Engine {
	api := r.Group("/api")
	RegisterTaskRoutes(api, cfg)
	return r
}

func RegisterTaskRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	taskServiceGrpcClient, _ := grpc.NewTaskServiceGrpcClient(cfg.Services.Task.GrpcAddr)
	handler := handlers.NewTaskHandler(taskServiceGrpcClient)
	handler.Register(rg.Group("/task"))
}
