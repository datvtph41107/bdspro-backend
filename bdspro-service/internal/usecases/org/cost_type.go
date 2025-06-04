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

type OrgCostTypeUc struct {
	common.OwnerUsecase[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}

func NewOrgCostTypeUc(repo org_repo.OrgCostTypeRepo) *OrgCostTypeUc {
	uc := &OrgCostTypeUc{
		OwnerUsecase: common.OwnerUsecase[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
			OwnerRepo: repo,
		},
	}
	uc.UC = uc
	return uc
}

func (h *OrgCostTypeUc) CheckPermission(c context.Context, id uint64) error {
	organizationId := _utils.GetOrganizationIdFromContext(c)
	assetCostType, err := h.OwnerRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if assetCostType.OwnerID != organizationId || assetCostType.OwnerType != enums.EOwnerTypeOrgnization {
		return &_routes.Except{
			Code:    403,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return nil
}

func (uc *OrgCostTypeUc) Create(c context.Context, entity *domain.AssetCostType) error {
	entity.OwnerID = _utils.GetOrganizationIdFromContext(c)
	entity.OwnerType = enums.EOwnerTypeOrgnization
	entity.IsCustom = true
	return uc.OwnerRepo.Create(c, entity)
}
