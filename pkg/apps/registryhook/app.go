package registryhook

import (
	"context"
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/registryhook/gc"
	"github.com/expectedsh/expected/pkg/apps/registryhook/registry"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/auth"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/expectedsh/expected/pkg/util/certs"
	"github.com/expectedsh/expected/pkg/util/cors"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	RegistryUrl string `envconfig:"registry_url" default:"http://localhost:5000"`
	Certs       certs.Config
	GcConfig    gc.Config
	Gc          *gc.GarbageCollector
}

func (s *App) Name() string {
	return "registryhook"
}

func (s *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
		stan.NewFromEnv(),
		auth.NewFromEnv(),
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
	s.Gc = gc.New(context.Background(), &gc.Config{
		OlderThan: s.GcConfig.OlderThan,
		Interval:  s.GcConfig.Interval,
		Limit:     s.GcConfig.Limit,
	})
	return nil
}

func (s *App) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/hook", Hook)

	if err := cors.ApplyMiddleware(router); err != nil {
		return err
	}

	if err := apps.HandleQueueSubscription(
		stan.SubjectImageDelete, s.Name(), s.DeleteImage,
		apps.StanQueueGroupOptions(s.Name())...); err != nil {
		return err
	}

	if err := apps.HandleQueueSubscription(
		stan.SubjectImageDeleteLayer, s.Name(), s.DeleteImageLayer,
		apps.StanQueueGroupOptions(s.Name())...); err != nil {
		return err
	}

	s.Gc.Run()
	return apps.HandleHTTP(router)
}
