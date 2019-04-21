package apps

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
	current                      App
	ErrAppAlreadyStarted         = errors.New("app already started")
	ErrHttpHandlerAlreadyDefined = errors.New("app http handler already defined")
	ErrGRPCHandlerAlreadyDefined = errors.New("app grpc handler already defined")
)

func Current() App {
	return current
}

func CurrentMode() Mode {
	if value := GetEnvOrDefault("MODE", string(ModeDevelopment)); value == string(ModeProduction) {
		return ModeProduction
	}
	return ModeDevelopment
}

func GetEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func start(app App) {
	if CurrentMode() == ModeProduction {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}

	if srvs := app.RequiredServices(); len(srvs) > 0 {
		for _, service := range srvs {
			logrus.Infof("registering service %v", service.Name())
			services.Register(service)
		}
		logrus.Info("starting services")
		services.Start()
	}

	logrus.Info("starting the app")
	if err := app.Configure(); err != nil {
		logrus.WithError(err).Fatal("failed to configure app")
	}
	if err := app.Run(); err != nil {
		logrus.WithError(err).Fatal("an error occurred while the app running")
	}
}

func Start(app App) error {
	if current != nil {
		return ErrAppAlreadyStarted
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go handleStop(ch)

	start(app)
	return nil
}
