package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
)

type SharingAccessRepo interface {
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	GetByUserAndProduct(userId, productId uint64) (*domain.SharingAccess, error)
	GetByTargetAndProduct(userId, target, productId uint64) (*domain.SharingAccess, error)
	GetAll() ([]domain.SharingAccess, error)
	OwnerSearchProduct(c context.Context, productId uint64, fromType enums.EOwnerOf, dto *dto.SharingAccessSearch) (*[]domain.SharingAccess, error)
	OwnerSearchAsset(c context.Context, assetId uint64, fromType enums.EOwnerOf, dto *dto.SharingAccessSearch) (*[]domain.SharingAccess, error)
	UpdateCommission(c context.Context,
		profileId uint64,
		dto dto.CommissionUpdate,
	) error
	GetProductAccessByOrgID(ctx context.Context, orgID int64, dto dto.SharingAccessSearch) ([]*domain.SharingAccess, error)
	BulkSave(c context.Context, id uint64, dto dto.SharingAccessBulk) (*[]domain.SharingAccess, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.SharingAccess, error)

	// Đếm số sharing access theo productId
	CountSharingAccessByProductID(ctx context.Context, productID uint64) (uint32, error)
}
