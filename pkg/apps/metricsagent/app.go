package metricsagent

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/docker"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"google.golang.org/grpc"
)

type App struct{}

func (a *App) Name() string {
	return "metricsagent"
}

func (a *App) RequiredServices() []services.Service {
	return []services.Service{}
}
func (a *App) ConfigureGRPC(server *grpc.Server) {
	protocol.RegisterMetricsServer(server, a)
}

func (a *App) Configure() error {
	return docker.Init()
}

func (a *App) Run() error {
	go run()

	return apps.HandleGRPC(a)
}
