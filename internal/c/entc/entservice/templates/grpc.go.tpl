// @file: internal/transport/server/grpc.go

package server

import (
	"context"

	"{{.Module}}/internal/config"
    "{{.Module}}/internal/transport/interceptor"
	"{{.Module}}/internal/transport/service"
	pb "{{.ProtoPackage}}"

	"github.com/syralon/coconut/proto/syralon/coconut/field"
	"github.com/syralon/coconut/transport/grpc"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func NewGRPCServer(c *config.Config, services *service.Services) *grpc.Server {
	srv := grpc.NewServer(&c.GRPC)

	srv.WithUnaryInterceptor(interceptor.LogServerInterceptor())
	srv.WithUnaryInterceptor(func(ctx context.Context, req any, info *stdgrpc.UnaryServerInfo, handler stdgrpc.UnaryHandler) (resp any, err error) {
		if err = field.Bind(ctx, req.(proto.Message)); err != nil { // bind header or metadata value into message
			return resp, err
		}
		return handler(ctx, req)
	})

	srv.Register(func(srv *stdgrpc.Server) {
		// grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
        {{ range .Services }}pb.Register{{.}}ServiceServer(srv, services.{{.}})
        {{ end }}
	})

	srv.WithOTELHandler()

	return srv
}
