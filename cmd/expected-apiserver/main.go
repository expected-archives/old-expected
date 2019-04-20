package main

import (
	"github.com/expectedsh/expected/pkg/apiserver"
	"github.com/expectedsh/expected/pkg/app"
	_ "github.com/lib/pq"
)

func main() {
	if err := app.Start(&apiserver.App{}); err != nil {
		panic(err)
	}
}
