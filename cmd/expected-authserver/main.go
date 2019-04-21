package main

import (
	"github.com/expectedsh/expected/pkg/app"
	"github.com/expectedsh/expected/pkg/authserver"
)

type Config struct {
}

func main() {
	if err := app.Start(&authserver.App{}); err != nil {
		panic(err)
	}
}
