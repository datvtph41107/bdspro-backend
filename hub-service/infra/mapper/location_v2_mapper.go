package mapper

import (
	_utils "common/utils"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type LocationV2Mapper struct{}

func NewLocationV2Mapper() *LocationV2Mapper {
	return &LocationV2Mapper{}
}

// ProvinceV2ToProto converts ProvinceV2 entity to protobuf message
func (m *LocationV2Mapper) ProvinceV2ToProto(entity *domain.ProvinceV2) *hubpb.ProvinceV2 {
	if entity == nil {
		return nil
	}

	return &hubpb.ProvinceV2{
		Id:           entity.ID,
		Name:         entity.Name,
		Code:         int32(entity.Code),
		Codename:     entity.Codename,
		DivisionType: entity.DivisionType,
		PhoneCode:    int32(entity.PhoneCode),
		CreatedAt:    _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

// ProvincesV2ToProto converts slice of ProvinceV2 entities to protobuf messages
func (m *LocationV2Mapper) ProvincesV2ToProto(entities []domain.ProvinceV2) []*hubpb.ProvinceV2 {
	if entities == nil {
		return []*hubpb.ProvinceV2{}
	}

	result := make([]*hubpb.ProvinceV2, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.ProvinceV2ToProto(&entity))
	}
	return result
}

// WardV2ToProto converts WardV2 entity to protobuf message
func (m *LocationV2Mapper) WardV2ToProto(entity *domain.WardV2) *hubpb.WardV2 {
	if entity == nil {
		return nil
	}

	return &hubpb.WardV2{
		Id:   entity.ID,
		Name: entity.Name,
		// Code:          int32(entity.Code),
		// Codename:      entity.Codename,
		// DivisionType:  entity.DivisionType,
		// ShortCodename: entity.ShortCodename,
		// ProvinceCode:  int32(entity.ProvinceCode),
		// CreatedAt:     _utils.FormatTimeToString(entity.CreatedAt),
		// UpdatedAt:     _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

// WardsV2ToProto converts slice of WardV2 entities to protobuf messages
func (m *LocationV2Mapper) WardsV2ToProto(entities []domain.WardV2) []*hubpb.WardV2 {
	if entities == nil {
		return []*hubpb.WardV2{}
	}

	result := make([]*hubpb.WardV2, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.WardV2ToProto(&entity))
	}
	return result
}

// ProvinceV2WithWardsToProto converts ProvinceV2 with wards to protobuf message
func (m *LocationV2Mapper) ProvinceV2WithWardsToProto(entity *domain.ProvinceV2) *hubpb.ProvinceV2WithWards {
	if entity == nil {
		return nil
	}

	return &hubpb.ProvinceV2WithWards{
		Province: m.ProvinceV2ToProto(entity),
		Wards:    m.WardsV2ToProto(entity.Wards),
	}
}

// LocationSearchResultV2ToProto converts LocationSearchResultV2 entity to protobuf message
func (m *LocationV2Mapper) LocationSearchResultV2ToProto(entity *domain.LocationSearchResultV2) *hubpb.LocationSearchResultV2 {
	if entity == nil {
		return nil
	}

	return &hubpb.LocationSearchResultV2{
		WardId:       entity.WardID,
		WardName:     entity.WardName,
		WardType:     entity.WardType,
		ProvinceId:   entity.ProvinceID,
		ProvinceName: entity.ProvinceName,
		ProvinceType: entity.ProvinceType,
		FullAddress:  entity.FullAddress,
	}
}

// LocationSearchResultsV2ToProto converts slice of LocationSearchResultV2 entities to protobuf messages
func (m *LocationV2Mapper) LocationSearchResultsV2ToProto(entities []domain.LocationSearchResultV2) []*hubpb.LocationSearchResultV2 {
	if entities == nil {
		return []*hubpb.LocationSearchResultV2{}
	}

	result := make([]*hubpb.LocationSearchResultV2, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.LocationSearchResultV2ToProto(&entity))
	}
	return result
}
