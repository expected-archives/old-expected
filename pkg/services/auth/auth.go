package auth

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
	client protocol.AuthClient
}

type Config struct {
	Addr string `envconfig:"addr" default:"localhost:4001"`
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
	if err := envconfig.Process("AUTH", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}

	return New(config)
}

func (srv *Service) Name() string {
	return "auth"
}

func (srv *Service) Start() error {
	conn, err := grpc.Dial(srv.config.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	srv.conn = conn
	srv.client = protocol.NewAuthClient(conn)
	return nil
}

func (srv *Service) Stop() error {
	return srv.conn.Close()
}

func (srv *Service) Started() bool {
	return srv.conn != nil && srv.conn.GetState() == connectivity.Ready
}

func (srv *Service) Client() protocol.AuthClient {
	return srv.client
}
