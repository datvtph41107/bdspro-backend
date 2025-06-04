package org_repo

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
)

type OrgAssetCostRepo interface {
	common.IOwnerRepo[domain.AssetCost, *dto.AssetCostSearchDTO]
}
