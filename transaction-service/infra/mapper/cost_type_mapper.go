package mapper

import (
	_models "common/models"
	transactionpb "pb/types/transaction"
	"tx/internal/domain/entity"
	"tx/internal/enums"
)

type CostTypeMapper struct{}

func NewCostTypeMapper() *CostTypeMapper {
	return &CostTypeMapper{}
}

func (m *CostTypeMapper) CostTypeModelToProto(model *entity.CostType) *transactionpb.CostTypeDTO {
	return &transactionpb.CostTypeDTO{
		Id:          model.ID,
		TypeName:    model.TypeName,
		IsCustom:    model.IsCustom,
		CostType:    uint32(model.CostType),
		Description: model.Description,
	}
}

func (m *CostTypeMapper) CostTypeProtoToModel(proto *transactionpb.CostTypeDTO) *entity.CostType {
	return &entity.CostType{
		BaseEntity:  _models.BaseEntity{ID: proto.Id},
		TypeName:    proto.TypeName,
		IsCustom:    proto.IsCustom,
		CostType:    enums.ECostType(proto.CostType),
		Description: proto.Description,
	}
}
