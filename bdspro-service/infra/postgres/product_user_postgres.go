package postgres

import (
	"bdspro/internal/domain"
	_utils "common/utils"
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

func (r *ProductUserPostgres) UpdatePriceID(
	ctx context.Context,
	productID uint64,
	profileID uint64,
	priceID uint64,
) error {
	return r.db.WithContext(ctx).
		Model(&domain.ProductUser{}).
		Where(
			"product_id = ? AND profile_id = ? AND deleted_at IS NULL",
			productID, profileID,
		).
		Update("price_id", priceID).
		Error
}

func (r *ProductUserPostgres) CreateOwner(ctx context.Context, productID uint64, profileID uint64) (*domain.ProductUser, error) {
	originId := _utils.GetOriginIdFromContext(ctx)
	productOfUser := &domain.ProductUser{
		ProductID: productID,
		ProfileID: profileID,
		IsOwner:   true,

		OriginProfileID: &originId,
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

func (r *ProductUserPostgres) GetByProductAndOriginProfile(
	ctx context.Context,
	productID, originProfileId uint64,
) (*domain.ProductUser, error) {
	var pu domain.ProductUser
	err := GetDB(ctx, r.db).
		Where(
			"product_id = ? AND origin_profile_id = ? AND deleted_at IS NULL",
			productID, originProfileId,
		).
		First(&pu).Error
	if err != nil {
		return nil, err
	}

	return &pu, nil
}

func (r *ProductUserPostgres) GetByDistributeID(ctx context.Context, distributeID uint64) ([]*domain.ProductUser, error) {
	var productOfUsers []*domain.ProductUser
	err := r.db.WithContext(ctx).
		Where("distribute_id = ? AND deleted_at IS NULL", distributeID).
		Find(&productOfUsers).Error
	if err != nil {
		return nil, err
	}
	return productOfUsers, nil
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
