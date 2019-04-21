package imageserver

import (
	"context"
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/nats"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/certs"
	"github.com/expectedsh/expected/pkg/util/cors"
	"github.com/expectedsh/expected/pkg/util/registry"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	RegistryUrl string `envconfig:"registry_url" default:"http://localhost:5000"`
	Certs       certs.Config
	GcConfig    apps.Config
	Gc          *apps.GarbageCollector
}

func (s *App) Name() string {
	return "registryhook"
}

func (s *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
		nats.NewFromEnv(),
	}
}

func (s *App) Configure() error {
	if err := envconfig.Process("", s); err != nil {
		return err
	}
	if err := certs.Init(s.Certs); err != nil {
		return err
	}
	registry.Init(s.RegistryUrl)
	s.Gc = apps.New(context.Background(), &apps.Config{
		OlderThan: s.GcConfig.OlderThan,
		Interval:  s.GcConfig.Interval,
		Limit:     s.GcConfig.Limit,
	})
	return nil
}

func (s *App) ConfigureGRPC(server *grpc.Server) {
	protocol.RegisterRegistryHookServer(server, s)
}

func must(err error) {
	logrus.WithError(err).Fatal("unable to start grpc server")
}

func (s *App) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/hook", apps.Hook)

	if err := cors.ApplyMiddleware(router); err != nil {
		return err
	}

	s.Gc.Run()
	go must(apps.HandleGRPC(s))
	return apps.HandleHTTP(router)
}
