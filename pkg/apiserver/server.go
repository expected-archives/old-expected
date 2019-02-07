package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	router.HandleFunc("/oauth/github", s.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", s.OAuthGithubCallback).Methods("GET")
	v1 := router.PathPrefix("/v1").Subrouter()
	{
		v1.Use(s.corsMiddleware, s.authMiddleware)
		v1.HandleFunc("/account", s.GetAccount).Methods("GET")
		v1.HandleFunc("/containers", s.GetContainers).Methods("GET")
		v1.HandleFunc("/containers", s.CreateContainer).Methods("POST")
	}
	return http.ListenAndServe(s.Addr, router)
}
