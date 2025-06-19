package gateway

import (
	"github.com/distributedmarketplace/internal/gateway/config"
	"github.com/distributedmarketplace/internal/gateway/handler"
	"github.com/distributedmarketplace/internal/gateway/service"
	grpcClient "github.com/distributedmarketplace/internal/gateway/service/grpc"
	httpClient "github.com/distributedmarketplace/internal/gateway/service/http"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the API router
func SetupRouter(config *config.Config) *gin.Engine {
	r := gin.Default()

	var serviceClient service.TaskServiceClient
	var err error

	// Use gRPC or HTTP client based on configuration
	if config.UseGrpc {
		// Use the TaskServiceGrpcClient to communicate with the task service via gRPC
		serviceClient, err = grpcClient.NewTaskServiceGrpcClient(config.TaskServiceGrpcAddr)
		if err != nil {
			// Fall back to HTTP if the gRPC connection fails
			serviceClient = httpClient.NewTaskServiceClient(config.TaskServiceURL)
		}
	} else {
		// Use the TaskServiceHttpClient to communicate with the task service via HTTP
		serviceClient = httpClient.NewTaskServiceClient(config.TaskServiceURL)
	}
	// Set up handlers
	taskHandler := handler.NewTaskHandler(serviceClient)

	// Register task routes
	r.POST("/tasks", taskHandler.SubmitTask)
	r.GET("/tasks/:id", taskHandler.GetTaskStatus)

	return r
}
