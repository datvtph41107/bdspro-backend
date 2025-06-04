package org_usecase

import (
	"bdspro/internal/common"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	org_repo "bdspro/internal/repo/org"
	_routes "common/routes"
	_utils "common/utils"
	"context"
)

type OrgAssetCostUc struct {
	common.OwnerUsecase[domain.AssetCost, *dto.AssetCostSearchDTO]
}

func NewOrgAssetCostUc(repo org_repo.OrgAssetCostRepo) *OrgAssetCostUc {
	uc := &OrgAssetCostUc{
		OwnerUsecase: common.OwnerUsecase[domain.AssetCost, *dto.AssetCostSearchDTO]{
			OwnerRepo: repo,
		},
	}
	uc.UC = uc
	return uc
}

func (h *OrgAssetCostUc) CheckPermission(c context.Context, id uint64) error {
	organizationId := _utils.GetOrganizationIdFromContext(c)
	assetCost, err := h.OwnerRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if assetCost.OwnerID != organizationId || assetCost.OwnerType != enums.EOwnerTypeOrgnization {
		return &_routes.Except{
			Code:    403,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return nil
}

func (uc *OrgAssetCostUc) Create(c context.Context, entity *domain.AssetCost) error {
	entity.OwnerID = _utils.GetOrganizationIdFromContext(c)
	entity.OwnerType = enums.EOwnerTypeOrgnization
	return uc.OwnerRepo.Create(c, entity)
}
