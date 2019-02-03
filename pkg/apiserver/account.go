package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
	"net/http"
)

func (s *ApiServer) Account(w http.ResponseWriter, r *http.Request) {
	response.SingleResource(w, "account", session.GetAccount(r))
}
