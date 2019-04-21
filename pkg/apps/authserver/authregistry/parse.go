package authregistry

import (
	"fmt"
	"net/http"
)

type Labels map[string][]string

type Request struct {
	Login   string
	Token   string
	Service string
	Scopes  []Scope
	Labels  Labels
}

func ParseRequest(r *http.Request) (req *Request, err error) {
	req = &Request{
		Scopes: make([]Scope, 0),
	}

	user, token, hasAuth := r.BasicAuth()
	if !hasAuth {
		return nil, fmt.Errorf("require basic account")
	}

	req.Token = token
	req.Login = user

	account := r.FormValue("account")
	if account != req.Login {
		return nil, fmt.Errorf("user and account are not the same (login %q vs account %q)", req.Login, account)
	}

	req.Service = r.FormValue("service")
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("invalid form value")
	}

	if err = parseScope(r.FormValue("scope"), req, r.Form); err != nil {
		return nil, err
	}
	return req, nil
}
