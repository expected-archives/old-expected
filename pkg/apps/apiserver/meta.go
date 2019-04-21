package apiserver

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/models/images"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *app.App) GetTags(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	tags, err := containers.FindTagsByNamespaceID(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to get tags")
		apps.ErrorInternal(w)
		return
	}
	if tags == nil {
		tags = []string{}
	}
	apps.Resource(w, "tags", tags)
}

func (a *app.App) GetImagesName(w http.ResponseWriter, r *http.Request) {
	account := apps.GetAccount(r)
	names, err := images.FindImagesName(r.Context(), account.ID)
	if err != nil {
		logrus.WithError(err).WithField("account", account.ID).Errorln("unable to get names")
		apps.ErrorInternal(w)
		return
	}
	if names == nil {
		names = []string{}
	}
	apps.Resource(w, "names", names)
}
