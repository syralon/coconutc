// @file: internal/bootstrap/wire.go

//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package bootstrap

import (
	"{{.Module}}/internal/config"
	"{{.Module}}/internal/infra"
	"{{.Module}}/internal/transport"

	"github.com/google/wire"
	"github.com/syralon/coconut"
)

func NewApp(config *config.Config) (*coconut.App, func(), error) {
	panic(wire.Build(
		infra.ProviderSet,
		transport.ProviderSet,
		newApp,
	))
}
