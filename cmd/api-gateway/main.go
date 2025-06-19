package main

import (
	"fmt"
	"github.com/distributedmarketplace/internal/gateway"
	"github.com/distributedmarketplace/internal/gateway/config"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	router := gateway.SetupRouter(&cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting API Gateway on %s (env: %s)", addr, cfg.Env)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
