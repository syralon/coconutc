// @internal/infra/tx/txn.go

package tx

import (
	"context"

	"{{.Module}}/ent"
	"{{.Module}}/internal/domain/txn"
	"github.com/syralon/coconut/proto/syralon/coconut/errors"
)

type txContextKey struct{}

func withContext(ctx context.Context, tx *ent.Client) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func fromContext(ctx context.Context) (*ent.Client, bool) {
	val := ctx.Value(txContextKey{})
	if val == nil {
		return nil, false
	}
	tx, ok := val.(*ent.Client)
	return tx, ok
}

type Repository struct {
	c *ent.Client
}

func NewRepository(c *ent.Client) *Repository {
	return &Repository{c: c}
}

func (rep *Repository) client(ctx context.Context) *ent.Client {
	if tx, ok := fromContext(ctx); ok {
		return tx
	}
	return rep.c
}

{{ range .Services }}
func (rep *Repository) {{.}}(ctx context.Context) *ent.{{.}}Client {
	return rep.client(ctx).{{.}}
}
{{ end }}


func (rep *Repository) Tx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := rep.c.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(errors.Recovery(recover()))
		if err == nil {
			err = tx.Commit()
		} else {
			err = errors.Join(err, tx.Rollback())
		}
	}()
	return fn(ctx)
}
