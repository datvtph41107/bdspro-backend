package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_models "common/models"
	pb_bdspro "pb/types/bdspro"
)

type CostTypeMapper struct {
}

func NewCostTypeMapper() *CostTypeMapper {
	return &CostTypeMapper{}
}

func (m *CostTypeMapper) CostTypeToPb(costType *domain.AssetCostType) *pb_bdspro.AssetCostTypeDTO {
	return &pb_bdspro.AssetCostTypeDTO{
		Id:          costType.ID,
		TypeName:    costType.TypeName,
		IsCustom:    costType.IsCustom,
		Type:        uint32(costType.Type),
		Description: costType.Description,
	}
}

func (m *CostTypeMapper) CostTypeToPbList(costTypes []domain.AssetCostType) []*pb_bdspro.AssetCostTypeDTO {
	pbCostTypes := make([]*pb_bdspro.AssetCostTypeDTO, len(costTypes))
	for i, costType := range costTypes {
		pbCostTypes[i] = m.CostTypeToPb(&costType)
	}
	return pbCostTypes
}

func (m *CostTypeMapper) CostTypePbToDomain(pbCostType *pb_bdspro.AssetCostTypeDTO) *domain.AssetCostType {
	return &domain.AssetCostType{
		BaseEntity: _models.BaseEntity{
			ID: pbCostType.Id,
		},
		TypeName:    pbCostType.TypeName,
		IsCustom:    pbCostType.IsCustom,
		Type:        enums.ECostType(pbCostType.Type),
		Description: pbCostType.Description,
	}
}
