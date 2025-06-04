package user_repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type UserCostTypeRepo interface {
	common.IOwnerRepo[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}
