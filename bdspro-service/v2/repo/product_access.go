package repo

import (
	"bdspro/v2/internal/domain"
)

type ProductAccessRepo interface {
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	GetByUserAndProduct(userId, productId uint64) (*domain.ProductAccess, error)
	GetAll() ([]domain.ProductAccess, error)
}
