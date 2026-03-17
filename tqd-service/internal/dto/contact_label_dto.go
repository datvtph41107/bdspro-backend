package dto

import (
	_utils "common/utils"
	"tqd/internal/domain"
	"tqd/internal/enums"

	_dto "common/domain/dto"
	_entity "common/domain/entity"
)

// ContactLabelDTO represents the data transfer object for contact label
type ContactLabelDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	Icon         string `json:"icon"`
	IsActive     bool   `json:"isActive"`
	SortOrder    int32  `json:"sortOrder"`
	ContactCount int32  `json:"contactCount"`
	IsSystem     bool   `json:"isSystem"`
	Status       int    `json:"status"`
	Type         int    `json:"type"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateContactLabelRequestDTO represents the request DTO for creating contact label
type CreateContactLabelRequestDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Code        string `json:"code" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"max=7"`
	Icon        string `json:"icon" validate:"max=100"`
	IsActive    bool   `json:"isActive"`
	SortOrder   int    `json:"sortOrder"`
	IsSystem    bool   `json:"isSystem"`
	Type        int    `json:"type"`
}

// UpdateContactLabelRequestDTO represents the request DTO for updating contact label
type UpdateContactLabelRequestDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Code        string `json:"code" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"max=7"`
	Icon        string `json:"icon" validate:"max=100"`
	IsActive    bool   `json:"isActive"`
	SortOrder   int    `json:"sortOrder"`
	IsSystem    bool   `json:"isSystem"`
	Type        int    `json:"type"`
}

// ListContactLabelsRequestDTO represents the request DTO for listing contact labels
type ListContactLabelsRequestDTO struct {
	Pagable  _dto.Pagable
	Search   string `json:"search"`
	IsActive *bool  `json:"isActive"`
	IsSystem *bool  `json:"isSystem"`
}

// ListContactLabelsResponseDTO represents the response DTO for listing contact labels
type ListContactLabelsResponseDTO struct {
	Data  []ContactLabelDTO `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

// ToDomain converts ContactLabelDTO to domain ContactLabel
func (d *ContactLabelDTO) ToDomain() *domain.ContactLabel {
	return &domain.ContactLabel{
		BaseEntity: _entity.BaseEntity{
			ID:        d.ID,
			CreatedAt: _utils.ParseStringToTime(d.CreatedAt),
			UpdatedAt: _utils.ParseStringToTime(d.UpdatedAt),
		},
		Name:         d.Name,
		Code:         d.Code,
		Description:  d.Description,
		Color:        d.Color,
		Icon:         d.Icon,
		IsActive:     d.IsActive,
		SortOrder:    int(d.SortOrder),
		ContactCount: int(d.ContactCount),
		IsSystem:     d.IsSystem,
		Status:       enums.ContactLabelStatus(d.Status),
		Type:         enums.ContactLabelType(d.Type),
	}
}

// FromDomain converts domain ContactLabel to ContactLabelDTO
func (d *ContactLabelDTO) FromDomain(contactLabel *domain.ContactLabel) {
	if contactLabel == nil {
		return
	}
	d.ID = contactLabel.ID
	d.Name = contactLabel.Name
	d.Code = contactLabel.Code
	d.Description = contactLabel.Description
	d.Color = contactLabel.Color
	d.Icon = contactLabel.Icon
	d.IsActive = contactLabel.IsActive
	d.SortOrder = int32(contactLabel.SortOrder)
	d.ContactCount = int32(contactLabel.ContactCount)
	d.IsSystem = contactLabel.IsSystem
	d.Status = int(contactLabel.Status)
	d.Type = int(contactLabel.Type)
	d.CreatedAt = _utils.FormatTimeToString(contactLabel.CreatedAt)
	d.UpdatedAt = _utils.FormatTimeToString(contactLabel.UpdatedAt)
}
