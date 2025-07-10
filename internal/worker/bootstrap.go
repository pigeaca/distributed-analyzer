package bootstrap

import (
	"github.com/distributedmarketplace/internal/worker/config"
	app "github.com/distributedmarketplace/pkg/application"
)

func StartApplication(cfg config.Config) error {
	runner := app.NewApplicationRunner()
	return runner.StartBlocking()
}
