package main

import (
	"context"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/expectedsh/expected/pkg/registryhook"
	"github.com/expectedsh/expected/pkg/registryhook/auth/token"
	"github.com/expectedsh/expected/pkg/registryhook/gc"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Addr  string `envconfig:"addr" default:":3000"`
	Certs struct {
		PublicKey  string `envconfig:"public_key" default:"./certs/server.crt"`
		PrivateKey string `envconfig:"private_key" default:"./certs/server.key"`
	}
	Gc struct {
		OlderThan time.Duration `envconfig:"gc_older_than" default:"1 hour"`
		Interval  time.Duration `envconfig:"gc_interval" default:"1 hour"`
		Limit     int64         `envconfig:"gc_limit" default:"100"`
	}
}

func main() {
	logrus.Infoln("processing environment configuration")
	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		logrus.WithError(err).Fatalln("unable to parse environment variables")
	}

	services.Register(postgres.NewFromEnv())
	services.Start()
	defer services.Stop()

	if err := accounts.InitDB(services.Postgres().Client()); err != nil {
		logrus.Fatal(err)
	}
	if err := images.InitDB(services.Postgres().Client()); err != nil {
		logrus.Fatal(err)
	}

	token.Init(config.Certs.PublicKey, config.Certs.PrivateKey)

	gc.New(context.Background(), &gc.Options{
		OlderThan: config.Gc.OlderThan,
		Interval:  config.Gc.Interval,
		Limit:     config.Gc.Limit,
	}).Run()

	logrus.Infoln("starting api server")
	server := registryhook.New(config.Addr)

	logrus.Infof("listening on %v", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start api server")
	}
}