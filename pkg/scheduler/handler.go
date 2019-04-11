package scheduler

import (
	"context"
	"fmt"
	"github.com/expectedsh/expected/pkg/models/containers"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
)

func ContainerDeploymentHandler(b []byte) error {
	message := &protocol.ContainerDeploymentRequest{}
	if err := proto.Unmarshal(b, message); err != nil {
		return err
	}
	container, err := containers.FindByID(context.Background(), message.Id)
	if err != nil || container == nil {
		logrus.WithField("id", message.Id).Errorln("failed to find container")
	}
	fmt.Println(container.Name)
	return nil
}
