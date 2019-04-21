package authregistry

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/apps"
	"net/url"
	"sort"
	"strings"
)

type Scope struct {
	Type    string
	Name    string
	Actions []string
}

func parseScope(value string, req *apps.Request, form url.Values) error {
	if value != "" {
		for _, scopeStr := range form["scope"] {
			parts := strings.Split(scopeStr, ":")
			var scope Scope
			switch len(parts) {
			case 3:
				scope = Scope{
					Type:    parts[0],
					Name:    parts[1],
					Actions: strings.Split(parts[2], ","),
				}
			case 4:
				scope = Scope{
					Type:    parts[0],
					Name:    parts[1] + ":" + parts[2],
					Actions: strings.Split(parts[3], ","),
				}
			default:
				return fmt.Errorf("invalid scope: %q", scopeStr)
			}
			sort.Strings(scope.Actions)
			req.Scopes = append(req.Scopes, scope)
		}
	}
	return nil
}
