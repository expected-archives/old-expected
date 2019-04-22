package apps

import (
	"context"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

var (
	httpServer    *http.Server
	grpcServer    *grpc.Server
	subscriptions []*nats.Subscription
)

type GRPCConfigurer interface {
	ConfigureGRPC(server *grpc.Server)
}

func handleStop(ch chan os.Signal) {
	<-ch
	logrus.Info("stopping the app")
	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if err := sub.Unsubscribe(); err != nil {
				logrus.WithField("subject", sub.Subject).WithError(err).Error("failed to unsubscribe")
			}
		}
	}
	if httpServer != nil {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logrus.WithError(err).Error("failed to shutdown http server")
		}
	}
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	logrus.Info("stopping services")
	services.Stop()
}

func HandleHTTP(h http.Handler) error {
	if httpServer != nil {
		return ErrHttpHandlerAlreadyDefined
	}
	httpServer = &http.Server{
		Handler: h,
		Addr:    GetEnvOrDefault("HTTP_ADDR", ":3000"),
	}
	logrus.Infof("http server listening on %v", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func HandleGRPC(configurer GRPCConfigurer) error {
	if grpcServer != nil {
		return ErrGRPCHandlerAlreadyDefined
	}
	grpcServer = grpc.NewServer()
	configurer.ConfigureGRPC(grpcServer)
	addr := GetEnvOrDefault("GRPC_ADDR", ":4000")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logrus.Infof("grpc server listening on %v", addr)
	if err := grpcServer.Serve(listener); err != nil {
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
