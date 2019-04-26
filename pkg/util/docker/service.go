package docker

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/models/plans"
	"github.com/expectedsh/expected/pkg/util"
	"io"
)

func ServiceFindByName(ctx context.Context, name string) (*swarm.Service, error) {
	args := filters.NewArgs()
	args.Add("name", name)
	services, err := util.cli.ServiceList(ctx, types.ServiceListOptions{
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

func ServiceCreate(ctx context.Context, container *containers.Container) error {
	plan, err := plans.FindPlanByID(ctx, container.PlanID)
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

	_, err = util.cli.ServiceCreate(ctx, swarm.ServiceSpec{
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
				Env:   util.convertEnv(container.Environment),
			},
			Resources: &swarm.ResourceRequirements{
				Limits:       resources,
				Reservations: resources,
			},
			LogDriver: &swarm.Driver{
				Name: "json-file",
				Options: map[string]string{
					"max-size": "10m",
					"max-file": "3",
				},
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

func ServiceRemove(ctx context.Context, container *containers.Container) error {
	return util.cli.ServiceRemove(ctx, container.ID)
}

func ServiceGetLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	return util.cli.ServiceLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Details:    true,
		Timestamps: true,
	})
}
