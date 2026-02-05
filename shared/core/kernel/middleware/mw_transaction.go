package middleware

import (
	"context"
	"core/kernel/container"
	core_ctx "core/kernel/ctx"
	"fmt"

	"gorm.io/gorm"
)

// TransactionMiddleware: starts a db tx and stores it in context for downstream
func TransactionMiddleware(c *container.Container) Middleware {
	return MiddlewareFunc(func(ctx context.Context, payload interface{}, next func(context.Context) error) error {
		// resolve db from container
		dbi, err := c.Resolve("db")
		if err != nil {
			// no db configured -> just continue
			return next(ctx)
		}
		db, ok := dbi.(*gorm.DB)
		if !ok {
			return fmt.Errorf("db binding invalid")
		}
		var errTx error
		// start tx
		tx := db.Begin()
		if tx.Error != nil {
			return tx.Error
		}
		ctxWithTx := core_ctx.WithTx(ctx, tx)
		// run next
		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback()
				panic(r)
			}
		}()
		if errTx = next(ctxWithTx); errTx != nil {
			_ = tx.Rollback()
			return errTx
		}
		// commit
		if err := tx.Commit().Error; err != nil {
			_ = tx.Rollback()
			return err
		}
		return nil
	})
}
