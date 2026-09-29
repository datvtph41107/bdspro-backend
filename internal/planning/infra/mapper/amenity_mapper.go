package mapper

import (
	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

type AmenityMapper struct{}

func NewAmenityMapper() *AmenityMapper {
	return &AmenityMapper{}
}

func (m *AmenityMapper) ToDTO(amenity *domain.Amenity) *dto.AmenityDTO {
	if amenity == nil {
		return nil
	}

	return &dto.AmenityDTO{
		ID:          amenity.ID,
		Name:        amenity.Name,
		Description: amenity.Description,
		Icon:        amenity.Icon,
		Category:    amenity.Category,
		SortOrder:   amenity.SortOrder,
		IsActive:    amenity.IsActive,
		CreatedAt:   _utils.FormatTimeToString(amenity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(amenity.UpdatedAt),
	}
}

func (m *AmenityMapper) ToDomain(dto *dto.AmenityDTO) *domain.Amenity {
	if dto == nil {
		return nil
	}

	return &domain.Amenity{
		BaseEntity: _entity.BaseEntity{
			ID:        dto.ID,
			CreatedAt: _utils.ParseStringToTime(dto.CreatedAt),
			UpdatedAt: _utils.ParseStringToTime(dto.UpdatedAt),
		},
		Name:        dto.Name,
		Description: dto.Description,
		Icon:        dto.Icon,
		Category:    dto.Category,
		SortOrder:   dto.SortOrder,
		IsActive:    dto.IsActive,
	}
}

func (m *AmenityMapper) ToDomainFromCreateRequest(req *dto.CreateAmenityRequestDTO) *domain.Amenity {
	if req == nil {
		return nil
	}

	return &domain.Amenity{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Category:    req.Category,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}
}

func (m *AmenityMapper) ToDomainFromUpdateRequest(req *dto.UpdateAmenityRequestDTO, id uint64) *domain.Amenity {
	if req == nil {
		return nil
	}

	return &domain.Amenity{
		BaseEntity: _entity.BaseEntity{
			ID: id,
		},
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Category:    req.Category,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}
}

// ToDTOList converts slice of domain Amenity to slice of AmenityDTO
func (m *AmenityMapper) ToDTOList(amenities []*domain.Amenity) []*dto.AmenityDTO {
	if amenities == nil {
		return nil
	}
	result := make([]*dto.AmenityDTO, len(amenities))
	for i, amenity := range amenities {
		result[i] = m.ToDTO(amenity)
	}
	return result
}

// ToDomainList converts slice of AmenityDTO to slice of domain Amenity
func (m *AmenityMapper) ToDomainList(dtos []*dto.AmenityDTO) []*domain.Amenity {
	if dtos == nil {
		return nil
	}
	result := make([]*domain.Amenity, len(dtos))
	for i, dto := range dtos {
		result[i] = m.ToDomain(dto)
	}
	return result
}

// ToProto converts domain Amenity to proto Amenity
func (m *AmenityMapper) ToProto(amenity *domain.Amenity) *tqdpb.Amenity {
	if amenity == nil {
		return nil
	}

	return &tqdpb.Amenity{
		Id:          amenity.ID,
		Name:        amenity.Name,
		Description: amenity.Description,
		Icon:        amenity.Icon,
		// Category:    amenity.Category, // Add to proto if needed
		// SortOrder:   int32(amenity.SortOrder), // Add to proto if needed
		// IsActive:    amenity.IsActive, // Add to proto if needed
		CreatedAt: _utils.FormatTimeToString(amenity.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(amenity.UpdatedAt),
	}
}

// ToProtoList converts slice of domain Amenity to slice of proto Amenity
func (m *AmenityMapper) ToProtoList(amenities []*domain.Amenity) []*tqdpb.Amenity {
	if amenities == nil {
		return nil
	}
	result := make([]*tqdpb.Amenity, len(amenities))
	for i, amenity := range amenities {
		result[i] = m.ToProto(amenity)
	}
	return result
}
