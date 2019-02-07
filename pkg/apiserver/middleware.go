package apiserver

import (
	"net/http"

	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
	"github.com/expectedsh/expected/pkg/models"
)

func (s *ApiServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			w.WriteHeader(http.StatusOK)
		}
	})
}

func (s *ApiServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.ErrorForbidden(w)
			return
		}
		account, err := models.Accounts.GetByApiKey(r.Context(), header)
		if err != nil {
			response.ErrorInternal(w)
			return
		}
		if account == nil {
			response.ErrorForbidden(w)
			return
		}
		session.SetAccount(r, account)
		next.ServeHTTP(w, r)
	})
}
