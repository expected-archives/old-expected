package registry

import (
	"github.com/expectedsh/expected/pkg/authserver/authregistry"
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
	return authregistry.Generate(authregistry.RequestFromDaemon{
		Login:   "admin",
		Service: "registry",
	}, []authregistry.AuthorizedScope{
		{
			Scope: authregistry.Scope{
				Type: "repository",
				Name: repository,
			},
			AuthorizedActions: []string{"pull", "push", "delete"},
		},
	})
}
