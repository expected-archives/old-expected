package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type ApiServer struct {
	Addr         string
	Secret       string
	Admin        string
	DashboardURL string
	OAuth        *oauth2.Config
}

func New(addr, secret, admin, dashboardUrl, githubClientId, githubClientSecret string) *ApiServer {
	return &ApiServer{
		Addr:         addr,
		Secret:       secret,
		DashboardURL: dashboardUrl,
		OAuth: &oauth2.Config{
			ClientID:     githubClientId,
			ClientSecret: githubClientSecret,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user", "user:email"},
		},
		Admin: admin,
	}
}

func (s *ApiServer) Start() error {
	router := mux.NewRouter()

	router.Use(s.authMiddleware)

	router.HandleFunc("v1/account", s.GetAccount).Methods("GET")
	router.HandleFunc("v1/account/sync", s.SyncAccount).Methods("POST")
	router.HandleFunc("v1/account/regenerate_apikey", s.RegenerateAPIKeyAccount).Methods("POST")

	router.HandleFunc("v1/containers", s.GetContainers).Methods("GET")
	router.HandleFunc("v1/containers", s.CreateContainer).Methods("POST")

	if err := corsMiddleware(router); err != nil {
		return err
	}
	return http.ListenAndServe(s.Addr, router)
}
