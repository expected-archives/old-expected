package image

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Service struct {
	config *Config
	conn   *grpc.ClientConn
	client protocol.ImageClient
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
	if err := envconfig.Process("IMAGE", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}

	return New(config)
}

func (srv *Service) Name() string {
	return "image"
}

func (srv *Service) Start() error {
	var conn *grpc.ClientConn
	var err error
	if apps.CurrentMode() == apps.ModeProduction {
		conn, err = grpc.Dial(srv.config.Addr)
	} else {
		conn, err = grpc.Dial(srv.config.Addr, grpc.WithInsecure())
	}
	if err != nil {
		return err
	}
	srv.conn = conn
	srv.client = protocol.NewImageClient(conn)
	return nil
}

func (srv *Service) Stop() error {
	return srv.conn.Close()
}

func (srv *Service) Started() bool {
	return srv.conn != nil && srv.conn.GetState() == connectivity.Ready
}

func (srv *Service) Client() protocol.ImageClient {
	return srv.client
}
