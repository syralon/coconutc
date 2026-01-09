// @file: internal/transport/provider.go

package transport

import (
	"{{.Module}}/internal/transport/server"
	"{{.Module}}/internal/transport/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewServices,
	
	server.NewGRPCServer,
	server.NewGatewayServer,
)
