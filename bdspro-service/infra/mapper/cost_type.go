package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

func CostTypeToPb(costType *domain.AssetCostType) *pb_bdspro.AssetCostTypeDTO {
	return &pb_bdspro.AssetCostTypeDTO{
		Id:       costType.ID,
		TypeName: costType.TypeName,
		IsCustom: costType.IsCustom,
		Type:     uint32(costType.Type),
	}
}

func CostTypeToPbList(costTypes []domain.AssetCostType) []*pb_bdspro.AssetCostTypeDTO {
	pbCostTypes := make([]*pb_bdspro.AssetCostTypeDTO, len(costTypes))
	for i, costType := range costTypes {
		pbCostTypes[i] = CostTypeToPb(&costType)
	}
	return pbCostTypes
}
