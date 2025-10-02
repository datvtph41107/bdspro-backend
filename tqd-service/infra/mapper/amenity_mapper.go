package mapper

import (
	_entity "common/domain/entity"
	_utils "common/utils"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

type AmenityMapper struct{}

func NewAmenityMapper() *AmenityMapper {
	return &AmenityMapper{}
}

// ToDTO converts domain Amenity to AmenityDTO
func (m *AmenityMapper) ToDTO(amenity *domain.Amenity) *dto.AmenityDTO {
	if amenity == nil {
		return nil
	}

	return &dto.AmenityDTO{
		ID:          amenity.ID,
		Name:        amenity.Name,
		Description: amenity.Description,
		Icon:        amenity.Icon,
		CreatedAt:   _utils.FormatTimeToString(amenity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(amenity.UpdatedAt),
	}
}

// ToDomain converts AmenityDTO to domain Amenity
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
	}
}

// ToDomainFromCreateRequest converts CreateAmenityRequestDTO to domain Amenity
func (m *AmenityMapper) ToDomainFromCreateRequest(req *dto.CreateAmenityRequestDTO) *domain.Amenity {
	if req == nil {
		return nil
	}

	return &domain.Amenity{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}
}

// ToDomainFromUpdateRequest converts UpdateAmenityRequestDTO to domain Amenity
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
	}
}

// ToDTOList converts slice of domain Amenity to slice of AmenityDTO
func (m *AmenityMapper) ToDTOList(amenities []*domain.Amenity) []*dto.AmenityDTO {
	result := make([]*dto.AmenityDTO, len(amenities))
	for i, amenity := range amenities {
		result[i] = m.ToDTO(amenity)
	}
	return result
}
