package config

import "github.com/pigeaca/DistributedMarketplace/libs/common/config"

// Config holds the API gateway configuration
type Config struct {
	// Server settings
	config.ServerConfig `yaml:",inline"`

	// Services hold the configuration for all-services
	// WorkerManagement settings
	WorkerManagement WorkerManagementConfig `yaml:"worker_management"`
}

// WorkerManagementConfig holds worker management configuration
type WorkerManagementConfig struct {
	// HeartbeatInterval specifies the worker heartbeat interval
	// Can be set via WORKER_HEARTBEAT_INTERVAL environment variable
	// Default: 30s
	HeartbeatInterval string `envconfig:"WORKER_HEARTBEAT_INTERVAL" default:"30s" yaml:"heartbeat_interval"`

	// Timeout specifying the worker timeout
	// Can be set via WORKER_TIMEOUT environment variable
	// Default: 60s
	Timeout string `envconfig:"WORKER_TIMEOUT" default:"60s" yaml:"timeout"`
}
