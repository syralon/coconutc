// @file: cmd/{{ .Module | basepath }}/main.go

package main

import (
	"context"
	"os"

	"{{.Module}}/internal/bootstrap"
	"{{.Module}}/internal/config"
	"{{.Module}}/version"

	"github.com/google/gops/agent"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		version.Show()
		os.Exit(0)
	}

	if err := agent.Listen(agent.Options{}); err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	cfg, err := config.Read(ctx)
	if err != nil {
		panic(err)
	}
	app, cancel, err := bootstrap.NewApp(cfg)
	if err != nil {
		panic(err)
	}
	defer cancel()
	if err = app.Run(ctx); err != nil {
		panic(err)
	}
}
