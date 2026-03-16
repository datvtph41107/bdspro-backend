package mapper

import (
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"

	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
)

// ContactLabelMapper handles mapping between domain and DTO
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
		SortOrder:    contactLabel.SortOrder,
		ContactCount: contactLabel.ContactCount,
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
		SortOrder:    dto.SortOrder,
		ContactCount: dto.ContactCount,
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
		IsSystem:     req.IsSystem,
		Status:       enums.ContactLabelStatusActive,
		Type:         enums.ContactLabelType(req.Type),
		ContactCount: 0,
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
	result := make([]dto.ContactLabelDTO, len(contactLabels))
	for i, contactLabel := range contactLabels {
		result[i] = *m.ToDTO(&contactLabel)
	}
	return result
}

// ToProto converts domain ContactLabel to proto ContactLabel
func (m *ContactLabelMapper) ToProto(contactLabel *dto.ContactLabelDTO) *tqdpb.ContactLabel {
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
		SortOrder:    contactLabel.SortOrder,
		ContactCount: int32(contactLabel.ContactCount),
		IsSystem:     contactLabel.IsSystem,
		Status:       int32(contactLabel.Status),
		Type:         int32(contactLabel.Type),
		CreatedAt:    contactLabel.CreatedAt,
		UpdatedAt:    contactLabel.UpdatedAt,
	}
}
