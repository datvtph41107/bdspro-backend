package repository

import (
	"context"

	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// PropertyRepository is the property-facing store on the v4 DSN.
type PropertyRepository interface {
	ListMine(ctx context.Context, originID uint64, page, size uint32) ([]*entity.PropertyLineage, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.PropertyLineage, error)
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
	err := q.Preload("Info").Preload("Location").
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
	err := r.db.WithContext(ctx).Preload("Info").Preload("Location").
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
