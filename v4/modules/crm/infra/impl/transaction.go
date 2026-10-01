package impl

import (
	"context"

	"gorm.io/gorm"
)

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	if defaultDB == nil {
		return nil
	}
	return defaultDB.WithContext(ctx)
}
