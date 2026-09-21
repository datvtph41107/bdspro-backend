package admin_postgres

import (
	"bdspro/internal"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	_db "common/db"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/admin.AdminPropertyTypeRepo
type AdminPropertyTypePostgres struct {
	*crud.CrudRepo[domain.PropertyType]
}

func NewAdminPropertyTypePostgres(db *_db.TransactionRepo) *AdminPropertyTypePostgres {
	x := &AdminPropertyTypePostgres{
		CrudRepo: &crud.CrudRepo[domain.PropertyType]{},
	}
	x.CrudRepo.Init(x, db)
	return x
}

func (r *AdminPropertyTypePostgres) QueryDomain(query *gorm.DB, dto *dto.TextSearchRequest) *gorm.DB {
	if dto.Text != "" {
		query = query.Where("property_type.name ILIKE ?", "%"+dto.Text+"%")
	}
	query = query.Where("property_type.deleted_at is NULL")
	return query
}

func (r *AdminPropertyTypePostgres) GetList(c context.Context, pagable _dto.IPagable) ([]domain.PropertyType, int64, error) {
	dto, ok := pagable.(*dto.TextSearchRequest)
	if !ok {
		return nil, 0, _errors.ReturnError(service.InvalidPagable)
	}

	var propertyTypes []domain.PropertyType
	var total int64

	query := r.GetDB(c).
		Debug().
		Model(domain.PropertyType{}).
		Order("property_type.created_at asc").
		Where("property_type.deleted_at is NULL")

	query = r.QueryDomain(query, dto)
	err := query.
		Order("property_type.updated_at desc").
		Limit(pagable.GetLimit()).
		Offset(pagable.GetOffset()).
		Find(&propertyTypes).Error

	if err != nil {
		return nil, 0, err
	}

	err = query.Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	return propertyTypes, total, nil
}
