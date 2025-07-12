// Package config provides common configuration structures used across services
package configloader

// ServerConfig holds common server configuration settings
type ServerConfig struct {
	// Port specifies the HTTP server port
	// Can be set via PORT environment variable
	Port string `envconfig:"PORT" yaml:"port"`

	// GrpcPort specifies the gRPC server port
	// Can be set via GRPC_PORT environment variable
	GrpcPort string `envconfig:"GRPC_PORT" yaml:"grpc_port"`

	// Env specifies the environment (development, staging, production)
	// Can be set via ENV environment variable
	// Default: development
	Env string `envconfig:"ENV" default:"development" yaml:"env"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	// Host specifying the database host
	// Can be set via DB_HOST environment variable
	// Default: localhost
	Host string `envconfig:"DB_HOST" default:"localhost" yaml:"host"`

	// Port specifies the database port
	// Can be set via DB_PORT environment variable
	// Default: 5432
	Port int `envconfig:"DB_PORT" default:"5432" yaml:"port"`

	// User specifies the database user
	// Can be set via DB_USER environment variable
	// Default: postgres
	User string `envconfig:"DB_USER" default:"postgres" yaml:"user"`

	// Password specifies the database password
	// Can be set via DB_PASSWORD environment variable
	// Default: postgres
	Password string `envconfig:"DB_PASSWORD" default:"postgres" yaml:"password"`

	// Name specifies the database name
	// Can be set via DB_NAME environment variable
	Name string `envconfig:"DB_NAME" yaml:"name"`

	// SSLMode specifies the database SSL mode
	// Can be set via DB_SSL_MODE environment variable
	// Default: disable
	SSLMode string `envconfig:"DB_SSL_MODE" default:"disable" yaml:"ssl_mode"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	// Level specifies the log level
	// Can be set via LOG_LEVEL environment variable
	// Default: info
	Level string `envconfig:"LOG_LEVEL" default:"info" yaml:"level"`

	// Format specifying the log format
	// Can be set via LOG_FORMAT environment variable
	// Default: json
	Format string `envconfig:"LOG_FORMAT" default:"json" yaml:"format"`
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Topics []string `envconfig:"KAFKA_TOPICS" yaml:"topics"`
	// Brokers specifies the Kafka brokers
	// Can be set via KAFKA_BROKERS environment variable
	// Default: localhost:9092
	Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092" yaml:"brokers"`

	// GroupID specifies the Kafka consumer group ID
	// Can be set via KAFKA_GROUP_ID environment variable
	GroupID string `envconfig:"KAFKA_GROUP_ID" yaml:"group_id"`
}

// ServiceConnectionConfig holds service connection details
type ServiceConnectionConfig struct {
	// Url specifies the Url of the service for HTTP communication
	Url string `koanf:"url" yaml:"url"`

	// GrpcAddr specifies the gRPC address of the service
	// Can be set via <SERVICE>_GRPC_ADDR environment variable
	GrpcAddr string `koanf:"grpc_addr" yaml:"grpc_addr"`
}
