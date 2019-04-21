package apiserver

import (
	"context"
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (a *app.App) ListContainers(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	ctrs, err := containers.FindContainersByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).Errorln("unable to get containers")
		apps.ErrorInternal(w)
		return
	}
	if ctrs == nil {
		ctrs = []*containers.Container{}
	}
	apps.Resource(w, "containers", ctrs)
}

func (a *app.App) CreateContainer(w http.ResponseWriter, r *http.Request) {
	form := &apps.CreateContainer{}
	account := apps.GetAccount(r)
	if err := apps.ParseBody(r, form); err != nil {
		apps.ErrorBadRequest(w, "Invalid json payload.", nil)
		return
	}
	if errors := form.Validate(r.Context(), account.ID); len(errors) > 0 {
		apps.ErrorBadRequest(w, "Invalid form.", errors)
		return
	}
	container, err := containers.CreateContainer(r.Context(), form.Name, form.Image, form.PlanID,
		form.Environment, form.Tags, account.ID)
	if err != nil {
		logrus.WithError(err).Errorln("unable to create container")
		apps.ErrorInternal(w)
		return
	}
	endpoint := strings.ReplaceAll(container.ID, "-", "") + ".ctr.expected.sh"
	if _, err := containers.CreateEndpoint(context.Background(), container, endpoint, true); err != nil {
		logrus.WithError(err).Errorln("unable to create default container endpoint")
		apps.ErrorInternal(w)
		return
	}
	apps.Resource(w, "container", container)
}

func (a *app.App) GetContainer(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	if err != nil {
		logrus.
			WithField("name", name).
			WithField("action", "get").
			WithError(err).
			Errorln("unable to get container")
		apps.ErrorInternal(w)
		return
	}
	if ctr == nil {
		apps.ErrorNotFound(w)
		return
	}
	apps.Resource(w, "container", ctr)
}

func (a *app.App) StartContainer(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	log := logrus.WithField("name", name).WithField("action", "start")
	if err != nil {
		log.WithError(err).Error("unable to get container")
		apps.ErrorInternal(w)
		return
	}
	if ctr == nil {
		apps.ErrorNotFound(w)
		return
	}
	if _, err := apps.RequestChangeContainerState(r.Context(), ctr.ID, protocol.State_START); err != nil {
		log.WithError(err).Error("unable to request container state change")
		apps.ErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *app.App) StopContainer(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	log := logrus.WithField("name", name).WithField("action", "stop")
	if err != nil {
		log.WithError(err).Error("unable to get container")
		apps.ErrorInternal(w)
		return
	}
	if ctr == nil {
		apps.ErrorNotFound(w)
		return
	}
	if _, err := apps.RequestChangeContainerState(r.Context(), ctr.ID, protocol.State_STOP); err != nil {
		log.WithError(err).Error("unable to request container state change")
		apps.ErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
