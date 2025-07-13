package http

import (
	"distributed-analyzer/services/api-gateway/internal/config"
	"distributed-analyzer/services/api-gateway/internal/http/handlers"
	"distributed-analyzer/services/api-gateway/internal/service/grpc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) *gin.Engine {
	api := r.Group("/api")
	RegisterTaskRoutes(api, cfg)
	return r
}

func RegisterTaskRoutes(rg *gin.RouterGroup, cfg *config.Config) {
	taskServiceGrpcClient, _ := grpc.NewTaskServiceGrpcClient(cfg.Services.Task.GRPCAddr)
	handler := handlers.NewTaskHandler(taskServiceGrpcClient)
	handler.Register(rg.Group("/task"))
}
