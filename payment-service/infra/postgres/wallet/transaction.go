package walletpostgres

import (
	"context"

	walletuc "payment/internal/usecase/wallet"

	"gorm.io/gorm"
)

type txContextKey struct{}

type transaction struct {
	tx *gorm.DB
}

func (t *transaction) Commit() error {
	if t == nil || t.tx == nil {
		return gorm.ErrInvalidDB
	}
	return t.tx.Commit().Error
}

func (t *transaction) Rollback() error {
	if t == nil || t.tx == nil {
		return gorm.ErrInvalidDB
	}
	return t.tx.Rollback().Error
}

func (t *transaction) Context(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if t == nil || t.tx == nil {
		return ctx
	}
	return context.WithValue(ctx, txContextKey{}, t.tx.WithContext(ctx))
}

func dbFromContext(ctx context.Context, root *gorm.DB) *gorm.DB {
	if ctx != nil {
		if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
			return tx.WithContext(ctx)
		}
	}
	if root == nil {
		return nil
	}
	return root.WithContext(ctx)
}

var _ walletuc.Transaction = (*transaction)(nil)
