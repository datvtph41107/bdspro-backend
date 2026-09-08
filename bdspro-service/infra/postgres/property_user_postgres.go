package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PropertyUserPostgres struct {
	db *gorm.DB
}

func NewPropertyUserPostgres(db *gorm.DB) repo.PropertyUserRepo {
	return &PropertyUserPostgres{db: db}
}

func (r *PropertyUserPostgres) CreateOwner(ctx context.Context, propertyId uint64, profileID uint64) (*domain.PropertyUser, error) {
	originId := _utils.GetOriginIdFromContext(ctx)
	propertyOfUser := &domain.PropertyUser{
		PropertyLineageID: &propertyId,
		OwnerOriginID:     profileID,
		OwnerAt:           time.Now(),

		OriginProfileID: &originId,
	}
	err := GetDB(ctx, r.db).Create(propertyOfUser).Error
	if err != nil {
		return nil, err
	}
	return propertyOfUser, nil
}

func (r *PropertyUserPostgres) GetByLineageAndOwner(
	ctx context.Context,
	propertyId uint64,
	ownerOriginID uint64,
) (*domain.PropertyUser, error) {
	var propertyUser domain.PropertyUser
	err := GetDB(ctx, r.db).
		Where("property_lineage_id = ? AND owner_origin_id = ? AND deleted_at IS NULL", propertyId, ownerOriginID).
		First(&propertyUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &propertyUser, nil
}

func (r *PropertyUserPostgres) GetByLineageAndOwnerUnscoped(
	ctx context.Context,
	propertyId uint64,
	ownerOriginID uint64,
) (*domain.PropertyUser, error) {
	var propertyUser domain.PropertyUser
	err := GetDB(ctx, r.db).
		Unscoped().
		Where("property_lineage_id = ? AND owner_origin_id = ?", propertyId, ownerOriginID).
		First(&propertyUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &propertyUser, nil
}

func (r *PropertyUserPostgres) RestoreByID(ctx context.Context, id uint64) error {
	return GetDB(ctx, r.db).
		Unscoped().
		Model(&domain.PropertyUser{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (r *PropertyUserPostgres) UpdateFieldsByLineageAndOwner(
	ctx context.Context,
	propertyId uint64,
	ownerOriginID uint64,
	fields map[string]any,
) error {
	if len(fields) == 0 {
		return nil
	}

	db := GetDB(ctx, r.db)
	var propertyUser domain.PropertyUser
	err := db.
		Where("property_lineage_id = ? AND owner_origin_id = ? AND deleted_at IS NULL", propertyId, ownerOriginID).
		First(&propertyUser).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		originId := _utils.GetOriginIdFromContext(ctx)
		propertyUser = domain.PropertyUser{
			PropertyLineageID: &propertyId,
			OwnerOriginID:     ownerOriginID,
			OwnerAt:           time.Now(),
			OriginProfileID:   &originId,
		}
		if err := db.Create(&propertyUser).Error; err != nil {
			return err
		}
	}

	return db.Model(&propertyUser).Updates(fields).Error
}
