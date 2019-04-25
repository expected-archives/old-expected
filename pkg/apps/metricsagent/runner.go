package metricsagent

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/collector"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/docker"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/stats"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

func run() {
	for {
		containers, err := docker.GetContainers(context.Background())

		if err != nil {
			logrus.WithError(err).Error("can't get containers")
			time.Sleep(10 * time.Second)
			continue
		}

		group := sync.WaitGroup{}
		group.Add(len(containers))

		packet := make([]stats.Stats, 0)

		for _, ctr := range containers {
			go func(ctr types.Container) {
				st, err := getDockerStats(ctr)
				if err != nil {
					logrus.WithError(err).Error(10 * time.Second)
					group.Done()
				}

				packet = append(packet, stats.FromDockerStats(*st, uuid.New()))
				group.Done()
			}(ctr)
		}

		group.Wait()
		collector.AddPackets(packet)

		time.Sleep(10 * time.Second)
	}
}

func getDockerStats(ctr types.Container) (*types.StatsJSON, error) {
	response, err := docker.GetStats(context.Background(), ctr.ID)
	if err != nil {
		logrus.WithError(err).Error("can't get stats")
		return nil, err
	}
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var v *types.StatsJSON

	if err := dec.Decode(&v); err != nil {
		dec = json.NewDecoder(io.MultiReader(dec.Buffered(), response.Body))

		if err == io.EOF {
			logrus.WithError(err).Error("can't get stats")
			return nil, err
		}
	}

	return v, nil
}
