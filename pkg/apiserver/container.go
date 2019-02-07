package apiserver

import (
	"encoding/json"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"io/ioutil"
	"net/http"
)

type createContainer struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Size        int               `json:"size"`
	Tags        []string          `json:"tags"`
	Environment map[string]string `json:"environment"`
}

func (s *ApiServer) GetContainers(w http.ResponseWriter, r *http.Request) {
	// account := session.GetAccount(r)
}

func (s *ApiServer) CreateContainer(w http.ResponseWriter, r *http.Request) {
	form := &createContainer{}
	//account := session.GetAccount(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ErrorInternal(w)
		return
	}
	if err = json.Unmarshal(b, form); err != nil {
		response.ErrorBadRequest(w, "Invalid json payload.")
		return
	}
	response.SingleResource(w, "form", form)
}
