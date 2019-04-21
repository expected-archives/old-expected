package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/expectedsh/expected/pkg/util/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func (a *App) GetAccount(w http.ResponseWriter, r *http.Request) {
	response.Resource(w, "account", request.GetAccount(r))
}

func (a *App) SyncAccount(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	token := &oauth2.Token{
		AccessToken: account.GithubAccessToken,
		TokenType:   "bearer",
	}
	user, err := github.GetUser(r.Context(), token)
	if err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to get github user")
		response.ErrorInternal(w)
		return
	}
	account.Name = user.Name
	account.AvatarURL = user.AvatarUrl
	email, err := github.GetPrimaryEmail(r.Context(), token)
	if err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to get github user email")
		response.ErrorInternal(w)
		return
	}
	account.Email = email.Email
	if err = accounts.UpdateAccount(r.Context(), account); err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to update an account")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "account", account)
}

func (a *App) RegenerateAPIKeyAccount(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	account.RegenerateAPIKey()
	if err := accounts.UpdateAccount(r.Context(), account); err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to regenerate api key")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "account", account)
}
