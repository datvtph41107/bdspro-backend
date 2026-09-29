package repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type AssetCostTypeRepo interface {
	common.IOwnerRepo[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}
