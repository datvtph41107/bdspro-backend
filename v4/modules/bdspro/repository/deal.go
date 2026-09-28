package repository

import (
	"context"
	"strings"

	"userv4/modules/bdspro/entity"

	"gorm.io/gorm"
)

// DealRepository is the deal-facing store on the v4 DSN.
type DealRepository interface {
	ListForUser(ctx context.Context, userID uint64, page, size uint32, text string) ([]*entity.Deal, uint32, error)
	ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Deal, uint32, error)
}

type dealPostgres struct {
	db *gorm.DB
}

func NewDealPostgres(db *gorm.DB) DealRepository {
	return &dealPostgres{db: db}
}

func (r *dealPostgres) ListForUser(ctx context.Context, userID uint64, page, size uint32, text string) ([]*entity.Deal, uint32, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Joins("JOIN deal_members dm ON dm.deal_id = deals.id AND dm.deleted_at IS NULL").
		Where("dm.member_id = ?", userID).
		Where("deals.deleted_at IS NULL")
	if strings.TrimSpace(text) != "" {
		like := "%" + strings.TrimSpace(text) + "%"
		q = q.Where("deals.name ILIKE ?", like)
	}

	var total int64
	if err := q.Distinct("deals.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	offset := int(page) * limit

	var list []*entity.Deal
	err := q.Preload("Members").
		Order("deals.updated_at DESC NULLS LAST").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, uint32(total), err
}

func (r *dealPostgres) ListByProductID(ctx context.Context, productID uint64, page, size uint32) ([]*entity.Deal, uint32, error) {
	if r.db == nil || productID == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Joins("JOIN deal_products dp ON dp.deal_id = deals.id AND dp.deleted_at IS NULL").
		Where("dp.product_id = ? AND deals.deleted_at IS NULL", productID)

	var total int64
	if err := q.Distinct("deals.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	if limit > 999 {
		limit = 999
	}
	offset := int(page) * limit

	var list []*entity.Deal
	err := q.Preload("Members").
		Order("deals.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, uint32(total), err
}
