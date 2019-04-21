package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/nats"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/cors"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Secret       string `envconfig:"secret" default:"changeme"`
	DashboardURL string `envconfig:"dashboard_url"`
}

func (a *App) Name() string {
	return "apiserver"
}

func (a *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
		nats.NewFromEnv(),
	}
}

func (a *App) Configure() error {
	return envconfig.Process("", a)
}

func (a *App) Run() error {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	{
		v1.Use(a.authMiddleware)

		v1.HandleFunc("/account", a.GetAccount).Methods("GET")
		v1.HandleFunc("/account/sync", a.SyncAccount).Methods("POST")
		v1.HandleFunc("/account/regenerate_apikey", a.RegenerateAPIKeyAccount).Methods("POST")

		v1.HandleFunc("/containers", a.ListContainers).Methods("GET")
		v1.HandleFunc("/containers", a.CreateContainer).Methods("POST")
		v1.HandleFunc("/containers/{name}", a.GetContainer).Methods("GET")
		v1.HandleFunc("/containers/{name}/start", a.StartContainer).Methods("POST")
		v1.HandleFunc("/containers/{name}/stop", a.StopContainer).Methods("POST")

		v1.HandleFunc("/images", a.ListImages).Methods("GET")
		v1.HandleFunc("/images/{name}:{tag}", a.GetImage).Methods("GET")
		v1.HandleFunc("/images/{id}", a.DeleteImage).Methods("DELETE")

		v1.HandleFunc("/plans", a.ListPlans).Methods("GET")
		v1.HandleFunc("/plans/{id}", a.GetPlan).Methods("GET")

		v1.HandleFunc("/meta/tags", a.GetTags).Methods("GET")
		v1.HandleFunc("/meta/images", a.GetImagesName).Methods("GET")
	}
	if err := cors.ApplyMiddleware(router); err != nil {
		return err
	}
	return apps.HandleHTTP(router)
}
