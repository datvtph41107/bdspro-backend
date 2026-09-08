package mapper

import (
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"

	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
)

// ContactLabelMapper handles mapping between domain, DTO and proto
type ContactLabelMapper struct{}

// NewContactLabelMapper creates a new ContactLabelMapper
func NewContactLabelMapper() *ContactLabelMapper {
	return &ContactLabelMapper{}
}

// ToDTO converts domain ContactLabel to ContactLabelDTO
func (m *ContactLabelMapper) ToDTO(contactLabel *domain.ContactLabel) *dto.ContactLabelDTO {
	if contactLabel == nil {
		return nil
	}

	return &dto.ContactLabelDTO{
		ID:           contactLabel.ID,
		Name:         contactLabel.Name,
		Code:         contactLabel.Code,
		Description:  contactLabel.Description,
		Color:        contactLabel.Color,
		Icon:         contactLabel.Icon,
		IsActive:     contactLabel.IsActive,
		SortOrder:    int32(contactLabel.SortOrder),
		ContactCount: int32(contactLabel.ContactCount),
		IsSystem:     contactLabel.IsSystem,
		Status:       int(contactLabel.Status),
		Type:         int(contactLabel.Type),
		CreatedAt:    _utils.FormatTimeToString(contactLabel.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(contactLabel.UpdatedAt),
	}
}

// ToDomain converts ContactLabelDTO to domain ContactLabel
func (m *ContactLabelMapper) ToDomain(dto *dto.ContactLabelDTO) *domain.ContactLabel {
	if dto == nil {
		return nil
	}

	return &domain.ContactLabel{
		BaseEntity: _entity.BaseEntity{
			ID:        dto.ID,
			CreatedAt: _utils.ParseStringToTime(dto.CreatedAt),
			UpdatedAt: _utils.ParseStringToTime(dto.UpdatedAt),
		},
		Name:         dto.Name,
		Code:         dto.Code,
		Description:  dto.Description,
		Color:        dto.Color,
		Icon:         dto.Icon,
		IsActive:     dto.IsActive,
		SortOrder:    int(dto.SortOrder),
		ContactCount: int(dto.ContactCount),
		IsSystem:     dto.IsSystem,
		Status:       enums.ContactLabelStatus(dto.Status),
		Type:         enums.ContactLabelType(dto.Type),
	}
}

// ToDomainFromCreateRequest converts CreateContactLabelRequestDTO to domain ContactLabel
func (m *ContactLabelMapper) ToDomainFromCreateRequest(req *dto.CreateContactLabelRequestDTO) *domain.ContactLabel {
	if req == nil {
		return nil
	}

	return &domain.ContactLabel{
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		Color:        req.Color,
		Icon:         req.Icon,
		IsActive:     req.IsActive,
		SortOrder:    req.SortOrder,
		ContactCount: 0,
		IsSystem:     req.IsSystem,
		Status:       enums.ContactLabelStatusActive,
		Type:         enums.ContactLabelType(req.Type),
	}
}

// ToDomainFromUpdateRequest converts UpdateContactLabelRequestDTO to domain ContactLabel
func (m *ContactLabelMapper) ToDomainFromUpdateRequest(req *dto.UpdateContactLabelRequestDTO, id uint64) *domain.ContactLabel {
	if req == nil {
		return nil
	}

	return &domain.ContactLabel{
		BaseEntity: _entity.BaseEntity{
			ID: id,
		},
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		IsSystem:    req.IsSystem,
		Type:        enums.ContactLabelType(req.Type),
	}
}

// ToDTOList converts slice of domain ContactLabel to slice of ContactLabelDTO
func (m *ContactLabelMapper) ToDTOList(contactLabels []domain.ContactLabel) []dto.ContactLabelDTO {
	if contactLabels == nil {
		return nil
	}
	result := make([]dto.ContactLabelDTO, len(contactLabels))
	for i, contactLabel := range contactLabels {
		dto := m.ToDTO(&contactLabel)
		if dto != nil {
			result[i] = *dto
		}
	}
	return result
}

// ToDomainList converts slice of ContactLabelDTO to slice of domain ContactLabel
func (m *ContactLabelMapper) ToDomainList(dtos []dto.ContactLabelDTO) []*domain.ContactLabel {
	if dtos == nil {
		return nil
	}
	result := make([]*domain.ContactLabel, len(dtos))
	for i, dto := range dtos {
		result[i] = m.ToDomain(&dto)
	}
	return result
}

// ToProto converts ContactLabelDTO to proto ContactLabel
func (m *ContactLabelMapper) ToProto(dto *dto.ContactLabelDTO) *tqdpb.ContactLabel {
	if dto == nil {
		return nil
	}

	return &tqdpb.ContactLabel{
		Id:           dto.ID,
		Name:         dto.Name,
		Code:         dto.Code,
		Description:  dto.Description,
		Color:        dto.Color,
		Icon:         dto.Icon,
		IsActive:     dto.IsActive,
		SortOrder:    dto.SortOrder,
		ContactCount: dto.ContactCount,
		IsSystem:     dto.IsSystem,
		Status:       int32(dto.Status),
		Type:         int32(dto.Type),
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

// ToProtoFromDomain converts domain ContactLabel directly to proto ContactLabel
func (m *ContactLabelMapper) ToProtoFromDomain(contactLabel *domain.ContactLabel) *tqdpb.ContactLabel {
	if contactLabel == nil {
		return nil
	}

	return &tqdpb.ContactLabel{
		Id:           contactLabel.ID,
		Name:         contactLabel.Name,
		Code:         contactLabel.Code,
		Description:  contactLabel.Description,
		Color:        contactLabel.Color,
		Icon:         contactLabel.Icon,
		IsActive:     contactLabel.IsActive,
		SortOrder:    int32(contactLabel.SortOrder),
		ContactCount: int32(contactLabel.ContactCount),
		IsSystem:     contactLabel.IsSystem,
		Status:       int32(contactLabel.Status),
		Type:         int32(contactLabel.Type),
		CreatedAt:    _utils.FormatTimeToString(contactLabel.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(contactLabel.UpdatedAt),
	}
}

// ToProtoList converts slice of ContactLabelDTO to slice of proto ContactLabel
func (m *ContactLabelMapper) ToProtoList(dtos []dto.ContactLabelDTO) []*tqdpb.ContactLabel {
	if dtos == nil {
		return nil
	}
	result := make([]*tqdpb.ContactLabel, len(dtos))
	for i, dto := range dtos {
		result[i] = m.ToProto(&dto)
	}
	return result
}
