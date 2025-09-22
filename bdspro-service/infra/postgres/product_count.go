package postgres

import (
	base_enum "base/enum"
	"bdspro/internal/domain"
	"context"
)

// CountByOwner đếm số sản phẩm theo ownerOf và ownerId
func (r *GormUserProductRepo) CountByOwner(ctx context.Context, ownerOf base_enum.EOwnerOf, ownerId uint64) (uint32, error) {
	var count int64

	query := r.DB.WithContext(ctx).Model(&domain.Product{}).
		Where("deleted_at IS NULL and owner_id = ? and owner_of = ?", ownerId, ownerOf)

	err := query.Debug().Count(&count).Error
	if err != nil {
		return 0, err
	}

	return uint32(count), nil
}
