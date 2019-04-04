package registryserver

import (
	"encoding/json"
	"fmt"
	"github.com/expectedsh/expected/pkg/registryserver/auth"
	"github.com/expectedsh/expected/pkg/registryserver/auth/token"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Auth(response http.ResponseWriter, request *http.Request) {
	req, err := auth.ParseRequest(request)
	if err != nil {
		logrus.WithField("Parsing request fail", err).Error()
		http.Error(response, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}

	account, err := auth.Authenticate(req.Login, req.Token)
	if account == nil || err != nil {
		logrus.WithField("Authenticate fail", err).Error()
		http.Error(response, fmt.Sprintf("Authentication failed: %s", err), http.StatusUnauthorized)
		response.Header()["WWW-Authenticate"] = []string{fmt.Sprintf(`Basic realm="%s"`, token.Issuer)}
		return
	}

	authorizedScopes, err := auth.Authorize(*account, req.Scopes)
	if err != nil {
		logrus.WithField("Authorize fail", err).Error()
		http.Error(response, fmt.Sprintf("Authorization failed: %s", err), http.StatusInternalServerError)
		return
	}

	tok, err := token.Generate(*req, authorizedScopes)
	if err != nil {
		logrus.WithField("Generating token fail", err).Error()
		http.Error(response, fmt.Sprintf("Token generation failed: %s", err), http.StatusInternalServerError)
		return
	}

	result, _ := json.Marshal(&map[string]string{"token": tok})
	response.Header().Set("Content-Type", "application/json")
	_, _ = response.Write(result)
}
