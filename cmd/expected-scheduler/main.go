package main

import (
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/nats"
	"github.com/expectedsh/expected/pkg/services/postgres"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("initializing services")
	services.Register(nats.NewFromEnv())
	services.Register(postgres.NewFromEnv())
	services.Start()
	defer services.Stop()

	logrus.Infoln("starting scheduler")

	if err := scheduler.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start scheduler")
	}
}
