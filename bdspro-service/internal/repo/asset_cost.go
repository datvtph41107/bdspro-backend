package repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
)

type AssetCostRepo interface {
	common.IOwnerRepo[domain.AssetCost, *dto.AssetCostSearchDTO]
	Find(c context.Context, profileId uint64, ownerType enums.EOwnerType, dto *dto.AssetCostSearchDTO) ([]domain.AssetCostQuery, int64, error)
}
