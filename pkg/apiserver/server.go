package apiserver

import (
	"github.com/expectedsh/expected/pkg/util/cors"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type ApiServer struct {
	Addr         string
	Secret       string
	DashboardURL string
	OAuth        *oauth2.Config
}

func New(addr, secret, dashboardUrl string) *ApiServer {
	return &ApiServer{
		Addr:         addr,
		Secret:       secret,
		DashboardURL: dashboardUrl,
	}
}

func (s *ApiServer) Start() error {
	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	{
		v1.Use(s.authMiddleware)

		v1.HandleFunc("/account", s.GetAccount).Methods("GET")
		v1.HandleFunc("/account/sync", s.SyncAccount).Methods("POST")
		v1.HandleFunc("/account/regenerate_apikey", s.RegenerateAPIKeyAccount).Methods("POST")

		v1.HandleFunc("/containers", s.ListContainers).Methods("GET")
		v1.HandleFunc("/containers", s.CreateContainer).Methods("POST")

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
	return http.ListenAndServe(s.Addr, router)
}
