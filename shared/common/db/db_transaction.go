package _db

import (
	"context"

	"gorm.io/gorm"
)

type IDBTransaction interface {
	Begin() *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type DBTransaction struct {
	DB *gorm.DB
}

func NewDBTransaction(db *gorm.DB) *DBTransaction {
	return &DBTransaction{DB: db}
}

func (t *DBTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.DB.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

func (t *DBTransaction) StartTransaction(ctx context.Context) context.Context {
	tx := t.DB.Begin()
	return context.WithValue(ctx, "tx", tx)
}

func (t *DBTransaction) RollbackTransaction(ctx context.Context) error {
	return t.GetDB(ctx).Rollback().Error
}

func (t *DBTransaction) CommitTransaction(ctx context.Context) error {
	return t.GetDB(ctx).Commit().Error
}

func (t *DBTransaction) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return t.DB
}
