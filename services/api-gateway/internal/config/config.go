package config

import "github.com/pigeaca/DistributedMarketplace/libs/common/config"

// Config holds the API gateway configuration
type Config struct {

	// Port and Env settings
	Port string `koanf:"port" default:"8081" yaml:"port"`
	Env  string `koanf:"env" default:"development" yaml:"env"`

	// Services holds the configuration for all services
	Services struct {
		// Task service configuration
		// Default Url: http://localhost:8082
		// Default GrpcAddr: localhost:9082
		Task config.ServiceConnectionConfig `koanf:"task" yaml:"task"`

		// Scheduler service configuration
		// Default Url: http://localhost:8083
		// Default GrpcAddr: localhost:9083
		Scheduler config.ServiceConnectionConfig `koanf:"scheduler" yaml:"scheduler"`

		// Billing service configuration
		// Default Url: http://localhost:8084
		// Default GrpcAddr: localhost:9084
		Billing config.ServiceConnectionConfig `koanf:"billing" yaml:"billing"`
	} `koanf:"services" yaml:"services"`
}
