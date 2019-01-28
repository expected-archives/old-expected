package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/github"
	"github.com/expectedsh/expected/pkg/models"
	"golang.org/x/oauth2"
	"net/http"
)

func (s *ApiServer) OAuthGithub(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.OAuth.AuthCodeURL("", oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
}

func (s *ApiServer) OAuthGithubCallback(w http.ResponseWriter, r *http.Request) {
	token, err := s.OAuth.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		response.ErrorBadRequest(w, "Invalid oauth code.")
		return
	}
	user, err := github.GetUser(r.Context(), token)
	if err != nil {
		response.ErrorInternal(w, "Unable to retrieve your github data.")
		return
	}
	account, err := models.Accounts.GetByGithubID(r.Context(), user.ID)
	if err != nil {
		response.ErrorInternal(w, "Unable to retrieve your account.")
		return
	}
	if account == nil {
		email, err := github.GetPrimaryEmail(r.Context(), token)
		if err != nil {
			response.ErrorInternal(w, "Unable to retrieve your github data.")
			return
		}
		account = &models.Account{}
		account.Name = user.Name
		account.Email = email.Email
		account.AvatarUrl = user.AvatarUrl
		account.GithubID = user.ID
		account.GithubAccessToken = token.AccessToken
		account.GithubRefreshToken = token.RefreshToken
		account.Admin = s.Admin == user.Login
		if err = models.Accounts.Create(r.Context(), account); err != nil {
			response.ErrorInternal(w, "Unable to create your account.")
			return
		}
		response.SingleResource(w, "user", user)
	} else {
		response.SingleResource(w, "user", user)
	}
}
