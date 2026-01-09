// @file: internal/config/config.go

package config

import (
	"context"

	"github.com/syralon/coconut/configuration"
	"github.com/syralon/coconut/pkg/etcdutil"
	"github.com/syralon/coconut/transport/gateway"
	"github.com/syralon/coconut/transport/grpc"
)

type Config struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway"`
	GRPC    grpc.Config    `json:"grpc"    yaml:"grpc"`

	ETCD etcdutil.Config `json:"etcd"      yaml:"etcd"`

	Database Database `json:"database" yaml:"database"`
}

type Database struct {
	Driver string `json:"driver" yaml:"driver"`
	DSN    string `json:"dsn"    yaml:"dsn"`
}

func Read(ctx context.Context) (*Config, error) {
	c := new(Config)
	if err := configuration.Read(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
