package authserver

import (
	"context"
	"github.com/expectedsh/expected/pkg/apps/authserver/authregistry"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/sirupsen/logrus"
	"time"
)

func (App) GenerateToken(ctx context.Context, r *protocol.GenerateTokenRequest) (*protocol.GenerateTokenReply, error) {
	s, err := authregistry.Generate(authregistry.Request{
		Login:   "admin",
		Service: "registry",
	}, []authregistry.AuthorizedScope{
		{
			Scope: authregistry.Scope{
				Type: "repository",
				Name: r.Image,
			},
			AuthorizedActions: []string{"pull"},
		},
	}, time.Duration(r.Duration))

	if err != nil {
		logrus.WithError(err).Error("failed to generate token")
		return nil, err
	}

	return &protocol.GenerateTokenReply{
		Token: s,
	}, nil
}
