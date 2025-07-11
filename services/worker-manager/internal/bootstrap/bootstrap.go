package bootstrap

import (
	app "github.com/pigeaca/DistributedMarketplace/libs/application"
	"github.com/pigeaca/DistributedMarketplace/services/worker-manager/internal/config"
	"log"
)

func StartApplication(cfg *config.Config) {
	runner := app.NewApplicationRunner()
	if err := runner.Start(); err != nil {
		log.Println("Error while starting application", err)
	}
}
