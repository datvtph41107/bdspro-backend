package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type BusinessDomainRepository interface {
	Create(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error)
	FindById(ctx context.Context, id uint32) (*entity.BusinessDomain, error)
	FindByCode(ctx context.Context, code string) (*entity.BusinessDomain, error)
	FindAll(ctx context.Context) ([]*entity.BusinessDomain, error)
	FindActive(ctx context.Context) ([]*entity.BusinessDomain, error)
	Update(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error)
	Delete(ctx context.Context, id uint32) error
	FindByIds(ctx context.Context, ids []uint32) ([]*entity.BusinessDomain, error)
}
