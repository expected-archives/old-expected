package app

import (
	"github.com/expectedsh/expected/pkg/app/services"
	"github.com/expectedsh/expected/pkg/app/services/postgres"
	"github.com/expectedsh/expected/pkg/backoff"
	"github.com/sirupsen/logrus"
)

var (
	Name             = "default"
	requiredServices = map[string]services.Service{}
)

func Service(service string) services.Service {
	return requiredServices[service]
}

func AddService(service services.Service) {
	requiredServices[service.Name()] = service
}

func Start() {
	logrus.WithField("app", Name).Info("starting the app")
	for _, service := range requiredServices {
		entry := logrus.WithField("service", service.Name())
		if err := backoff.New("starting service", service.Start, entry).Execute(); err != nil {
			logrus.WithField("service", service.Name()).WithError(err).Fatalln("unable to start this service")
		}
	}
}

func Stop() {
	logrus.WithField("app", Name).Info("stopping the app")
	for _, service := range requiredServices {
		logrus.WithField("service", service.Name()).Info("stopping service")
		if err := service.Stop(); err != nil {
			logrus.WithField("service", service.Name()).WithError(err).Fatalln("unable to stop this service")
		}
	}
}

func Postgres() *postgres.Service {
	return Service(services.Postgres).(*postgres.Service)
}
