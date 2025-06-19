package config

type Config struct {
	// Port specifies the HTTP server port
	// Can be set via PORT environment variable
	// Default: 8083
	Port string `envconfig:"PORT" default:"8083"`

	// GrpcPort specifies the gRPC server port
	// Can be set via GRPC_PORT environment variable
	// Default: 9083
	GrpcPort string `envconfig:"GRPC_PORT" default:"9083"`

	// Env specifies the environment (development, staging, production)
	// Can be set via ENV environment variable
	// Default: development
	Env string `envconfig:"ENV" default:"development"`

	// KafkaBrokers specifies the Kafka brokers
	// Can be set via KAFKA_BROKERS environment variable
	// Default: localhost:9092
	KafkaBrokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`

	// KafkaGroupID specifies the Kafka consumer group ID
	// Can be set via KAFKA_GROUP_ID environment variable
	// Default: scheduler-service
	KafkaGroupID string `envconfig:"KAFKA_GROUP_ID" default:"scheduler-service"`
}
