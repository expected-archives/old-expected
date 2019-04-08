package main

import (
	"github.com/expectedsh/expected/pkg/authserver"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DashboardURL string   `envconfig:"dashboard_url"`
	Secret       string   `envconfig:"secret" default:"changeme"`
	Admins       []string `envconfig:"admin"`
	Addr         string   `envconfig:"addr" default:":3000"`

	Certs struct {
		PublicKey  string `envconfig:"public_key" default:"./certs/server.crt"`
		PrivateKey string `envconfig:"private_key" default:"./certs/server.key"`
	}

	Github struct {
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

	logrus.Infoln("starting auth server")
	//
	authserver.New(config.Addr, config.Secret,
		config.DashboardURL, config.Github.ClientID, config.Github.ClientSecret, config.Admins)
}
