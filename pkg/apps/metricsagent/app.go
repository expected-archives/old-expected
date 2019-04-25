package metricsagent

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/docker"
	"github.com/expectedsh/expected/pkg/services"
	"io/ioutil"
)

type App struct{}

func (a *App) Name() string {
	return "metricsagent"
}

func (a *App) RequiredServices() []services.Service {
	return []services.Service{}
}

func (a *App) Configure() error {
	return nil
}

func (a *App) Run() error {
	containers, err := docker.GetContainers(context.Background())
	if err != nil {
		return err
	}

	target := containers[0]

	containerStats, _ := docker.GetStats(context.Background(), target.ID)

	bytes, err := ioutil.ReadAll(containerStats.Body)
	fmt.Println(string(bytes))
	return err
}
