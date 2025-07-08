package configloader

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

func LoadApplicationConfig[T any](prefix string) T {
	var cfg T
	if err := envconfig.Process(prefix, &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}
