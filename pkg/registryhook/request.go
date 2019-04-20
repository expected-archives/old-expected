package registryhook

import (
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"time"
)

var queue *amqp.Queue

func initQueue(ch *amqp.Channel) error {
	if queue == nil {
		q, err := ch.QueueDeclare("images", true, false, false, false, nil)
		if err != nil {
			return err
		}
		queue = &q
	}
	return nil
}

func RequestDeleteImage(id string) error {
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
			"Message-Type": "ImageDeleteRequest",
		},
		Body: &protocol.ImageDeleteRequest{
			Id: id,
		},
	})
}

func RequestTokenRegistry(imageName string, tokenDuration time.Duration) (token *protocol.ImageTokenResponse, err error) {
	ch, err := services.RabbitMQ().Client().Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	corrId := uuid.New().String()

	q, err := ch.QueueDeclare(
		"", false, false,
		true, false, nil,
	)
	if err != nil {
		return nil, err
	}

	err = services.RabbitMQ().Publish(ch, "", "images", rabbitmq.Message{
		CorrelationId: corrId,
		ReplyTo:       q.Name,
		Headers: amqp.Table{
			"Message-Type": "ImageTokenRequest",
		},
		Body: &protocol.ImageTokenRequest{
			ImageName:     imageName,
			TokenDuration: int64(tokenDuration),
		},
	})
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, "", true,
		false, false, false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	message := <-msgs

	resp := protocol.ImageTokenResponse{}
	err = proto.Unmarshal(message.Body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
