package repo

import (
	"bdspro/internal/domain"
	"context"
)

type ProductUserRepo interface {
	CreateOwner(ctx context.Context, productID uint64, profileID uint64) (*domain.ProductUser, error)
	GetByProductID(ctx context.Context, productID uint64) ([]*domain.ProductUser, error)
	GetByProfileID(ctx context.Context, profileID uint64) ([]*domain.ProductUser, error)
	GetByProductAndProfile(ctx context.Context, productID, profileID uint64) (*domain.ProductUser, error)
	GetByDistributeID(ctx context.Context, distributeID uint64) ([]*domain.ProductUser, error)
	Update(ctx context.Context, productOfUser *domain.ProductUser) (*domain.ProductUser, error)
	Delete(ctx context.Context, id uint64) error
	DeleteByProductID(ctx context.Context, productID uint64) error

	UpdatePriceID(ctx context.Context, productID uint64, profileID uint64, priceID uint64) error
	GetByProductAndOriginProfile(ctx context.Context, productID, originProfileId uint64) (*domain.ProductUser, error)
	GetByOriginProfile(ctx context.Context, originProfileId uint64) ([]*domain.ProductUser, error)
	UpdateArchivedStatus(ctx context.Context, productID uint64, originProfileId uint64, archived bool) error
	GetProfileIDsByProduct(ctx context.Context, productID uint64) ([]uint64, error)
	GetUserProductsWithTimestamp(ctx context.Context, profileID uint64) (map[uint64]int64, error)
}
