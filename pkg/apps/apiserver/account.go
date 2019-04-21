package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/expectedsh/expected/pkg/util/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func (a *app.App) GetAccount(w http.ResponseWriter, r *http.Request) {
	apps.Resource(w, "account", apps.GetAccount(r))
}

func (a *app.App) SyncAccount(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	token := &oauth2.Token{
		AccessToken: account.GithubAccessToken,
		TokenType:   "bearer",
	}
	user, err := github.GetUser(r.Context(), token)
	if err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to get github user")
		apps.ErrorInternal(w)
		return
	}
	account.Name = user.Name
	account.AvatarURL = user.AvatarUrl
	email, err := github.GetPrimaryEmail(r.Context(), token)
	if err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to get github user email")
		apps.ErrorInternal(w)
		return
	}
	account.Email = email.Email
	if err = accounts.UpdateAccount(r.Context(), account); err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to update an account")
		apps.ErrorInternal(w)
		return
	}
	apps.Resource(w, "account", account)
}

func (a *app.App) RegenerateAPIKeyAccount(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	account.RegenerateAPIKey()
	if err := accounts.UpdateAccount(r.Context(), account); err != nil {
		logrus.WithField("account", account.ID).WithError(err).Errorln("unable to regenerate api key")
		apps.ErrorInternal(w)
		return
	}
	apps.Resource(w, "account", account)
}
