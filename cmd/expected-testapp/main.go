package main

import (
	"github.com/expectedsh/expected/pkg/app"
	"github.com/expectedsh/expected/pkg/app/services/postgres"
	"github.com/sirupsen/logrus"
)

func main() {
	app.AddService(postgres.NewFromEnv())
	app.Start()
	defer app.Stop()

	logrus.Infoln("ok")
}
