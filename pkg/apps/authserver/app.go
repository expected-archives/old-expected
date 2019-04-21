package authserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/certs"
	"github.com/expectedsh/expected/pkg/util/github"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
)

type App struct {
	DashboardURL string   `envconfig:"dashboard_url"`
	Secret       string   `envconfig:"secret" default:"changeme"`
	Addr         string   `envconfig:"addr" default:":3000"`
	Admins       []string `envconfig:"addr" default:"remicaumette,alexisviscox"`

	Github github.Config
	Certs  certs.Config

	OAuth *oauth2.Config
}

func (a *App) Name() string {
	return "authserver"
}

func (a *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
	}
}

func (a *App) Configure() error {
	if err := envconfig.Process("", a); err != nil {
		return err
	}
	a.OAuth = &oauth2.Config{
		ClientID:     a.Github.ClientID,
		ClientSecret: a.Github.ClientSecret,
		Endpoint:     gh.Endpoint,
		Scopes:       []string{"user", "user:email"},
	}
	if err := certs.Init(a.Certs); err != nil {
		return err
	}
	return nil
}

func (a *App) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/oauth/github", a.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", a.OAuthGithubCallback).Methods("GET")

	router.HandleFunc("/auth/registry", a.AuthRegistry).Methods("GET")

	return apps.HandleHTTP(router)
}
