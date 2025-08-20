package admin_postgres

import (
	"bdspro/infra/postgre"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	_dto "common/dto"
	_errors "common/errors"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/admin.AdminProductRepo
type AdminProductPostgres struct {
	crud.BaseRepo[domain.Product]
	ProductRepo *postgre.GormUserProductRepo
}

func NewAdminProductPostgres(db *gorm.DB, productRepo *postgre.GormUserProductRepo) *AdminProductPostgres {
	return &AdminProductPostgres{
		BaseRepo:    crud.BaseRepo[domain.Product]{DB: db},
		ProductRepo: productRepo,
	}
}

func (r *AdminProductPostgres) GetList(c context.Context, pagable _dto.IPagable) ([]domain.Product, int64, error) {
	dto, ok := pagable.(*dto.ProductSearchRequest)
	if !ok {
		return nil, 0, _errors.ReturnError(400, "invalid pagable")
	}

	var products []domain.Product
	var total int64

	query := r.DB.
		Debug().
		Model(domain.Product{}).
		Select("products.*, products.area - COALESCE(child.total_area, 0) AS avaiable_area").
		Joins(`LEFT JOIN (SELECT parent_id, SUM(area) AS total_area 
					FROM products WHERE deleted_at is null GROUP BY parent_id) 
					AS child ON products.id = child.parent_id`).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		// Preload("Amenities", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		// Preload("SaleTransaction", "deleted_at IS NULL").
		// Preload("RentTransaction", "deleted_at IS NULL").
		Where("products.deleted_at is NULL")

	query = r.ProductRepo.QueryDomain(query, dto)
	err := query.
		Order("updated_at desc").
		Limit(pagable.GetLimit()).
		Offset(pagable.GetOffset()).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	err = query.Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
