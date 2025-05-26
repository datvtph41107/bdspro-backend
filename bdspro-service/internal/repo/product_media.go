package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type ProductMediaRepo interface {
	Create(access *domain.ProductMedia) error
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	UpdateProductMediaItems(productID uint64, mediaItems []dto.MediaItem) error
}
