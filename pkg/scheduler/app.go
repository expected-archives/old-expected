package scheduler

import (
	"github.com/expectedsh/expected/pkg/app"
	"github.com/expectedsh/expected/pkg/scheduler/docker"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/nats"
	"github.com/expectedsh/expected/pkg/services/postgres"
)

type App struct{}

func (s *App) Name() string {
	return "scheduler"
}

func (s *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
		nats.NewFromEnv(),
	}
}

func (s *App) Configure() error {
	if err := docker.Init(); err != nil {
		return err
	}
	return nil
}

func (s *App) Run() error {
	if err := app.HandleSubscription("container:change-state", handleContainerChangeState); err != nil {
		return err
	}
	return nil
}
