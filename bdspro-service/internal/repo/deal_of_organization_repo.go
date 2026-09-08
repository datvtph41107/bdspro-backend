package repo

import (
	"bdspro/internal/domain"
	"context"
)

type DealOfOrganizationRepo interface {
	Create(ctx context.Context, dealOfOrganization *domain.DealOfOrganization) error
	GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]domain.DealOfOrganization, int64, error)
	GetByDealID(ctx context.Context, dealID uint64) (*domain.DealOfOrganization, error)
	DeleteByDealID(ctx context.Context, dealID uint64) error
}
