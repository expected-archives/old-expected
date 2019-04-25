package main

import (
	"github.com/expectedsh/expected/pkg/apps"
	"github.com/expectedsh/expected/pkg/apps/metricsagent"
)

func main() {
	if err := apps.Start(&metricsagent.App{}); err != nil {
		panic(err)
	}
}
