package dto

import (
	"encoding/json"
	"time"
	"tqd/internal/domain"

	_dto "common/domain/dto"
	_entity "common/domain/entity"
)

// DirectorySourceDTO represents the data transfer object for directory source
type DirectorySourceDTO struct {
	_entity.BaseEntity
	ID             uint64     `json:"id"`
	Name           string     `json:"name"`
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	CategoryID     *uint64    `json:"categoryId"`
	Type           string     `json:"type"`
	Icon           string     `json:"icon"`
	Color          string     `json:"color"`
	IsActive       bool       `json:"isActive"`
	SortOrder      uint32     `json:"sortOrder"`
	ExpectedAmount int64      `json:"expectedAmount"`
	ActualAmount   int64      `json:"actualAmount"`
	Frequency      string     `json:"frequency"`
	StartDate      *time.Time `json:"startDate"`
	IsRecurring    bool       `json:"isRecurring"`
	PaymentMethod  string     `json:"paymentMethod"`
	Notes          string     `json:"notes"`
	Tags           []string   `json:"tags"`
	RawTags        string     `json:"-" gorm:"column:tags"` // JSON string array

	Category domain.DirectoryCategory `json:"category"`
}

// ListDirectorySourcesRequestDTO represents the request DTO for listing directory sources
type ListDirectorySourcesRequestDTO struct {
	_dto.Pagable
	Search      string  `json:"search"`
	Category    *string `json:"category"`
	Type        *string `json:"type"`
	IsActive    *bool   `json:"isActive"`
	IsRecurring *bool   `json:"isRecurring"`
}

// ListDirectorySourcesResponseDTO represents the response DTO for listing directory sources
type ListDirectorySourcesResponseDTO struct {
	Data  []DirectorySourceDTO `json:"data"`
	Total int64                `json:"total"`
}

// ToDomain converts DirectorySourceDTO to domain DirectorySource
func (d *DirectorySourceDTO) ToDomain() *domain.DirectorySource {
	// Convert tags to JSON string
	tagsJSON, _ := json.Marshal(d.Tags)

	return &domain.DirectorySource{
		BaseEntity: _entity.BaseEntity{
			ID: d.ID,
		},
		Name:        d.Name,
		Code:        d.Code,
		Description: d.Description,
		// Category:       d.Category,
		CategoryID:     d.CategoryID,
		Type:           d.Type,
		Icon:           d.Icon,
		Color:          d.Color,
		IsActive:       d.IsActive,
		SortOrder:      d.SortOrder,
		ExpectedAmount: d.ExpectedAmount,
		ActualAmount:   d.ActualAmount,
		Frequency:      d.Frequency,
		StartDate:      d.StartDate,
		IsRecurring:    d.IsRecurring,
		PaymentMethod:  d.PaymentMethod,
		Notes:          d.Notes,
		RawTags:        string(tagsJSON),
	}
}

// FromDomain converts domain DirectorySource to DirectorySourceDTO
func (d *DirectorySourceDTO) FromDomain(directorySource *domain.DirectorySource) {
	d.ID = directorySource.ID
	d.Name = directorySource.Name
	d.Code = directorySource.Code
	d.Description = directorySource.Description
	d.Category = directorySource.Category
	d.CategoryID = directorySource.CategoryID
	d.Type = directorySource.Type
	d.Icon = directorySource.Icon
	d.Color = directorySource.Color
	d.IsActive = directorySource.IsActive
	d.SortOrder = directorySource.SortOrder
	d.ExpectedAmount = directorySource.ExpectedAmount
	d.ActualAmount = directorySource.ActualAmount
	d.Frequency = directorySource.Frequency
	d.StartDate = directorySource.StartDate
	d.IsRecurring = directorySource.IsRecurring
	d.PaymentMethod = directorySource.PaymentMethod
	d.Notes = directorySource.Notes

	// Parse tags from JSON string
	var tags []string
	json.Unmarshal([]byte(directorySource.RawTags), &tags)
	d.Tags = tags
}
