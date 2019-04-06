package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func corsMiddleware(router *mux.Router) error {
	routes := make(map[string][]string)
	if err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			return nil
		}
		path, _ := route.GetPathTemplate()
		if len(path) == 0 {
			return nil
		}
		for _, method := range methods {
			routes[path] = append(routes[path], method)
		}
		return nil
	}); err != nil {
		return err
	}
	for route := range routes {
		router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}).Methods("OPTIONS")
	}
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(routes[r.URL.Path], ","))
			w.Header().Add("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			next.ServeHTTP(w, r)
		})
	})
	return nil
}

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
		session.SetAccount(r, account)
		next.ServeHTTP(w, r)
	})
}
