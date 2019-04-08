package services

import (
	"github.com/expectedsh/expected/pkg/services/etcd"
	"github.com/expectedsh/expected/pkg/services/postgres"
	"github.com/expectedsh/expected/pkg/util/backoff"
	"github.com/sirupsen/logrus"
	"time"
)

var registeredServices = map[string]Service{}

type Service interface {
	Name() string
	Start() error
	Stop() error
	Started() bool
}

func Get(service string) Service {
	return registeredServices[service]
}

func Register(service Service) {
	registeredServices[service.Name()] = service
}

func Start() {
	for _, service := range registeredServices {
		entry := logrus.WithField("service", service.Name())
		if err := backoff.New("starting service", service.Start, entry).Execute(); err != nil {
			logrus.WithField("service", service.Name()).WithError(err).Fatalln("unable to start this service")
		}
	}
}

func Stop() {
	for _, service := range registeredServices {
		logrus.WithField("service", service.Name()).Info("stopping service")
		if err := service.Stop(); err != nil {
			logrus.WithField("service", service.Name()).WithError(err).Fatalln("unable to stop this service")
		}
	}
}

func Postgres() *postgres.Service {
	return Get("postgres").(*postgres.Service)
}

func Etcd() *etcd.Service {
	return Get("etcd").(*etcd.Service)
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	})
}
