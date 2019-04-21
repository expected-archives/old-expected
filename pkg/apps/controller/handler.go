package controller

import (
	"context"
	"errors"
	"github.com/expectedsh/expected/pkg/apps/controller/docker"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/sirupsen/logrus"
)

func (App) ChangeContainerState(ctx context.Context, r *protocol.ChangeContainerStateRequest) (*protocol.ChangeContainerStateReply, error) {
	container, err := containers.FindContainerByID(ctx, r.Id)
	log := logrus.WithField("id", r.Id).WithField("request", r.RequestedState.String())
	log.Info("new container change request received")

	if err == nil && container == nil {
		err = errors.New("container not found")
	}
	if err != nil {
		log.WithError(err).Error("failed to find container")
		return nil, err
	}

	service, err := docker.ServiceFindByName(container.ID)
	if err != nil {
		log.WithError(err).Error("failed to find container service")
		return nil, err
	}

	if (service == nil && r.RequestedState == protocol.State_STOP) ||
		(service != nil && r.RequestedState == protocol.State_START) {
		log.Info("service current state is already to desired state")
		return &protocol.ChangeContainerStateReply{}, nil
	}

	if r.RequestedState == protocol.State_START {
		log.Info("creating the service")
		if err := docker.ServiceCreate(container); err != nil {
			log.WithError(err).Error("failed to create container service")
			return nil, err
		}
	}

	if r.RequestedState == protocol.State_STOP {
		log.Info("removing the service")
		if err := docker.ServiceRemove(container); err != nil {
			log.WithError(err).Error("failed to remove container service")
			return nil, err
		}
	}

	return &protocol.ChangeContainerStateReply{}, nil
}
