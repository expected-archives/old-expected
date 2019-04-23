package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/apiserver/request"
	"github.com/expectedsh/expected/pkg/apps/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/auth"
	"github.com/expectedsh/expected/pkg/services/controller"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/expectedsh/expected/pkg/util/cors"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"net/http"
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
		stan.NewFromEnv(),
	}
}

func (s *App) Configure() error {
	return envconfig.Process("", s)
}

func (s *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.ErrorForbidden(w)
			return
		}
		account, err := accounts.FindAccountByAPIKey(r.Context(), header)
		if err != nil {
			logrus.WithError(err).Errorln("unable to find account")
			response.ErrorInternal(w)
			return
		}
		if account == nil {
			response.ErrorForbidden(w)
			return
		}
		request.SetAccount(r, account)
		next.ServeHTTP(w, r)
	})
}

func (s *App) Run() error {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	{
		v1.Use(s.AuthMiddleware)

		v1.HandleFunc("/account", s.GetAccount).Methods("GET")
		v1.HandleFunc("/account/sync", s.SyncAccount).Methods("POST")
		v1.HandleFunc("/account/regenerate_apikey", s.RegenerateAPIKeyAccount).Methods("POST")

		v1.HandleFunc("/containers", s.ListContainers).Methods("GET")
		v1.HandleFunc("/containers", s.CreateContainer).Methods("POST")
		containers := v1.PathPrefix("/containers/{name}").Subrouter()
		{
			containers.Use(s.ContainerMiddleware)

			containers.HandleFunc("/", s.GetContainer).Methods("GET")
			containers.HandleFunc("/logs", s.GetContainerLogs).Methods("GET")
			containers.HandleFunc("/start", s.StartContainer).Methods("POST")
			containers.HandleFunc("/stop", s.StopContainer).Methods("POST")
		}

		v1.HandleFunc("/images", s.ListImages).Methods("GET")
		v1.HandleFunc("/images/{name}:{tag}", s.GetImage).Methods("GET")
		v1.HandleFunc("/images/{name}:{tag}/{digest}", s.DeleteImage).Methods("DELETE")

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
