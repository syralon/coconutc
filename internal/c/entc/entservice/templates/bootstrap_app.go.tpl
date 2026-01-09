// @file: internal/bootstrap/app.go

package bootstrap

import (
	"github.com/syralon/coconut"
	"github.com/syralon/coconut/pkg/etcdutil"
	"github.com/syralon/coconut/transport"
	"github.com/syralon/coconut/transport/gateway"
	"github.com/syralon/coconut/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func newApp(client *clientv3.Client, grpcServer *grpc.Server, gatewayServer *gateway.Server) (*coconut.App, error) {
	app := coconut.NewApp(
		coconut.WithHooks(
			transport.Logger(),
			transport.Registry(etcdutil.NewRegistry(client)),
		),
	)
	app.Add(grpcServer, gatewayServer)
	return app, nil
}