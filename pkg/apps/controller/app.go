package controller

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/expectedsh/expected/pkg/util"
	"google.golang.org/grpc"
)

type App struct{}

func (s *App) Name() string {
	return "controller"
}

func (s *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
		stan.NewFromEnv(),
	}
}

func (s *App) Configure() error {
	if err := util.Init(); err != nil {
		return err
	}
	return nil
}

func (s *App) ConfigureGRPC(server *grpc.Server) {
	protocol.RegisterControllerServer(server, s)
}

func (s *App) Run() error {
	return apps.HandleGRPC(s)
}
