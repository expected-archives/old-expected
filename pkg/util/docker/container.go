package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func GetStats(ctx context.Context, id string) (types.ContainerStats, error) {
	return cli.ContainerStats(ctx, id, false)
}

func GetContainers(ctx context.Context) ([]types.Container, error) {
	args := filters.NewArgs()
	args.Add("status", "running")

	return cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: args,
	})
}
