package db

import (
	"context"

	coreCTX "core/kernel/ctx"

	"gorm.io/gorm"
)

func WithTransaction(ctx context.Context, db *gorm.DB, fn func(ctx context.Context) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	ctx = coreCTX.WithTx(ctx, tx)

	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
