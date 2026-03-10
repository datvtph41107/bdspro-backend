package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyStatisticRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyStatisticRepo(db *gorm.DB) repo.PropertyStatisticRepo {
	return &PostgrePropertyStatisticRepo{DB: db}
}

func (r *PostgrePropertyStatisticRepo) Save(ctx context.Context, entity *domain.PropertyStatistic) error {
	return GetDB(ctx, r.DB).Save(entity).Error
}
