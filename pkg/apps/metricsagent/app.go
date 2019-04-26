package metricsagent

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/docker"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/stan"
)

type App struct{}

func (a *App) Name() string {
	return "metricsagent"
}

func (a *App) RequiredServices() []services.Service {
	return []services.Service{
		stan.NewFromEnv(),
	}
}

func (a *App) Configure() error {
	return docker.Init()
}

func (a *App) Run() error {
	return apps.HandleRunner(a.run)
}
