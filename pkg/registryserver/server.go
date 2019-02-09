package registryserver

import (
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

	router.HandleFunc("/registry/auth", Auth).Methods("GET")
	router.HandleFunc("/registry/hook", Hook).Methods("POST")

	return http.ListenAndServe(s.Addr, router)
}
