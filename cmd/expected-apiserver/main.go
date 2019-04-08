package main

import (
	"github.com/expectedsh/expected/pkg/apiserver"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr         string `envconfig:"addr" default:":3000"`
	Secret       string `envconfig:"secret" default:"changeme"`
	Admin        string `envconfig:"admin"`
	DashboardURL string `envconfig:"dashboard_url"`
	Github       struct {
		ClientID     string `envconfig:"client_id"`
		ClientSecret string `envconfig:"client_secret"`
	}
}

func main() {
	logrus.Infoln("processing environment configuration")
	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		logrus.WithError(err).Fatalln("unable to parse environment variables")
	}

	logrus.Infoln("initializing services")
	services.Register(postgres.NewFromEnv())
	services.Start()
	defer services.Stop()

	logrus.Infoln("starting api server")
	server := apiserver.New(config.Addr, config.Secret, config.Admin,
		config.DashboardURL, config.Github.ClientID, config.Github.ClientSecret)

	logrus.Infof("listening on %v", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start api server")
	}
}
