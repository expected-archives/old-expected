package consul

import (
	"github.com/hashicorp/consul/api"
)

type Service struct {
	config *Config
	client *api.Client

	stopped bool
}

type Config struct {
	AppName string `envconfig:"appname"`
	Address string `envconfig:"address" default:"localhost:8500"`
}

func (s *Service) Name() string {
	return "etcd"
}

func (s *Service) Start() error {
	client, err := api.NewClient(&api.Config{Address: s.config.Address})
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Started() bool {
	return s.client != nil
}

func (s *Service) Client() *api.Client {
	return s.client
}

func (s *Service) Config() *Config {
	return s.config
}
