package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps/apiserver/request"
	"github.com/expectedsh/expected/pkg/apps/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.ErrorForbidden(w)
			return
		}
		account, err := accounts.FindAccountByAPIKey(r.Context(), header)
		if err != nil {
			logrus.WithError(err).Errorln("unable to find account")
			response.ErrorInternal(w)
			return
		}
		if account == nil {
			response.ErrorForbidden(w)
			return
		}
		request.SetAccount(r, account)
		next.ServeHTTP(w, r)
	})
}
