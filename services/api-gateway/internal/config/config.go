package config

import configloader "distributed-analyzer/libs/config"

type Config struct {
	configloader.ServerConfig `yaml:",inline"`

	Services ServicesConfig `yaml:"services"`
}

type ServicesConfig struct {
	Task      ServiceConnectionConfig `yaml:"task"`
	Scheduler ServiceConnectionConfig `yaml:"scheduler"`
	Billing   ServiceConnectionConfig `yaml:"billing"`
}

type ServiceConnectionConfig struct {
	URL      string `yaml:"url"       env:"{PREFIX}_URL"       env-default:"http://localhost:8080"`
	GRPCAddr string `yaml:"grpc_addr" env:"{PREFIX}_GRPC_ADDR" env-default:"localhost:9090"`
}
