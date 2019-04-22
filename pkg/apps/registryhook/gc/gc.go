package gc

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Config struct {
	Interval  time.Duration `envconfig:"older_than" default:"1h"` // Interval between each run of the gc
	OlderThan time.Duration `envconfig:"limit" default:"1h"`      // OlderThan define since when the layers should not be used
	Limit     int64         `envconfig:"interval" default:"100"`  // Limit define how many layers the gc can handle per run
}

type GarbageCollector struct {
	mutex  sync.Mutex
	logger *logrus.Entry
	ctx    context.Context

	Config *Config
}

func New(ctx context.Context, config *Config) *GarbageCollector {
	return &GarbageCollector{
		ctx:    ctx,
		mutex:  sync.Mutex{},
		Config: config,
		logger: logrus.WithField("task", "garbage-collector"),
	}
}

func (gc *GarbageCollector) Run() {
	go func() {
		for {
			gc.process()
			time.Sleep(gc.Config.Interval)
		}
	}()
}

func (gc *GarbageCollector) process() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	layersDeleted := 0

	gc.logger.Info("start")

	layers, err := images.FindUnusedLayers(gc.ctx, gc.Config.OlderThan, gc.Config.Limit)
	if err != nil {
		gc.logger.WithError(err).Error("can't find unused layers, skip this garbage collector")
	} else {
		for _, layer := range layers {
			err := services.Stan().Publish(stan.SubjectImageDeleteLayer, &protocol.DeleteImageLayerEvent{
				Repository: layer.Repository,
				Digest:     layer.Digest,
			})

			if err == nil {
				layersDeleted++
			}
		}
	}

	gc.logger.
		WithField("layers-in-deletion", fmt.Sprintf("%d/%d", layersDeleted, len(layers))).
		Info("end")
}
