package hook

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/distribution/notifications"
	"github.com/expectedsh/expected/pkg/apps/imageserver/registry"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/sirupsen/logrus"
)

// onPush is idempotent.
// It is not supposed to create edge effects by replaying the same event.
func onPush(ctx context.Context, event notifications.Event) error {
	namespaceId, name, err := parseRepository(event.Target.Repository)
	if err != nil {
		return nil
	}

	digest := event.Target.Digest.String()

	log := logrus.NewEntry(logrus.StandardLogger()).
		WithField("task", "registry-hook").
		WithField("event", "push").
		WithField("repo", fmt.Sprintf("%s/%s", namespaceId, name)).
		WithField("digest", digest)

	log.Info()

	if err := defaultLayer(ctx, event, digest); err != nil {
		log.WithField("err", err).Error("can't insert default layers")
		return err
	}

	if event.Target.Tag != "" {
		log = log.WithField("tag", event.Target.Tag)

		account, err := accounts.FindAccountByID(ctx, namespaceId)
		if err != nil {
			log.WithError(err).Error("finding account by id")
			return err
		}

		if account == nil {
			log.WithError(err).Error("can't find account")
			return errors.New("can't find account")
		}

		image, err := images.FindImageByInfos(ctx, namespaceId, name, event.Target.Tag, digest)
		if err != nil {
			log.WithError(err).Error("finding image by infos", err)
			return err
		}

		// insert image if not exist
		if image == nil {
			image, err = images.CreateImage(
				ctx,
				name,
				digest,
				account.ID, // todo change this with the real
				event.Target.Tag,
			)
			if err != nil {
				log.WithError(err).WithField("image", image).Error("creating image")
				return err
			}
		}

		log = log.WithField("image", image)

		// get layers by calling the registry manifest
		layers := registry.GetLayers(event.Target.Repository, digest, event.Target.Size)
		if layers == nil {
			log.WithError(err).Error("getting layers")
			return fmt.Errorf("can't get layers with digest %s and repo %s", digest, event.Target.Repository)
		}

		// insert layers and many to many relation with image id <-> layer digest
		err = insertLayers(layers, image.ID)
		if err != nil {
			log.WithError(err).Error("inserting layers")
			return err
		}
	}
	return nil
}

// defaultLayer create a layer if not exist with 0 to count.
// So that if the user stops pushing (ctrl-c or networks problem) the image,
// the garbage collector can delete the unused layer.
func defaultLayer(ctx context.Context, event notifications.Event, digest string) error {
	layer, err := images.FindLayerByDigest(ctx, digest)
	if err != nil {
		return err
	} else if layer == nil {
		_, err := images.CreateLayer(ctx, event.Target.Repository, digest, event.Target.Size)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertLayers will insert layers to table layers and image_layer.
// This will set all layers to the normal count (old + 1).
func insertLayers(layers []images.Layer, imageId string) error {
	err := images.CreateLayers(context.Background(), layers)
	if err != nil {
		return err
	}

	err = images.CreateImageLayer(context.Background(), layers, imageId)
	if err != nil {
		return err
	}
	return nil
}
