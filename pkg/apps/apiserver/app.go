package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/auth"
	"github.com/expectedsh/expected/pkg/services/controller"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/cors"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Secret       string `envconfig:"secret" default:"changeme"`
	DashboardURL string `envconfig:"dashboard_url"`
}

func (s *App) Name() string {
	return "apiserver"
}

func (s *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
		controller.NewFromEnv(),
		auth.NewFromEnv(),
	}
}

func (s *App) Configure() error {
	return envconfig.Process("", s)
}

func (s *App) Run() error {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	{
		v1.Use(s.authMiddleware)

		v1.HandleFunc("/account", s.GetAccount).Methods("GET")
		v1.HandleFunc("/account/sync", s.SyncAccount).Methods("POST")
		v1.HandleFunc("/account/regenerate_apikey", s.RegenerateAPIKeyAccount).Methods("POST")

		v1.HandleFunc("/containers", s.ListContainers).Methods("GET")
		v1.HandleFunc("/containers", s.CreateContainer).Methods("POST")
		v1.HandleFunc("/containers/{name}", s.GetContainer).Methods("GET")
		v1.HandleFunc("/containers/{name}/start", s.StartContainer).Methods("POST")
		v1.HandleFunc("/containers/{name}/stop", s.StopContainer).Methods("POST")

		v1.HandleFunc("/images", s.ListImages).Methods("GET")
		v1.HandleFunc("/images/{name}:{tag}", s.GetImage).Methods("GET")
		v1.HandleFunc("/images/{id}", s.DeleteImage).Methods("DELETE")

		v1.HandleFunc("/plans", s.ListPlans).Methods("GET")
		v1.HandleFunc("/plans/{id}", s.GetPlan).Methods("GET")

		v1.HandleFunc("/meta/tags", s.GetTags).Methods("GET")
		v1.HandleFunc("/meta/images", s.GetImagesName).Methods("GET")
	}
	if err := cors.ApplyMiddleware(router); err != nil {
		return err
	}
	return apps.HandleHTTP(router)
}
