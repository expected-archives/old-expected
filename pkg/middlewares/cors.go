package middlewares

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func ApplyCORS(router *mux.Router) error {
	routes := make(map[string][]string)
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
	})
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
