package admin_postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	case_archived "common/case/archived"
	"common/case/crud"
	_db "common/db"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/admin.AdminProductRepo
type AdminProductPostgres struct {
	*crud.CrudRepo[domain.Product]
	*case_archived.ArchivedRepo[domain.Product]
	ProductRepo repo.ProductRepo
}

func NewAdminProductPostgres(db *_db.TransactionRepo, productRepo repo.ProductRepo) *AdminProductPostgres {
	x := &AdminProductPostgres{
		CrudRepo:     &crud.CrudRepo[domain.Product]{},
		ArchivedRepo: &case_archived.ArchivedRepo[domain.Product]{},
		ProductRepo:  productRepo,
	}

	x.CrudRepo.Init(x, db)
	x.ArchivedRepo.Init(x, db)

	return x
}

func (r *AdminProductPostgres) GetList(c context.Context, pagable _dto.IPagable) ([]domain.Product, int64, error) {
	dto, ok := pagable.(*dto.ProductSearchRequest)
	if !ok {
		return nil, 0, _errors.ReturnError(400, "invalid pagable")
	}
	// r.WithTransaction(c, func(c context.Context) error {
	// 	return nil
	// })
	var products []domain.Product
	var total int64

	query := r.CrudRepo.GetDB(c).
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
		// Order("updated_at desc").
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

func (r *AdminProductPostgres) GetDetail(c context.Context, id uint64) (*domain.Product, error) {
	var product domain.Product
	err := r.CrudRepo.GetDB(c).
		Preload("Price", func(db *gorm.DB) *gorm.DB {
			return db.Order("product_price.created_at asc")
		}).
		Preload("Province", "deleted_at IS NULL").
		Preload("District", "deleted_at IS NULL").
		Preload("Ward", "deleted_at IS NULL").
		Preload("DocType", "deleted_at IS NULL").
		Preload("PropertyType", "deleted_at IS NULL").
		Preload("MediaList", "deleted_at IS NULL").
		Where("id = ? and deleted_at is null", id).
		First(&product).Error
	return &product, err
}
