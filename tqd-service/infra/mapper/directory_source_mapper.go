package mapper

import (
	"encoding/json"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/dto"

	_entity "common/domain/entity"
	_utils "common/utils"
)

// DirectorySourceMapper handles mapping between domain and DTO
type DirectorySourceMapper struct{}

// NewDirectorySourceMapper creates a new DirectorySourceMapper
func NewDirectorySourceMapper() *DirectorySourceMapper {
	return &DirectorySourceMapper{}
}

// ToDTO converts domain DirectorySource to DirectorySourceDTO
func (m *DirectorySourceMapper) ToDTO(directorySource *domain.DirectorySource) *dto.DirectorySourceDTO {
	if directorySource == nil {
		return nil
	}

	// Parse tags from JSON string
	var tags []string
	json.Unmarshal([]byte(directorySource.RawTags), &tags)

	return &dto.DirectorySourceDTO{
		ID:             directorySource.ID,
		Name:           directorySource.Name,
		Code:           directorySource.Code,
		Description:    directorySource.Description,
		Category:       directorySource.Category,
		Type:           directorySource.Type,
		Icon:           directorySource.Icon,
		Color:          directorySource.Color,
		IsActive:       directorySource.IsActive,
		SortOrder:      directorySource.SortOrder,
		ExpectedAmount: directorySource.ExpectedAmount,
		ActualAmount:   directorySource.ActualAmount,
		Frequency:      directorySource.Frequency,
		StartDate:      directorySource.StartDate,
		IsRecurring:    directorySource.IsRecurring,
		PaymentMethod:  directorySource.PaymentMethod,
		Notes:          directorySource.Notes,
		Tags:           tags,
		BaseEntity: _entity.BaseEntity{
			ID:        directorySource.ID,
			CreatedAt: directorySource.CreatedAt,
			UpdatedAt: directorySource.UpdatedAt,
		},
	}
}

// ToDomain converts DirectorySourceDTO to domain DirectorySource
func (m *DirectorySourceMapper) ToDomain(dto *dto.DirectorySourceDTO) *domain.DirectorySource {
	if dto == nil {
		return nil
	}

	// Convert tags to JSON string
	tagsJSON, _ := json.Marshal(dto.Tags)

	return &domain.DirectorySource{
		BaseEntity: _entity.BaseEntity{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		Name:           dto.Name,
		Code:           dto.Code,
		Description:    dto.Description,
		Category:       dto.Category,
		Type:           dto.Type,
		Icon:           dto.Icon,
		Color:          dto.Color,
		IsActive:       dto.IsActive,
		SortOrder:      dto.SortOrder,
		ExpectedAmount: dto.ExpectedAmount,
		ActualAmount:   dto.ActualAmount,
		Frequency:      dto.Frequency,
		StartDate:      dto.StartDate,
		IsRecurring:    dto.IsRecurring,
		PaymentMethod:  dto.PaymentMethod,
		Notes:          dto.Notes,
		RawTags:        string(tagsJSON),
	}
}

// ToDTOList converts slice of domain DirectorySource to slice of DirectorySourceDTO
func (m *DirectorySourceMapper) ToDTOList(directorySources []domain.DirectorySource) []dto.DirectorySourceDTO {
	result := make([]dto.DirectorySourceDTO, len(directorySources))
	for i, directorySource := range directorySources {
		result[i] = *m.ToDTO(&directorySource)
	}
	return result
}

// ToProto converts DirectorySourceDTO to proto DirectorySource
func (m *DirectorySourceMapper) ToProto(dto *dto.DirectorySourceDTO) *tqdpb.DirectorySource {
	if dto == nil {
		return nil
	}

	return &tqdpb.DirectorySource{
		Id:          dto.ID,
		Name:        dto.Name,
		Code:        dto.Code,
		Description: dto.Description,
		// Category:       dto.Category,
		Type:           dto.Type,
		Icon:           dto.Icon,
		Color:          dto.Color,
		IsActive:       dto.IsActive,
		SortOrder:      dto.SortOrder,
		ExpectedAmount: dto.ExpectedAmount,
		ActualAmount:   dto.ActualAmount,
		Frequency:      dto.Frequency,
		StartDate:      _utils.FormatTimeToString(dto.StartDate),
		IsRecurring:    dto.IsRecurring,
		PaymentMethod:  dto.PaymentMethod,
		Notes:          dto.Notes,
		Tags:           dto.Tags,
		CreatedAt:      _utils.FormatTimeToString(dto.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(dto.UpdatedAt),
	}
}

// ToProtoList converts slice of DirectorySourceDTO to slice of proto DirectorySource
func (m *DirectorySourceMapper) ToProtoList(dtos []dto.DirectorySourceDTO) []*tqdpb.DirectorySource {
	result := make([]*tqdpb.DirectorySource, len(dtos))
	for i, dto := range dtos {
		result[i] = m.ToProto(&dto)
	}
	return result
}

// FromProto converts proto DirectorySource to DirectorySourceDTO
func (m *DirectorySourceMapper) FromProto(proto *tqdpb.DirectorySource) *dto.DirectorySourceDTO {
	if proto == nil {
		return nil
	}

	return &dto.DirectorySourceDTO{
		ID:          proto.Id,
		Name:        proto.Name,
		Code:        proto.Code,
		Description: proto.Description,
		// Category:       proto.Category,
		Type:           proto.Type,
		Icon:           proto.Icon,
		Color:          proto.Color,
		IsActive:       proto.IsActive,
		SortOrder:      proto.SortOrder,
		ExpectedAmount: proto.ExpectedAmount,
		ActualAmount:   proto.ActualAmount,
		Frequency:      proto.Frequency,
		StartDate:      _utils.ParseStringToTime(proto.StartDate),
		IsRecurring:    proto.IsRecurring,
		PaymentMethod:  proto.PaymentMethod,
		Notes:          proto.Notes,
		Tags:           proto.Tags,
		BaseEntity: _entity.BaseEntity{
			ID:        proto.Id,
			CreatedAt: _utils.ParseStringToTime(proto.CreatedAt),
			UpdatedAt: _utils.ParseStringToTime(proto.UpdatedAt),
		},
	}
}
