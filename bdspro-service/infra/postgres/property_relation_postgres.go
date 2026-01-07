package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyRelationRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyRelationRepo(db *gorm.DB) repo.PropertyRelationRepo {
	return &PostgrePropertyRelationRepo{
		DB: db,
	}
}

func (r *PostgrePropertyRelationRepo) Create(ctx context.Context, relation *domain.PropertyRelation) error {
	return r.DB.WithContext(ctx).Create(relation).Error
}

func (r *PostgrePropertyRelationRepo) Update(ctx context.Context, relation *domain.PropertyRelation) error {
	return r.DB.WithContext(ctx).
		Model(&domain.PropertyRelation{}).
		Where("property_id = ?", relation.PropertyID).
		Updates(relation).Error
}

func (r *PostgrePropertyRelationRepo) GetByPropertyID(ctx context.Context, propertyID uint64) (*domain.PropertyRelation, error) {
	var relation domain.PropertyRelation
	err := r.DB.WithContext(ctx).
		Where("property_id = ?", propertyID).
		First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &relation, nil
}

func (r *PostgrePropertyRelationRepo) GetByAssetID(ctx context.Context, assetID uint64) (*domain.PropertyRelation, error) {
	var relation domain.PropertyRelation
	err := r.DB.WithContext(ctx).
		Where("asset_id = ?", assetID).
		First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &relation, nil
}

func (r *PostgrePropertyRelationRepo) GetByProductID(ctx context.Context, productID uint64) (*domain.PropertyRelation, error) {
	var relation domain.PropertyRelation
	err := r.DB.WithContext(ctx).
		Where("primary_product_id = ?", productID).
		First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &relation, nil
}

func (r *PostgrePropertyRelationRepo) GetCountsByPropertyID(ctx context.Context, propertyID uint64) (productCount uint32, assetCount uint32, postCount uint32, err error) {
	// Đếm số product (primary_product_id)
	var pCount int64
	err = r.DB.WithContext(ctx).Model(&domain.PropertyRelation{}).
		Where("property_id = ? AND primary_product_id IS NOT NULL", propertyID).
		Count(&pCount).Error
	if err != nil {
		return 0, 0, 0, err
	}
	productCount = uint32(pCount)

	// Đếm số asset (asset_id)
	var aCount int64
	err = r.DB.WithContext(ctx).Model(&domain.PropertyRelation{}).
		Where("property_id = ? AND asset_id IS NOT NULL", propertyID).
		Count(&aCount).Error
	if err != nil {
		return 0, 0, 0, err
	}
	assetCount = uint32(aCount)

	// Đếm số post thông qua products liên quan đến property này
	// Lấy primary_product_id từ property_relation
	var relation domain.PropertyRelation
	err = r.DB.WithContext(ctx).
		Where("property_id = ?", propertyID).
		First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Không có relation thì trả về 0
			return productCount, assetCount, 0, nil
		}
		return 0, 0, 0, err
	}

	// Nếu có primary_product_id, đếm posts của product đó
	if relation.PrimaryProductID != nil {
		var postCountInt int64
		err = r.DB.WithContext(ctx).Table("posts").
			Where("product_id = ? AND deleted_at IS NULL", *relation.PrimaryProductID).
			Count(&postCountInt).Error
		if err != nil {
			return 0, 0, 0, err
		}
		postCount = uint32(postCountInt)
	}

	return productCount, assetCount, postCount, nil
}
