package repository

import (
	"context"

	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// PropertyRepository is the property-facing store on the v4 DSN.
type PropertyRepository interface {
	ListMine(ctx context.Context, originID uint64, page, size uint32) ([]*entity.PropertyLineage, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.PropertyLineage, error)
	GetDetail(ctx context.Context, id, ownerOriginID uint64) (*entity.PropertyLineageDetail, error)
	UserCreate(ctx context.Context, originID uint64, bundle *entity.PropertyCreateBundle) (*entity.PropertyLineage, error)
}

type propertyPostgres struct {
	db *gorm.DB
}

func NewPropertyPostgres(db *gorm.DB) PropertyRepository {
	return &propertyPostgres{db: db}
}

func (r *propertyPostgres) ListMine(ctx context.Context, originID uint64, page, size uint32) ([]*entity.PropertyLineage, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.PropertyLineage{}).
		Joins("JOIN property_user pu ON pu.property_lineage_id = property_lineage.id AND pu.deleted_at IS NULL").
		Where("pu.owner_origin_id = ?", originID).
		Where("property_lineage.deleted_at IS NULL")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	offset := int(page) * limit

	var list []*entity.PropertyLineage
	err := q.Preload("Info").
		Preload("Info.PropertyType").
		Preload("Info.Avatar").
		Preload("Location").
		Preload("Location.Province").
		Preload("Location.Ward").
		Preload("LandInfo").
		Preload("BuildingInfo").
		Preload("Statistic").
		Preload("Identify").
		Preload("MediaList", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Order("property_lineage.updated_at DESC NULLS LAST").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, total, err
}

func (r *propertyPostgres) GetByID(ctx context.Context, id uint64) (*entity.PropertyLineage, error) {
	if r.db == nil {
		return nil, nil
	}
	var p entity.PropertyLineage
	err := r.db.WithContext(ctx).
		Preload("Info", "deleted_at IS NULL").
		Preload("Info.PropertyType", "deleted_at IS NULL").
		Preload("Info.Avatar", "deleted_at IS NULL").
		Preload("Location", "deleted_at IS NULL").
		Preload("Location.Province").
		Preload("Location.Ward").
		Preload("LandInfo", "deleted_at IS NULL").
		Preload("BuildingInfo", "deleted_at IS NULL").
		Preload("Statistic", "deleted_at IS NULL").
		Preload("Identify", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
