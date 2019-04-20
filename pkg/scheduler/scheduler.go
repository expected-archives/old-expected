package scheduler

import (
	"github.com/expectedsh/expected/pkg/scheduler/docker"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

var subscriptions []*nats.Subscription

func subscribe(subject string, h nats.Handler) error {
	sub, err := services.NATS().Client().Subscribe(subject, h)
	if err != nil {
		return err
	}
	subscriptions = append(subscriptions, sub)
	return nil
}

func Start() error {
	if err := docker.Init(); err != nil {
		return err
	}
	if err := subscribe("container:change-state", handleContainerChangeState); err != nil {
		return err
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	logrus.Info("unsubscribing to all subjects")
	for _, sub := range subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			logrus.WithField("subject", sub.Subject).WithError(err).Error("failed to unsubscribe")
		}
	}
	time.Sleep(2 * time.Second)
	return nil
}
