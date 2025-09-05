package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DealOfOrganizationRepo interface {
	Create(ctx context.Context, dealOfOrganization *entity.DealOfOrganization) error
	GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]entity.DealOfOrganization, int64, error)
	GetByDealID(ctx context.Context, dealID uint64) (*entity.DealOfOrganization, error)
	DeleteByDealID(ctx context.Context, dealID uint64) error
}
