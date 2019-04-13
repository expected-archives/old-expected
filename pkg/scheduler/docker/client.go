package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var cli *client.Client

func Init() error {
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	cli = c
	return nil
}

func Client() *client.Client {
	return cli
}

func FindServiceByName(name string) (*swarm.Service, error) {
	args := filters.NewArgs()
	args.Add("name", name)
	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{
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

func CreateService(service swarm.ServiceSpec) (types.ServiceCreateResponse, error) {
	return cli.ServiceCreate(context.Background(), swarm.ServiceSpec{

	}, types.ServiceCreateOptions{})
}
