package main

import (
	"context"
	"github.com/expectedsh/expected/pkg/registryhook"
	"github.com/expectedsh/expected/pkg/registryhook/gc"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	"github.com/expectedsh/expected/pkg/util/certs"
	"github.com/expectedsh/expected/pkg/util/registry"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr        string `envconfig:"addr" default:":3000"`
	RegistryUrl string `envconfig:"registry_url" default:"http://localhost:5000"`

	Certs certs.Config
	Gc    gc.Config
}

func main() {
	logrus.Infoln("processing environment configuration")
	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		logrus.WithError(err).Fatalln("unable to parse environment variables")
	}

	services.Register(postgres.NewFromEnv())
	services.Register(rabbitmq.NewFromEnv())
	services.Start()
	defer services.Stop()

	certs.Init(config.Certs)
	registry.Init(config.RegistryUrl)

	gc.New(context.Background(), &gc.Config{
		OlderThan: config.Gc.OlderThan,
		Interval:  config.Gc.Interval,
		Limit:     config.Gc.Limit,
	}).Run()

	go registryhook.Start()
	logrus.Infoln("starting api server")
	server := registryhook.New(config.Addr)

	logrus.Infof("listening on %v", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start api server")
	}
}
