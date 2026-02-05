package contract

import "context"

type DB interface {
	Transaction(fn func(ctx context.Context) error) error
}
