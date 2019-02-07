package main

import (
	"database/sql"
	"time"

	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/apiserver"
	"github.com/expectedsh/expected/pkg/containers"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr         string `envconfig:"addr" default:":3000"`
	Secret       string `envconfig:"secret" default:"changeme"`
	Admin        string `envconfig:"admin"`
	DashboardURL string `envconfig:"dashboard_url"`
	Postgres     struct {
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

func initDB(addr string, connMaxLifetime time.Duration, maxIdleConns, maxOpenConns int) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	if err = accounts.InitDB(db); err != nil {
		return nil, err
	}
	if err = containers.InitDB(db); err != nil {
		return nil, err
	}
	return db, err
}

func main() {
	logrus.Infoln("processing environment configuration")
	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		logrus.WithError(err).Fatalln("unable to parse environment variables")
	}

	logrus.Infoln("initializing the database")
	db, err := initDB(config.Postgres.Addr, config.Postgres.ConnMaxLifetime,
		config.Postgres.MaxIdleConns, config.Postgres.MaxOpenConns)
	if err != nil {
		logrus.WithError(err).Fatalln("unable to init the database")
	}
	defer db.Close()

	logrus.Infoln("starting api server")
	server := apiserver.New(config.Addr, config.Secret, config.Admin,
		config.DashboardURL, config.Github.ClientID, config.Github.ClientSecret)

	logrus.Infof("listening on %v\n", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start api server")
	}
}
