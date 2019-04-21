package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *app.App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			apps.ErrorForbidden(w)
			return
		}
		account, err := accounts.FindAccountByAPIKey(r.Context(), header)
		if err != nil {
			logrus.WithError(err).Errorln("unable to find account")
			apps.ErrorInternal(w)
			return
		}
		if account == nil {
			apps.ErrorForbidden(w)
			return
		}
		apps.SetAccount(r, account)
		next.ServeHTTP(w, r)
	})
}
