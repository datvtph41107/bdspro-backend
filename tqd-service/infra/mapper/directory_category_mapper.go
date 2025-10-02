package mapper

import (
	"tqd/internal/domain"
)

// DirectoryCategoryMapper handles mapping between domain and proto
type DirectoryCategoryMapper struct{}

// NewDirectoryCategoryMapper creates a new DirectoryCategoryMapper
func NewDirectoryCategoryMapper() *DirectoryCategoryMapper {
	return &DirectoryCategoryMapper{}
}

// ToProto converts domain DirectoryCategory to proto DirectoryCategory
func (m *DirectoryCategoryMapper) ToProto(directoryCategory *domain.DirectoryCategory) map[string]interface{} {
	if directoryCategory == nil {
		return nil
	}

	return map[string]interface{}{
		"id":             directoryCategory.ID,
		"name":           directoryCategory.Name,
		"code":           directoryCategory.Code,
		"description":    directoryCategory.Description,
		"icon":           directoryCategory.Icon,
		"color":          directoryCategory.Color,
		"isActive":       directoryCategory.IsActive,
		"sortOrder":      directoryCategory.SortOrder,
		"level":          directoryCategory.Level,
		"path":           directoryCategory.Path,
		"directoryCount": directoryCategory.DirectoryCount,
		"createdAt":      directoryCategory.CreatedAt,
		"updatedAt":      directoryCategory.UpdatedAt,
	}
}

// ToProtoList converts slice of domain DirectoryCategory to slice of proto maps
func (m *DirectoryCategoryMapper) ToProtoList(categories []domain.DirectoryCategory) []map[string]interface{} {
	result := make([]map[string]interface{}, len(categories))
	for i, category := range categories {
		result[i] = m.ToProto(&category)
	}
	return result
}
