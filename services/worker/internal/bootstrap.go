package bootstrap

import (
	app "github.com/pigeaca/DistributedMarketplace/libs/application"
	"github.com/pigeaca/DistributedMarketplace/services/worker/internal/config"
)

func StartApplication(cfg *config.Config) error {
	runner := app.NewApplicationRunner()
	return runner.Start()
}
