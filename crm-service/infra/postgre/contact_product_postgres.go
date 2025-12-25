package postgre

import (
	_dto "common/domain/dto"
	"context"
	"crm/infra/impl"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// @bind: crm/internal/repo.ContactProductRepo
type ContactProductPostgre struct {
	DB *gorm.DB
}

func NewContactProductPostgre(DB *gorm.DB) *ContactProductPostgre {
	return &ContactProductPostgre{DB: DB}
}

// SyncContactProductsCTE cập nhật tb_contact_product chỉ với 1 câu SQL
func (repo ContactProductPostgre) Bulk(ctx context.Context, productIds []uint64, contactID uint64) error {
	const q = `
WITH new_data AS (
    SELECT unnest($1::bigint[]) AS product_id
), deleted AS (
    DELETE FROM tb_contact_product
    WHERE contact_id = $2 AND product_id NOT IN (SELECT product_id FROM new_data)
), inserted AS (
    INSERT INTO tb_contact_product (contact_id, product_id)
    SELECT $3, product_id FROM new_data
    ON CONFLICT (contact_id, product_id) DO NOTHING
)
SELECT 1;
`
	return impl.GetDB(ctx, repo.DB).
		Exec(q, pq.Array(productIds), contactID, contactID).
		Error
}

func (repo ContactProductPostgre) GetProductIds(ctx context.Context, contactID uint64) ([]uint64, error) {
	var productIds []uint64
	err := impl.GetDB(ctx, repo.DB).
		Model(&struct {
			ProductID uint64 `gorm:"column:product_id"`
		}{}).
		Table("tb_contact_product").
		Select("product_id").
		Where("contact_id = ?", contactID).
		Pluck("product_id", &productIds).
		Error
	return productIds, err
}

func (repo ContactProductPostgre) GetProductIdsWithPagination(ctx context.Context, contactID uint64, pagable _dto.Pagable) ([]uint64, int64, error) {
	var productIds []uint64
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()

	// Count total
	var total int64
	countQuery := impl.GetDB(ctx, repo.DB).
		Table("tb_contact_product").
		Where("contact_id = ?", contactID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := impl.GetDB(ctx, repo.DB).
		Model(&struct {
			ProductID uint64 `gorm:"column:product_id"`
		}{}).
		Table("tb_contact_product").
		Select("product_id").
		Where("contact_id = ?", contactID).
		Offset(offset).
		Limit(limit).
		Pluck("product_id", &productIds).
		Error

	return productIds, total, err
}
