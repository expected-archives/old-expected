package controller

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/nats"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"google.golang.org/grpc"
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
	if err := apps.Init(); err != nil {
		return err
	}
	return nil
}

func (s *App) ConfigureGRPC(server *grpc.Server) {
	protocol.RegisterSchedulerServer(server, s)
}

func (s *App) Run() error {
	return apps.HandleGRPC(s)
}
