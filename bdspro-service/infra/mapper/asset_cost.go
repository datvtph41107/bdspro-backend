package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_utils "common/utils"
	pb_bdspro "pb/types/bdspro"
)

func AssetCostToPbList(assetCosts []domain.AssetCostQuery) []*pb_bdspro.AssetCostDTO {
	pbAssetCosts := make([]*pb_bdspro.AssetCostDTO, len(assetCosts))
	for i, assetCost := range assetCosts {
		pbAssetCosts[i] = AssetCostToPb(&assetCost)
	}
	return pbAssetCosts
}

func AssetCostToPb(assetCost *domain.AssetCostQuery) *pb_bdspro.AssetCostDTO {
	return &pb_bdspro.AssetCostDTO{
		Id:          assetCost.ID,
		AssetId:     assetCost.AssetID,
		Amount:      assetCost.Amount,
		Date:        _utils.FormatTimeToString(assetCost.Date),
		Description: assetCost.Description,
		// OwnerType:    uint32(assetCost.OwnerType),
		Type:         uint32(assetCost.Type),
		TypeName:     enums.CostTypeMap[assetCost.Type],
		CostTypeId:   assetCost.CostTypeID,
		CostTypeName: assetCost.CostTypeName,
	}
}
