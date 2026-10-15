package repository

import (
	"context"

	infpostgres "userv4/infra/postgres"
	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

type ProductPriceRepository interface {
	GetByID(ctx context.Context, id uint64) (*entity.ProductPrice, error)
	Create(ctx context.Context, row *entity.ProductPrice) error
	UpdatePrice(ctx context.Context, priceID *uint64, row *entity.ProductPrice) error
}

type productPricePostgres struct {
	db *gorm.DB
}

func NewProductPricePostgres(db *gorm.DB) ProductPriceRepository {
	return &productPricePostgres{db: db}
}

func (r *productPricePostgres) dbCtx(ctx context.Context) *gorm.DB {
	return infpostgres.GetDB(ctx, r.db)
}

func (r *productPricePostgres) GetByID(ctx context.Context, id uint64) (*entity.ProductPrice, error) {
	var row entity.ProductPrice
	if err := r.dbCtx(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *productPricePostgres) Create(ctx context.Context, row *entity.ProductPrice) error {
	return r.dbCtx(ctx).Create(row).Error
}

func (r *productPricePostgres) UpdatePrice(ctx context.Context, priceID *uint64, row *entity.ProductPrice) error {
	if priceID == nil || row == nil {
		return nil
	}
	return r.dbCtx(ctx).Model(&entity.ProductPrice{}).Where("id = ?", *priceID).Updates(row).Error
}
