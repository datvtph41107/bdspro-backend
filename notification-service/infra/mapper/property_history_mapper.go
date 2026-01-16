package mapper

import (
	_enum "common/domain/enum"
	_utils "common/utils"

	"notification/internal/domain"
	"notification/internal/dto"
	notificationpb "pb/types/notification"
)

type PropertyHistoryMapper struct{}

func NewPropertyHistoryMapper() *PropertyHistoryMapper {
	return &PropertyHistoryMapper{}
}

func (m *PropertyHistoryMapper) CreateProtoToDTO(
	req *notificationpb.CreatePropertyHistoryRequest,
) *dto.CreatePropertyHistoryDTO {

	if req == nil {
		return nil
	}

	return &dto.CreatePropertyHistoryDTO{
		SubjectID:   req.SubjectId,
		SubjectType: _enum.SubjectTypeProperty(req.SubjectType),
		Action:      _enum.PropertyAction(req.Action),
		Description: req.Description,
		ActorID:     req.ActorId,
	}
}

func (m *PropertyHistoryMapper) SearchProtoToDTO(
	req *notificationpb.GetPropertyHistoryRequest,
) dto.PropertyHistoryQueryDTO {

	return dto.PropertyHistoryQueryDTO{
		SubjectID:   &req.SubjectId,
		SubjectType: &req.SubjectType,
	}
}

func (m *PropertyHistoryMapper) EntityToItem(
	entity *domain.PropertyHistoryEntity,
) *notificationpb.PropertyHistoryItem {

	if entity == nil {
		return nil
	}

	return &notificationpb.PropertyHistoryItem{
		Id:        entity.ID,
		SubjectId: entity.SubjectID,
		// SubjectType: entity.SubjectType.String(),
		Action:      entity.Action.String(),
		Description: entity.Description,
		ActorId:     entity.ActorID,
		OccurredAt:  _utils.FormatTimeToString(entity.OccurredAt),
		CreatedAt:   _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

func (m *PropertyHistoryMapper) EntitiesToListResponse(
	entities []domain.PropertyHistoryEntity,
	total int64,
	page, size int32,
) *notificationpb.GetPropertyHistoryResponse {

	items := make([]*notificationpb.PropertyHistoryItem, 0, len(entities))

	for _, e := range entities {
		eCopy := e // 🧠 best practice (tránh pointer reuse bug)
		items = append(items, m.EntityToItem(&eCopy))
	}

	return &notificationpb.GetPropertyHistoryResponse{
		Data:  items,
		Total: total,
	}
}

func (m *PropertyHistoryMapper) EntityToDetailResponse(
	entity *domain.PropertyHistoryEntity,
) *notificationpb.GetPropertyHistoryByIDResponse {

	return &notificationpb.GetPropertyHistoryByIDResponse{
		Data: m.EntityToItem(entity), // ✅ dùng thẳng
	}
}
