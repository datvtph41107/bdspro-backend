package provider

import (
	"context"
	"core/kernel/contract"
	core_ctx "core/kernel/ctx"

	"gorm.io/gorm"
)

type gormDB struct {
	db *gorm.DB
}

func NewGormDB(db *gorm.DB) contract.DB {
	return &gormDB{db: db}
}

func (g *gormDB) Transaction(fn func(ctx context.Context) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		ctx := core_ctx.WithTx(context.Background(), tx)
		return fn(ctx)
	})
}
