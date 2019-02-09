package auth

import (
	"context"
	"errors"
	"github.com/expectedsh/expected/pkg/accounts"
	"github.com/expectedsh/expected/pkg/containers"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthorizedScope struct {
	Scope
	AuthorizedActions []string
}

func Authorize(account accounts.Account, scopes []Scope) ([]AuthorizedScope, error) {
	authorizedScopes := make([]AuthorizedScope, 0)

	ctrs, err := containers.FindByOwnerID(context.Background(), account.ID)
	if err != nil {
		return nil, err
	}

	for _, scope := range scopes {
		logrus.Info("authorization for %s actions: %v type: %v", scope.Name, scope.Actions, scope.Type)

		namespace, image, err := resource(scope.Name)

		if err != nil || namespace != account.ID || !hasAuthorization(image, ctrs) {
			continue
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
		return "", "", errors.New("should be url.com/repo/image")
	}
	return scopeSplit[0], scopeSplit[1], nil
}

func hasAuthorization(image string, ctrs []*containers.Container) bool {
	for _, container := range ctrs {
		if image == container.Name {
			return true
		}
	}
	return false
}
