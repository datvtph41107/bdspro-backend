package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type ProductMediaRepo interface {
	Create(access *domain.ProductMedia) error
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	UpdateProductMediaItems(ctx context.Context, productID uint64, mediaItems []dto.MediaItem) ([]uint64, error)
	GetByProductID(ctx context.Context, productID uint64) ([]domain.ProductMedia, error)
}
