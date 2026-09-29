package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
)

// DirectoryCategoryMapper handles mapping between domain and proto
type DirectoryCategoryMapper struct{}

// NewDirectoryCategoryMapper creates a new DirectoryCategoryMapper
func NewDirectoryCategoryMapper() *DirectoryCategoryMapper {
	return &DirectoryCategoryMapper{}
}

// ToProto converts domain DirectoryCategory to proto DirectoryCategory
func (m *DirectoryCategoryMapper) ToProto(directoryCategory *domain.DirectoryCategory) *tqdpb.DirectoryCategory {
	if directoryCategory == nil {
		return nil
	}

	return &tqdpb.DirectoryCategory{
		Id:             directoryCategory.ID,
		Name:           directoryCategory.Name,
		Code:           directoryCategory.Code,
		Description:    directoryCategory.Description,
		Icon:           directoryCategory.Icon,
		Color:          directoryCategory.Color,
		IsActive:       directoryCategory.IsActive,
		SortOrder:      directoryCategory.SortOrder,
		Level:          directoryCategory.Level,
		Path:           directoryCategory.Path,
		DirectoryCount: directoryCategory.DirectoryCount,
		CreatedAt:      _utils.FormatTimeToString(directoryCategory.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(directoryCategory.UpdatedAt),
		CreatedBy:      *directoryCategory.CreatedBy,
		UpdatedBy:      *directoryCategory.UpdatedBy,
	}
}

// ToProtoList converts slice of domain DirectoryCategory to slice of proto DirectoryCategory
func (m *DirectoryCategoryMapper) ToProtoList(categories []domain.DirectoryCategory) []*tqdpb.DirectoryCategory {
	result := make([]*tqdpb.DirectoryCategory, len(categories))
	for i, category := range categories {
		result[i] = m.ToProto(&category)
	}
	return result
}
