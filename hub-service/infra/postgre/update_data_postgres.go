package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
)

// @bind: hub/internal/repo.IUpdateDataRepo
type UpdateDataRepo struct {
	*_db.TransactionRepo
}

func NewUpdateDataRepo(db *_db.TransactionRepo) *UpdateDataRepo {
	return &UpdateDataRepo{TransactionRepo: db}
}

func (r *UpdateDataRepo) Upsert(c context.Context, entity *domain.UpdateData) error {
	var existing domain.UpdateData
	err := r.GetDB(c).
		Where("owner_id = ? AND resource_type = ? AND resource_id = ?", entity.OwnerID, entity.ResourceType, entity.ResourceID).
		First(&existing).Error
	if err == nil {
		existing.UpdatedAt = entity.UpdatedAt
		return r.GetDB(c).Model(&existing).Update("updated_at", entity.UpdatedAt).Error
	}
	return r.GetDB(c).Create(entity).Error
}

func (r *UpdateDataRepo) DeleteByOwnerResourceId(c context.Context, ownerID uint64, resourceType domain.ESyncResource, resourceID uint64) error {
	return r.GetDB(c).
		Where("owner_id = ? AND resource_type = ? AND resource_id = ?", ownerID, resourceType, resourceID).
		Delete(&domain.UpdateData{}).Error
}

func (r *UpdateDataRepo) GetChangedIdsSince(c context.Context, ownerID uint64, resourceType domain.ESyncResource, lastSync int64, limit int) ([]uint64, error) {
	var ids []uint64
	err := r.GetDB(c).Model(&domain.UpdateData{}).
		Select("resource_id").
		Where("owner_id = ? AND resource_type = ? AND updated_at > ?", ownerID, resourceType, lastSync).
		Order("updated_at DESC").
		Limit(limit).
		Pluck("resource_id", &ids).Error
	return ids, err
}

func (r *UpdateDataRepo) FlushOld(c context.Context, ownerID uint64, resourceType domain.ESyncResource, keepLimit int64) error {
	if keepLimit <= 0 {
		return nil
	}
	// Lấy updated_at của bản ghi thứ keepLimit (từ mới đến cũ)
	var cutoff int64
	err := r.GetDB(c).Model(&domain.UpdateData{}).
		Select("updated_at").
		Where("owner_id = ? AND resource_type = ?", ownerID, resourceType).
		Order("updated_at DESC").
		Offset(int(keepLimit)).
		Limit(1).
		Pluck("updated_at", &cutoff).Error
	if err != nil || cutoff == 0 {
		return err
	}
	// Xóa các bản ghi có updated_at < cutoff
	return r.GetDB(c).
		Where("owner_id = ? AND resource_type = ? AND updated_at < ?", ownerID, resourceType, cutoff).
		Delete(&domain.UpdateData{}).Error
}
