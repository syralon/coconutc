// @file: internal/infra/dependency/ent.go

package dependency

import (
	"{{.Module}}/ent"
	"{{.Module}}/internal/config"

)

func NewEnt(c *config.Config) (*ent.Client, func(), error) {
	client, err := ent.Open(c.Database.Driver, c.Database.DSN)
	if err != nil {
		return nil, nil, err
	}
	return client, func() { _ = client.Close() }, nil
}
