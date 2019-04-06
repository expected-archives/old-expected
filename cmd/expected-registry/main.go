package main

import (
	"context"
	"database/sql"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/containers"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/expectedsh/expected/pkg/registryhook"
	"github.com/expectedsh/expected/pkg/registryhook/auth/token"
	"github.com/expectedsh/expected/pkg/registryhook/gc"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/registrycli"
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
	if err = images.InitDB(db); err != nil {
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
		OlderThan: time.Second,
		Interval:  time.Minute * 5,
		Limit:     10,
	}).Run()

	status, e := registrycli.DeleteManifest(
		"b431e1c9-3b04-42bd-83f5-47c05e49c70c",
		"registry",
		"sha256:b1165286043f2745f45ea637873d61939bff6d9a59f76539d6228abf79f87774")
	logrus.Info("status", status.String(), e)

	logrus.Infoln("starting api server")
	server := registryhook.New(config.Addr)

	logrus.Infof("listening on %v", config.Addr)
	if err := server.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start api server")
	}
}
