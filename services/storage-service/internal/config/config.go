package config

import (
	commonConfig "github.com/distributedmarketplace/internal/common/config"
)

type Config struct {
	// Server settings
	commonConfig.ServerConfig `yaml:",inline"`

	// Storage settings
	Storage StorageConfig `yaml:"storage"`

	// Files settings
	Files FilesConfig `yaml:"files"`

	// Cache settings
	Cache CacheConfig `yaml:"cache"`

	// Log settings
	Log commonConfig.LogConfig `yaml:"log"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	// Endpoint specifies the MinIO/S3 endpoint
	// Can be set via STORAGE_ENDPOINT environment variable
	// Default: localhost:9000
	Endpoint string `envconfig:"STORAGE_ENDPOINT" default:"localhost:9000" yaml:"endpoint"`

	// AccessKey specifies the MinIO/S3 access key
	// Can be set via STORAGE_ACCESS_KEY environment variable
	// Default: minioadmin
	AccessKey string `envconfig:"STORAGE_ACCESS_KEY" default:"minioadmin" yaml:"access_key"`

	// SecretKey specifies the MinIO/S3 secret key
	// Can be set via STORAGE_SECRET_KEY environment variable
	// Default: minioadmin
	SecretKey string `envconfig:"STORAGE_SECRET_KEY" default:"minioadmin" yaml:"secret_key"`

	// UseSSL specifies whether to use SSL for MinIO/S3 connections
	// Can be set via STORAGE_USE_SSL environment variable
	// Default: false
	UseSSL bool `envconfig:"STORAGE_USE_SSL" default:"false" yaml:"use_ssl"`

	// Bucket specifies the MinIO/S3 bucket
	// Can be set via STORAGE_BUCKET environment variable
	// Default: distributed-analyzer
	Bucket string `envconfig:"STORAGE_BUCKET" default:"distributed-analyzer" yaml:"bucket"`

	// Region specifies the MinIO/S3 region
	// Can be set via STORAGE_REGION environment variable
	// Default: us-east-1
	Region string `envconfig:"STORAGE_REGION" default:"us-east-1" yaml:"region"`
}

// FilesConfig holds file-related configuration
type FilesConfig struct {
	// MaxSize specifies the maximum file size
	// Can be set via FILES_MAX_SIZE environment variable
	// Default: 100MB
	MaxSize string `envconfig:"FILES_MAX_SIZE" default:"100MB" yaml:"max_size"`

	// AllowedTypes specifies the allowed file types
	// Can be set via FILES_ALLOWED_TYPES environment variable
	// Default: .go,.mod,.sum,.txt,.json,.yaml,.yml
	AllowedTypes []string `envconfig:"FILES_ALLOWED_TYPES" default:".go,.mod,.sum,.txt,.json,.yaml,.yml" yaml:"allowed_types"`

	// TempDir specifies the temporary directory for file processing
	// Can be set via FILES_TEMP_DIR environment variable
	// Default: /tmp/storage-service
	TempDir string `envconfig:"FILES_TEMP_DIR" default:"/tmp/storage-service" yaml:"temp_dir"`
}

// CacheConfig holds cache-related configuration
type CacheConfig struct {
	// Enabled specifies whether caching is enabled
	// Can be set via CACHE_ENABLED environment variable
	// Default: true
	Enabled bool `envconfig:"CACHE_ENABLED" default:"true" yaml:"enabled"`

	// TTL specifies the cache time-to-live
	// Can be set via CACHE_TTL environment variable
	// Default: 1h
	TTL string `envconfig:"CACHE_TTL" default:"1h" yaml:"ttl"`

	// MaxSize specifies the maximum cache size
	// Can be set via CACHE_MAX_SIZE environment variable
	// Default: 1GB
	MaxSize string `envconfig:"CACHE_MAX_SIZE" default:"1GB" yaml:"max_size"`
}
