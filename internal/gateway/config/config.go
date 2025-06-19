package config

// Config holds the API gateway configuration
type Config struct {
	// Port specifies the HTTP server port
	// Can be set via PORT environment variable
	// Default: 8081
	Port string `envconfig:"PORT" default:"8081"`

	// Env specifies the environment (development, staging, production)
	// Can be set via ENV environment variable
	// Default: development
	Env string `envconfig:"ENV" default:"development"`

	// TaskServiceURL specifies the URL of the task service for HTTP communication
	// Can be set via TASK_SERVICE_URL environment variable
	// Default: http://localhost:8082
	TaskServiceURL string `envconfig:"TASK_SERVICE_URL" default:"http://localhost:8082"`

	// TaskServiceGrpcAddr specifies the address of the task service for gRPC communication
	// Can be set via TASK_SERVICE_GRPC_ADDR environment variable
	// Default: localhost:9082
	TaskServiceGrpcAddr string `envconfig:"TASK_SERVICE_GRPC_ADDR" default:"localhost:9082"`

	// SchedulerServiceURL specifies the URL of the scheduler service for HTTP communication
	// Can be set via SCHEDULER_SERVICE_URL environment variable
	// Default: http://localhost:8083
	SchedulerServiceURL string `envconfig:"SCHEDULER_SERVICE_URL" default:"http://localhost:8083"`

	// SchedulerServiceGrpcAddr specifies the address of the scheduler service for gRPC communication
	// Can be set via SCHEDULER_SERVICE_GRPC_ADDR environment variable
	// Default: localhost:9083
	SchedulerServiceGrpcAddr string `envconfig:"SCHEDULER_SERVICE_GRPC_ADDR" default:"localhost:9083"`

	// BillingServiceURL specifies the URL of the billing service for HTTP communication
	// Can be set via BILLING_SERVICE_URL environment variable
	// Default: http://localhost:8084
	BillingServiceURL string `envconfig:"BILLING_SERVICE_URL" default:"http://localhost:8084"`

	// BillingServiceGrpcAddr specifies the address of the billing service for gRPC communication
	// Can be set via BILLING_SERVICE_GRPC_ADDR environment variable
	// Default: localhost:9084
	BillingServiceGrpcAddr string `envconfig:"BILLING_SERVICE_GRPC_ADDR" default:"localhost:9084"`

	// UseGrpc specifies whether to use gRPC for service-to-service communication
	// Can be set via USE_GRPC environment variable
	// Default: true
	UseGrpc bool `envconfig:"USE_GRPC" default:"true"`
}
