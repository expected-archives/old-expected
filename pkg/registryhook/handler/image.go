package handler

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/authserver/authregistry"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/util/registry"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type ImageDletete struct {
	Logger *logrus.Entry
}

func (ImageDletete) Name() string {
	return "ImageDeleteRequest"
}

func (i ImageDletete) Handle(msg amqp.Delivery) error {
	ctx := context.Background()
	request := protocol.ImageDeleteRequest{}
	if err := proto.Unmarshal(msg.Body, &request); err != nil {
		return err
	}
	logger := i.Logger.WithField("image-id", request.Id)
	img, err := images.FindImageByID(ctx, request.Id)
	if err != nil {
		logger.WithError(err).Error("unable find image")
		return err
	}
	if img == nil {
		return nil
	}

	logger = i.Logger.
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag).
		WithField("id", img.ID).
		WithField("tag", img.Tag)

	logger.Info("image delete request from rabbitmq")

	layers, err := images.FindLayersByImageId(ctx, img.ID)
	if err != nil {
		logger.WithError(err).Error("finding layers by img id")
		return err
	}
	// deleting relations between img and layers
	if err := images.DeleteImageLayerByImageID(ctx, img.ID); err != nil {
		logger.WithError(err).Error("deleting image_layer rows by img id")
		return err
	}
	for _, layer := range layers {
		// If layer is again referenced and unfortunately the repository property is the one that
		// the registry delete, another repository is choose.
		// Else the layer update_at property is updated to be garbage collected.
		layerLog := logger.WithField("digest", layer.Digest)
		if cnt, err := images.FindLayerCountReferences(ctx, layer.Digest); err != nil {
			layerLog.WithError(err).Error("finding layer count references")
			return err
		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceID, img.Name) {
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
	status, err := registry.DeleteManifest(fmt.Sprintf("%s/%s", img.NamespaceID, img.Name), img.Digest)
	if err != nil || status == registry.DeleteStatusUnknown {
		logger.WithError(err).Error("can't delete manifest")
		return err
	}
	// deleting the img at the end to be sure all actions above has been executed
	if err := images.DeleteImageByID(ctx, img.ID); err != nil {
		logger.WithError(err).Error("deleting image")
		return err
	}
	return nil
}

type ImageToken struct {
}

func (ImageToken) Handle(msg amqp.Delivery) error {
	request := protocol.ImageTokenRequest{}
	if err := proto.Unmarshal(msg.Body, &request); err != nil {
		return err
	}
	s, err := authregistry.Generate(authregistry.Request{
		Login:   "admin",
		Service: "registry",
	}, []authregistry.AuthorizedScope{
		{
			Scope: authregistry.Scope{
				Type: "repository",
				Name: request.ImageName,
			},
			AuthorizedActions: []string{"pull"},
		},
	}, ((time.Hour*24)*365)*10)
	if err != nil {
		return err
	}
	ch, err := services.RabbitMQ().Client().Channel()
	if err != nil {
		return err
	}

	resp, err := proto.Marshal(&protocol.ImageTokenResponse{Token: s})
	if err != nil {
		return err
	}

	return ch.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
		ContentType:   "application/vnd.google.protobuf",
		CorrelationId: msg.CorrelationId,
		Body:          resp,
	})
}

func (ImageToken) Name() string {
	return "ImageTokenRequest"
}
