package gc

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/expectedsh/expected/pkg/util/registrycli"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Options struct {
	Interval  time.Duration // Interval define the interval between each run of the gc
	Limit     int64         // Limit define how many element max the gc can handle per run
	OlderThan time.Duration // OlderThan define since when the element should not be used
}

type GarbageCollector struct {
	mutex  sync.Mutex
	logger *logrus.Entry
	ctx    context.Context

	Options *Options
}

func New(ctx context.Context, options *Options) *GarbageCollector {
	return &GarbageCollector{
		ctx:     ctx,
		mutex:   sync.Mutex{},
		Options: options,
		logger:  logrus.WithField("service", "garbage-collector"),
	}
}

func (gc *GarbageCollector) Run() {
	go func() {
		for {
			gc.process()
			time.Sleep(gc.Options.Interval)
		}
	}()
}

func (gc *GarbageCollector) process() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	now := time.Now()
	layersDeleted := 0

	gc.logger.Info("starting garbage collection")

	layers, err := images.FindUnusedLayers(gc.ctx, gc.Options.OlderThan, gc.Options.Limit)
	if err != nil {
		gc.logger.WithError(err).Error("can't find unused layers, skip this garbage collector")
	} else {
		for _, layer := range layers {

			// re-compute manually the count to ensure the layer is unused
			count, err := images.FindActualLayerCount(gc.ctx, layer.Digest)
			if err != nil {
				gc.logger.WithError(err).
					WithField("digest", layer.Digest).
					Error("layer count can't be computed, skip this layer")
				continue
			}

			if count <= 0 {

				// delete the layer by calling the registry
				deleteStatus, err := registrycli.DeleteLayer(layer.OriginRepo, layer.Digest)
				if err != nil {
					gc.logger.
						WithError(err).
						WithField("digest", layer.Digest).
						Error("layer can't be deleted with the registry")
					continue
				}

				// the layer has not be delete in the registry
				if deleteStatus == registrycli.Unknown || deleteStatus == registrycli.NotFound {
					gc.logger.
						WithField("digest", layer.Digest).
						WithField("delete-status", deleteStatus.String()).
						Warn("layer delete status is incoherent")
				}

				// delete the layer in the database
				err = images.DeleteLayer(gc.ctx, layer.Digest)
				if err != nil {
					gc.logger.
						WithField("digest", layer.Digest).
						WithField("delete-status", deleteStatus.String()).
						Error("can't delete layer in postgres")
				}

				// layer has be delete at least in the database or in the registry
				if err == nil || deleteStatus == registrycli.Deleted {
					gc.logger.WithField("digest", layer.Digest).Info("layer deleted")
					layersDeleted++
				}

			} else {

				// incoherence founded, update the count attribute of the layer
				if err := images.UpdateLayer(gc.ctx, layer.Digest); err != nil {
					gc.logger.WithField("err", err).
						WithField("digest", layer.Digest).
						Warn("layer can't be updated")
				}
			}
		}
	}

	gc.logger.
		WithField("gc-duration", time.Now().Sub(now).String()).
		WithField("layers-deleted", fmt.Sprintf("%d/%d", layersDeleted, len(layers))).
		Info("end of garbage collection")
}
