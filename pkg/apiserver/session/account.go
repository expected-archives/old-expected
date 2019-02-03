package session

import (
	"github.com/expectedsh/expected/pkg/models"
	"github.com/gorilla/context"
	"net/http"
)

func GetAccount(r *http.Request) *models.Account {
	account, ok := context.GetOk(r, "account")
	if !ok {
		return nil
	}
	return account.(*models.Account)
}

func SetAccount(r *http.Request, account *models.Account) {
	context.Set(r, "account", account)
}
