package main

import (
	"fmt"
	"github.com/expectedsh/expected/pkg/registryserver/auth"
	"github.com/expectedsh/expected/pkg/registryserver/auth/token"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 4 || os.Getenv("EMAIL") == "" {
		exit()
	}
	token.Init("./certs/server.crt", "./certs/server.key")
	s, _ := token.Generate(auth.RequestFromDaemon{
		Login:   os.Getenv("EMAIL"),
		Service: "registry",
	}, []auth.AuthorizedScope{
		{
			Scope: auth.Scope{
				Type: os.Args[1],
				Name: os.Args[2],
			},
			AuthorizedActions: strings.Split(os.Args[3], ","),
		},
	})
	fmt.Println(s)
}

func exit() {
	logrus.Infoln("environment variables required: $EMAIL")
	logrus.Infoln("usage: generate-token <type> <resource> <actions>")
	logrus.Infoln("example: EMAIL=Alexis.viscogliosi@outlook.fr generate-token repository hello-world \"pull,push,delete\"")
	os.Exit(0)
}
