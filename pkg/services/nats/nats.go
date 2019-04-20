package nats

import (
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats/encoders/protobuf"
	"github.com/sirupsen/logrus"
)

type Service struct {
	config *Config
	conn   *nats.EncodedConn
}

type Config struct {
	Addr string `envconfig:"addr" default:"nats://localhost:4222"`
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		conn:   nil,
	}
}

func NewFromEnv() *Service {
	config := &Config{}
	if err := envconfig.Process("NATS", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}
	return New(config)
}

func (srv *Service) Name() string {
	return "nats"
}

func (srv *Service) Start() error {
	c, err := nats.Connect(srv.config.Addr)
	if err != nil {
		return err
	}
	conn, err := nats.NewEncodedConn(c, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return err
	}
	srv.conn = conn
	return conn.LastError()
}

func (srv *Service) Stop() error {
	srv.conn.Close()
	return srv.conn.LastError()
}

func (srv *Service) Started() bool {
	return srv.conn != nil && srv.conn.Conn.IsConnected()
}

func (srv *Service) Client() *nats.EncodedConn {
	return srv.conn
}
