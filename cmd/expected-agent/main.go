package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/agent"
)

func main() {
	if err := apps.Start(&agent.App{}); err != nil {
		panic(err)
	}
}
