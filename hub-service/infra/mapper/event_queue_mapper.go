package mapper

import (
	_utils "common/utils"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type EventQueueMapper struct{}

func NewEventQueueMapper() *EventQueueMapper {
	return &EventQueueMapper{}
}

// EntityToPb chuyển từ Entity sang Protobuf DTO
func (m *EventQueueMapper) EntityToPb(entity *domain.EventQueueEntity) *hubpb.EventQueueDTO {
	if entity == nil {
		return nil
	}

	pb := &hubpb.EventQueueDTO{
		Id:         entity.ID,
		EventType:  entity.EventType,
		EventName:  entity.EventName,
		Status:     uint32(entity.Status),
		RetryCount: entity.RetryCount,
		MaxRetries: entity.MaxRetries,
		CreatedAt:  _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:  _utils.FormatTimeToString(entity.UpdatedAt),
	}

	if entity.EventData != nil {
		pb.EventData = entity.EventData
	}

	if entity.Metadata != nil {
		pb.Metadata = entity.Metadata
	}

	if entity.ProfileID != nil {
		pb.ProfileId = entity.ProfileID
	}

	if entity.OrganizationID != nil {
		pb.OrganizationId = entity.OrganizationID
	}

	if entity.ScheduledAt != nil {
		scheduledAtStr := _utils.FormatTimeToString(entity.ScheduledAt)
		pb.ScheduledAt = &scheduledAtStr
	}

	if entity.ProcessedAt != nil {
		processedAtStr := _utils.FormatTimeToString(entity.ProcessedAt)
		pb.ProcessedAt = &processedAtStr
	}

	if entity.ErrorMessage != nil {
		pb.ErrorMessage = entity.ErrorMessage
	}

	if entity.CreatedBy != nil {
		pb.CreatedBy = entity.CreatedBy
	}

	if entity.UpdatedBy != nil {
		pb.UpdatedBy = entity.UpdatedBy
	}

	return pb
}

// ListEntityToPb chuyển từ List Entity sang List Protobuf DTO
func (m *EventQueueMapper) ListEntityToPb(entities []domain.EventQueueEntity) []*hubpb.EventQueueDTO {
	result := make([]*hubpb.EventQueueDTO, 0, len(entities))

	for _, entity := range entities {
		result = append(result, m.EntityToPb(&entity))
	}

	return result
}
