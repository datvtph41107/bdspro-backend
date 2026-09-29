package mapper

import (
	_utils "common/utils"
	"hub/internal/domain"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type LocationMapper struct{}

func NewLocationMapper() *LocationMapper {
	return &LocationMapper{}
}

// ProvinceToProto converts Province entity to protobuf message
func (m *LocationMapper) ProvinceToProto(entity *domain.Province) *hubpb.Province {
	if entity == nil {
		return nil
	}

	return &hubpb.Province{
		Id:        entity.ID,
		Name:      entity.Name,
		Type:      int32(entity.Type),
		TypeText:  entity.TypeText,
		Slug:      entity.Slug,
		CreatedAt: _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

// ProvincesToProto converts slice of Province entities to protobuf messages
func (m *LocationMapper) ProvincesToProto(entities []domain.Province) []*hubpb.Province {
	if entities == nil {
		return []*hubpb.Province{}
	}

	result := make([]*hubpb.Province, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.ProvinceToProto(&entity))
	}
	return result
}

// DistrictToProto converts District entity to protobuf message
func (m *LocationMapper) DistrictToProto(entity *domain.District) *hubpb.District {
	if entity == nil {
		return nil
	}

	return &hubpb.District{
		Id:         entity.ID,
		Name:       entity.Name,
		ProvinceId: entity.ProvinceID,
		Type:       int32(entity.Type),
		TypeText:   entity.TypeText,
		CreatedAt:  _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:  _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

// DistrictsToProto converts slice of District entities to protobuf messages
func (m *LocationMapper) DistrictsToProto(entities []domain.District) []*hubpb.District {
	if entities == nil {
		return []*hubpb.District{}
	}

	result := make([]*hubpb.District, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.DistrictToProto(&entity))
	}
	return result
}

// LocationSearchResultToProto converts LocationSearchResult to protobuf message
func (m *LocationMapper) LocationSearchResultToProto(entity *domain.LocationSearchResult) *hubpb.LocationSearchResult {
	if entity == nil {
		return nil
	}

	return &hubpb.LocationSearchResult{
		DistrictId:   entity.DistrictID,
		DistrictName: entity.DistrictName,
		DistrictType: entity.DistrictType,
		ProvinceId:   entity.ProvinceID,
		ProvinceName: entity.ProvinceName,
		ProvinceType: entity.ProvinceType,
		FullAddress:  entity.FullAddress,
	}
}

// LocationSearchResultsToProto converts slice of LocationSearchResult to protobuf messages
func (m *LocationMapper) LocationSearchResultsToProto(entities []domain.LocationSearchResult) []*hubpb.LocationSearchResult {
	if entities == nil {
		return []*hubpb.LocationSearchResult{}
	}

	result := make([]*hubpb.LocationSearchResult, 0, len(entities))
	for _, entity := range entities {
		result = append(result, m.LocationSearchResultToProto(&entity))
	}
	return result
}

// LocationInfoToProto converts LocationInfo to protobuf Location message
func (m *LocationMapper) LocationInfoToProto(info *_usecase.LocationInfo) *hubpb.Location {
	if info == nil {
		return nil
	}

	location := &hubpb.Location{
		Id:       info.ID,
		Name:     info.Name,
		Type:     int32(info.Type),
		TypeText: info.TypeText,
		Level:    int32(info.Level),
	}

	if info.ParentID != "" {
		location.ParentId = &info.ParentID
	}

	return location
}

// LocationInfosToProto converts slice of LocationInfo to protobuf Location messages
func (m *LocationMapper) LocationInfosToProto(infos []_usecase.LocationInfo) []*hubpb.Location {
	if infos == nil {
		return []*hubpb.Location{}
	}

	result := make([]*hubpb.Location, 0, len(infos))
	for _, info := range infos {
		result = append(result, m.LocationInfoToProto(&info))
	}
	return result
}
