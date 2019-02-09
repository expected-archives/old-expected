package auth

import (
	"context"
	"errors"
	"github.com/expectedsh/expected/pkg/accounts"
	"log"
)

func Authenticate(login, token string) (*accounts.Account, error) {
	account, e := accounts.FindByAPIKey(context.Background(), token)
	if e != nil {
		return account, e
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.Email == login {
		return account, nil
	}
	log.Printf("failed authentication with login %q and token %q\n", login, token)
	return nil, errors.New("bad credentials")
}