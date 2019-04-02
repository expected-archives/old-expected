package apiserver

import (
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
)

func (s *ApiServer) GetAccount(w http.ResponseWriter, r *http.Request) {
	response.Resource(w, "account", session.GetAccount(r))
}

func (s *ApiServer) SyncAccount(w http.ResponseWriter, r *http.Request) {

	response.Resource(w, "account", session.GetAccount(r))
}

func (s *ApiServer) RegenerateAPIKeyAccount(w http.ResponseWriter, r *http.Request) {
	account := session.GetAccount(r)
	account.RegenerateAPIKey()
	if err := accounts.Update(r.Context(), account); err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to regenerate api key")
	}
	response.Resource(w, "account", account)
}
