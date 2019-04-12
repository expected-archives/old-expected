package main

import (
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/docker"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("initializing services")
	services.Register(rabbitmq.NewFromEnv())
	services.Register(postgres.NewFromEnv())
	services.Register(docker.NewFromEnv())
	services.Start()
	defer services.Stop()

	logrus.Infoln("starting scheduler")

	if err := scheduler.Start(); err != nil {
		logrus.WithError(err).Fatalln("unable to start scheduler")
	}
}
