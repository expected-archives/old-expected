package app

import (
	"context"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var (
	httpServer    *http.Server
	subscriptions []*nats.Subscription
)

func handleStop(ch chan os.Signal) {
	<-ch
	logrus.Info("stopping the app")
	if httpServer != nil {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logrus.WithError(err).Error("failed to shutdown http server")
		}
	}
	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if err := sub.Unsubscribe(); err != nil {
				logrus.WithField("subject", sub.Subject).WithError(err).Error("failed to unsubscribe")
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func HandleHTTP(h http.Handler) error {
	if httpServer != nil {
		return ErrHttpHandlerAlreadyDefined
	}
	httpServer = &http.Server{
		Handler: h,
		Addr:    GetEnvOrDefault("ADDR", ":3000"),
	}
	logrus.Infof("listening on %v", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func HandleSubscription(subject string, h nats.Handler) error {
	sub, err := services.NATS().Client().Subscribe(subject, h)
	if err != nil {
		return err
	}
	subscriptions = append(subscriptions, sub)
	logrus.WithField("subject", subject).Debug("handling new subscription")
	return nil
}
