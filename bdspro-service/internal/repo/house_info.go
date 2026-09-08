package repo

import (
	"bdspro/internal/domain"
	"context"
)

type HouseInfoRepo interface {
	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
	Create(access *domain.HouseInfo) error
	UpdateHouseInfo(ctx context.Context, houseInfo *domain.HouseInfo) error
}
