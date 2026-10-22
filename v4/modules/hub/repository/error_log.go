package repository

import (
	"context"

	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

// ErrorLogRepository persists mobile/app error logs.
type ErrorLogRepository interface {
	Insert(ctx context.Context, log *entity.ErrorLog) error
}

type errorLogPostgres struct{ db *gorm.DB }

func NewErrorLogPostgres(db *gorm.DB) ErrorLogRepository {
	return &errorLogPostgres{db: db}
}

func (r *errorLogPostgres) Insert(ctx context.Context, log *entity.ErrorLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
