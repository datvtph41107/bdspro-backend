package mapper

import (
	"encoding/json"
	"time"
	"tqd/internal/domain"
	"tqd/internal/dto"

	_entity "common/domain/entity"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
)

// DirectorySupplierMapper handles mapping between domain, DTO and proto
type DirectorySupplierMapper struct{}

// NewDirectorySupplierMapper creates a new DirectorySupplierMapper
func NewDirectorySupplierMapper() *DirectorySupplierMapper {
	return &DirectorySupplierMapper{}
}

// ToDTO converts domain DirectorySupplier to DirectorySupplierDTO
func (m *DirectorySupplierMapper) ToDTO(supplier *domain.DirectorySupplier) *dto.DirectorySupplierDTO {
	if supplier == nil {
		return nil
	}

	// Parse tags from JSON string
	var tags []string
	if supplier.RawTags != "" {
		json.Unmarshal([]byte(supplier.RawTags), &tags)
	}

	// Parse categories from JSON string
	var categories []string
	if supplier.RawCategories != "" {
		json.Unmarshal([]byte(supplier.RawCategories), &categories)
	}

	return &dto.DirectorySupplierDTO{
		ID:              supplier.ID,
		Name:            supplier.Name,
		Code:            supplier.Code,
		Description:     supplier.Description,
		ContactPerson:   supplier.ContactPerson,
		Phone:           supplier.Phone,
		Email:           supplier.Email,
		Address:         supplier.Address,
		Website:         supplier.Website,
		TaxCode:         supplier.TaxCode,
		BusinessLicense: supplier.BusinessLicense,
		Rating:          supplier.Rating,
		IsActive:        supplier.IsActive,
		ServiceCount:    supplier.ServiceCount,
		ContractCount:   supplier.ContractCount,
		TotalValue:      supplier.TotalValue,
		JoinDate:        supplier.JoinDate,
		LastContactDate: supplier.LastContactDate,
		Notes:           supplier.Notes,
		Tags:            tags,
		Categories:      categories,
		CreatedAt:       supplier.CreatedAt,
		UpdatedAt:       supplier.UpdatedAt,
		CreatedBy:       supplier.CreatedBy,
		UpdatedBy:       supplier.UpdatedBy,
	}
}

// ToDomain converts DirectorySupplierDTO to domain DirectorySupplier
func (m *DirectorySupplierMapper) ToDomain(dto *dto.DirectorySupplierDTO) *domain.DirectorySupplier {
	if dto == nil {
		return nil
	}

	// Convert tags to JSON string
	var tagsJSON string
	if dto.Tags != nil {
		tagsBytes, _ := json.Marshal(dto.Tags)
		tagsJSON = string(tagsBytes)
	}

	// Convert categories to JSON string
	var categoriesJSON string
	if dto.Categories != nil {
		categoriesBytes, _ := json.Marshal(dto.Categories)
		categoriesJSON = string(categoriesBytes)
	}

	return &domain.DirectorySupplier{
		BaseEntity: _entity.BaseEntity{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		Name:            dto.Name,
		Code:            dto.Code,
		Description:     dto.Description,
		ContactPerson:   dto.ContactPerson,
		Phone:           dto.Phone,
		Email:           dto.Email,
		Address:         dto.Address,
		Website:         dto.Website,
		TaxCode:         dto.TaxCode,
		BusinessLicense: dto.BusinessLicense,
		Rating:          dto.Rating,
		IsActive:        dto.IsActive,
		ServiceCount:    dto.ServiceCount,
		ContractCount:   dto.ContractCount,
		TotalValue:      dto.TotalValue,
		JoinDate:        dto.JoinDate,
		LastContactDate: dto.LastContactDate,
		Notes:           dto.Notes,
		RawTags:         tagsJSON,
		RawCategories:   categoriesJSON,
	}
}

// ToDomainFromCreateRequest converts CreateDirectorySupplierRequestDTO to domain DirectorySupplier
func (m *DirectorySupplierMapper) ToDomainFromCreateRequest(req *dto.CreateDirectorySupplierRequestDTO) (*domain.DirectorySupplier, error) {
	if req == nil {
		return nil, nil
	}

	// Parse join date if provided
	var joinDate *time.Time
	if req.JoinDate != nil && *req.JoinDate != "" {
		joinDate = _utils.ParseStringToTime(*req.JoinDate)
	}

	// Parse last contact date if provided
	var lastContactDate *time.Time
	if req.LastContactDate != nil && *req.LastContactDate != "" {
		lastContactDate = _utils.ParseStringToTime(*req.LastContactDate)
	}

	// Convert tags to JSON string
	var tagsJSON string
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, err
		}
		tagsJSON = string(tagsBytes)
	}

	// Convert categories to JSON string
	var categoriesJSON string
	if req.Categories != nil {
		categoriesBytes, err := json.Marshal(req.Categories)
		if err != nil {
			return nil, err
		}
		categoriesJSON = string(categoriesBytes)
	}

	return &domain.DirectorySupplier{
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		ContactPerson:   req.ContactPerson,
		Phone:           req.Phone,
		Email:           req.Email,
		Address:         req.Address,
		Website:         req.Website,
		TaxCode:         req.TaxCode,
		BusinessLicense: req.BusinessLicense,
		Rating:          req.Rating,
		IsActive:        req.IsActive,
		Notes:           req.Notes,
		RawTags:         tagsJSON,
		RawCategories:   categoriesJSON,
		JoinDate:        joinDate,
		LastContactDate: lastContactDate,
		ServiceCount:    0,
		ContractCount:   0,
		TotalValue:      0,
	}, nil
}

// ToDomainFromUpdateRequest converts UpdateDirectorySupplierRequestDTO to domain DirectorySupplier
func (m *DirectorySupplierMapper) ToDomainFromUpdateRequest(req *dto.UpdateDirectorySupplierRequestDTO, id uint64) (*domain.DirectorySupplier, error) {
	if req == nil {
		return nil, nil
	}

	// Parse join date if provided
	var joinDate *time.Time
	if req.JoinDate != nil && *req.JoinDate != "" {
		joinDate = _utils.ParseStringToTime(*req.JoinDate)
	}

	// Parse last contact date if provided
	var lastContactDate *time.Time
	if req.LastContactDate != nil && *req.LastContactDate != "" {
		lastContactDate = _utils.ParseStringToTime(*req.LastContactDate)
	}

	// Convert tags to JSON string
	var tagsJSON string
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, err
		}
		tagsJSON = string(tagsBytes)
	}

	// Convert categories to JSON string
	var categoriesJSON string
	if req.Categories != nil {
		categoriesBytes, err := json.Marshal(req.Categories)
		if err != nil {
			return nil, err
		}
		categoriesJSON = string(categoriesBytes)
	}

	return &domain.DirectorySupplier{
		BaseEntity: _entity.BaseEntity{
			ID: id,
		},
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		ContactPerson:   req.ContactPerson,
		Phone:           req.Phone,
		Email:           req.Email,
		Address:         req.Address,
		Website:         req.Website,
		TaxCode:         req.TaxCode,
		BusinessLicense: req.BusinessLicense,
		Rating:          req.Rating,
		IsActive:        req.IsActive,
		Notes:           req.Notes,
		RawTags:         tagsJSON,
		RawCategories:   categoriesJSON,
		JoinDate:        joinDate,
		LastContactDate: lastContactDate,
	}, nil
}

// ToDTOList converts slice of domain DirectorySupplier to slice of DirectorySupplierDTO
func (m *DirectorySupplierMapper) ToDTOList(suppliers []domain.DirectorySupplier) []dto.DirectorySupplierDTO {
	if suppliers == nil {
		return nil
	}
	result := make([]dto.DirectorySupplierDTO, len(suppliers))
	for i, supplier := range suppliers {
		dto := m.ToDTO(&supplier)
		if dto != nil {
			result[i] = *dto
		}
	}
	return result
}

// ToProto converts DirectorySupplierDTO to proto DirectorySupplier
func (m *DirectorySupplierMapper) ToProto(dto *dto.DirectorySupplierDTO) *tqdpb.DirectorySupplier {
	if dto == nil {
		return nil
	}

	return &tqdpb.DirectorySupplier{
		Id:              dto.ID,
		Name:            dto.Name,
		Code:            dto.Code,
		Description:     dto.Description,
		ContactPerson:   dto.ContactPerson,
		Phone:           dto.Phone,
		Email:           dto.Email,
		Address:         dto.Address,
		Website:         dto.Website,
		TaxCode:         dto.TaxCode,
		BusinessLicense: dto.BusinessLicense,
		Rating:          dto.Rating,
		IsActive:        dto.IsActive,
		ServiceCount:    dto.ServiceCount,
		ContractCount:   dto.ContractCount,
		TotalValue:      dto.TotalValue,
		JoinDate:        _utils.FormatTimeToString(dto.JoinDate),
		LastContactDate: _utils.FormatTimeToString(dto.LastContactDate),
		Notes:           dto.Notes,
		Tags:            dto.Tags,
		Categories:      dto.Categories,
		CreatedAt:       _utils.FormatTimeToString(dto.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(dto.UpdatedAt),
		CreatedBy:       uint64Value(dto.CreatedBy),
		UpdatedBy:       uint64Value(dto.UpdatedBy),
	}
}

// FromProto converts proto DirectorySupplier to DirectorySupplierDTO
func (m *DirectorySupplierMapper) FromProto(proto *tqdpb.DirectorySupplier) (*dto.DirectorySupplierDTO, error) {
	if proto == nil {
		return nil, nil
	}

	// Parse join date if provided
	var joinDate *time.Time
	if proto.JoinDate != "" {
		joinDate = _utils.ParseStringToTime(proto.JoinDate)
	}

	// Parse last contact date if provided
	var lastContactDate *time.Time
	if proto.LastContactDate != "" {
		lastContactDate = _utils.ParseStringToTime(proto.LastContactDate)
	}

	return &dto.DirectorySupplierDTO{
		ID:              proto.Id,
		Name:            proto.Name,
		Code:            proto.Code,
		Description:     proto.Description,
		ContactPerson:   proto.ContactPerson,
		Phone:           proto.Phone,
		Email:           proto.Email,
		Address:         proto.Address,
		Website:         proto.Website,
		TaxCode:         proto.TaxCode,
		BusinessLicense: proto.BusinessLicense,
		Rating:          proto.Rating,
		IsActive:        proto.IsActive,
		ServiceCount:    proto.ServiceCount,
		ContractCount:   proto.ContractCount,
		TotalValue:      proto.TotalValue,
		JoinDate:        joinDate,
		LastContactDate: lastContactDate,
		Notes:           proto.Notes,
		Tags:            proto.Tags,
		Categories:      proto.Categories,
		CreatedAt:       _utils.ParseStringToTime(proto.CreatedAt),
		UpdatedAt:       _utils.ParseStringToTime(proto.UpdatedAt),
		CreatedBy:       &proto.CreatedBy,
		UpdatedBy:       &proto.UpdatedBy,
	}, nil
}

// Helper function to convert *uint64 to uint64 (with default 0)
func uint64Value(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// ToResponse converts slice of DirectorySupplierDTO to slice of proto DirectorySupplier
func (m *DirectorySupplierMapper) ToResponse(responses *dto.DirectorySupplierDTO) []*dto.DirectorySupplierDTO {
	if responses == nil {
		return nil
	}

	return []*dto.DirectorySupplierDTO{responses}
}
