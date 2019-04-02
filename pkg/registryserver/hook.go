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

			//todo add error actions

			//todo change this
			account, err := accounts.FindByID(context.Background(), getUserId(v.Target.Repository))
			if err != nil {
				// todo add a way to know there is a error here
				continue
			}

			//
			layer, _ := images.FindLayerByDigest(context.Background(), v.Target.Digest.String())
			if layer == nil {
				_, _ = images.CreateLayer(context.Background(), v.Target.Digest.String(), v.Target.Size)
			}

			if v.Target.Tag != "" {

				// todo check if image exist

				img, _ := images.Create(
					context.Background(),
					getName(v.Target.Repository),
					account.ID,
					v.Target.Digest.String(),
					account.ID, // namespace id , now its the user id
					v.Target.Tag,
				)

				registerManifest(account.Email, img.ID, v.Target.Repository, v.Target.Digest.String())
				// todo check total size of the image
			}
		}
	}
}

func registerManifest(email, imageId, repo, digest string) {
	manifest := registrycli.GetManifest("http://localhost:5000", email, repo, digest)

	if manifest == nil {
		// todo err
		return
	}

	// add layer digest
	for _, layer := range manifest.Layers {
		dig := layer.Digest.String()
		_, _ = images.CreateImageLayer(context.Background(), imageId, dig)
		_ = images.LayerIncrement(context.Background(), dig)
	}

	// add final digest
	_, _ = images.CreateImageLayer(context.Background(), imageId, digest)
	_ = images.LayerIncrement(context.Background(), digest)

	// add config digest
	_, _ = images.CreateImageLayer(context.Background(), imageId, manifest.Config.Digest.String())
	_ = images.LayerIncrement(context.Background(), manifest.Config.Digest.String())
}

func getName(repo string) string {
	return strings.Split(repo, "/")[1]
}

func getUserId(repo string) string {
	return strings.Split(repo, "/")[0]
}
