package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type AssetCostMapper struct {
}

func NewAssetCostMapper() *AssetCostMapper {
	return &AssetCostMapper{}
}

func (m *AssetCostMapper) AssetCostToPbList(assetCosts []domain.AssetCost) []*bdspropb.AssetCostDTO {
	pbAssetCosts := make([]*bdspropb.AssetCostDTO, len(assetCosts))
	for i, assetCost := range assetCosts {
		pbAssetCosts[i] = m.AssetCostToPb(&assetCost)
	}
	return pbAssetCosts
}

func (m *AssetCostMapper) AssetCostToPb(assetCost *domain.AssetCost) *bdspropb.AssetCostDTO {
	result := &bdspropb.AssetCostDTO{
		Id:          assetCost.ID,
		AssetId:     assetCost.AssetID,
		Amount:      assetCost.Amount,
		Date:        _utils.FormatTimeToString(assetCost.Date),
		Description: assetCost.Description,
		OwnerType:   uint32(assetCost.OwnerType),
		Type:        uint32(assetCost.Type),
		TypeName:    enums.CostTypeMap[assetCost.Type],
		CostTypeId:  assetCost.CostTypeID,
	}

	if assetCost.CostType != nil {
		result.CostTypeName = assetCost.CostType.TypeName
	}

	return result
}

func (m *AssetCostMapper) AssetCostPbToDomain(pbAssetCost *bdspropb.AssetCostDTO) *domain.AssetCost {
	res := &domain.AssetCost{
		BaseEntity: _models.BaseEntity{
			ID: pbAssetCost.Id,
		},
		AssetID:     pbAssetCost.AssetId,
		Amount:      pbAssetCost.Amount,
		Date:        _utils.ParseStringToTime(pbAssetCost.Date),
		Description: pbAssetCost.Description,
		OwnerType:   enums.EOwnerOf(pbAssetCost.OwnerType),
		Type:        enums.ECostType(pbAssetCost.Type),
		CostTypeID:  pbAssetCost.CostTypeId,
	}

	if pbAssetCost.OwnerId != nil {
		res.OwnerID = *pbAssetCost.OwnerId
	}

	return res
}
