package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

func ServiceFindByName(name string) (*swarm.Service, error) {
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

func ServiceCreate(service swarm.ServiceSpec) (types.ServiceCreateResponse, error) {
	return cli.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
}
