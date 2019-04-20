package main

import (
	"github.com/expectedsh/expected/pkg/app"
	"github.com/expectedsh/expected/pkg/scheduler"
	_ "github.com/lib/pq"
)

func main() {
	if err := app.Start(&scheduler.App{}); err != nil {
		panic(err)
	}
}
