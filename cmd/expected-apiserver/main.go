package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/apiserver"
	_ "github.com/lib/pq"
)

func main() {
	if err := apps.Start(&apiserver.App{}); err != nil {
		panic(err)
	}
}
