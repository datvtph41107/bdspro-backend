package org_repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type OrgCostTypeRepo interface {
	common.IOwnerRepo[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}
