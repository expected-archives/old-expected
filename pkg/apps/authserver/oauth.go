package authserver

import (
	"context"
	"crypto/tls"
	"github.com/expectedsh/expected/pkg/apps/apiserver/response"
	"github.com/expectedsh/expected/pkg/util/github"
	"net/http"

	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (s *App) OAuthGithub(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.OAuth.AuthCodeURL("", oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
}

func (s *App) OAuthGithubCallback(w http.ResponseWriter, r *http.Request) {
	// todo remove/improve this
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, client)

	token, err := s.OAuth.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		response.ErrorBadRequest(w, "Invalid oauth code.", nil)
		return
	}

	user, err := github.GetUser(ctx, token)
	if err != nil {
		logrus.WithError(err).Errorln("unable to retrieve your github data")
		response.ErrorInternal(w)
		return
	}

	account, err := accounts.FindAccountByGithubID(ctx, user.ID)
	if err != nil {
		logrus.WithError(err).Errorln("unable to retrieve your account")
		response.ErrorInternal(w)
		return
	}

	if account == nil {
		email, err := github.GetPrimaryEmail(ctx, token)

		if err != nil {
			logrus.WithError(err).Errorln("unable to retrieve your github data")
			response.ErrorInternal(w)
			return
		}

		if account, err = accounts.CreateAccount(ctx, user.Name, email.Email, user.AvatarUrl, user.ID,
			token.AccessToken, s.isAdmin(user.Login)); err != nil {
			logrus.WithError(err).Errorln("unable to create your account")
			response.ErrorInternal(w)
			return
		}
	} else {
		account.GithubAccessToken = token.AccessToken

		if err = accounts.UpdateAccount(ctx, account); err != nil {
			logrus.WithField("account", account.ID).WithError(err).Errorln("unable to update your account")
			response.ErrorInternal(w)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Path:  "/",
		Value: account.APIKey,
	})
	http.Redirect(w, r, s.DashboardURL, http.StatusTemporaryRedirect)
}

func (s App) isAdmin(username string) bool {
	for _, name := range s.Admins {
		if name == username {
			return true
		}
	}
	return false
}
