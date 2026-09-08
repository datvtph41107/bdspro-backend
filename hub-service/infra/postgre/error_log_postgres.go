package postgres

import (
	"context"

	"hub/internal/domain"
	"hub/internal/repo"

	_db "common/db"
)

type ErrorLogPostgres struct {
	db *_db.TransactionRepo
}

func NewErrorLogRepo(db *_db.TransactionRepo) repo.IErrorLogRepo {
	return &ErrorLogPostgres{db: db}
}

func (r *ErrorLogPostgres) Insert(ctx context.Context, log *domain.ErrorLog) error {
	return r.db.GetDB(ctx).WithContext(ctx).Create(log).Error
}
