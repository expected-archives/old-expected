package rabbitmq

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Service struct {
	config *Config
	conn   *amqp.Connection
}

type Config struct {
	Addr string `envconfig:"addr" default:"amqp://expected:expected@localhost/expected"`
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		conn:   nil,
	}
}

func NewFromEnv() *Service {
	config := &Config{}
	if err := envconfig.Process("RABBITMQ", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}
	return New(config)
}

func (srv *Service) Name() string {
	return "rabbitmq"
}

func (srv *Service) Start() error {
	conn, err := amqp.Dial(srv.config.Addr)
	if err != nil {
		return err
	}
	srv.conn = conn
	return nil
}

func (srv *Service) Stop() error {
	return srv.conn.Close()
}

func (srv *Service) Started() bool {
	return srv.conn != nil && !srv.conn.IsClosed()
}

func (srv *Service) Client() *amqp.Connection {
	return srv.conn
}
