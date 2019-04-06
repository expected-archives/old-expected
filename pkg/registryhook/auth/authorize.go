package auth

import (
	"errors"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthorizedScope struct {
	Scope
	AuthorizedActions []string
}

func Authorize(account accounts.Account, scopes []Scope) ([]AuthorizedScope, error) {
	authorizedScopes := make([]AuthorizedScope, 0)

	for _, scope := range scopes {
		logrus.Infof("authorization for %s actions: %v type: %v", scope.Name, scope.Actions, scope.Type)

		if !account.Admin {
			namespace, _, err := resource(scope.Name)
			if err != nil || namespace != account.ID {
				continue
			}
		}

		authorizedScopes = append(authorizedScopes, AuthorizedScope{
			Scope:             scope,
			AuthorizedActions: []string{"pull", "push"},
		})
	}
	return authorizedScopes, nil
}

func resource(scopeName string) (namespace string, image string, err error) {
	scopeSplit := strings.Split(scopeName, "/")
	if len(scopeSplit) != 2 {
		return "", "", errors.New("should be url/repo/image")
	}
	return scopeSplit[0], scopeSplit[1], nil
}
