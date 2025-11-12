package mapper

import (
	_models "common/models"
	_utils "common/utils"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type ApiKeyMapper struct{}

func NewApiKeyMapper() *ApiKeyMapper {
	return &ApiKeyMapper{}
}

func (m *ApiKeyMapper) ProtoToEntity(req *hubpb.ApiKeyProto) *domain.ApiKeyEntity {
	if req == nil {
		return nil
	}

	entity := &domain.ApiKeyEntity{
		BaseEntity: _models.BaseEntity{
			ID: req.GetId(),
		},
		Name:        req.GetName(),
		AppName:     req.GetAppName(),
		Description: req.GetDescription(),
		ApiKey:      req.GetApiKey(),
		ExpiredAt:   _utils.ParseStringToTime(req.GetExpiredAt()),
	}
	return entity
}

func (m *ApiKeyMapper) EntityToProto(entity *domain.ApiKeyEntity) *hubpb.ApiKeyProto {
	if entity == nil {
		return nil
	}

	return &hubpb.ApiKeyProto{
		Id:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		ApiKey:      entity.ApiKey,
		AppName:     entity.AppName,
		ExpiredAt:   _utils.FormatTimeToString(entity.ExpiredAt),
		CreatedAt:   _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

func (m *ApiKeyMapper) EntitiesToProtos(entities []domain.ApiKeyEntity) []*hubpb.ApiKeyProto {
	if len(entities) == 0 {
		return []*hubpb.ApiKeyProto{}
	}

	result := make([]*hubpb.ApiKeyProto, 0, len(entities))
	for i := range entities {
		e := entities[i]
		result = append(result, m.EntityToProto(&e))
	}
	return result
}
