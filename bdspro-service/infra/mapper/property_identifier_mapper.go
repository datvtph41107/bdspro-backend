// infra/mapper/identifier_mapper.go
package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type PropertyIdentifierMapper struct{}

func NewIdentifierMapper() *PropertyIdentifierMapper {
	return &PropertyIdentifierMapper{}
}

func (m *PropertyIdentifierMapper) ToDTO(identifier *domain.CountryIdentifier) *dto.IdentifierDTO {
	if identifier == nil {
		return nil
	}

	return &dto.IdentifierDTO{
		ID:             identifier.ID,
		LandParcelCode: identifier.LandParcelCode,
		ProjectCode:    identifier.ProjectCode,
		ProvinceID:     identifier.ProvinceID,
		DistrictID:     identifier.DistrictID,
		WardID:         identifier.WardID,
		Latitude:       identifier.Latitude,
		Longitude:      identifier.Longitude,
		Type:           identifier.Type,
		LegalStatus:    identifier.LegalStatus,
		CurrentOwnerID: identifier.CurrentOwnerID,
		CreatedAt:      _utils.FormatTimeToStringCustom(&identifier.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToStringCustom(&identifier.UpdatedAt),
	}
}

func (m *PropertyIdentifierMapper) ToDTOList(identifiers []*domain.CountryIdentifier) []*dto.IdentifierDTO {
	if identifiers == nil {
		return nil
	}

	result := make([]*dto.IdentifierDTO, len(identifiers))
	for i, identifier := range identifiers {
		result[i] = m.ToDTO(identifier)
	}
	return result
}

func (m *PropertyIdentifierMapper) ToDomain(dto *dto.IdentifierDTO) *domain.CountryIdentifier {
	if dto == nil {
		return nil
	}

	identifier := &domain.CountryIdentifier{
		ID:             dto.ID,
		LandParcelCode: dto.LandParcelCode,
		ProjectCode:    dto.ProjectCode,
		ProvinceID:     dto.ProvinceID,
		DistrictID:     dto.DistrictID,
		WardID:         dto.WardID,
		Latitude:       dto.Latitude,
		Longitude:      dto.Longitude,
		Type:           dto.Type,
		LegalStatus:    dto.LegalStatus,
		CurrentOwnerID: dto.CurrentOwnerID,
	}

	// Parse time strings to time.Time
	if createdAt := _utils.ParseStringToTimeCustom(dto.CreatedAt); createdAt != nil {
		identifier.CreatedAt = *createdAt
	}
	if updatedAt := _utils.ParseStringToTimeCustom(dto.UpdatedAt); updatedAt != nil {
		identifier.UpdatedAt = *updatedAt
	}

	return identifier
}

func (m *PropertyIdentifierMapper) FromProtoToDTO(proto *bdspropb.Identifier) *dto.IdentifierDTO {
	if proto == nil {
		return nil
	}

	return &dto.IdentifierDTO{
		ID:             proto.Id,
		LandParcelCode: proto.LandParcelCode,
		ProjectCode:    proto.ProjectCode,
		ProvinceID:     proto.ProvinceId,
		DistrictID:     proto.DistrictId,
		WardID:         proto.WardId,
		Latitude:       proto.Latitude,
		Longitude:      proto.Longitude,
		Type:           proto.Type,
		LegalStatus:    proto.LegalStatus,
		CurrentOwnerID: proto.CurrentOwnerId,
		// CreatedAt và UpdatedAt sẽ được xử lý ở handler layer
	}
}

func (m *PropertyIdentifierMapper) ToProto(dto *dto.IdentifierDTO) *bdspropb.Identifier {
	if dto == nil {
		return nil
	}

	return &bdspropb.Identifier{
		Id:             dto.ID,
		LandParcelCode: dto.LandParcelCode,
		ProjectCode:    dto.ProjectCode,
		ProvinceId:     dto.ProvinceID,
		DistrictId:     dto.DistrictID,
		WardId:         dto.WardID,
		Latitude:       dto.Latitude,
		Longitude:      dto.Longitude,
		Type:           dto.Type,
		LegalStatus:    dto.LegalStatus,
		CurrentOwnerId: dto.CurrentOwnerID,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}

func (m *PropertyIdentifierMapper) ToProtoList(dto []*dto.IdentifierDTO) []*bdspropb.Identifier {
	if dto == nil {
		return nil
	}

	result := make([]*bdspropb.Identifier, len(dto))
	for i, dto := range dto {
		result[i] = m.ToProto(dto)
	}
	return result
}

func (m *PropertyIdentifierMapper) SearchCriteriaFromProto(proto *bdspropb.SearchCriteria) *dto.SearchCriteriaDTO {
	if proto == nil {
		return nil
	}

	return &dto.SearchCriteriaDTO{
		ProvinceID:  proto.ProvinceId,
		DistrictID:  proto.DistrictId,
		WardID:      proto.WardId,
		Type:        proto.Type,
		LegalStatus: proto.LegalStatus,
		OwnerID:     proto.OwnerId,
		Keyword:     proto.Keyword,
		CreatedFrom: proto.CreatedFrom,
		CreatedTo:   proto.CreatedTo,
		RadiusKm:    proto.RadiusKm,
		CenterLat:   proto.CenterLat,
		CenterLng:   proto.CenterLng,
	}
}

func (m *PropertyIdentifierMapper) PaginationFromProto(proto *bdspropb.IdentifierPagination) *dto.PaginationDTO {
	if proto == nil {
		return &dto.PaginationDTO{Page: 1, Limit: 10} // Default values
	}

	return &dto.PaginationDTO{
		Page:  proto.Page,
		Limit: proto.Limit,
	}
}
