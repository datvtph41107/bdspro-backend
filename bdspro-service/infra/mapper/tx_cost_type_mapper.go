package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	_models "common/models"
	bdspropb "pb/types/bdspro"
)

type TxCostTypeMapper struct{}

func NewTxCostTypeMapper() *TxCostTypeMapper {
	return &TxCostTypeMapper{}
}

func (m *TxCostTypeMapper) TxCostTypeModelToProto(model *domain.TxCostType) *bdspropb.CostTypeDTO {
	return &bdspropb.CostTypeDTO{
		Id:          model.ID,
		TypeName:    model.TypeName,
		IsCustom:    model.IsCustom,
		CostType:    uint32(model.CostType),
		Description: model.Description,
	}
}

func (m *TxCostTypeMapper) TxCostTypeProtoToModel(proto *bdspropb.CostTypeDTO) *domain.TxCostType {
	return &domain.TxCostType{
		BaseEntity:  _models.BaseEntity{ID: proto.Id},
		TypeName:    proto.TypeName,
		IsCustom:    proto.IsCustom,
		CostType:    enums.TxCostType(proto.CostType),
		Description: proto.Description,
	}
}
