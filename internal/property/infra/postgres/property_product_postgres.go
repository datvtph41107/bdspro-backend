package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyProductRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyProductRepo(db *gorm.DB) repo.PropertyProductRepo {
	return &PostgrePropertyProductRepo{
		DB: db,
	}
}

func (r *PostgrePropertyProductRepo) Create(ctx context.Context, propertyProduct *domain.PropertyProduct) error {
	return r.DB.WithContext(ctx).Create(propertyProduct).Error
}

func (r *PostgrePropertyProductRepo) GetByProductID(ctx context.Context, productID uint64) ([]domain.PropertyProduct, error) {
	var propertyProducts []domain.PropertyProduct
	err := r.DB.WithContext(ctx).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Find(&propertyProducts).Error
	return propertyProducts, err
}

func (r *PostgrePropertyProductRepo) GetByPropertyID(ctx context.Context, propertyID uint64) ([]domain.PropertyProduct, error) {
	var propertyProducts []domain.PropertyProduct
	err := r.DB.WithContext(ctx).
		Where("property_id = ? AND deleted_at IS NULL", propertyID).
		Find(&propertyProducts).Error
	return propertyProducts, err
}
