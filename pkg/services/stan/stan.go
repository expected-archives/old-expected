package stan

import (
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/go-nats-streaming"
	"github.com/sirupsen/logrus"
)

type Service struct {
	config   *Config
	conn     stan.Conn
	clientId string
}

type Config struct {
	Addr      string `envconfig:"addr" default:"nats://localhost:4222"`
	ClusterId string `envconfig:"cluster_id" default:"test-cluster"`
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		conn:   nil,
	}
}

func NewFromEnv() *Service {
	config := &Config{}
	if err := envconfig.Process("STAN", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}
	return New(config)
}

func (srv *Service) Name() string {
	return "stan"
}

func (srv *Service) Start() error {
	srv.clientId = uuid.New().String()
	c, err := stan.Connect(srv.config.ClusterId, srv.clientId, stan.NatsURL(srv.config.Addr))
	if err != nil {
		return err
	}
	srv.conn = c
	return nil
}

func (srv *Service) Stop() error {
	if srv.conn != nil {
		return srv.conn.Close()
	}
	return nil
}

func (srv *Service) Started() bool {
	return srv.conn != nil && srv.conn.NatsConn().IsConnected()
}

func (srv *Service) Client() stan.Conn {
	return srv.conn
}

func (srv *Service) Publish(subject string, message proto.Message) error {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return srv.conn.Publish(subject, bytes)
}

func (srv *Service) PublishAsync(subject string, message proto.Message, ah stan.AckHandler) (string, error) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return "", err
	}
	return srv.conn.PublishAsync(subject, bytes, ah)
}
