package docker

import (
	"github.com/docker/docker/client"
)

var cli *client.Client

func Init() error {
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	cli = c
	return nil
}

func Client() *client.Client {
	return cli
}
