package mapper

import (
	_models "common/domain/entity"
	_utils "common/utils"
	"mdn4/modules/user/domain/user"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

type PurposeUseMapper struct{}

func NewPurposeUseMapper() *PurposeUseMapper {
	return &PurposeUseMapper{}
}

func (m *PurposeUseMapper) ProtoToEntity(req *userpb.PurposeUse) *user.PurposeUseEntity {
	return &user.PurposeUseEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.GetId(),
		},
		Name:          req.GetName(),
		Description:   req.GetDescription(),
		NumberProfile: req.GetNumberProfile(),
	}
}

func (m *PurposeUseMapper) EntityToProto(entity *user.PurposeUseEntity) *userpb.PurposeUse {
	if entity == nil {
		return nil
	}
	result := &userpb.PurposeUse{
		Id:            entity.ID,
		Name:          entity.Name,
		NumberProfile: entity.NumberProfile,
		CreatedAt:     _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:     _utils.FormatTimeToString(entity.UpdatedAt),
	}
	if entity.Description != "" {
		result.Description = &entity.Description
	}
	return result
}

func (m *PurposeUseMapper) EntitiesToProtos(entities []user.PurposeUseEntity) []*userpb.PurposeUse {
	result := make([]*userpb.PurposeUse, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.EntityToProto(&entity))
	}
	return result
}

func (m *PurposeUseMapper) EntitiesToItemProtos(entities []user.PurposeUseEntity) []*sharepb.ItemV3Proto {
	result := make([]*sharepb.ItemV3Proto, 0, len(entities))
	for _, entity := range entities {
		result = append(result, &sharepb.ItemV3Proto{
			Id:   entity.ID,
			Name: entity.Name,
		})
	}
	return result
}
