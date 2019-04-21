package registry

import (
	"github.com/expectedsh/expected/pkg/apps"
	"time"
)

var registryUrl = "http://localhost:5000"

func Init(url string) {
	registryUrl = url
}

type DeleteStatus int

const (
	DeleteStatusDeleted DeleteStatus = iota
	DeleteStatusNotFound
	DeleteStatusUnknown
)

func (d DeleteStatus) String() string {
	switch d {
	case DeleteStatusDeleted:
		return "deleted"
	case DeleteStatusNotFound:
		return "not found"
	}
	return "unknown"
}

func newToken(repository string) (string, error) {
	return apps.Generate(apps.Request{
		Login:   "admin",
		Service: "registry",
	}, []apps.AuthorizedScope{
		{
			Scope: apps.Scope{
				Type: "repository",
				Name: repository,
			},
			AuthorizedActions: []string{"pull", "push", "delete"},
		},
	}, time.Hour)
}
