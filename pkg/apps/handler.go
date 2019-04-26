package apps

import (
	"context"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/nats-io/go-nats-streaming"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	httpServer *http.Server
	grpcServer *grpc.Server

	cancellers    []context.CancelFunc
	subscriptions []stan.Subscription
)

type GRPCConfigurer interface {
	ConfigureGRPC(server *grpc.Server)
}

func handleStop(ch chan os.Signal) {
	<-ch
	logrus.Info("stopping the app")
	if len(cancellers) > 0 {
		for _, c := range cancellers {
			c()
		}
	}
	if len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if err := sub.Close(); err != nil {
				logrus.WithError(err).Error("failed to unsubscribe")
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

func HandleRunner(runner func(ctx context.Context) error) error {
	ctx, canceller := context.WithCancel(context.Background())
	cancellers = append(cancellers, canceller)
	return runner(ctx)
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

func StanQueueGroupOptions(name string) []stan.SubscriptionOption {
	return []stan.SubscriptionOption{
		stan.AckWait(3 * time.Second),
		stan.DeliverAllAvailable(),
		stan.MaxInflight(5),
		stan.SetManualAckMode(),
		stan.DurableName(name),
	}
}

func HandleSubscription(subject string, h stan.MsgHandler, opts ...stan.SubscriptionOption) error {
	sub, err := services.Stan().Client().Subscribe(subject, h, opts...)
	if err != nil {
		return err
	}
	subscriptions = append(subscriptions, sub)
	logrus.WithField("subject", subject).Debug("handling new subscription")
	return nil
}

func HandleQueueSubscription(subject string, qgroup string, h stan.MsgHandler, opts ...stan.SubscriptionOption) error {
	sub, err := services.Stan().Client().QueueSubscribe(subject, qgroup, h, opts...)
	if err != nil {
		return err
	}
	subscriptions = append(subscriptions, sub)
	logrus.WithField("subject", subject).Debug("handling new queue subscription")
	return nil
}
