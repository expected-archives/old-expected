package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/controller"
	_ "github.com/lib/pq"
)

func main() {
	if err := apps.Start(&controller.App{}); err != nil {
		panic(err)
	}
}
