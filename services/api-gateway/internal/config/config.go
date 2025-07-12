package config

import (
	"distributed-analyzer/libs/config"
)

// Config holds the API gateway configuration
type Config struct {
	ServerConfig configloader.ServerConfig `koanf:"server" yaml:"server"`

	// Services hold the configuration for all services
	Services struct {
		// Task service configuration
		// Default Url: http://localhost:8082
		// Default GrpcAddr: localhost:9082
		Task configloader.ServiceConnectionConfig `koanf:"task" yaml:"task"`

		// Scheduler service configuration
		// Default Url: http://localhost:8083
		// Default GrpcAddr: localhost:9083
		Scheduler configloader.ServiceConnectionConfig `koanf:"scheduler" yaml:"scheduler"`

		// Billing service configuration
		// Default Url: http://localhost:8084
		// Default GrpcAddr: localhost:9084
		Billing configloader.ServiceConnectionConfig `koanf:"billing" yaml:"billing"`
	} `koanf:"services" yaml:"services"`
}
