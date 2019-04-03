package registryserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/notifications"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/images"
	"github.com/expectedsh/expected/pkg/util/registrycli"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func Hook(res http.ResponseWriter, req *http.Request) {

	logrus.Infoln("Hook incoming")

	if req.Body == nil {
		http.Error(res, "ignoring request. Required non-empty request body", http.StatusOK)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != notifications.EventsMediaType {
		http.Error(res, fmt.Sprintf("ignoring request. Required mimetype is %q but got %q", notifications.EventsMediaType, contentType), http.StatusOK)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var envelope notifications.Envelope
	err := decoder.Decode(&envelope)
	if err != nil {
		http.Error(res, fmt.Sprintf("failed to decode envelope: %s", err), http.StatusBadRequest)
		return
	}

	err = processNotifications(envelope)
	if err != nil {
		http.Error(res, fmt.Sprintf("failed to process notifications: %s", err), http.StatusInternalServerError)
		return
	}
}

// processNotifications is idem potent.
func processNotifications(envelope notifications.Envelope) error {
	for _, v := range envelope.Events {
		if v.Action == "push" {

			// create layer if not exist with 0 to count
			digest := v.Target.Digest.String()
			layer, err := images.FindLayerByDigest(context.Background(), digest)
			if err != nil {
				return err
			} else if layer == nil {
				_, err := images.CreateLayer(context.Background(), digest, v.Target.Size)
				if err != nil {
					return err
				}
			}

			if v.Target.Tag != "" {

				account, err := accounts.FindByID(context.Background(), getNamespaceID(v.Target.Repository))
				if err != nil {
					return err
				}

				image, err := images.FindImageByInfos(context.Background(), v.Target.Repository, v.Target.Tag, digest)
				if err != nil {
					logrus.Trace(err)
					return err
				}

				// insert image if not exist
				if image == nil {
					image, err = images.Create(
						context.Background(),
						getName(v.Target.Repository),
						digest,
						account.ID, // todo change this with the real
						v.Target.Tag,
					)
					if err != nil {
						logrus.Trace("can't create image", image, err)
						return err
					}
				}

				// get layers by calling the registry manifest
				layers := getLayers(account.Email, v.Target.Repository, digest, v.Target.Size)
				if layers == nil {
					logrus.Trace("can't get layers", image)
					return fmt.Errorf("can't get layers with digest %s and repo %s", digest, v.Target.Repository)
				}

				// insert layers and many to many relation with image id <-> layer digest
				err = insertLayers(layers, image.ID)
				if err != nil {
					logrus.Trace("can't insert layers", image, err)
					return err
				}
			}
		}
	}
	return nil
}

// getLayers call the registry to get all fs layers for a given digest and repo.
func getLayers(email, repo, digest string, size int64) []images.Layer {
	manifest := registrycli.GetManifest("http://localhost:5000", email, repo, digest)

	if manifest == nil {
		return nil
	}

	var layers []images.Layer

	// add layer digest
	for _, layer := range manifest.Layers {
		layers = append(layers, images.Layer{
			Digest:    layer.Digest.String(),
			Size:      layer.Size,
			Count:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	layers = append(layers, images.Layer{
		Digest:    digest,
		Size:      size,
		Count:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	layers = append(layers, images.Layer{
		Digest:    manifest.Config.Digest.String(),
		Size:      manifest.Config.Size,
		Count:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return layers
}

// insertLayers will insert layers to table layers and image_layer.
func insertLayers(layers []images.Layer, imageId string) error {
	err := images.InsertLayers(context.Background(), layers, imageId)
	if err != nil {
		return err
	}
	err = images.InsertImageLayer(context.Background(), layers, imageId)
	if err != nil {
		return err
	}
	return nil
}

func getName(repo string) string {
	return strings.Split(repo, "/")[1]
}

func getNamespaceID(repo string) string {
	return strings.Split(repo, "/")[0]
}
