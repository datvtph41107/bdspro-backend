package repository

import (
	"context"

	infpostgres "userv4/infra/postgres"
	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

type ProductUserRepository interface {
	GetByDistributeID(ctx context.Context, distributeID uint64) ([]*entity.ProductUser, error)
}

type productUserPostgres struct {
	db *gorm.DB
}

func NewProductUserPostgres(db *gorm.DB) ProductUserRepository {
	return &productUserPostgres{db: db}
}

func (r *productUserPostgres) dbCtx(ctx context.Context) *gorm.DB {
	return infpostgres.GetDB(ctx, r.db)
}

func (r *productUserPostgres) GetByDistributeID(ctx context.Context, distributeID uint64) ([]*entity.ProductUser, error) {
	var rows []*entity.ProductUser
	err := r.dbCtx(ctx).
		Where("distribute_id = ? AND deleted_at IS NULL", distributeID).
		Find(&rows).Error
	return rows, err
}
