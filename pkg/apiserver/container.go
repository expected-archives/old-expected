package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/sirupsen/logrus"
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
	account := request.GetAccount(r)
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

func (s *ApiServer) GetContainerPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := containers.FindPlans(r.Context())
	if err != nil {
		logrus.WithError(err).Errorln("unable to get container plans")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "plans", plans)
}

func (s *ApiServer) CreateContainer(w http.ResponseWriter, r *http.Request) {
	form := &createContainer{}
	account := request.GetAccount(r)
	if err := request.ParseBody(r, form); err != nil {
		response.ErrorBadRequest(w, "Invalid json payload.", nil)
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
	if err = scheduler.RequestDeployment(container.ID); err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to send container deployment request")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "container", container)
}
