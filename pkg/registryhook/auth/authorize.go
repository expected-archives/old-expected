package auth

import (
	"errors"
	"fmt"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/google/uuid"
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

		namespace, _, err := resource(scope.Name)

		// check the semantic of the repository
		if err != nil {
			continue
		}

		// check if namespace is an UUID
		_, err = uuid.Parse(namespace)
		if err != nil {
			fmt.Println("error not an uuid", namespace)
			continue
		}

		if !account.Admin {
			if namespace != account.ID {
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
