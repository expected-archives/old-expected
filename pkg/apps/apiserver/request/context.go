package request

import (
	"github.com/expectedsh/expected/pkg/models/containers"
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

func GetContainer(r *http.Request) *containers.Container {
	container, ok := context.GetOk(r, "container")
	if !ok {
		return nil
	}
	return container.(*containers.Container)
}

func SetContainer(r *http.Request, container *containers.Container) {
	context.Set(r, "container", container)
}
