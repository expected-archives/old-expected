package metrics

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/expectedsh/expected/pkg/apps/agent/docker"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

var ui64 = uint64(0)

func Runner(ctx context.Context) error {
	for {
		waveStartedAt := time.Now()
		ui64++

		select {
		case <-ctx.Done():
			return nil
		default:
			containers, err := docker.GetContainers(ctx)
			if err != nil {
				logrus.WithError(err).Error("can't get containers")
				break
			}

			group := sync.WaitGroup{}
			group.Add(len(containers))

			// list of metrics for each containers running in this host
			// at this moment.
			metricList := make([]Metric, 0)

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
					data := FromDockerStats(*st, uuid.New())
					metricList = append(metricList, data)
					group.Done()
				}(ctr)
			}

			group.Wait()
			ingest(metricList)
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
