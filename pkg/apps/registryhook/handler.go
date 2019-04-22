package registryhook

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps/registryhook/registry"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/sirupsen/logrus"
)

func (App) DeleteImage(msg *stan.Msg) {
	ctx := context.Background()
	img := &protocol.DeleteImageEvent{}
	err := proto.Unmarshal(msg.Data, img)
	if err != nil {
		logrus.WithField("data", string(msg.Data)).WithField("nats-subject", msg.Subject).
			WithError(err).
			Error("can't unmarshall proto")
		if err := msg.Ack(); err != nil {
			return
		}
	}

	log := logrus.WithField("image-id", img.Id).
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceId, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag).
		WithField("id", img.Id).
		WithField("tag", img.Tag).
		WithField("nats-subject", msg.Subject)

	layers, err := images.FindLayersByImageId(ctx, img.Id)
	if err != nil {
		log.WithError(err).Error("finding layers by image id")
		return
	}
	// deleting relations between img and layers
	if err := images.DeleteImageLayerByImageID(ctx, img.Id); err != nil {
		log.WithError(err).Error("deleting image_layer rows by img id")
		return
	}
	for _, layer := range layers {
		// If layer is again referenced and unfortunately the repository property is the one that
		// the registry delete, another repository is choose.
		// Else the layer update_at property is updated to be garbage collected.
		layerLog := log.WithField("digest", layer.Digest)
		if cnt, err := images.FindLayerCountReferences(ctx, layer.Digest); err != nil {
			layerLog.WithError(err).Error("finding layer count references")
			return
		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceId, img.Name) {
			if err := images.UpdateLayerRepository(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating repository of layer")
				return
			}
		} else {
			if err := images.UpdateLayer(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating layer")
				return
			}
		}
	}

	// deleting the img at the end to be sure all actions above has been executed
	if err := images.DeleteImageByID(ctx, img.Id); err != nil {
		log.WithError(err).Error("deleting image")
		return
	}

	// Ack before deleting the manifest, if the manifest deletion fail, we ensure to the client that his image
	// has been deleted, so we need to check internally if there is some manifest "alone".
	if err := msg.Ack(); err != nil {
		log.WithError(err).Error("can't ACK")
		return
	}

	msg.Reset()
	status, err := registry.DeleteManifest(fmt.Sprintf("%s/%s", img.NamespaceId, img.Name), img.Digest)
	if err != nil || status == registry.DeleteStatusUnknown {
		log.WithError(err).Error("deleting manifest")
		return
	}
}
