package hook

import (
	"context"
	"fmt"
	"github.com/docker/distribution/notifications"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/sirupsen/logrus"
)

// onDelete is call when
func onDelete(ctx context.Context, event notifications.Event) error {
	digest := event.Target.Digest.String()

	fmt.Println("digest ->", digest, event.Target.Tag)

	namespaceId, name, err := parseRepository(event.Target.Repository)
	if err != nil {
		// intended to not block the registry because the repository syntax is not respected
		return nil
	}

	log := logrus.NewEntry(logrus.StandardLogger()).
		WithField("service", "registry-hook").
		WithField("event", "delete").
		WithField("repo", fmt.Sprintf("%s/%s", namespaceId, name)).
		WithField("digest", digest)

	if event.Target.Tag != "" {
		image, err := images.FindImageByInfos(ctx, namespaceId, name, event.Target.Tag, digest)
		if err != nil {
			log.WithError(err).Error("finding image by infos")
			return err
		}

		if image == nil {
			// intended to not block the registry because the image is not in our database
			log.Warn("image not found")
			return nil
		}

		log = log.WithField("id", image.ID).WithField("tag", image.Tag)
		log.Info()

		layers, err := images.FindLayersByImageId(ctx, image.ID)
		if err != nil {
			log.WithError(err).Error("finding layers by image id")
			return err
		}

		// deleting relations between image and layers
		if err := images.DeleteImageLayerByImageID(ctx, image.ID); err != nil {
			log.WithError(err).Error("deleting image_layer rows by image id")
			return err
		}

		for _, layer := range layers {

			// If layer is again referenced and unfortunately the repository property is the one that
			// the registry delete, another repository is choose.
			// Else the layer update_at property is updated to be garbage collected.

			layerLog := log.WithField("digest", layer.Digest)

			if cnt, err := images.FindLayerCountReferences(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("finding layer count references")
				return err
			} else if cnt != 0 && layer.Repository == event.Target.Repository {
				if err := images.UpdateLayerRepository(ctx, layer.Digest); err != nil {
					layerLog.WithError(err).Error("updating repository of layer")
					return err
				}
			} else {
				if err := images.UpdateLayer(ctx, layer.Digest); err != nil {
					layerLog.WithError(err).Error("updating layer")
					return err
				}
			}
		}

		// deleting the image at the end to be sure all actions above has been executed
		if err := images.DeleteImage(ctx, image.ID); err != nil {
			log.WithError(err).Error("deleting image")
			return err
		}
	}
	return nil
}
