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

func NewBaseRepo(db *gorm.DB) *BaseRepo {
	return &BaseRepo{
		TransactionRepo: &TransactionRepo{db: db},
	}
}

func (t *BaseRepo) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return t.db.WithContext(ctx)
}

type ITransactionRepo interface {
	Begin() *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (t *TransactionRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, "tx", tx))
	})
}

func (t *TransactionRepo) StartTransaction(ctx context.Context) context.Context {
	tx := t.db.Begin()
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
	return t.db.WithContext(ctx)
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return defaultDB
}
