package apiserver

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
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

func cors(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func (s *ApiServer) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/oauth/github", s.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", s.OAuthGithubCallback).Methods("GET")
	router.HandleFunc("/account", cors(s.Account)).Methods("GET")
	return http.ListenAndServe(s.Addr, router)
}
