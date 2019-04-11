package docker

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Service struct {
	client *client.Client
}

func NewFromEnv() *Service {
	cli, err := client.NewEnvClient()
	if err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}
	return &Service{client: cli}
}

func (srv *Service) Name() string {
	return "docker"
}

func (srv *Service) Start() error {
	return nil
}

func (srv *Service) Stop() error {
	return srv.client.Close()
}

func (srv *Service) Started() bool {
	return srv.client != nil
}

func (srv *Service) Client() *client.Client {
	return srv.client
}
