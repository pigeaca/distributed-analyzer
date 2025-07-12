package configloader

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LoadApplicationConfig loads and validates config.
// It applies defaults, then YAML, then ENV.
func LoadApplicationConfig[T any](prefix string) T {
	var cfg T

	configPath := resolveConfigPath(prefix)
	log.Printf("Loading config from: %s", configPath)

	// Read a config file (optional)
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read environment variables: %v", err)
	}

	printFinalConfig(cfg)

	return cfg
}

func printFinalConfig[T any](cfg T) {
	out, err := yaml.Marshal(cfg)
	if err != nil {
		log.Printf("Failed to marshal config to YAML: %v", err)
		return
	}
	log.Println("Configuration:")
	log.Println(string(out))
}

func resolveConfigPath(prefix string) string {
	// 1. CONFIG_PATH env
	if custom := os.Getenv("CONFIG_PATH"); custom != "" {
		return custom
	}

	// 2. Try to resolve from known locations
	filePrefix := strings.ReplaceAll(prefix, "_", "-")
	candidates := []string{
		filepath.Join(findConfigDir(), fmt.Sprintf("%s.yml", prefix)),
		filepath.Join(findConfigDir(), fmt.Sprintf("%s.yml", filePrefix)),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	log.Fatalf("no config file found for prefix %s", prefix)
	return ""
}

func findConfigDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}

	for i := 0; i < 5; i++ {
		configs := filepath.Join(dir, "configs")
		if stat, err := os.Stat(configs); err == nil && stat.IsDir() {
			return configs
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	log.Fatalf("configs directory not found near working directory")
	return ""
}
