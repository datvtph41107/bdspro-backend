package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_models "common/models"
	bdspropb "pb/types/bdspro"
)

type CostTypeMapper struct {
}

func NewCostTypeMapper() *CostTypeMapper {
	return &CostTypeMapper{}
}

func (m *CostTypeMapper) CostTypeToPb(costType *domain.AssetCostType) *bdspropb.AssetCostTypeDTO {
	return &bdspropb.AssetCostTypeDTO{
		Id:          costType.ID,
		TypeName:    costType.TypeName,
		IsCustom:    costType.IsCustom,
		Type:        uint32(costType.Type),
		Description: costType.Description,
	}
}

func (m *CostTypeMapper) CostTypeToPbList(costTypes []domain.AssetCostType) []*bdspropb.AssetCostTypeDTO {
	pbCostTypes := make([]*bdspropb.AssetCostTypeDTO, len(costTypes))
	for i, costType := range costTypes {
		pbCostTypes[i] = m.CostTypeToPb(&costType)
	}
	return pbCostTypes
}

func (m *CostTypeMapper) CostTypePbToDomain(pbCostType *bdspropb.AssetCostTypeDTO) *domain.AssetCostType {
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
