package metricsagent

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/docker"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/ingester"
	"github.com/expectedsh/expected/pkg/apps/metricsagent/metrics"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

var ui64 = uint64(0)

func (App) run(ctx context.Context) {
	for {
		waveStartedAt := time.Now()
		ui64++

		select {
		case <-ctx.Done():
			return
		default:
			containers, err := docker.GetContainers(ctx)
			if err != nil {
				logrus.WithError(err).Error("can't get containers")
				break
			}

			group := sync.WaitGroup{}
			group.Add(len(containers))

			// list of metrics for each container at this moment.
			metricList := make([]metrics.Metric, 0)

			for _, ctr := range containers {
				go func(ctr types.Container) {
					st, err := getDockerStats(ctx, ctr)
					if err != nil {
						logrus.WithError(err).Error(10 * time.Second)
						group.Done()
						return
					}

					// translate docker stats to our metric structure.
					// todo change uuid.New() with uuid of the container
					data := metrics.FromDockerStats(*st, uuid.New())
					metricList = append(metricList, data)
					group.Done()
				}(ctr)
			}

			group.Wait()
			ingester.Ingest(metricList)
		}
		logrus.WithField("wave", ui64).WithField("duration", time.Now().Sub(waveStartedAt).Round(time.Millisecond).String()).Infof("metrics processed")
		time.Sleep(10 * time.Second)
	}
}

func getDockerStats(ctx context.Context, ctr types.Container) (*types.StatsJSON, error) {
	response, err := docker.GetStats(ctx, ctr.ID)
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
