package repo

import (
	"bdspro/v2/internal/domain"
	"bdspro/v2/internal/dto"
)

type ProductMediaRepo interface {
	Create(access *domain.ProductMedia) error
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	UpdateProductMediaItems(productID uint64, mediaItems []dto.MediaItem) error
}
