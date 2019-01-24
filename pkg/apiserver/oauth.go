package apiserver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httputil"
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
	client := oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://api.github.com/user")

	// /user/emails get email
	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
