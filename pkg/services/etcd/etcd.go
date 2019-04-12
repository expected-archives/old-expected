package etcd

import (
	client "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/connectivity"
)

type Service struct {
	config  *Config
	client  *client.Client
	stopped bool
}

type Config struct {
	AppName   string   `envconfig:"appname"`
	Addresses []string `envconfig:"addresses" default:"localhost:2379"`
}

func (s *Service) Name() string {
	return "etcd"
}

func (s *Service) Start() error {
	etcd, err := client.New(client.Config{Endpoints: s.config.Addresses})
	if err != nil {
		return err
	}
	s.client = etcd
	return nil
}

func (s *Service) Stop() error {
	return s.client.Close()
}

func (s *Service) Started() bool {
	return s.client != nil && s.client.ActiveConnection().GetState() == connectivity.Ready
}

func (s *Service) Client() *client.Client {
	return s.client
}

func (s *Service) Config() *Config {
	return s.config
}
