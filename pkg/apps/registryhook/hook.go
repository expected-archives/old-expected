package registryhook

import (
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/notifications"
	"github.com/expectedsh/expected/pkg/apps/registryhook/hook"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Hook(res http.ResponseWriter, req *http.Request) {

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

	err = hook.Handle(envelope)
	if err != nil {
		logrus.WithField("http-hook-response", err).Error()
		http.Error(res, fmt.Sprintf("failed to process notifications: %s", err), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(200)
}
