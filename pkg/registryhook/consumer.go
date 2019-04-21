package registryhook

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/authserver/authregistry"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/registry"
	"github.com/sirupsen/logrus"
	"time"
)

func handleDeleteImage(subject, reply string, r *protocol.DeleteImageRequest) {
	log := logrus.WithField("image-id", r.Id)
	img, err := images.FindImageByID(context.Background(), r.Id)
	if err != nil || img == nil {
		log.WithError(err).Error("unable find image")
		return
	}
	log = log.
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag).
		WithField("id", img.ID).
		WithField("tag", img.Tag)

	layers, err := images.FindLayersByImageId(context.Background(), img.ID)
	if err != nil {
		log.WithError(err).Error("finding layers by img id")
		return
	}
	// deleting relations between img and layers
	if err := images.DeleteImageLayerByImageID(context.Background(), img.ID); err != nil {
		log.WithError(err).Error("deleting image_layer rows by img id")
		return
	}
	for _, layer := range layers {
		// If layer is again referenced and unfortunately the repository property is the one that
		// the registry delete, another repository is choose.
		// Else the layer update_at property is updated to be garbage collected.
		layerLog := log.WithField("digest", layer.Digest)
		if cnt, err := images.FindLayerCountReferences(context.Background(), layer.Digest); err != nil {
			layerLog.WithError(err).Error("finding layer count references")
			return
		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceID, img.Name) {
			if err := images.UpdateLayerRepository(context.Background(), layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating repository of layer")
				return
			}
		} else {
			if err := images.UpdateLayer(context.Background(), layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating layer")
				return
			}
		}
	}
	status, err := registry.DeleteManifest(fmt.Sprintf("%s/%s", img.NamespaceID, img.Name), img.Digest)
	if err != nil || status == registry.DeleteStatusUnknown {
		log.WithError(err).Error("can't delete manifest")
		return
	}
	// deleting the img at the end to be sure all actions above has been executed
	if err := images.DeleteImageByID(context.Background(), img.ID); err != nil {
		log.WithError(err).Error("deleting image")
		return
	}

	if err := services.NATS().Client().PublishRequest(subject, reply, &protocol.DeleteImageReply{}); err != nil {
		logrus.WithError(err).Error("failed to send response")
		return
	}
}

func handleGenerateToken(subject, reply string, r *protocol.GenerateTokenRequest) {
	s, err := authregistry.Generate(authregistry.Request{
		Login:   "admin",
		Service: "registry",
	}, []authregistry.AuthorizedScope{
		{
			Scope: authregistry.Scope{
				Type: "repository",
				Name: r.Image,
			},
			AuthorizedActions: []string{"pull"},
		},
	}, ((time.Hour*24)*365)*10)
	if err != nil {
		logrus.WithError(err).Error("failed to generate token")
		return
	}
	if err := services.NATS().Client().PublishRequest(subject, reply, &protocol.GenerateTokenReply{
		Token: s,
	}); err != nil {
		logrus.WithError(err).Error("failed to send response")
	}
}
