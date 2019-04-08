package main

import (
	"github.com/expectedsh/expected/pkg/authserver"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/certs"
	"github.com/expectedsh/expected/pkg/util/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DashboardURL string   `envconfig:"dashboard_url"`
	Secret       string   `envconfig:"secret" default:"changeme"`
	Admins       []string `envconfig:"admin"`
	Addr         string   `envconfig:"addr" default:":3000"`

	Github github.Config
	Certs  certs.Config
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

	certs.Init(config.Certs)

	logrus.Infoln("starting auth server")
	//
	server := authserver.New(config.Addr, config.Secret, config.DashboardURL, config.Github, config.Admins)

	logrus.Infof("listening on %v", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatal("unable to start auth server")
	}
}
