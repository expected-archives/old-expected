package authserver

import (
	"encoding/json"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (app.App) AuthRegistry(response http.ResponseWriter, request *http.Request) {
	req, err := apps.ParseRequest(request)
	if err != nil {
		logrus.WithField("Parsing request fail", err).Error()
		http.Error(response, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}

	account, err := apps.Authenticate(req.Login, req.Token)
	if account == nil || err != nil {
		logrus.WithField("Authenticate fail", err).Error()
		http.Error(response, fmt.Sprintf("Authentication failed: %s", err), http.StatusUnauthorized)
		response.Header()["WWW-Authenticate"] = []string{fmt.Sprintf(`Basic realm="%s"`, apps.Issuer)}
		return
	}

	authorizedScopes, err := apps.Authorize(*account, req.Scopes)
	if err != nil {
		logrus.WithField("Authorize fail", err).Error()
		http.Error(response, fmt.Sprintf("Authorization failed: %s", err), http.StatusInternalServerError)
		return
	}

	tok, err := apps.Generate(*req, authorizedScopes, time.Hour)
	if err != nil {
		logrus.WithField("Generating token fail", err).Error()
		http.Error(response, fmt.Sprintf("Token generation failed: %s", err), http.StatusInternalServerError)
		return
	}

	result, _ := json.Marshal(&map[string]string{"token": tok})
	response.Header().Set("Content-Type", "application/json")
	_, _ = response.Write(result)
}
