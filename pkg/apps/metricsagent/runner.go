package metricsagent

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/collector"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/docker"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/metrics"
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

		// list of metrics for each container at this moment.
		packets := make([][]byte, 0)

		for _, ctr := range containers {
			go func(ctr types.Container) {
				st, err := getDockerStats(ctr)
				if err != nil {
					logrus.WithError(err).Error(10 * time.Second)
					group.Done()
				}

				// translate docker stats to our bytes marshalled metric structure.
				data, err := metrics.FromDockerStats(*st, uuid.New()).MarshalBinary()

				if err != nil {
					logrus.WithError(err).Error(10 * time.Second)
					group.Done()
				}

				packets = append(packets, data)
				group.Done()
			}(ctr)
		}

		group.Wait()
		// lock mutex
		collector.AddPackets(packets)

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
