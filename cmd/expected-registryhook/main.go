package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/registryhook"
	_ "github.com/lib/pq"
)

func main() {
	if err := apps.Start(&registryhook.App{}); err != nil {
		panic(err)
	}
}
