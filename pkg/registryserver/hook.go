package registryserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/notifications"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/images"
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

	processNotifications(envelope)
}

func processNotifications(envelope notifications.Envelope) {
	for _, v := range envelope.Events {
		if v.Action == "push" {
			account, err := accounts.FindByID(context.Background(), getUserId(v.Target.Repository))
			if err != nil {
				// todo add a way to know there is a error here
				continue
			}
			image, err := images.Create(
				context.Background(),
				getName(v.Target.Repository),
				account.ID,
				v.Target.Digest.String(),
				v.Target.Repository,
				v.Target.Tag,
				v.Target.Size,
			)
			if err != nil {
				// todo add a way to know there is a error here
				continue
			}
			logrus.Infof("%v+", image)
		}
	}
}

func getName(repo string) string {
	return strings.Split(repo, "/")[1]
}

func getUserId(repo string) string {
	return strings.Split(repo, "/")[0]
}
