package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
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
	return GetDB(ctx, r.DB).Create(relation).Error
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

func (r *PostgrePropertyRelationRepo) GetAllByPropertyID(ctx context.Context, propertyID uint64) ([]domain.PropertyRelation, error) {
	var relations []domain.PropertyRelation
	err := r.DB.WithContext(ctx).
		Where("property_id = ?", propertyID).
		Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *PostgrePropertyRelationRepo) GetByAssetID(ctx context.Context, assetID uint64) (*domain.PropertyRelation, error) {
	var relation domain.PropertyRelation
	err := r.DB.WithContext(ctx).
		Where("relation_id = ? AND relation_type = ?", assetID, enums.ERelationTypeAsset).
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
		Where("relation_id = ? AND relation_type = ?", productID, enums.ERelationTypeProduct).
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
	// Đếm số product (relation_type = 10)
	var pCount int64
	err = r.DB.WithContext(ctx).Model(&domain.PropertyRelation{}).
		Where("property_id = ? AND relation_type = ? AND relation_id IS NOT NULL", propertyID, enums.ERelationTypeProduct).
		Count(&pCount).Error
	if err != nil {
		return 0, 0, 0, err
	}
	productCount = uint32(pCount)

	// Đếm số asset (relation_type = 20)
	var aCount int64
	err = r.DB.WithContext(ctx).Model(&domain.PropertyRelation{}).
		Where("property_id = ? AND relation_type = ? AND relation_id IS NOT NULL", propertyID, enums.ERelationTypeAsset).
		Count(&aCount).Error
	if err != nil {
		return 0, 0, 0, err
	}
	assetCount = uint32(aCount)

	// Đếm số post thông qua products liên quan đến property này
	// Lấy relation_id từ property_relation với relation_type = product
	var relation domain.PropertyRelation
	err = r.DB.WithContext(ctx).
		Where("property_id = ? AND relation_type = ?", propertyID, enums.ERelationTypeProduct).
		First(&relation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Không có relation thì trả về 0
			return productCount, assetCount, 0, nil
		}
		return 0, 0, 0, err
	}

	// Nếu có relation_id (product ID), đếm posts của product đó
	if relation.RelationID != nil {
		var postCountInt int64
		err = r.DB.WithContext(ctx).Table("posts").
			Where("product_id = ? AND deleted_at IS NULL", *relation.RelationID).
			Count(&postCountInt).Error
		if err != nil {
			return 0, 0, 0, err
		}
		postCount = uint32(postCountInt)
	}

	return productCount, assetCount, postCount, nil
}
