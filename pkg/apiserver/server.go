package apiserver

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
	"strings"
)

type ApiServer struct {
	Addr       string
	OAuth      *oauth2.Config
	OAuthState string
}

func New(addr, githubClientId, githubClientSecret string) *ApiServer {
	return &ApiServer{
		Addr: addr,
		OAuth: &oauth2.Config{
			ClientID:     githubClientId,
			ClientSecret: githubClientSecret,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user", "user:email"},
		},
		OAuthState: strings.Replace(uuid.New().String(), "-", "", -1),
	}
}

func (s *ApiServer) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/oauth/github", s.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", s.OAuthGithubCallback).Methods("GET")

	return http.ListenAndServe(s.Addr, router)
}
