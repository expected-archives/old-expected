package authserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/certs"
	"github.com/expectedsh/expected/pkg/util/github"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
	"google.golang.org/grpc"
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

func (s *App) Name() string {
	return "authserver"
}

func (s *App) RequiredServices() []services.Service {
	return []services.Service{
		postgres.NewFromEnv(),
	}
}

func (s *App) Configure() error {
	if err := envconfig.Process("", s); err != nil {
		return err
	}
	s.OAuth = &oauth2.Config{
		ClientID:     s.Github.ClientID,
		ClientSecret: s.Github.ClientSecret,
		Endpoint:     gh.Endpoint,
		Scopes:       []string{"user", "user:email"},
	}
	if err := certs.Init(s.Certs); err != nil {
		return err
	}
	return nil
}

func (s *App) ConfigureGRPC(server *grpc.Server) {
	protocol.RegisterAuthServer(server, s)
}

func (s *App) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/oauth/github", s.OAuthGithub).Methods("GET")
	router.HandleFunc("/oauth/github/callback", s.OAuthGithubCallback).Methods("GET")

	router.HandleFunc("/auth/registry", s.AuthRegistry).Methods("GET")

	go func() {
		if err := apps.HandleGRPC(s); err != nil {
			logrus.Fatal(err)
			return
		}
	}()

	return apps.HandleHTTP(router)
}
