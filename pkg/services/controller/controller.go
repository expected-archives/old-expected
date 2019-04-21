package controller

import (
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Service struct {
	config *Config
	conn   *grpc.ClientConn
	client protocol.ControllerClient
}

type Config struct {
	Addr string `envconfig:"addr" default:"localhost:3000"`
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		conn:   nil,
		client: nil,
	}
}

func NewFromEnv() *Service {
	config := &Config{}
	if err := envconfig.Process("CONTROLLER", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}

	return New(config)
}

func (srv *Service) Name() string {
	return "controller"
}

func (srv *Service) Start() error {
	conn, err := grpc.Dial(srv.config.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	srv.conn = conn
	srv.client = protocol.NewControllerClient(conn)
	return nil
}

func (srv *Service) Stop() error {
	return srv.conn.Close()
}

func (srv *Service) Started() bool {
	return srv.conn != nil && srv.conn.GetState() == connectivity.Ready
}

func (srv *Service) Client() protocol.ControllerClient {
	return srv.client
}
