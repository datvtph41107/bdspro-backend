package transformer

import (
	_models "common/models"
	"organization/internal/domain/entity"
	"organization/internal/enums"
	pb "pb/types/organization"
)

type CostTypeMapper struct{}

func NewCostTypeMapper() *CostTypeMapper {
	return &CostTypeMapper{}
}

func (m *CostTypeMapper) CostTypeModelToProto(model *entity.CostType) *pb.CostTypeDTO {
	return &pb.CostTypeDTO{
		Id:          model.ID,
		TypeName:    model.TypeName,
		IsCustom:    model.IsCustom,
		CostType:    uint32(model.CostType),
		Description: model.Description,
	}
}

func (m *CostTypeMapper) CostTypeProtoToModel(proto *pb.CostTypeDTO) *entity.CostType {
	return &entity.CostType{
		BaseEntity:  _models.BaseEntity{ID: proto.Id},
		TypeName:    proto.TypeName,
		IsCustom:    proto.IsCustom,
		CostType:    enums.ECostType(proto.CostType),
		Description: proto.Description,
	}
}
