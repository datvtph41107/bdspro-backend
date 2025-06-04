package user_repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type UserAssetCostRepo interface {
	common.IOwnerRepo[domain.AssetCost, *dto.AssetCostSearchDTO]
}
