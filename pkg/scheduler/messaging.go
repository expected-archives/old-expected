package scheduler

import (
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	"github.com/streadway/amqp"
)

var queue *amqp.Queue

func initQueue(ch *amqp.Channel) error {
	if queue == nil {
		q, err := ch.QueueDeclare("containers", true, false, false, false, nil)
		if err != nil {
			return err
		}
		queue = &q
	}
	return nil
}

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
		Body: &protocol.ContainerDeploymentRequest{
			Id: id,
		},
	})
}
