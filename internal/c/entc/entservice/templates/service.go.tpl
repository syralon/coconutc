// @internal/transport/service/service.go

package service

import (
    {{ range .Services }}"{{$.Module}}/internal/transport/service/{{ . | toLower }}service"
    {{ end }}
    "github.com/google/wire"
)

type Services struct {
    {{ range .Services }}{{.}} *{{ . | toLower }}service.{{.}}Service
    {{ end }}
}

func NewServices(
    {{ range .Services }}{{ . | toLower }}Repo repository.{{.}}Repository,
    {{ end }}
) *Services {
    return &Services{
        {{ range .Services }}{{.}}: {{ . | toLower }}service.New{{.}}Service({{.|toLower}}Repo),
        {{ end }}
    }
}

var ProviderSet = wire.NewSet(
    NewServices,
)
