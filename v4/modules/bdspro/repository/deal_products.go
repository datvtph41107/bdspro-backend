package repository

import (
	"context"
	"time"

	"userv4/modules/bdspro/entity"
)

func (r *dealPostgres) AddDealProducts(ctx context.Context, dealID uint64, productIDs []uint64, createdBy uint64) error {
	if r.db == nil || dealID == 0 || len(productIDs) == 0 {
		return nil
	}
	ids := dedupeUint64(productIDs)
	if len(ids) == 0 {
		return nil
	}

	var existing []entity.DealProduct
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("deal_id = ? AND product_id IN ?", dealID, ids).
		Find(&existing).Error; err != nil {
		return err
	}

	active := make(map[uint64]struct{}, len(existing))
	softDeleted := make(map[uint64]uint64, len(existing)) // productID -> row id
	for _, row := range existing {
		if row.DeletedAt.Valid {
			softDeleted[row.ProductID] = row.ID
			continue
		}
		active[row.ProductID] = struct{}{}
	}

	now := time.Now()
	toCreate := make([]entity.DealProduct, 0, len(ids))
	for _, pid := range ids {
		if _, ok := active[pid]; ok {
			continue
		}
		if rowID, ok := softDeleted[pid]; ok {
			updates := map[string]interface{}{
				"deleted_at": nil,
				"updated_at": now,
			}
			if createdBy > 0 {
				updates["created_by"] = createdBy
				updates["updated_by"] = createdBy
			}
			if err := r.db.WithContext(ctx).Unscoped().
				Model(&entity.DealProduct{}).
				Where("id = ?", rowID).
				Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		row := entity.DealProduct{
			DealID:    dealID,
			ProductID: pid,
		}
		if createdBy > 0 {
			row.CreatedBy = &createdBy
			row.UpdatedBy = &createdBy
		}
		toCreate = append(toCreate, row)
	}
	if len(toCreate) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&toCreate).Error
}

func (r *dealPostgres) RemoveDealProductsByCreator(ctx context.Context, dealID uint64, productIDs []uint64, createdBy uint64) (int64, error) {
	if r.db == nil || dealID == 0 || createdBy == 0 || len(productIDs) == 0 {
		return 0, nil
	}
	ids := dedupeUint64(productIDs)
	result := r.db.WithContext(ctx).
		Where("deal_id = ? AND product_id IN ? AND created_by = ? AND deleted_at IS NULL", dealID, ids, createdBy).
		Delete(&entity.DealProduct{})
	return result.RowsAffected, result.Error
}

func dedupeUint64(in []uint64) []uint64 {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(in))
	out := make([]uint64, 0, len(in))
	for _, id := range in {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
