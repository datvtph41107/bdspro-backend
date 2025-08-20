package mapper

import (
	"bdspro/internal/domain"
	_utils "common/utils"

	pb_bdspro "pb/types/bdspro"
)

type PropertyTypeMapper struct {
}

func NewPropertyTypeMapper() *PropertyTypeMapper {
	return &PropertyTypeMapper{}
}

func (m *PropertyTypeMapper) MapPropertyTypePbList(propertyTypes []domain.PropertyType) []*pb_bdspro.PropertyType {
	propertyTypesPb := make([]*pb_bdspro.PropertyType, len(propertyTypes))
	for i, propertyType := range propertyTypes {
		propertyTypesPb[i] = m.MapPropertyTypePb(&propertyType)
	}
	return propertyTypesPb
}

func (m *PropertyTypeMapper) MapPropertyTypePb(propertyType *domain.PropertyType) *pb_bdspro.PropertyType {
	return &pb_bdspro.PropertyType{
		Id:        propertyType.ID,
		Name:      propertyType.Name,
		CreatedAt: _utils.FormatTimeToString(propertyType.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(propertyType.UpdatedAt),
		Active:    propertyType.Active,
	}
}

func (m *PropertyTypeMapper) PropertyTypePbToDomain(propertyType *pb_bdspro.PropertyType) *domain.PropertyType {
	return &domain.PropertyType{
		ID:     propertyType.Id,
		Name:   propertyType.Name,
		Active: propertyType.Active,
	}
}
