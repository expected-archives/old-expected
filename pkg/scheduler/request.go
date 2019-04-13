package scheduler

import (
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	"github.com/streadway/amqp"
)

func RequestDeployment(id string) error {
	ch, err := services.RabbitMQ().Client().Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err = initQueue(ch); err != nil {
		return err
	}
	return services.RabbitMQ().Publish(ch, "", queue.Name, rabbitmq.Message{
		DeliveryMode: amqp.Persistent,
		Headers: amqp.Table{
			"Message-Type": "ContainerDeploymentRequest",
		},
		Body: &protocol.ContainerDeploymentRequest{
			Id: id,
		},
	})
}
