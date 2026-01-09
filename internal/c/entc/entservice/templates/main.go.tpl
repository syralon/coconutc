// @file: cmd/{{ .Module | basepath }}/main.go

package main

import (
	"context"

	"github.com/syralon/coconut-example/internal/bootstrap"
	"github.com/syralon/coconut-example/internal/config"
)

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
