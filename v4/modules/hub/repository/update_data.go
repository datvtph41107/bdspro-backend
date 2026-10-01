package repository

import (
	"context"

	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

// UpdateDataRepository manages update_data persistence for sync tracking.
type UpdateDataRepository interface {
	Upsert(ctx context.Context, ud *entity.UpdateData) error
	GetChangedIDsSince(ctx context.Context, ownerID uint64, rt entity.ESyncResource, lastSync int64, limit int) ([]uint64, error)
	FlushOld(ctx context.Context, ownerID uint64, rt entity.ESyncResource, keepLimit int64) error
}

type updateDataPostgres struct{ db *gorm.DB }

func NewUpdateDataPostgres(db *gorm.DB) UpdateDataRepository {
	return &updateDataPostgres{db: db}
}

func (r *updateDataPostgres) Upsert(ctx context.Context, ud *entity.UpdateData) error {
	var existing entity.UpdateData
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND resource_type = ? AND resource_id = ?", ud.OwnerID, ud.ResourceType, ud.ResourceID).
		First(&existing).Error
	if err == nil {
		return r.db.WithContext(ctx).Model(&existing).Update("updated_at", ud.UpdatedAt).Error
	}
	return r.db.WithContext(ctx).Create(ud).Error
}

func (r *updateDataPostgres) GetChangedIDsSince(ctx context.Context, ownerID uint64, rt entity.ESyncResource, lastSync int64, limit int) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).Model(&entity.UpdateData{}).
		Select("resource_id").
		Where("owner_id = ? AND resource_type = ? AND updated_at > ?", ownerID, rt, lastSync).
		Order("updated_at DESC").Limit(limit).Pluck("resource_id", &ids).Error
	return ids, err
}

func (r *updateDataPostgres) FlushOld(ctx context.Context, ownerID uint64, rt entity.ESyncResource, keepLimit int64) error {
	if keepLimit <= 0 {
		return nil
	}
	var cutoff int64
	err := r.db.WithContext(ctx).Model(&entity.UpdateData{}).
		Select("updated_at").
		Where("owner_id = ? AND resource_type = ?", ownerID, rt).
		Order("updated_at DESC").Offset(int(keepLimit)).Limit(1).
		Pluck("updated_at", &cutoff).Error
	if err != nil || cutoff == 0 {
		return err
	}
	return r.db.WithContext(ctx).
		Where("owner_id = ? AND resource_type = ? AND updated_at < ?", ownerID, rt, cutoff).
		Delete(&entity.UpdateData{}).Error
}
