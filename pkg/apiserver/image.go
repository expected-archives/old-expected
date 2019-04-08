package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) GetImages(w http.ResponseWriter, r *http.Request) {
	account := session.GetAccount(r)
	imagesStats, err := images.FindImagesStatsByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Error("unable to get images list")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "images", imagesStats)
}

func (s *ApiServer) GetImage(w http.ResponseWriter, r *http.Request) {
	account := session.GetAccount(r)
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

// todo use rabbitmq
//func (s *ApiServer) DeleteImage(w http.ResponseWriter, r *http.Request) {
//	account := session.GetAccount(r)
//	imageId := mux.Vars(r)["id"]
//	img, err := images.FindImageByID(r.Context(), imageId)
//	if err != nil {
//		logrus.WithError(err).WithField("account", account.ID).Error("deleting img: unable find img")
//		response.ErrorInternal(w)
//		return
//	}
//	if img == nil {
//		response.ErrorNotFound(w)
//		return
//	}
//	if img.NamespaceID != account.ID {
//		response.ErrorForbidden(w)
//		return
//	}
//	log := logrus.NewEntry(logrus.StandardLogger()).
//		WithField("service", "delete-image").
//		WithField("repo", fmt.Sprintf("%s/%s", img.NamespaceID, img.Name)).
//		WithField("digest", img.Digest).
//		WithField("tag", img.Tag)
//
//	log = log.WithField("id", img.ID).WithField("tag", img.Tag)
//	log.Info()
//
//	layers, err := images.FindLayersByImageId(r.Context(), img.ID)
//	if err != nil {
//		log.WithError(err).Error("finding layers by img id")
//		response.ErrorInternal(w)
//		return
//	}
//
//	// deleting relations between img and layers
//	if err := images.DeleteImageLayerByImageID(r.Context(), img.ID); err != nil {
//		log.WithError(err).Error("deleting image_layer rows by img id")
//		response.ErrorInternal(w)
//		return
//	}
//
//	for _, layer := range layers {
//
//		// If layer is again referenced and unfortunately the repository property is the one that
//		// the registry delete, another repository is choose.
//		// Else the layer update_at property is updated to be garbage collected.
//
//		layerLog := log.WithField("digest", layer.Digest)
//
//		if cnt, err := images.FindLayerCountReferences(r.Context(), layer.Digest); err != nil {
//			layerLog.WithError(err).Error("finding layer count references")
//			response.ErrorInternal(w)
//			return
//		} else if cnt != 0 && layer.Repository == fmt.Sprintf("%s/%s", img.NamespaceID, img.Name) {
//			if err := images.UpdateLayerRepository(r.Context(), layer.Digest); err != nil {
//				layerLog.WithError(err).Error("updating repository of layer")
//				response.ErrorInternal(w)
//				return
//			}
//		} else {
//			if err := images.UpdateLayer(r.Context(), layer.Digest); err != nil {
//				layerLog.WithError(err).Error("updating layer")
//				response.ErrorInternal(w)
//				return
//			}
//		}
//	}
//
//	status, err := registry.DeleteManifest(fmt.Sprintf("%s/%s", img.NamespaceID, img.Name), img.Digest)
//	if err != nil || status == registry.DeleteStatusUnknown {
//		log.WithError(err).Error("can't delete manifest")
//		return
//	}
//
//	// deleting the img at the end to be sure all actions above has been executed
//	if err := images.DeleteImage(r.Context(), img.ID); err != nil {
//		log.WithError(err).Error("deleting img")
//		response.ErrorInternal(w)
//		return
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//	return
//}
