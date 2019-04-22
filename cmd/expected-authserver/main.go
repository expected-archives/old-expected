package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/authserver"
)

type Config struct {
}

func main() {
	if err := apps.Start(&authserver.App{}); err != nil {
		panic(err)
	}
}
