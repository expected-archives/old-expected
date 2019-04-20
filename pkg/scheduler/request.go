package scheduler

import (
	"context"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"time"
)

func RequestChangeContainerState(parent context.Context, id string, requestedState protocol.State) (*protocol.ChangeContainerStateReply, error) {
	req := &protocol.ChangeContainerStateRequest{
		Id:             id,
		RequestedState: requestedState,
	}
	reply := &protocol.ChangeContainerStateReply{}
	ctx, _ := context.WithTimeout(parent, time.Second)
	if err := services.NATS().Client().RequestWithContext(ctx, "container:change-state", req, reply); err != nil {
		return nil, err
	}
	return reply, nil
}
