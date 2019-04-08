package authserver

import (
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type ApiServer struct {
	Addr         string
	Secret       string
	Admins       []string
	DashboardURL string
	OAuth        *oauth2.Config
}

func New(addr, secret, dashboardUrl, githubClientId, githubClientSecret string, admin []string) *ApiServer {
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
		Admins: admin,
	}
}

func (s *ApiServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/oauth/github", s.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", s.OAuthGithubCallback).Methods("GET")
}

func (s ApiServer) isAdmin(username string) bool {
	for _, name := range s.Admins {
		if name == username {
			return true
		}
	}
	return false
}
