package apiserver

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/expectedsh/expected/pkg/util/sse"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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

var server = sse.New()

func (s *ApiServer) LogContainer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("___ HELLO WORLD ___")
	//account := request.GetAccount(r)
	containerId := "lol"
	//container, err := containers.FindByID(r.Context(), containerId)
	//if err != nil {
	//	response.ErrorInternal(w)
	//	return
	//}
	//if container == nil {
	//	response.ErrorNotFound(w)
	//	return
	//}
	//if container.OwnerID != account.ID {
	//	response.ErrorForbidden(w)
	//	return
	//}

	if !server.StreamExists(containerId) {
		go sendLogs(server)
	}
	server.CreateStream(containerId)
	if server.StreamExists(containerId) {
		fmt.Println("____ STREAM EXIST _____")
	}
	server.HTTPHandler(containerId, w, r)
	fmt.Println("EXIT STREAM ______")
	server.RemoveStream(containerId)
}

func sendLogs(server *sse.Server) {
	i := 0
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("_____ SEND LOGS _____")
		server.Publish("lol", &sse.Event{
			ID:   []byte(strconv.Itoa(i)),
			Data: []byte(strconv.Itoa(i) + " - " + strconv.Itoa(rand.Int())),
		})
		i++
	}
}
