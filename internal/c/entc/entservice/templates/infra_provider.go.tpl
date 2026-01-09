// @file: internal/infra/provider.go

package infra

import (
	"{{.Module}}/internal/infra/data"
	"{{.Module}}/internal/infra/dependency"
	"{{.Module}}/internal/infra/tx"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	dependency.NewEnt,
	dependency.NewETCDClient,

	tx.NewRepository,
    {{ range .Services }}
    data.New{{.}}Repository,{{ end }}
    {{ range .Services }}
    data.To{{.}}Repository,{{ end }}
)
