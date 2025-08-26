package _db

import (
	"context"

	"gorm.io/gorm"
)

type IRepo interface {
	ITransactionRepo
	GetDB(ctx context.Context) *gorm.DB
}

type BaseRepo struct {
	*TransactionRepo
}

func (t *BaseRepo) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return t.DB
}

type ITransactionRepo interface {
	Begin() *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionRepo struct {
	DB *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{DB: db}
}

func (t *TransactionRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.DB.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

func (t *TransactionRepo) StartTransaction(ctx context.Context) context.Context {
	tx := t.DB.Begin()
	return context.WithValue(ctx, "tx", tx)
}

func (t *TransactionRepo) RollbackTransaction(ctx context.Context) error {
	return t.GetDB(ctx).Rollback().Error
}

func (t *TransactionRepo) CommitTransaction(ctx context.Context) error {
	return t.GetDB(ctx).Commit().Error
}

func (t *TransactionRepo) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return t.DB
}
