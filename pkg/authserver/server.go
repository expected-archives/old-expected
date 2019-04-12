package authserver

import (
	"github.com/expectedsh/expected/pkg/util/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
	"net/http"
)

type AuthServer struct {
	Addr         string
	Secret       string
	Admins       []string
	DashboardURL string
	OAuth        *oauth2.Config
}

func New(addr, secret, dashboardUrl string, githubConfig github.Config) *AuthServer {
	return &AuthServer{
		Addr:         addr,
		Secret:       secret,
		DashboardURL: dashboardUrl,
		OAuth: &oauth2.Config{
			ClientID:     githubConfig.ClientID,
			ClientSecret: githubConfig.ClientSecret,
			Endpoint:     gh.Endpoint,
			Scopes:       []string{"user", "user:email"},
		},
		Admins: []string{"remicaumette", "alexisvisco"},
	}
}

func (s *AuthServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/oauth/github", s.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", s.OAuthGithubCallback).Methods("GET")

	router.HandleFunc("/auth/registry", s.AuthRegistry).Methods("GET")

	return http.ListenAndServe(s.Addr, router)
}

func (s AuthServer) isAdmin(username string) bool {
	for _, name := range s.Admins {
		if name == username {
			return true
		}
	}
	return false
}
