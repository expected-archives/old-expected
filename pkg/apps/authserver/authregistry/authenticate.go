package authregistry

import (
	"context"
	"errors"
	"github.com/expectedsh/expected/pkg/models/accounts"
	"strings"
)

func Authenticate(login, token string) (*accounts.Account, error) {
	account, e := accounts.FindAccountByAPIKey(context.Background(), token)
	if e != nil {
		return account, e
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	if strings.ToLower(account.Email) != strings.ToLower(login) {
		return nil, errors.New("bad credentials")
	}
	return account, nil
}
