package registrycli

import (
	"github.com/expectedsh/expected/pkg/registryhook/auth"
	"github.com/expectedsh/expected/pkg/registryhook/auth/token"
)

// todo change this in the future
const registryUrl = "http://localhost:5000"

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
	return token.Generate(auth.RequestFromDaemon{
		Login:   "admin",
		Service: "registry",
	}, []auth.AuthorizedScope{
		{
			Scope: auth.Scope{
				Type: "repository",
				Name: repository,
			},
			AuthorizedActions: []string{"pull", "push", "delete"},
		},
	})
}
