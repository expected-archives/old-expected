package authserver

import (
	"github.com/expectedsh/expected/pkg/util/github"
	"net/http"

	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (s *ApiServer) OAuthGithub(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.OAuth.AuthCodeURL("", oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
}

func (s *ApiServer) OAuthGithubCallback(w http.ResponseWriter, r *http.Request) {
	token, err := s.OAuth.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		response.ErrorBadRequest(w, "Invalid oauth code.", nil)
		return
	}
	user, err := github.GetUser(r.Context(), token)
	if err != nil {
		logrus.WithError(err).Errorln("unable to retrieve your github data")
		response.ErrorInternal(w)
		return
	}

	account, err := accounts.FindByGithubID(r.Context(), user.ID)
	if err != nil {
		logrus.WithError(err).Errorln("unable to retrieve your account")
		response.ErrorInternal(w)
		return
	}
	if account == nil {
		email, err := github.GetPrimaryEmail(r.Context(), token)
		if err != nil {
			logrus.WithError(err).Errorln("unable to retrieve your github data")
			response.ErrorInternal(w)
			return
		}
		if account, err = accounts.Create(r.Context(), user.Name, email.Email, user.AvatarUrl, user.ID,
			token.AccessToken, s.isAdmin(user.Login)); err != nil {
			logrus.WithError(err).Errorln("unable to create your account")
			response.ErrorInternal(w)
			return
		}
	} else {
		account.GithubAccessToken = token.AccessToken
		if err = accounts.Update(r.Context(), account); err != nil {
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
