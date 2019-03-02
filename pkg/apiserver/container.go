package apiserver

import (
	"encoding/json"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/apiserver/session"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type createContainer struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Memory      int               `json:"memory"`
	Tags        []string          `json:"tags"`
	Environment map[string]string `json:"environment"`
}

func (s *ApiServer) GetContainers(w http.ResponseWriter, r *http.Request) {
	account := session.GetAccount(r)
	ctrs, err := containers.FindByOwnerID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to get containers")
		response.ErrorInternal(w)
		return
	}
	if ctrs == nil {
		ctrs = []*containers.Container{}
	}
	response.Resource(w, "containers", ctrs)
}

func (s *ApiServer) CreateContainer(w http.ResponseWriter, r *http.Request) {
	form := &createContainer{}
	account := session.GetAccount(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ErrorInternal(w)
		return
	}
	if err = json.Unmarshal(b, form); err != nil {
		response.ErrorBadRequest(w, "Invalid json payload.")
		return
	}
	// todo check form
	container, err := containers.Create(r.Context(), form.Name, form.Image, form.Memory,
		form.Environment, form.Tags, account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to create container")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "container", container)
}
