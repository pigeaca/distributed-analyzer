package configloader

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LoadApplicationConfig loads the configuration and then validates it.
// It returns the loaded and validated configuration.
// If validation fails, it panics with an error message.
func LoadApplicationConfig[T any](prefix string) T {
	// Load configuration
	config := loadConfig[T](prefix)

	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		panic(fmt.Sprintf("configuration validation failed: %v", err))
	}

	return config
}

func loadConfig[T any](prefix string) T {
	var cfg T
	k := koanf.New(".")

	loadYamlConfig(k, prefix)

	// Override with environment variables
	k.Load(env.Provider(prefix+"_", ".", func(s string) string {
		//convert MYAPP_SERVER_PORT -> server.port
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil)

	// Unmarshal into typed struct
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("failed to unmarshal config into struct: %v", err)
	}

	return cfg
}

func loadYamlConfig(k *koanf.Koanf, prefix string) {
	filePrefix := strings.ReplaceAll(prefix, "_", "-")

	configPaths := []string{
		filepath.Join("configs", fmt.Sprintf("%s.yml", prefix)),
		filepath.Join("configs", fmt.Sprintf("%s.yml", filePrefix)),
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
				log.Printf("Warning: failed to load config from %s: %v", path, err)
			} else {
				log.Printf("Loaded configuration from %s", path)
			}
			return
		}
	}

	log.Printf("Warning: no config file found for prefix %s", prefix)
}
