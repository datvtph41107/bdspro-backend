package mapper

import (
	"encoding/json"
	"tqd/internal/domain"

	_utils "common/utils"
	tqdpb "pb/types/tqd"
)

// DirectorySupplierMapper handles mapping between domain and response
type DirectorySupplierMapper struct{}

// NewDirectorySupplierMapper creates a new DirectorySupplierMapper
func NewDirectorySupplierMapper() *DirectorySupplierMapper {
	return &DirectorySupplierMapper{}
}

// ToResponse converts domain DirectorySupplier to response
func (m *DirectorySupplierMapper) ToResponse(supplier *domain.DirectorySupplier) map[string]interface{} {
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

	return map[string]interface{}{
		"id":              supplier.ID,
		"name":            supplier.Name,
		"code":            supplier.Code,
		"description":     supplier.Description,
		"contactPerson":   supplier.ContactPerson,
		"phone":           supplier.Phone,
		"email":           supplier.Email,
		"address":         supplier.Address,
		"website":         supplier.Website,
		"taxCode":         supplier.TaxCode,
		"businessLicense": supplier.BusinessLicense,
		"rating":          supplier.Rating,
		"isActive":        supplier.IsActive,
		"serviceCount":    supplier.ServiceCount,
		"contractCount":   supplier.ContractCount,
		"totalValue":      supplier.TotalValue,
		"joinDate":        _utils.FormatTimeToString(supplier.JoinDate),
		"lastContactDate": _utils.FormatTimeToString(supplier.LastContactDate),
		"notes":           supplier.Notes,
		"tags":            tags,
		"categories":      categories,
		"createdAt":       _utils.FormatTimeToString(supplier.CreatedAt),
		"updatedAt":       _utils.FormatTimeToString(supplier.UpdatedAt),
		"createdBy":       supplier.CreatedBy,
		"updatedBy":       supplier.UpdatedBy,
	}
}

// ToDomain converts request to domain DirectorySupplier
func (m *DirectorySupplierMapper) ToDomain(data map[string]interface{}) *domain.DirectorySupplier {
	if data == nil {
		return nil
	}

	supplier := &domain.DirectorySupplier{}

	if id, ok := data["id"].(float64); ok {
		supplier.ID = uint64(id)
	}
	if name, ok := data["name"].(string); ok {
		supplier.Name = name
	}
	if code, ok := data["code"].(string); ok {
		supplier.Code = code
	}
	if description, ok := data["description"].(string); ok {
		supplier.Description = description
	}
	if contactPerson, ok := data["contactPerson"].(string); ok {
		supplier.ContactPerson = contactPerson
	}
	if phone, ok := data["phone"].(string); ok {
		supplier.Phone = phone
	}
	if email, ok := data["email"].(string); ok {
		supplier.Email = email
	}
	if address, ok := data["address"].(string); ok {
		supplier.Address = address
	}
	if website, ok := data["website"].(string); ok {
		supplier.Website = website
	}
	if taxCode, ok := data["taxCode"].(string); ok {
		supplier.TaxCode = taxCode
	}
	if businessLicense, ok := data["businessLicense"].(string); ok {
		supplier.BusinessLicense = businessLicense
	}
	if rating, ok := data["rating"].(float64); ok {
		supplier.Rating = rating
	}
	if isActive, ok := data["isActive"].(bool); ok {
		supplier.IsActive = isActive
	}
	if serviceCount, ok := data["serviceCount"].(float64); ok {
		supplier.ServiceCount = uint32(serviceCount)
	}
	if contractCount, ok := data["contractCount"].(float64); ok {
		supplier.ContractCount = uint32(contractCount)
	}
	if totalValue, ok := data["totalValue"].(float64); ok {
		supplier.TotalValue = int64(totalValue)
	}
	if joinDate, ok := data["joinDate"].(string); ok {
		supplier.JoinDate = _utils.ParseStringToTime(joinDate)
	}
	if lastContactDate, ok := data["lastContactDate"].(string); ok {
		supplier.LastContactDate = _utils.ParseStringToTime(lastContactDate)
	}
	if notes, ok := data["notes"].(string); ok {
		supplier.Notes = notes
	}

	// Handle tags
	if tags, ok := data["tags"].([]interface{}); ok {
		tagStrings := make([]string, len(tags))
		for i, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				tagStrings[i] = tagStr
			}
		}
		if tagJSON, err := json.Marshal(tagStrings); err == nil {
			supplier.RawTags = string(tagJSON)
		}
	}

	// Handle categories
	if categories, ok := data["categories"].([]interface{}); ok {
		categoryStrings := make([]string, len(categories))
		for i, category := range categories {
			if categoryStr, ok := category.(string); ok {
				categoryStrings[i] = categoryStr
			}
		}
		if categoryJSON, err := json.Marshal(categoryStrings); err == nil {
			supplier.RawCategories = string(categoryJSON)
		}
	}

	return supplier
}

// ToListResponse converts list of domain DirectorySupplier to response
func (m *DirectorySupplierMapper) ToListResponse(suppliers []domain.DirectorySupplier, total int64) map[string]interface{} {
	data := make([]map[string]interface{}, len(suppliers))
	for i, supplier := range suppliers {
		data[i] = m.ToResponse(&supplier)
	}

	return map[string]interface{}{
		"data":  data,
		"total": total,
	}
}

// ToProto converts domain DirectorySupplier to proto DirectorySupplier
func (m *DirectorySupplierMapper) ToProto(supplier *domain.DirectorySupplier) *tqdpb.DirectorySupplier {
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

	return &tqdpb.DirectorySupplier{
		Id:              supplier.ID,
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
		JoinDate:        _utils.FormatTimeToString(supplier.JoinDate),
		LastContactDate: _utils.FormatTimeToString(supplier.LastContactDate),
		Notes:           supplier.Notes,
		Tags:            tags,
		Categories:      categories,
		CreatedAt:       _utils.FormatTimeToString(supplier.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(supplier.UpdatedAt),
		CreatedBy:       *supplier.CreatedBy,
		UpdatedBy:       *supplier.UpdatedBy,
	}
}
