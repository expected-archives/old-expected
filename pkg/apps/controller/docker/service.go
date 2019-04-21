package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/models/plans"
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

func ServiceCreate(container *containers.Container) error {
	plan, err := plans.FindPlanByID(context.Background(), container.PlanID)
	if err != nil {
		return err
	}
	if plan == nil || plan.Type != plans.TypeContainer {
		return errors.New("invalid plan id")
	}

	replicas := uint64(1)
	resources := &swarm.Resources{
		MemoryBytes: int64(plan.Metadata["memory"].(float64) * 1024 * 1024),
		NanoCPUs:    int64(plan.Metadata["cpu"].(float64) * 100000000),
	}

	_, err = cli.ServiceCreate(context.Background(), swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: container.ID,
			Labels: map[string]string{
				"traefik.enable":        "true",
				"traefik.domain":        "expected.sh",
				"traefik.frontend.rule": "Host:" + container.Endpoints[0].Endpoint,
				"traefik.port":          "80",
				//"traefik.docker.network": "private",
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: container.Image,
				Env:   convertEnv(container.Environment),
			},
			Resources: &swarm.ResourceRequirements{
				Limits:       resources,
				Reservations: resources,
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			},
		},
		Networks: []swarm.NetworkAttachmentConfig{
			{
				Target: "traefik-net",
			},
		},
	}, types.ServiceCreateOptions{})

	return err
}

func ServiceRemove(container *containers.Container) error {
	return cli.ServiceRemove(context.Background(), container.ID)
}
