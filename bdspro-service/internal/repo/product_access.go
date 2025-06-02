package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type ProductAccessRepo interface {
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	GetByUserAndProduct(userId, productId uint64) (*domain.ProductAccess, error)
	GetByTargetAndProduct(userId, target, productId uint64) (*domain.ProductAccess, error)
	GetAll() ([]domain.ProductAccess, error)
	OwnerSearchProduct(c context.Context, id uint64, dto *dto.ProductAccessSearch) (*[]domain.ProductAccessQuery, error)
	UpdateCommission(c context.Context,
		profileId uint64,
		dto dto.ProductAccessUpdate,
	) error
	GetByID(id uint64) (*domain.ProductAccess, error)
	Delete(c context.Context, id uint64) error
	// GetProductAccessByOrgID(ctx context.Context, orgID int64, dto dto.ProductAccessSearch) ([]*domain.ProductAccess, error)
	// BulkSave(c context.Context, id uint64, dto dto.ProductAccessBulkRequest) (*[]domain.ProductAccess, error)
}
