package apiserver

import (
	"github.com/expectedsh/expected/pkg/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *ApiServer) OAuthGithub(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.OAuth.AuthCodeURL(s.OAuthState, oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
}

func (s *ApiServer) OAuthGithubCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != s.OAuthState {
		http.Redirect(w, r, s.OAuth.AuthCodeURL(s.OAuthState, oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
		return
	}
	token, err := s.OAuth.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		logrus.WithError(err).Warningln("unable exchange oauth token")
	}
	ghUser, err := github.GetUser(r.Context(), token)
	if err != nil {
		logrus.WithError(err).Warningln("ghuser")
	}
	email, err := github.GetPrimaryEmail(r.Context(), token)
	if err != nil {
		logrus.WithError(err).Warningln("email")
	}
	logrus.WithField("ghUser", ghUser).WithField("emai", email).Infoln("ok")

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
