package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/registry/auth/token"
	"github.com/expectedsh/expected/pkg/registryhook/auth"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	Issuer     = "auth_registry"
	Expiration = time.Hour
)

func Generate(auth auth.RequestFromDaemon, scopes []auth.AuthorizedScope) (string, error) {
	now := time.Now().Unix()
	_, alg, err := privateKey.Sign(strings.NewReader("dummy"), 0)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %s", err)
	}

	header := token.Header{
		Type:       "JWT",
		SigningAlg: alg,
		KeyID:      publicKey.KeyID(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %s", err)
	}

	claims := token.ClaimSet{
		Issuer:     Issuer,
		Subject:    auth.Login,
		Audience:   auth.Service,
		NotBefore:  now - 10,
		IssuedAt:   now,
		Expiration: now + int64(Expiration),
		JWTID:      fmt.Sprintf("%d", rand.Int63()),
		Access:     scopeToResourceActions(scopes),
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %s", err)
	}

	payload := fmt.Sprintf("%s%s%s", toBase64(headerJSON), token.TokenSeparator, toBase64(claimsJSON))
	sig, sigAlg, err := privateKey.Sign(strings.NewReader(payload), 0)
	if err != nil || sigAlg != alg {
		return "", fmt.Errorf("failed to sign token: %s", err)
	}
	return fmt.Sprintf("%s%s%s", payload, token.TokenSeparator, toBase64(sig)), nil
}

func scopeToResourceActions(scopes []auth.AuthorizedScope) []*token.ResourceActions {
	var actions []*token.ResourceActions

	for _, a := range scopes {
		ra := &token.ResourceActions{
			Type:    a.Type,
			Name:    a.Name,
			Actions: a.AuthorizedActions,
		}
		if ra.Actions == nil {
			ra.Actions = []string{}
		}
		sort.Strings(ra.Actions)
		actions = append(actions, ra)
	}
	return actions
}

func toBase64(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
