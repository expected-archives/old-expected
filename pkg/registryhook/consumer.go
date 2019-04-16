package registryhook

import (
	"github.com/expectedsh/expected/pkg/registryhook/handler"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/rabbitmq"
	"github.com/sirupsen/logrus"
)

var handlers = []rabbitmq.MessageHandler{
	&handler.ImageDletete{Logger: logrus.WithField("task", "image-delete")},
}

func findHandler(name string) rabbitmq.MessageHandler {
	for _, h := range handlers {
		if h.Name() == name {
			return h
		}
	}
	return nil
}

func Start() error {
	ch, err := services.RabbitMQ().Client().Channel()
	if err != nil {
		return err
	}
	if err = initQueue(ch); err != nil {
		return err
	}
	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	for message := range messages {
		logrus.
			WithField("routing-key", message.RoutingKey).
			WithField("headers", message.Headers).
			Infoln("handling new message")
		messageType := message.Headers["Message-Type"]
		if messageType == nil {
			logrus.Warningln("invalid message, no message type provided")
			if err = message.Ack(false); err != nil {
				logrus.WithError(err).Errorln("unable to ack the message")
			}
			continue
		}
		h := findHandler(messageType.(string))
		if h == nil {
			logrus.
				WithField("message-type", messageType.(string)).
				Warningln("unhandled message type")
			if err = message.Nack(false, true); err != nil {
				logrus.WithError(err).Errorln("unable to nack the message")
			}
			continue
		}
		if err = h.Handle(message); err != nil {
			logrus.
				WithField("message-type", messageType.(string)).
				WithError(err).
				Errorln("failed to handle message")
			if err = message.Nack(false, true); err != nil {
				logrus.WithError(err).Errorln("unable to nack the message")
			}
		} else {
			if err = message.Ack(false); err != nil {
				logrus.WithError(err).Errorln("unable to ack the message")
			}
		}
	}
	return nil
}