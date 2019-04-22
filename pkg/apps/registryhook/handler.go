package registryhook

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps/registryhook/registry"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/sirupsen/logrus"
)

func (App) DeleteImage(msg *stan.Msg) {
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
		WithField("nats-subject", msg.Subject)

	status, err := registry.DeleteManifest(fmt.Sprintf("%s/%s", img.NamespaceId, img.Name), img.Digest)
	if err != nil || status == registry.DeleteStatusUnknown {
		log.WithError(err).Error("deleting manifest")
		return
	}

	if err := msg.Ack(); err != nil {
		log.WithError(err).Error("can't ACK")
		return
	}
}

func (App) DeleteImageLayer(msg *stan.Msg) {
	layer := &protocol.DeleteImageLayerEvent{}
	err := proto.Unmarshal(msg.Data, layer)
	if err != nil {
		logrus.WithField("data", string(msg.Data)).WithField("nats-subject", msg.Subject).
			WithError(err).
			Error("can't unmarshall proto")
		if err := msg.Ack(); err != nil {
			return
		}
	}

	ctx := context.Background()

	log := logrus.WithField("nats-subject", msg.Subject).
		WithField("digest", layer.Digest)

	log.Info("layer deletion requested by gc")

	// todo check layer count before deleting

	//// delete the layer by calling the registry
	//deleteStatus, err := registry.DeleteLayer(layer.Repository, layer.Digest)
	//if err != nil {
	//	log.WithError(err).Error("layer can't be deleted with the registry")
	//	return
	//}
	//
	//// the layer has not been deleted in the registry
	//if deleteStatus == registry.DeleteStatusUnknown || deleteStatus == registry.DeleteStatusNotFound {
	//	log.WithField("delete-status", deleteStatus.String()).
	//		Warn("layer delete status is incoherent")
	//} else {
	//	log.WithField("delete-status", deleteStatus.String()).Info("layer deleted in the registry")
	//}
	//
	//// deleting the layer in the database
	//err = images.DeleteLayerByDigest(ctx, layer.Digest)
	//if err != nil {
	//	log.WithError(err).Error("can't delete layer in postgres")
	//} else {
	//	log.Info("layer deleted in postgres")
	//}
	//
	//if err := images.DeleteImageByDigest(ctx, layer.Digest); err != nil {
	//	log.WithError(err).Error("can't delete image by digest")
	//	return
	//}
	//
	//if err := msg.Ack(); err != nil {
	//	log.WithError(err).Error("can't ACK")
	//	return
	//}
}
