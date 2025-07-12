package config

import (
	configloader "distributed-analyzer/libs/config"
)

// Config is the main configuration for the storage service
type Config struct {
	configloader.ServerConfig `yaml:",inline"`

	Storage StorageConfig          `yaml:"storage"`
	Files   FilesConfig            `yaml:"files"`
	Cache   CacheConfig            `yaml:"cache"`
	Log     configloader.LogConfig `yaml:"log"`
}

// StorageConfig holds MinIO/S3-related settings
type StorageConfig struct {
	Endpoint  string `yaml:"endpoint"    env:"STORAGE_ENDPOINT"    env-default:"localhost:9000"`
	AccessKey string `yaml:"access_key"  env:"STORAGE_ACCESS_KEY"  env-default:"minioadmin"`
	SecretKey string `yaml:"secret_key"  env:"STORAGE_SECRET_KEY"  env-default:"minioadmin"`
	UseSSL    bool   `yaml:"use_ssl"     env:"STORAGE_USE_SSL"     env-default:"false"`
	Bucket    string `yaml:"bucket"      env:"STORAGE_BUCKET"      env-default:"distributed-analyzer"`
	Region    string `yaml:"region"      env:"STORAGE_REGION"      env-default:"us-east-1"`
}

// FilesConfig holds file-upload configuration
type FilesConfig struct {
	MaxSize      string   `yaml:"max_size"       env:"FILES_MAX_SIZE"        env-default:"100MB"`
	AllowedTypes []string `yaml:"allowed_types"  env:"FILES_ALLOWED_TYPES"   env-default:".go,.mod,.sum,.txt,.json,.yaml,.yml"`
	TempDir      string   `yaml:"temp_dir"       env:"FILES_TEMP_DIR"        env-default:"/tmp/storage-service"`
}

// CacheConfig holds in-memory cache settings
type CacheConfig struct {
	Enabled bool   `yaml:"enabled"    env:"CACHE_ENABLED"    env-default:"true"`
	TTL     string `yaml:"ttl"        env:"CACHE_TTL"        env-default:"1h"`
	MaxSize string `yaml:"max_size"   env:"CACHE_MAX_SIZE"   env-default:"1GB"`
}
