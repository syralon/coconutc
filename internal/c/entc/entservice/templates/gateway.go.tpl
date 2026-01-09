// @file: internal/transport/server/gateway.go

package server

import (
	"{{.Module}}/internal/config"
	"{{.Module}}/internal/transport/middleware"
	"{{.Module}}/internal/transport/service"
	pb "{{.ProtoPackage}}"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/syralon/coconut/transport/gateway"
	cocomiddleware "github.com/syralon/coconut/transport/gateway/middleware"
)

func NewGatewayServer(c *config.Config, services *service.Services) *gateway.Server {
	srv := gateway.NewServer(&c.Gateway)
	srv.WithOptions(
		runtime.WithMiddlewares(
		    middleware.Logger(),
			cocomiddleware.Recovery(),
		),
	)
	if c.Gateway.Endpoint != "" {
		srv.RegisterEndpoint(
			c.Gateway.Endpoint,
			{{ range .GatewayServices }}pb.Register{{.}}ServiceHandlerFromEndpoint,
            {{ end }}
		)
	} else {
		srv.Register(
            {{ range .GatewayServices }}gateway.ServerRegister[pb.{{.}}ServiceServer](services.{{.}}, pb.Register{{.}}ServiceHandlerServer),
            {{ end }}
		)
	}
	return srv
}
