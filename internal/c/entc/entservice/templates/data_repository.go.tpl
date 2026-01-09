// @file: internal/infra/data/repository.go

package data

import (
    "{{.Module}}/internal/domain/repository"
)
{{ range .Services }}
func To{{.}}Repository(rep *{{.}}Repository) repository.{{.}}Repository {
    return rep
}
{{ end }}