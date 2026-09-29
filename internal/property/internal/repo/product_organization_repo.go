package repo

import (
	"bdspro/internal/domain"
	"context"
)

type ProductOrganizationRepo interface {
	CreateOwner(ctx context.Context, productID uint64, organizationID uint64) (*domain.ProductOrganization, error)
	GetByProductID(ctx context.Context, productID uint64) ([]*domain.ProductOrganization, error)
	GetByOrganizationID(ctx context.Context, organizationID uint64) ([]*domain.ProductOrganization, error)
	GetByProductAndOrganization(ctx context.Context, productID, organizationID uint64) (*domain.ProductOrganization, error)
	Update(ctx context.Context, productOfOrg *domain.ProductOrganization) (*domain.ProductOrganization, error)
	Delete(ctx context.Context, id uint64) error
	DeleteByProductID(ctx context.Context, productID uint64) error
}
