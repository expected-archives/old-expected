package apiserver

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps/apiserver/request"
	"github.com/expectedsh/expected/pkg/apps/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *App) ListImages(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	imageSummaries, err := images.FindImagesSummariesByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable to get images list")
		response.ErrorInternal(w)
		return
	}
	if imageSummaries == nil {
		imageSummaries = []*images.ImageSummary{}
	}
	response.Resource(w, "images", imageSummaries)
}

func (s *App) GetImage(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	name := mux.Vars(r)["name"]
	tag := mux.Vars(r)["tag"]
	image, err := images.FindImageDetail(r.Context(), account.ID, name, tag)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable find image detail")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "image", image)
}

func (s *App) DeleteImage(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)

	name := mux.Vars(r)["name"]
	tag := mux.Vars(r)["tag"]
	digest := mux.Vars(r)["digest"]

	// Getting image and erroring if the image is already deleted
	img, err := images.FindImageByInfos(r.Context(), account.ID, name, tag, digest)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable find image")
		response.ErrorInternal(w)
		return
	}
	if img == nil {
		response.Error(w, http.StatusConflict, "The image is deleted or in the process of being deleted.")
		return
	}
	if img.NamespaceID != account.ID {
		response.ErrorForbidden(w)
		return
	}

	log := logrus.WithField("image-id", img.ID).
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag).
		WithField("id", img.NamespaceID).
		WithField("endpoint", "image-delete")

	layers, err := images.FindLayersByImageId(r.Context(), img.ID)
	if err != nil {
		log.WithError(err).Error("finding layers by image id")
		response.ErrorInternal(w)
		return
	}
	// Deleting relations between img and layers
	if err := images.DeleteImageLayerByImageID(r.Context(), img.ID); err != nil {
		log.WithError(err).Error("deleting image_layer rows by img id")
		response.ErrorInternal(w)
		return
	}
	for _, layer := range layers {

		// If layer is again referenced and unfortunately the repository property is the one that
		// the registry delete, another repository is choose.
		// Else the layer update_at property is updated to be garbage collected.

		layerLog := logrus.WithField("digest", layer.Digest)
		if cnt, err := images.FindLayerCountReferences(context.Background(), layer.Digest); err != nil {
			layerLog.WithError(err).Error("finding layer count references")
			response.ErrorInternal(w)
			return
		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceID, img.Name) {
			if err := images.UpdateLayerRepository(r.Context(), layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating repository of layer")
				response.ErrorInternal(w)
				return
			}
		} else {
			if err := images.UpdateLayer(r.Context(), layer.Digest); err != nil {
				layerLog.WithError(err).Error("updating layer")
				response.ErrorInternal(w)
				return
			}
		}
	}

	// deleting the img at the end to be sure all actions above has been executed
	if err := images.DeleteImageByID(r.Context(), img.ID); err != nil {
		log.WithError(err).Error("deleting image")
		response.ErrorInternal(w)
		return
	}

	event := &protocol.DeleteImageEvent{
		Id: img.ID, NamespaceId: img.NamespaceID, Name: img.Name,
		Tag: img.Tag, Digest: img.Digest,
	}

	bytes, err := proto.Marshal(event)
	if err != nil {
		response.ErrorInternal(w)
		return
	}

	// Request a manifest deletion
	if err := services.Stan().Client().Publish(stan.SubjectImageDelete, bytes); err != nil {
		response.ErrorInternal(w)
		return
	}

	w.WriteHeader(202)
}
