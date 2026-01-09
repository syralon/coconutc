// @file: internal/domain/txn/txn.go

package txn

import "context"

type Txn interface {
	Tx(ctx context.Context, fn func(ctx context.Context) error) (err error)
}
