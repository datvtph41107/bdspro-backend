package mapper

import (
	_enum "common/domain/enum"
	_utils "common/utils"

	"notification/internal/domain"
	"notification/internal/dto"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
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
		Action:      _enum.PropertyAction(req.Action),
		Description: req.Description,
		ActorID:     req.ActorId,
	}
}

func (m *PropertyHistoryMapper) SearchProtoToDTO(
	req *notificationpb.GetPropertyHistoryRequest,
) *dto.PropertyHistoryQueryDTO {
	q := &dto.PropertyHistoryQueryDTO{}
	if req != nil {
		if req.SubjectId != 0 {
			q.SubjectID = &req.SubjectId
		}

		if req.Action != 0 {
			action := _enum.PropertyAction(req.Action)
			q.Action = &action
		}

		if req.ActorId != 0 {
			q.ActorID = &req.ActorId
		}

		if req.Keyword != "" {
			q.Keyword = &req.Keyword
		}
	}

	q.Page = req.Page
	q.Size = req.Size
	q.Sort = req.Sort

	return q
}

func (m *PropertyHistoryMapper) EntityToItem(
	entity *domain.PropertyHistory,
	actor *sharepb.ProfileItem,
) *notificationpb.PropertyHistoryItem {
	if entity == nil {
		return nil
	}

	item := &notificationpb.PropertyHistoryItem{
		Id:          entity.ID,
		SubjectId:   entity.SubjectID,
		Action:      entity.Action.String(),
		Description: entity.Description,
		OccurredAt:  _utils.FormatTimeToString(entity.OccurredAt),
		CreatedAt:   _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(entity.UpdatedAt),
	}

	// nilpointer data k có actorId ; k cần check
	if actor != nil {
		item.Actor = &notificationpb.ActorInfo{
			Id:       actor.Id,
			FullName: actor.FullName,
			Avatar:   actor.Avatar,
		}
	}

	return item
}

func (m *PropertyHistoryMapper) EntitiesToListResponse(
	entities []*domain.PropertyHistory,
	userMap map[uint64]*sharepb.ProfileItem,
	total int64,
	page, size int32,
) *notificationpb.GetPropertyHistoryResponse {
	items := make([]*notificationpb.PropertyHistoryItem, 0, len(entities))

	for _, e := range entities {
		items = append(items,
			m.EntityToItem(e, userMap[e.ActorID]),
		)
	}

	return &notificationpb.GetPropertyHistoryResponse{
		Data:  items,
		Total: total,
	}
}

func (m *PropertyHistoryMapper) EntityToDetailResponse(
	entity *domain.PropertyHistory,
	actor *userpb.ProfileResponse,
) *notificationpb.GetPropertyHistoryByIDResponse {
	profileItem := &sharepb.ProfileItem{
		Id:       actor.ProfileId,
		FullName: actor.FullName,
		Avatar:   actor.Avatar,
	}
	return &notificationpb.GetPropertyHistoryByIDResponse{
		Data: m.EntityToItem(entity, profileItem),
	}
}
