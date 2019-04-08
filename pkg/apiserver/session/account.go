package session

import (
	"net/http"

	"github.com/expectedsh/expected/pkg/models/accounts"
	"github.com/gorilla/context"
)

func GetAccount(r *http.Request) *accounts.Account {
	account, ok := context.GetOk(r, "account")
	if !ok {
		return nil
	}
	return account.(*accounts.Account)
}

func SetAccount(r *http.Request, account *accounts.Account) {
	context.Set(r, "account", account)
}
