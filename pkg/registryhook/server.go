package registryhook

import (
	"github.com/expectedsh/expected/pkg/util/cors"
	"github.com/gorilla/mux"
	"net/http"
)

type RegistryServer struct {
	Addr string
}

func New(addr string) *RegistryServer {
	return &RegistryServer{
		Addr: addr,
	}
}

func (s *RegistryServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/hook", Hook)

	if err := cors.ApplyMiddleware(router); err != nil {
		return err
	}
	return http.ListenAndServe(s.Addr, router)
}
