package apiserver

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/expectedsh/expected/pkg/registryhook"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) GetImages(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	imagesStats, err := images.FindImagesStatsByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable to get images list")
		response.ErrorInternal(w)
		return
	}
	if imagesStats == nil {
		imagesStats = []*images.Stats{}
	}
	response.Resource(w, "images", imagesStats)
}

func (s *ApiServer) GetImage(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	imageId := mux.Vars(r)["id"]
	image, err := images.FindImageByID(r.Context(), imageId)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable find image")
		response.ErrorInternal(w)
		return
	}
	if image == nil {
		response.ErrorNotFound(w)
		return
	}
	if image.NamespaceID != account.ID {
		response.ErrorForbidden(w)
		return
	}
	layers, err := images.FindLayersByImageId(r.Context(), imageId)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("get image: unable find layers of images")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "image", images.ImageDetail{Image: image, Layers: layers})
}

func (s *ApiServer) DeleteImage(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	imageId := mux.Vars(r)["id"]
	img, err := images.FindImageByID(r.Context(), imageId)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable find image")
		response.ErrorInternal(w)
		return
	}
	if img == nil {
		response.ErrorNotFound(w)
		return
	}
	if img.NamespaceID != account.ID {
		response.ErrorForbidden(w)
		return
	}
	log := logrus.NewEntry(logrus.StandardLogger()).
		WithField("task", "api-delete-image").
		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag)
	log = log.WithField("id", img.ID).WithField("tag", img.Tag)
	log.Info()

	if img.DeleteMode {
		response.Error(w, "The image is being deleted", http.StatusConflict)
		return
	}

	if err := images.UpdateImageDeleteMode(r.Context(), img.ID); err != nil {
		log.WithError(err).Error("can't update image into delete mode")
		response.ErrorInternal(w)
		return
	}

	if err := registryhook.RequestDeleteImage(imageId); err != nil {
		log.WithError(err).Error("can't publish delete message to rabbit")
		response.ErrorInternal(w)
		return
	}
	w.WriteHeader(202)
}
