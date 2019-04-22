package apiserver

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/apps/apiserver/request"
	"github.com/expectedsh/expected/pkg/apps/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (s *App) ListContainers(w http.ResponseWriter, r *http.Request) {
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

func (s *App) CreateContainer(w http.ResponseWriter, r *http.Request) {
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
	endpoint := strings.ReplaceAll(container.ID, "-", "") + ".ctr.expected.sh"
	if _, err := containers.CreateEndpoint(context.Background(), container, endpoint, true); err != nil {
		logrus.WithError(err).Errorln("unable to create default container endpoint")
		response.ErrorInternal(w)
		return
	}
	response.Resource(w, "container", container)
}

func (s *App) GetContainer(w http.ResponseWriter, r *http.Request) {
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

func (s *App) StartContainer(w http.ResponseWriter, r *http.Request) {
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
	if _, err := services.Controller().Client().ChangeContainerState(r.Context(), &protocol.ChangeContainerStateRequest{
		Id:             ctr.ID,
		RequestedState: protocol.State_START,
	}); err != nil {
		log.WithError(err).Error("unable to request container state change")
		response.ErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *App) StopContainer(w http.ResponseWriter, r *http.Request) {
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
	if _, err := services.Controller().Client().ChangeContainerState(r.Context(), &protocol.ChangeContainerStateRequest{
		Id:             ctr.ID,
		RequestedState: protocol.State_STOP,
	}); err != nil {
		log.WithError(err).Error("unable to request container state change")
		response.ErrorInternal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *App) GetContainerLogs(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	name := mux.Vars(r)["name"]
	ctr, err := containers.FindContainerByNameAndNamespaceID(r.Context(), name, account.ID)
	log := logrus.WithField("name", name).WithField("action", "logs")
	if err != nil {
		log.WithError(err).Error("unable to get container")
		response.ErrorInternal(w)
		return
	}
	if ctr == nil {
		response.ErrorNotFound(w)
		return
	}
	logs, err := services.Controller().Client().GetContainerLogs(r.Context(), &protocol.GetContainersLogsRequest{
		Id: ctr.ID,
	})
	if err != nil {
		log.WithError(err).Error("unable to request container logs")
		response.ErrorInternal(w)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)

	for {
		reply, err := logs.Recv()
		if err != nil {
			log.WithError(err).Error("failed to receive container logs reply")
			return
		}
		fmt.Println(reply)
		fmt.Fprintf(w, "event: %s\n", "message")
		fmt.Fprintf(w, "data: %s\n", reply.Line)
		flusher.Flush()
	}
}
