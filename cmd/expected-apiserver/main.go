package main

import (
	"github.com/expectedsh/expected/pkg/apiserver"
	"github.com/expectedsh/expected/pkg/models"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Addr     string `envconfig:"addr" default:":3000"`
	Admin    string `envconfig:"admin"`
	Secret   string `envconfig:"secret" default:"changeme"`
	Postgres struct {
		Addr            string        `envconfig:"addr" default:"postgres://postgres:postgres@localhost/postgres?sslmode=disable"`
		ConnMaxLifetime time.Duration `envconfig:"connmaxlifetime" default:"10m"`
		MaxIdleConns    int           `envconfig:"maxidleconns" default:"1"`
		MaxOpenConns    int           `envconfig:"maxopenconns" default:"2"`
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

	logrus.Infoln("initializing the database")
	if err := models.InitDB(config.Postgres.Addr, config.Postgres.ConnMaxLifetime, config.Postgres.MaxIdleConns,
		config.Postgres.MaxOpenConns); err != nil {
		logrus.WithError(err).Fatalln("unable to init the database")
	}

	logrus.Infoln("starting api server")
	server := apiserver.New(config.Addr, config.Secret, config.Github.ClientID, config.Github.ClientSecret,
		config.Admin)

	logrus.Infof("listening on %v\n", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start api server")
	}
}
