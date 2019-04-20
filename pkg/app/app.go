package app

import (
	"errors"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Mode string

type App interface {
	// Returns the app name.
	Name() string
	// Register required app services.
	RequiredServices() []services.Service
	// Used for custom app configuration.
	Configure() error
	// Call when your app starts.
	Run() error
}

const (
	ModeProduction  = "PRODUCTION"
	ModeDevelopment = "DEVELOPMENT"
)

var (
	current                      *App
	mode                         Mode = ModeDevelopment
	ErrAppAlreadyStarted              = errors.New("app already started")
	ErrHttpHandlerAlreadyDefined      = errors.New("app http handler already defined")
)

func Current() *App {
	return current
}

func CurrentMode() Mode {
	return mode
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func Start(app App) error {
	if current != nil {
		return ErrAppAlreadyStarted
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go handleStop(ch)

	if srvs := app.RequiredServices(); len(srvs) > 0 {
		for _, service := range srvs {
			logrus.Infof("registering service %v", service.Name())
			services.Register(service)
		}
		logrus.Info("starting services")
		services.Start()
	}

	logrus.Info("configuring the app")
	if err := app.Configure(); err != nil {
		logrus.WithError(err).Fatal("failed to configure app")
	}

	logrus.Info("starting the app")
	if err := app.Run(); err != nil {
		logrus.WithError(err).Fatal("an error occurred while the app running")
	}
	return nil
}
