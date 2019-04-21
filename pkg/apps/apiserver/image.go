package apiserver

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *app.App) ListImages(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	imagesStats, err := images.FindImagesSummariesByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable to get images list")
		apps.ErrorInternal(w)
		return
	}
	if imagesStats == nil {
		imagesStats = []*images.ImageSummary{}
	}
	apps.Resource(w, "images", imagesStats)
}

func (a *app.App) GetImage(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	name := mux.Vars(r)["name"]
	tag := mux.Vars(r)["tag"]
	image, err := images.FindImageDetail(r.Context(), account.ID, name, tag)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable find image detail")
		apps.ErrorInternal(w)
		return
	}
	apps.Resource(w, "image", image)
}

func (a *app.App) DeleteImage(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	id := mux.Vars(r)["id"]
	if _, err := uuid.Parse(id); err != nil {
		apps.ErrorBadRequest(w, "Invalid image id.", nil)
		return
	}
	img, err := images.FindImageByID(r.Context(), id)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable find image")
		apps.ErrorInternal(w)
		return
	}
	if img == nil {
		apps.ErrorNotFound(w)
		return
	}
	if img.NamespaceID != account.ID {
		apps.ErrorForbidden(w)
		return
	}
	log := logrus.
		WithField("task", "api-delete-image").
		WithField("repo", fmt.Sprintf("%v/%v", img.NamespaceID, img.Name)).
		WithField("digest", img.Digest).
		WithField("tag", img.Tag).
		WithField("id", img.ID).
		WithField("tag", img.Tag)
	if img.DeleteMode {
		apps.Error(w, http.StatusConflict, "The image is being deleted")
		return
	}
	if err := images.UpdateImageDeleteMode(r.Context(), img.ID); err != nil {
		log.WithError(err).Error("can't update image into delete mode")
		apps.ErrorInternal(w)
		return
	}
	if _, err := apps.RequestDeleteImage(r.Context(), id); err != nil {
		log.WithError(err).Error("can't publish delete message")
		apps.ErrorInternal(w)
		return
	}
	w.WriteHeader(202)
}
