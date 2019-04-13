package handler

import (
	"context"
	"github.com/docker/docker/api/types/swarm"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/scheduler/docker"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type DeploymentHandler struct{}

func (DeploymentHandler) Name() string {
	return "ContainerDeploymentRequest"
}

func (DeploymentHandler) Handle(m amqp.Delivery) error {
	message := &protocol.ContainerDeploymentRequest{}
	if err := proto.Unmarshal(m.Body, message); err != nil {
		return err
	}
	container, err := containers.FindByID(context.Background(), message.Id)
	if err != nil || container == nil {
		return err
	}
	service, err := docker.FindServiceByName(container.ID)
	if err != nil {
		return err
	}
	if service == nil {
		replicas := uint64(1)
		resources := &swarm.Resources{
			MemoryBytes: int64(container.Memory * 1024 * 1024),
			NanoCPUs:    int64(100000000 * 2),
		}
		response, err := docker.CreateService(swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: container.ID,
				Labels: map[string]string{
					//"traefik.enable": "true",
					//"traefik.port":   "80",
					//"traefik.docker.network": "private",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image: container.Image,
					Env:   convertEnv(container.Environment),
				},
				Resources: &swarm.ResourceRequirements{
					Limits:       resources,
					Reservations: resources,
				},
			},
			Mode: swarm.ServiceMode{
				Replicated: &swarm.ReplicatedService{
					Replicas: &replicas,
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target: "eweww",
				},
			},
		})
		if err != nil {
			return err
		}
		logrus.Infoln(response.ID)
	} else {

	}
	return nil
}
