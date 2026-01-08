// @internal/config/config.go

package config

import (
	"github.com/syralon/coconut/transport/gateway"
	"github.com/syralon/coconut/transport/grpc"
)

type Config struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway"`
	GRPC    grpc.Config    `json:"grpc"    yaml:"grpc"`

	Database Database `json:"database" yaml:"database"`
}

type Database struct {
	Driver string `json:"driver" yaml:"driver"`
	DSN    string `json:"dsn"    yaml:"dsn"`
}
