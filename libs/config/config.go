package configloader

// ServerConfig holds common server configuration settings
type ServerConfig struct {
	Port     string `yaml:"port"      env:"PORT"`
	GrpcPort string `yaml:"grpc_port" env:"GRPC_PORT"`
	Env      string `yaml:"env"       env:"ENV" env-default:"development"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"     env:"DB_HOST"     env-default:"localhost"`
	Port     int    `yaml:"port"     env:"DB_PORT"     env-default:"5432"`
	User     string `yaml:"user"     env:"DB_USER"     env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
	Name     string `yaml:"name"     env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"disable"`
}

type LogConfig struct {
	Level  string `yaml:"level"  env:"LOG_LEVEL"  env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"json"`
}

type KafkaConfig struct {
	Topics  []string `yaml:"topics"  env:"KAFKA_TOPICS"`
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	GroupID string   `yaml:"group_id" env:"KAFKA_GROUP_ID"`
}

type ServiceConnectionConfig struct {
	URL      string `yaml:"url"       env:"SERVICE_URL"`
	GRPCAddr string `yaml:"grpc_addr" env:"SERVICE_GRPC_ADDR"`
}
