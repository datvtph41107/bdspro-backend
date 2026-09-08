package dto

import (
	"encoding/json"
	"time"
	"tqd/internal/domain"

	_dto "common/domain/dto"
	_entity "common/domain/entity"
)

// DirectorySupplierDTO represents the data transfer object for directory supplier
type DirectorySupplierDTO struct {
	ID              uint64     `json:"id"`
	Name            string     `json:"name"`
	Code            string     `json:"code"`
	Description     string     `json:"description"`
	ContactPerson   string     `json:"contactPerson"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Address         string     `json:"address"`
	Website         string     `json:"website"`
	TaxCode         string     `json:"taxCode"`
	BusinessLicense string     `json:"businessLicense"`
	Rating          float64    `json:"rating"`
	IsActive        bool       `json:"isActive"`
	ServiceCount    uint32     `json:"serviceCount"`
	ContractCount   uint32     `json:"contractCount"`
	TotalValue      int64      `json:"totalValue"`
	JoinDate        *time.Time `json:"joinDate"`
	LastContactDate *time.Time `json:"lastContactDate"`
	Notes           string     `json:"notes"`
	Tags            []string   `json:"tags"`
	Categories      []string   `json:"categories"`
	CreatedAt       *time.Time `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	RawTags         string     `json:"rawTags"`
	RawCategories   string     `json:"rawCategories"`
	CreatedBy       *uint64    `json:"createdBy"`
	UpdatedBy       *uint64    `json:"updatedBy"`
}

// CreateDirectorySupplierRequestDTO represents the request DTO for creating directory supplier
type CreateDirectorySupplierRequestDTO struct {
	Name            string   `json:"name" validate:"required"`
	Code            string   `json:"code" validate:"required"`
	Description     string   `json:"description"`
	ContactPerson   string   `json:"contactPerson"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email" validate:"omitempty,email"`
	Address         string   `json:"address"`
	Website         string   `json:"website" validate:"omitempty,url"`
	TaxCode         string   `json:"taxCode"`
	BusinessLicense string   `json:"businessLicense"`
	Rating          float64  `json:"rating" validate:"min=0,max=5"`
	IsActive        bool     `json:"isActive"`
	Notes           string   `json:"notes"`
	Tags            []string `json:"tags"`
	Categories      []string `json:"categories"`
	JoinDate        *string  `json:"joinDate"`        // ISO string
	LastContactDate *string  `json:"lastContactDate"` // ISO string
}

// UpdateDirectorySupplierRequestDTO represents the request DTO for updating directory supplier
type UpdateDirectorySupplierRequestDTO struct {
	ID              uint64   `json:"id"`
	Name            string   `json:"name" validate:"required"`
	Code            string   `json:"code" validate:"required"`
	Description     string   `json:"description"`
	ContactPerson   string   `json:"contactPerson"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email" validate:"omitempty,email"`
	Address         string   `json:"address"`
	Website         string   `json:"website" validate:"omitempty,url"`
	TaxCode         string   `json:"taxCode"`
	BusinessLicense string   `json:"businessLicense"`
	Rating          float64  `json:"rating" validate:"min=0,max=5"`
	IsActive        bool     `json:"isActive"`
	Notes           string   `json:"notes"`
	Tags            []string `json:"tags"`
	Categories      []string `json:"categories"`
	JoinDate        *string  `json:"joinDate"`        // ISO string
	LastContactDate *string  `json:"lastContactDate"` // ISO string
}

// DirectorySupplierFilterDTO represents the filter DTO for listing directory suppliers
type DirectorySupplierFilterDTO struct {
	_dto.Pagable
	Search     string   `json:"search"`
	Categories []string `json:"categories"`
	IsActive   *bool    `json:"isActive"`
	MinRating  *float64 `json:"minRating"`
	MaxRating  *float64 `json:"maxRating"`
	SortBy     string   `json:"sortBy"`    // rating, name, created_at
	SortOrder  string   `json:"sortOrder"` // asc, desc
}

// ListDirectorySuppliersResponseDTO represents the response DTO for listing directory suppliers
type ListDirectorySuppliersResponseDTO struct {
	Data  []DirectorySupplierDTO `json:"data"`
	Total int64                  `json:"total"`
	Page  uint32                 `json:"page"`
	Size  uint32                 `json:"size"`
}

// ToDomain converts DirectorySupplierDTO to domain DirectorySupplier
func (d *DirectorySupplierDTO) ToDomain() *domain.DirectorySupplier {
	// Convert tags to JSON string
	var tagsJSON string
	if d.Tags != nil {
		tagsBytes, _ := json.Marshal(d.Tags)
		tagsJSON = string(tagsBytes)
	}

	// Convert categories to JSON string
	var categoriesJSON string
	if d.Categories != nil {
		categoriesBytes, _ := json.Marshal(d.Categories)
		categoriesJSON = string(categoriesBytes)
	}

	return &domain.DirectorySupplier{
		BaseEntity: _entity.BaseEntity{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:            d.Name,
		Code:            d.Code,
		Description:     d.Description,
		ContactPerson:   d.ContactPerson,
		Phone:           d.Phone,
		Email:           d.Email,
		Address:         d.Address,
		Website:         d.Website,
		TaxCode:         d.TaxCode,
		BusinessLicense: d.BusinessLicense,
		Rating:          d.Rating,
		IsActive:        d.IsActive,
		ServiceCount:    d.ServiceCount,
		ContractCount:   d.ContractCount,
		TotalValue:      d.TotalValue,
		JoinDate:        d.JoinDate,
		LastContactDate: d.LastContactDate,
		Notes:           d.Notes,
		RawTags:         tagsJSON,
		RawCategories:   categoriesJSON,
	}
}

// FromDomain converts domain DirectorySupplier to DirectorySupplierDTO
func (d *DirectorySupplierDTO) FromDomain(supplier *domain.DirectorySupplier) {
	d.ID = supplier.ID
	d.Name = supplier.Name
	d.Code = supplier.Code
	d.Description = supplier.Description
	d.ContactPerson = supplier.ContactPerson
	d.Phone = supplier.Phone
	d.Email = supplier.Email
	d.Address = supplier.Address
	d.Website = supplier.Website
	d.TaxCode = supplier.TaxCode
	d.BusinessLicense = supplier.BusinessLicense
	d.Rating = supplier.Rating
	d.IsActive = supplier.IsActive
	d.ServiceCount = supplier.ServiceCount
	d.ContractCount = supplier.ContractCount
	d.TotalValue = supplier.TotalValue
	d.JoinDate = supplier.JoinDate
	d.LastContactDate = supplier.LastContactDate
	d.Notes = supplier.Notes
	d.CreatedAt = supplier.CreatedAt
	d.UpdatedAt = supplier.UpdatedAt
	d.CreatedBy = supplier.CreatedBy
	d.UpdatedBy = supplier.UpdatedBy

	// Parse tags from JSON string
	if supplier.RawTags != "" {
		var tags []string
		json.Unmarshal([]byte(supplier.RawTags), &tags)
		d.Tags = tags
	}

	// Parse categories from JSON string
	if supplier.RawCategories != "" {
		var categories []string
		json.Unmarshal([]byte(supplier.RawCategories), &categories)
		d.Categories = categories
	}
}

// ToDomainFromCreateRequest converts CreateDirectorySupplierRequestDTO to domain DirectorySupplier
func (d *CreateDirectorySupplierRequestDTO) ToDomain() (*domain.DirectorySupplier, error) {
	// Convert tags to JSON string
	var tagsJSON string
	if d.Tags != nil {
		tagsBytes, err := json.Marshal(d.Tags)
		if err != nil {
			return nil, err
		}
		tagsJSON = string(tagsBytes)
	}

	// Convert categories to JSON string
	var categoriesJSON string
	if d.Categories != nil {
		categoriesBytes, err := json.Marshal(d.Categories)
		if err != nil {
			return nil, err
		}
		categoriesJSON = string(categoriesBytes)
	}

	return &domain.DirectorySupplier{
		Name:            d.Name,
		Code:            d.Code,
		Description:     d.Description,
		ContactPerson:   d.ContactPerson,
		Phone:           d.Phone,
		Email:           d.Email,
		Address:         d.Address,
		Website:         d.Website,
		TaxCode:         d.TaxCode,
		BusinessLicense: d.BusinessLicense,
		Rating:          d.Rating,
		IsActive:        d.IsActive,
		Notes:           d.Notes,
		RawTags:         tagsJSON,
		RawCategories:   categoriesJSON,
		ServiceCount:    0,
		ContractCount:   0,
		TotalValue:      0,
	}, nil
}

// ToDomainFromUpdateRequest converts UpdateDirectorySupplierRequestDTO to domain DirectorySupplier
func (d *UpdateDirectorySupplierRequestDTO) ToDomain(id uint64) (*domain.DirectorySupplier, error) {
	// Convert tags to JSON string
	var tagsJSON string
	if d.Tags != nil {
		tagsBytes, err := json.Marshal(d.Tags)
		if err != nil {
			return nil, err
		}
		tagsJSON = string(tagsBytes)
	}

	// Convert categories to JSON string
	var categoriesJSON string
	if d.Categories != nil {
		categoriesBytes, err := json.Marshal(d.Categories)
		if err != nil {
			return nil, err
		}
		categoriesJSON = string(categoriesBytes)
	}

	return &domain.DirectorySupplier{
		BaseEntity: _entity.BaseEntity{
			ID: id,
		},
		Name:            d.Name,
		Code:            d.Code,
		Description:     d.Description,
		ContactPerson:   d.ContactPerson,
		Phone:           d.Phone,
		Email:           d.Email,
		Address:         d.Address,
		Website:         d.Website,
		TaxCode:         d.TaxCode,
		BusinessLicense: d.BusinessLicense,
		Rating:          d.Rating,
		IsActive:        d.IsActive,
		Notes:           d.Notes,
		RawTags:         tagsJSON,
		RawCategories:   categoriesJSON,
	}, nil
}
