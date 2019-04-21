package scheduler

import (
	"context"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/scheduler/docker"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/sirupsen/logrus"
)

func handleContainerChangeState(subject, reply string, r *protocol.ChangeContainerStateRequest) {
	container, err := containers.FindContainerByID(context.Background(), r.Id)
	log := logrus.WithField("id", r.Id).WithField("request", r.RequestedState.String())
	log.Info("new container change request received")

	if err != nil || container == nil {
		log.WithError(err).Error("failed to find container")
		return
	}

	service, err := docker.ServiceFindByName(container.ID)
	if err != nil {
		log.WithError(err).Error("failed to find container service")
		return
	}

	if (service == nil && r.RequestedState == protocol.State_STOP) ||
		(service != nil && r.RequestedState == protocol.State_START) {
		log.Info("service current state is already to desired state")
		return
	}

	if r.RequestedState == protocol.State_START {
		log.Info("creating the service")
		if err := docker.ServiceCreate(container); err != nil {
			log.WithError(err).Error("failed to create container service")
			return
		}
	}

	if r.RequestedState == protocol.State_STOP {
		log.Info("removing the service")
		if err := docker.ServiceRemove(container); err != nil {
			log.WithError(err).Error("failed to remove container service")
			return
		}
	}

	if err := services.NATS().Client().PublishRequest(subject, reply, &protocol.ChangeContainerStateReply{}); err != nil {
		log.WithError(err).Error("failed to publish reply")
		return
	}
}
