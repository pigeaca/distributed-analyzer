package config

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`

	Services ServicesConfig `yaml:"services"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"      env:"SERVER_PORT"      env-default:"8080"`
	GRPCPort int    `yaml:"grpc_port" env:"SERVER_GRPC_PORT" env-default:"9090"`
	Env      string `yaml:"env"       env:"ENV"              env-default:"development"`
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
