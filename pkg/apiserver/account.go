package apiserver

import (
	"net/http"

	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
)

func (s *ApiServer) GetAccount(w http.ResponseWriter, r *http.Request) {
	response.SingleResource(w, "account", session.GetAccount(r))
}
