package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductUserRepo
type ProductUserPostgres struct {
	db *gorm.DB
}

func NewProductUserPostgres(db *gorm.DB) *ProductUserPostgres {
	return &ProductUserPostgres{db: db}
}

func (r *ProductUserPostgres) CreateOwner(ctx context.Context, productID uint64, profileID uint64) (*domain.ProductUser, error) {
	productOfUser := &domain.ProductUser{
		ProductID: productID,
		ProfileID: profileID,
		IsOwner:   true,
	}
	err := GetDB(ctx, r.db).Create(productOfUser).Error
	if err != nil {
		return nil, err
	}
	return productOfUser, nil
}

func (r *ProductUserPostgres) GetByProductID(ctx context.Context, productID uint64) ([]*domain.ProductUser, error) {
	var productOfUsers []*domain.ProductUser
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Find(&productOfUsers).Error
	if err != nil {
		return nil, err
	}
	return productOfUsers, nil
}

func (r *ProductUserPostgres) GetByProfileID(ctx context.Context, profileID uint64) ([]*domain.ProductUser, error) {
	var productOfUsers []*domain.ProductUser
	err := r.db.WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Find(&productOfUsers).Error
	if err != nil {
		return nil, err
	}
	return productOfUsers, nil
}

func (r *ProductUserPostgres) GetByProductAndProfile(ctx context.Context, productID, profileID uint64) (*domain.ProductUser, error) {
	var productOfUser domain.ProductUser
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND profile_id = ? AND deleted_at IS NULL", productID, profileID).
		First(&productOfUser).Error
	if err != nil {
		return nil, err
	}
	return &productOfUser, nil
}

func (r *ProductUserPostgres) Update(ctx context.Context, productOfUser *domain.ProductUser) (*domain.ProductUser, error) {
	err := r.db.WithContext(ctx).Save(productOfUser).Error
	if err != nil {
		return nil, err
	}
	return productOfUser, nil
}

func (r *ProductUserPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductUser{}, id).Error
}

func (r *ProductUserPostgres) DeleteByProductID(ctx context.Context, productID uint64) error {
	return r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Delete(&domain.ProductUser{}).Error
}
