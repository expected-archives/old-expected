package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.ErrorForbidden(w)
			return
		}
		account, err := accounts.FindByAPIKey(r.Context(), header)
		if err != nil {
			logrus.WithField("header", header).WithError(err).Errorln("unable to find account")
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
