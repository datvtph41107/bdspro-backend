package mapper

import (
	"bdspro/internal/domain"
	_models "common/models"
	_utils "common/utils"

	bdspropb "pb/types/bdspro"
)

type PropertyTypeMapper struct {
}

func NewPropertyTypeMapper() *PropertyTypeMapper {
	return &PropertyTypeMapper{}
}

func (m *PropertyTypeMapper) MapPropertyTypePbList(propertyTypes []domain.PropertyType) []*bdspropb.PropertyType {
	propertyTypesPb := make([]*bdspropb.PropertyType, len(propertyTypes))
	for i, propertyType := range propertyTypes {
		propertyTypesPb[i] = m.MapPropertyTypePb(&propertyType)
	}
	return propertyTypesPb
}

func (m *PropertyTypeMapper) MapPropertyTypePb(propertyType *domain.PropertyType) *bdspropb.PropertyType {
	return &bdspropb.PropertyType{
		Id:        propertyType.ID,
		Name:      propertyType.Name,
		CreatedAt: _utils.FormatTimeToString(propertyType.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(propertyType.UpdatedAt),
		Active:    propertyType.Active,
	}
}

func (m *PropertyTypeMapper) PropertyTypePbToDomain(propertyType *bdspropb.PropertyType) *domain.PropertyType {
	return &domain.PropertyType{
		BaseEntity: _models.BaseEntity{ID: propertyType.Id},
		Name:       propertyType.Name,
		Active:     propertyType.Active,
	}
}
