package main

import (
	"github.com/expectedsh/expected/pkg/app"
	"github.com/expectedsh/expected/pkg/registryhook"
	_ "github.com/lib/pq"
)

func main() {
	if err := app.Start(&registryhook.App{}); err != nil {
		panic(err)
	}
}
