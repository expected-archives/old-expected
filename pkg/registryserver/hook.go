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
	bytes, err := json.MarshalIndent(envelope, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}

	fmt.Println()
	fmt.Println()

	processNotifications(envelope)
}

func processNotifications(envelope notifications.Envelope) {
	for _, v := range envelope.Events {
		if v.Action == "push" {
			if v.Target.Tag != "" {
				image, err := images.FindImageByDigest(context.Background(), v.Target.Digest.String())

				if err != nil {
					logrus.Trace(err)
					continue
				}

				if image != nil {
					continue
				}

				account, err := accounts.FindByID(context.Background(), getUserId(v.Target.Repository))
				if err != nil {
					logrus.Trace(err)
					continue
				}

				// insert image
				image, err = images.Create(
					context.Background(),
					getName(v.Target.Repository),
					v.Target.Digest.String(),
					account.ID, // todo change this with the real
					v.Target.Tag,
				)
				if err != nil {
					logrus.Trace("can't create image", image, err)
					continue
				}

				// get layers by calling the registry manifest
				layers := GetLayers(account.Email, v.Target.Repository, v.Target.Digest.String(), v.Target.Size)
				if layers == nil {
					logrus.Trace("can't get layers", image, err)
					continue
				}

				// insert layers and many to many relation with image id <-> layer digest
				err = registerManifest(layers, image.ID)
				if err != nil {
					logrus.Trace("can't insert layers", image, err)
				}
			}
		}
	}
}

func GetLayers(email, repo, digest string, size int64) []images.Layer {
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

func registerManifest(layers []images.Layer, imageId string) error {
	err := images.InsertLayers(context.Background(), layers)
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

func getUserId(repo string) string {
	return strings.Split(repo, "/")[0]
}
