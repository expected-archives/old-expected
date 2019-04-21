package gc

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/util/registry"
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

	now := time.Now()
	layersDeleted := 0

	gc.logger.Info("start")

	layers, err := images.FindUnusedLayers(gc.ctx, gc.Config.OlderThan, gc.Config.Limit)
	if err != nil {
		gc.logger.WithError(err).Error("can't find unused layers, skip this garbage collector")
	} else {
		for _, layer := range layers {

			logger := gc.logger.WithField("digest", layer.Digest)

			// delete the layer by calling the registry
			deleteStatus, err := registry.DeleteLayer(layer.Repository, layer.Digest)
			if err != nil {
				logger.WithError(err).Error("layer can't be deleted with the registry")
				continue
			}

			// the layer has not been deleted in the registry
			if deleteStatus == registry.DeleteStatusUnknown || deleteStatus == registry.DeleteStatusNotFound {
				logger.WithField("delete-status", deleteStatus.String()).
					Warn("layer delete status is incoherent")
			} else {
				logger.WithField("delete-status", deleteStatus.String()).Info("layer deleted in the registry")
			}

			// deleting the layer in the database
			err = images.DeleteLayerByDigest(gc.ctx, layer.Digest)
			if err != nil {
				logger.WithError(err).Error("can't delete layer in postgres")
			} else {
				logger.Info("layer deleted in postgres")
			}

			if err := images.DeleteImageByDigest(gc.ctx, layer.Digest); err != nil {
				logger.WithError(err).Error("can't delete image by digest")
				return
			}

			// layer has be delete at least in the database or in the registry
			if err == nil || deleteStatus == registry.DeleteStatusDeleted {
				layersDeleted++
			}
		}
	}

	gc.logger.
		WithField("gc-duration", time.Now().Sub(now).String()).
		WithField("layers-deleted", fmt.Sprintf("%d/%d", layersDeleted, len(layers))).
		Info("end")
}
