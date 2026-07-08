package postgre

import (
	"context"
	"crm/infra/impl"
	"crm/internal/domain"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// @bind: crm/internal/repo.ProductCareRepo
type ProductCarePostgre struct {
	DB *gorm.DB
}

func NewProductCarePostgre(DB *gorm.DB) *ProductCarePostgre {
	return &ProductCarePostgre{DB: DB}
}

// SyncLeadProductsCTE cập nhật tb_product_care chỉ với 1 câu SQL
func (repo ProductCarePostgre) Bulk(ctx context.Context, productCares []uint64, leadID uint64) error {
	const q = `
WITH new_data AS (
    SELECT unnest($1::bigint[]) AS product_id
), deleted AS (
    DELETE FROM tb_product_care
    WHERE lead_id = $2 AND product_id NOT IN (SELECT product_id FROM new_data)
), inserted AS (
    INSERT INTO tb_product_care (lead_id, product_id)
    SELECT $3, product_id FROM new_data
    ON CONFLICT DO NOTHING
)
SELECT 1;
`
	return impl.GetDB(ctx, repo.DB).
		Exec(q, pq.Array(productCares), leadID, leadID).
		Error
}

func (repo ProductCarePostgre) Create(ctx context.Context, productCares []domain.ProductCareEntity) error {
	return impl.GetDB(ctx, repo.DB).
		Model(&domain.ProductCareEntity{}).
		Create(&productCares).Error
}