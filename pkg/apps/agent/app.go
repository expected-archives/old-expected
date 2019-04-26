package agent

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/agent/metrics"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/expectedsh/expected/pkg/util/docker"
)

type App struct{}

func (a *App) Name() string {
	return "agent"
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
	return apps.HandleRunner(metrics.Runner)
}
