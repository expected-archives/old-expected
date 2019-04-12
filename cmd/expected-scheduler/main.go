package main

import (
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("initializing services")
	services.Register(rabbitmq.NewFromEnv())
	services.Start()
	defer services.Stop()

}
