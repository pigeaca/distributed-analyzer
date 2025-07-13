package bootstrap

import (
	app "distributed-analyzer/libs/application"
	"distributed-analyzer/services/worker/internal/config"
)

func StartApplication(cfg *config.Config) {
	runner := app.NewApplicationRunner()
	runner.DefaultStart()
}
