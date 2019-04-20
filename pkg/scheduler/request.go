package scheduler

import (
	"context"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
)

const Subject = "containers"

func RequestChangeContainerState(ctx context.Context, id string, requestedState protocol.State) (*protocol.ChangeContainerStateReply, error) {
	req := &protocol.ChangeContainerStateRequest{
		Id:             id,
		RequestedState: requestedState,
	}
	reply := &protocol.ChangeContainerStateReply{}
	if err := services.NATS().Client().RequestWithContext(ctx, Subject, req, reply); err != nil {
		return nil, err
	}
	return reply, nil
}
