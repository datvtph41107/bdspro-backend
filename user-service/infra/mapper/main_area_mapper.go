package mapper

import (
	_models "common/models"
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/internal/models"
)

type MainAreaMapper struct{}

func NewMainAreaMapper() *MainAreaMapper {
	return &MainAreaMapper{}
}

func (m *MainAreaMapper) ProtoToEntity(req *userpb.MainArea) *models.MainAreaEntity {
	return &models.MainAreaEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.GetId(),
		},
		Name:          req.GetName(),
		Description:   req.GetDescription(),
		NumberProfile: req.GetNumberProfile(),
	}
}

func (m *MainAreaMapper) EntityToProto(entity *models.MainAreaEntity) *userpb.MainArea {
	return &userpb.MainArea{
		Id:            entity.ID,
		Name:          entity.Name,
		NumberProfile: entity.NumberProfile,
		CreatedAt:     _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:     _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

func (m *MainAreaMapper) EntitiesToProtos(entities []models.MainAreaEntity) []*userpb.MainArea {
	result := make([]*userpb.MainArea, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.EntityToProto(&entity))
	}
	return result
}

func (m *MainAreaMapper) EntitiesToItemProtos(entities []models.MainAreaEntity) []*sharepb.ItemV3Proto {
	result := make([]*sharepb.ItemV3Proto, 0, len(entities))
	for _, entity := range entities {
		result = append(result, &sharepb.ItemV3Proto{
			Id:   entity.ID,
			Name: entity.Name,
		})
	}
	return result
}
