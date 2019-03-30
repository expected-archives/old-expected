package main

import (
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/sirupsen/logrus"
)

func main() {
	services.Register(postgres.NewFromEnv())
	services.Start()
	defer services.Stop()

	logrus.Infoln("ok")
}
