package registryhook

import (
	"context"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"time"
)

const Subject = "images"

func RequestDeleteImage(parent context.Context, id string) (*protocol.DeleteImageReply, error) {
	req := &protocol.DeleteImageRequest{
		Id: id,
	}
	reply := &protocol.DeleteImageReply{}
	ctx, _ := context.WithTimeout(parent, time.Second)
	if err := services.NATS().Client().RequestWithContext(ctx, Subject, req, reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func RequestTokenRegistry(parent context.Context, image string, duration time.Duration) (token *protocol.GenerateTokenReply, err error) {
	req := &protocol.GenerateTokenRequest{
		Image:    image,
		Duration: duration.Nanoseconds(),
	}
	reply := &protocol.GenerateTokenReply{}
	ctx, _ := context.WithTimeout(parent, time.Second)
	if err := services.NATS().Client().RequestWithContext(ctx, Subject, req, reply); err != nil {
		return nil, err
	}
	return reply, nil
}
