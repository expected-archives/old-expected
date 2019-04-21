package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/imageserver"
	_ "github.com/lib/pq"
)

func main() {
	if err := apps.Start(&imageserver.App{}); err != nil {
		panic(err)
	}
}
