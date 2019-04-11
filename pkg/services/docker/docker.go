package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
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

func (srv *Service) FindService(ctx context.Context, name string) (*swarm.Service, error) {
	args := filters.NewArgs()
	args.Add("name", name)
	services, err := srv.client.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: args,
	})
	if err != nil {
		return nil, err
	}
	if len(services) > 0 {
		return &services[0], nil
	} else {
		return nil, nil
	}
}
