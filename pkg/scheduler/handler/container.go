package handler

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
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
	service, err := services.Docker().FindService(context.Background(), container.ID)
	if err != nil {
		return err
	}
	if service == nil {
		replicas := uint64(1)
		resources := &swarm.Resources{
			MemoryBytes: int64(container.Memory * 1024 * 1024),
			NanoCPUs:    int64(100000000 * 2),
		}
		response, err := services.Docker().Client().ServiceCreate(context.Background(), swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: container.ID,
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
			EndpointSpec: &swarm.EndpointSpec{
				Ports: []swarm.PortConfig{
					{
						Name:        "http",
						Protocol:    "tcp",
						TargetPort:  80,
						PublishMode: "ingress",
					},
				},
			},
		}, types.ServiceCreateOptions{})
		if err != nil {
			return err
		}
		logrus.Infoln(response.ID)
	} else {

	}
	return nil
}
