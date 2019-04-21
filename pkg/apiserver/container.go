package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *App) ListContainers(w http.ResponseWriter, r *http.Request) {
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

func (a *App) CreateContainer(w http.ResponseWriter, r *http.Request) {
	form := &request.CreateContainer{}
	account := request.GetAccount(r)
	if err := request.ParseBody(r, form); err != nil {
		response.ErrorBadRequest(w, "Invalid json payload.", nil)
		return
	}
	if errors := form.Validate(r.Context(), account.ID); len(errors) > 0 {
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
	response.Resource(w, "container", container)
}

func (a *App) GetContainer(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	if err != nil {
		logrus.
			WithField("name", name).
			WithField("action", "get").
			WithError(err).
			Errorln("unable to get container")
		response.ErrorInternal(w)
		return
	}
	if ctr == nil {
		response.ErrorNotFound(w)
		return
	}
	response.Resource(w, "container", ctr)
}

func (a *App) StartContainer(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	log := logrus.WithField("name", name).WithField("action", "start")
	if err != nil {
		log.WithError(err).Error("unable to get container")
		response.ErrorInternal(w)
		return
	}
	if ctr == nil {
		response.ErrorNotFound(w)
		return
	}
	if _, err := scheduler.RequestChangeContainerState(r.Context(), ctr.ID, protocol.State_START); err != nil {
		log.WithError(err).Error("unable to request container state change")
		response.ErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) StopContainer(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	log := logrus.WithField("name", name).WithField("action", "stop")
	if err != nil {
		log.WithError(err).Error("unable to get container")
		response.ErrorInternal(w)
		return
	}
	if ctr == nil {
		response.ErrorNotFound(w)
		return
	}
	if _, err := scheduler.RequestChangeContainerState(r.Context(), ctr.ID, protocol.State_STOP); err != nil {
		log.WithError(err).Error("unable to request container state change")
		response.ErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
