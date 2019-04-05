package hook

import (
	"context"
	"fmt"
	"github.com/docker/distribution/notifications"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/expectedsh/expected/pkg/util/registrycli"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

// onPush is idempotent.
// It is not supposed to create edge effects by replaying the same event.
func onPush(ctx context.Context, event notifications.Event) error {

	namespaceId, name, err := parseRepository(event.Target.Repository)
	if err != nil {
		return err
	}

	digest := event.Target.Digest.String()

	log := logrus.NewEntry(logrus.StandardLogger()).
		WithField("repo", fmt.Sprintf("%s/%s", namespaceId, name)).
		WithField("tag", event.Target.Tag).
		WithField("digest", digest).
		WithField("event", "push")

	log.Info("layer is being push")

	if err := defaultLayer(ctx, event, digest); err != nil {
		log.WithField("err", err).Error("can't insert default layers")
		return err
	}

	if event.Target.Tag != "" {

		account, err := accounts.FindByID(ctx, namespaceId)
		if err != nil {
			log.Error("can't find account")
			return err
		}

		image, err := images.FindImageByInfos(ctx, namespaceId, name, event.Target.Tag, digest)
		if err != nil {
			log.WithField("err", err).Error("can't find image by repo+tag with digest ", err)
			return err
		}

		// insert image if not exist
		if image == nil {
			image, err = images.Create(
				ctx,
				name,
				digest,
				account.ID, // todo change this with the real
				event.Target.Tag,
			)
			if err != nil {
				log.WithField("err", err).WithField("image", image).Error("can't create image")
				return err
			}
		}

		log = log.WithField("image", image)

		// get layers by calling the registry manifest
		layers := registrycli.GetLayers(event.Target.Repository, digest, event.Target.Size)
		if layers == nil {
			log.WithField("err", err).Error("can't get layers")
			return fmt.Errorf("can't get layers with digest %s and repo %s", digest, event.Target.Repository)
		}

		// insert layers and many to many relation with image id <-> layer digest
		err = insertLayers(layers, image.ID)
		if err != nil {
			log.WithField("err", err).Error("can't insert layers")
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
		_, err := images.CreateLayer(ctx, digest, event.Target.Size)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertLayers will insert layers to table layers and image_layer.
// This will set all layers to the normal count (old + 1).
func insertLayers(layers []images.Layer, imageId string) error {
	err := images.CreateLayers(context.Background(), layers, imageId)
	if err != nil {
		return err
	}

	err = images.CreateImageLayer(context.Background(), layers, imageId)
	if err != nil {
		return err
	}
	return nil
}

// parseRepository return the namespace id and the name of the image.
// Can throw an error only if the repository is malformed.
func parseRepository(repo string) (namespaceID, name string, err error) {
	str := strings.Split(repo, "/")
	if len(str) != 2 {
		return "", "", errors.New("repository is malformed")
	}
	return str[0], str[1], nil
}
