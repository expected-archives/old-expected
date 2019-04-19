package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *ApiServer) ListContainers(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	ctrs, err := containers.FindContainersByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).Errorln("unable to get containers")
		response.ErrorInternal(w)
		return
	}
	if ctrs == nil {
		ctrs = []*containers.Container{}
	}
	response.Resource(w, "containers", ctrs)
}

func (s *ApiServer) CreateContainer(w http.ResponseWriter, r *http.Request) {
	form := &request.CreateContainer{}
	account := request.GetAccount(r)
	if err := request.ParseBody(r, form); err != nil {
		response.ErrorBadRequest(w, "Invalid json payload.", nil)
		return
	}
	if errors := form.Validate(r.Context()); len(errors) > 0 {
		response.ErrorBadRequest(w, "Invalid form.", errors)
		return
	}
	container, err := containers.CreateContainer(r.Context(), form.Name, form.Image, form.PlanID,
		form.Environment, form.Tags, account.ID)
	if err != nil {
		logrus.WithError(err).Errorln("unable to create container")
		response.ErrorInternal(w)
		return
	}
	if err = scheduler.RequestCreateContainer(container.ID); err != nil {
		logrus.WithError(err).Errorln("unable to request container creation to scheduler")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "container", container)
}
