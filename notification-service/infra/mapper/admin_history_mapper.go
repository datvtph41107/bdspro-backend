package mapper

import (
	shared_dto "pb/dto"
	"time"

	"notification/internal/domain"
	"notification/internal/dto"
)

type AdminHistoryMapper struct{}

func NewAdminHistoryMapper() *AdminHistoryMapper {
	return &AdminHistoryMapper{}
}

// AdminHistoryToDTO chuyển đổi AdminHistoryEntity thành AdminHistoryDTO
func (m *AdminHistoryMapper) AdminHistoryToDTO(entity *domain.AdminHistoryEntity) *dto.AdminHistoryDTO {
	if entity == nil {
		return nil
	}

	dto := &dto.AdminHistoryDTO{
		ID:         entity.ID,
		TargetId:   entity.TargetId,
		TargetType: entity.TargetType,
		ActionType: entity.ActionType,
		Title:      entity.Title,
		Note:       []string(entity.Note),
		PreStage:   entity.PreStage,
		AfterStage: entity.AfterStage,
		AdminID:    entity.AdminID,
		OwnerID:    entity.OwnerID,
		OwnerType:  entity.OwnerOf,
		IsInternal: entity.IsInternal,
		AdminRole:  entity.AdminRole,
		IPAddress:  entity.IPAddress,
		UserAgent:  entity.UserAgent,
		CreatedAt:  entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  entity.UpdatedAt.Format(time.RFC3339),
	}

	// Thêm thông tin user nếu có
	if entity.AdminName != nil {
		dto.AdminName = entity.AdminName
	}
	if entity.AdminAvatar != nil {
		dto.AdminAvatar = entity.AdminAvatar
	}
	if entity.Status != nil {
		dto.Status = entity.Status
	}

	return dto
}

// AdminHistoryToSharedDTO chuyển đổi AdminHistoryEntity thành shared_dto.HistoryDTO
func (m *AdminHistoryMapper) AdminHistoryToSharedDTO(entity *domain.AdminHistoryEntity) *shared_dto.HistoryDTO {
	if entity == nil {
		return nil
	}

	return &shared_dto.HistoryDTO{
		Title:      entity.Title,
		Note:       []string(entity.Note),
		TargetId:   entity.TargetId,
		TargetType: entity.TargetType,
		ActionType: entity.ActionType,
		OwnerID:    entity.OwnerID,
		OwnerOf:    entity.OwnerOf,
	}
}

// CreateDTOToEntity chuyển đổi AdminHistoryCreateDTO thành AdminHistoryEntity
func (m *AdminHistoryMapper) CreateDTOToEntity(dto *dto.AdminHistoryCreateDTO) *domain.AdminHistoryEntity {
	if dto == nil {
		return nil
	}

	return &domain.AdminHistoryEntity{
		TargetId:   dto.TargetID,
		TargetType: dto.TargetType,
		ActionType: dto.ActionType,
		Title:      dto.Title,
		Note:       dto.Note,
		PreStage:   dto.PreStage,
		AfterStage: dto.AfterStage,
		AdminID:    dto.AdminID,
		OwnerID:    dto.OwnerID,
		OwnerOf:    dto.OwnerType,
		IsInternal: dto.IsInternal,
		AdminRole:  dto.AdminRole,
		IPAddress:  dto.IPAddress,
		UserAgent:  dto.UserAgent,
	}
}

// EntitiesToDTOs chuyển đổi slice AdminHistoryEntity thành slice AdminHistoryDTO
func (m *AdminHistoryMapper) EntitiesToDTOs(entities []domain.AdminHistoryEntity) []dto.AdminHistoryDTO {
	if entities == nil {
		return nil
	}

	dtos := make([]dto.AdminHistoryDTO, len(entities))
	for i, entity := range entities {
		dto := m.AdminHistoryToDTO(&entity)
		if dto != nil {
			dtos[i] = *dto
		}
	}

	return dtos
}
