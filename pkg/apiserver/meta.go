package apiserver

import (
	"github.com/expectedsh/expected/pkg/apiserver/request"
	"github.com/expectedsh/expected/pkg/apiserver/response"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *App) GetTags(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	tags, err := containers.FindTagsByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to get tags")
		response.ErrorInternal(w)
		return
	}
	if tags == nil {
		tags = []string{}
	}
	response.Resource(w, "tags", tags)
}

func (a *App) GetImagesName(w http.ResponseWriter, r *http.Request) {
	account := request.GetAccount(r)
	names, err := images.FindImagesName(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to get names")
		response.ErrorInternal(w)
		return
	}
	if names == nil {
		names = []string{}
	}
	response.Resource(w, "names", names)
}
