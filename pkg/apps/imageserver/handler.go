package imageserver

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps/authserver/authregistry"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/util/registry"
	"github.com/sirupsen/logrus"
	"time"
)

func (App) DeleteImage(ctx context.Context, r *protocol.DeleteImageRequest) (*protocol.DeleteImageReply, error) {
	log := logrus.WithField("image-id", r.Id)
	img, err := images.FindImageByID(ctx, r.Id)
	if err != nil || img == nil {
		log.WithError(err).Error("unable find image")
		return nil, err
	}
	log = log.
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag).
		WithField("id", img.ID).
		WithField("tag", img.Tag)

	layers, err := images.FindLayersByImageId(ctx, img.ID)
	if err != nil {
		log.WithError(err).Error("finding layers by image id")
		return nil, err
	}
	// deleting relations between img and layers
	if err := images.DeleteImageLayerByImageID(ctx, img.ID); err != nil {
		log.WithError(err).Error("deleting image_layer rows by img id")
		return nil, err
	}
	for _, layer := range layers {
		// If layer is again referenced and unfortunately the repository property is the one that
		// the registry delete, another repository is choose.
		// Else the layer update_at property is updated to be garbage collected.
		layerLog := log.WithField("digest", layer.Digest)
		if cnt, err := images.FindLayerCountReferences(ctx, layer.Digest); err != nil {
			layerLog.WithError(err).Error("finding layer count references")
			return nil, err
		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceID, img.Name) {
			if err := images.UpdateLayerRepository(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating repository of layer")
				return nil, err
			}
		} else {
			if err := images.UpdateLayer(ctx, layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating layer")
				return nil, err
			}
		}
	}

	status, err := registry.DeleteManifest(fmt.Sprintf("%s/%s", img.NamespaceID, img.Name), img.Digest)
	if err != nil || status == registry.DeleteStatusUnknown {
		log.WithError(err).Error("can't delete manifest")
		return nil, err
	}
	// deleting the img at the end to be sure all actions above has been executed
	if err := images.DeleteImageByID(ctx, img.ID); err != nil {
		log.WithError(err).Error("deleting image")
		return nil, err
	}

	return &protocol.DeleteImageReply{}, nil
}

func (App) GenerateToken(ctx context.Context, r *protocol.GenerateTokenRequest) (*protocol.GenerateTokenReply, error) {
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
	}, time.Duration(r.Duration))

	if err != nil {
		logrus.WithError(err).Error("failed to generate token")
		return nil, err
	}

	return &protocol.GenerateTokenReply{
		Token: s,
	}, nil
}
