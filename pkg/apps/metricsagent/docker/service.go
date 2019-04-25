package docker

import (
	"context"
	"github.com/docker/docker/api/types"
)

func GetStats(ctx context.Context, id string) (types.ContainerStats, error) {
	return cli.ContainerStats(ctx, id, false)
}

func GetContainers(ctx context.Context) ([]types.Container, error) {
	return cli.ContainerList(ctx, types.ContainerListOptions{})
}
