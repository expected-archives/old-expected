package apiserver

import (
	"github.com/docker/distribution/reference"
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/scheduler"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
)

var (
	// nameRegexp define the name constraint. Must start with
	// alpha numeric char, can contain _ or - chars.
	nameRegexp = regexp.MustCompile(`^[\w0-9]+[\w_\-0-9]$`)

	// tagRegexp define the tag constaint. Must start with alpha
	// char, can numericals character and contain . - _ characters.
	tagRegexp = regexp.MustCompile(`^[\w][\w0-9.-_]$`)

	// environementKeyRegexp define the shell posix standart to accept environement
	// key variable.
	environementKeyRegexp = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]$`)
)

type createContainer struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Memory      int               `json:"memory"`
	Tags        []string          `json:"tags"`
	Environment map[string]string `json:"environment"`
}

func (s createContainer) validate() (errors map[string]string) {
	errors = make(map[string]string)

	if !nameRegexp.MatchString(s.Name) {
		errors["name"] = "Name must start with alphanumerical character and can contain dash or underscore. "
	}
	if len(s.Name) < 3 || len(s.Name) > 32 {
		errors["name"] = "Name must be between 3 and 32 characters."
	}

	if !reference.ReferenceRegexp.MatchString(s.Image) {
		errors["image "] = "Invalid image name."
	}

	for _, tag := range s.Tags {
		if !tagRegexp.MatchString(tag) {
			errors["tags"] = "Tags must start with alpha characters and can contain numbers, dot, dash or underscore."
		}
		if len(tag) == 0 || len(tag) > 127 {
			errors["tags"] = "Tags must be between 1 and 217 characters."
		}
	}

	for key, value := range s.Environment {
		if !environementKeyRegexp.MatchString(key) {
			errors["environment"] = "Environment key must start with alpha characters or underscore and can contain " +
				"numericals values."
		}
		if len(key) == 0 || len(key) > 1024 {
			errors["environment"] = "Environment key must be between 1 and 1024 characters."
		}

		if len(value) > 32768 {
			errors["environment"] = "Environment value must be lesser than 32768 characters."
		}
	}

	return errors
}

func (s *ApiServer) GetContainers(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	ctrs, err := containers.FindContainerByOwnerID(r.Context(), account.ID)
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

	errors := form.validate()
	if len(errors) > 0 {
		response.ErrorBadRequest(w, "Invalid form.", errors)
		return
	}

	container, err := containers.CreateContainer(r.Context(), form.Name, form.Image, form.Memory,
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
