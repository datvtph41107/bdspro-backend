package mapper

import (
	"notification/internal/domain"
	notificationpb "pb/types/notification"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EntityToCreateDealHistoryRequest chuyển đổi từ entity sang proto request
func EntityToCreateDealHistoryRequest(entity *domain.DealHistoryEntity) *notificationpb.CreateDealHistoryRequest {
	return &notificationpb.CreateDealHistoryRequest{
		DealId:      entity.DealID,
		ActorId:     entity.ActorID,
		ActorName:   entity.ActorName,
		ActorAvatar: entity.ActorAvatar,
		ActionType:  entity.ActionType,
		ActionName:  entity.ActionName,
		Content:     entity.Content,
		Metadata:    convertMapToStruct(entity.Metadata),
	}
}

// RequestToEntity chuyển đổi từ proto request sang entity
func RequestToEntity(req *notificationpb.CreateDealHistoryRequest) *domain.DealHistoryEntity {
	return &domain.DealHistoryEntity{
		DealID:      req.DealId,
		ActorID:     req.ActorId,
		ActorName:   req.ActorName,
		ActorAvatar: req.ActorAvatar,
		ActionType:  req.ActionType,
		ActionName:  req.ActionName,
		Content:     req.Content,
		Metadata:    convertStructToMap(req.Metadata),
	}
}

// EntityToDealHistoryItem chuyển đổi từ entity sang proto item
func EntityToDealHistoryItem(entity *domain.DealHistoryEntity) *notificationpb.DealHistoryItem {
	return &notificationpb.DealHistoryItem{
		Id:          entity.ID,
		DealId:      entity.DealID,
		ActorId:     entity.ActorID,
		ActorName:   entity.ActorName,
		ActorAvatar: entity.ActorAvatar,
		ActionType:  entity.ActionType,
		ActionName:  entity.ActionName,
		Content:     entity.Content,
		Metadata:    convertMapToStruct(entity.Metadata),
		CreatedAt:   convertTimeToTimestamp(entity.CreatedAt),
		UpdatedAt:   convertTimeToTimestamp(entity.UpdatedAt),
	}
}

// EntitiesToDealHistoryResponse chuyển đổi từ entities sang proto response
func EntitiesToDealHistoryResponse(entities []*domain.DealHistoryEntity, total uint32, page, size int32) *notificationpb.GetDealHistoryResponse {
	items := make([]*notificationpb.DealHistoryItem, len(entities))
	for i, entity := range entities {
		items[i] = EntityToDealHistoryItem(entity)
	}

	return &notificationpb.GetDealHistoryResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  size,
	}
}

// EntityToDealHistoryByIDResponse chuyển đổi từ entity sang proto response
func EntityToDealHistoryByIDResponse(entity *domain.DealHistoryEntity) *notificationpb.GetDealHistoryByIDResponse {
	return &notificationpb.GetDealHistoryByIDResponse{
		Data: EntityToDealHistoryItem(entity),
	}
}

// ResponseToEntity chuyển đổi từ proto response sang entity
func ResponseToEntity(resp *notificationpb.CreateDealHistoryResponse) *domain.DealHistoryResponse {
	return &domain.DealHistoryResponse{
		Success: resp.Success,
		Message: resp.Message,
		ID:      resp.Id,
	}
}

// Helper functions
func convertMapToStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return &structpb.Struct{}
	}

	// Convert map to structpb.Struct
	// This is a simplified conversion - you might need more sophisticated handling
	return &structpb.Struct{}
}

func convertStructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return make(map[string]interface{})
	}
	return s.AsMap()
}

func convertTimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

// RequestToDealHistoryRequest chuyển đổi từ proto request sang domain request
func RequestToDealHistoryRequest(req *notificationpb.CreateDealHistoryRequest) *domain.DealHistoryRequest {
	return &domain.DealHistoryRequest{
		DealID:      req.DealId,
		ActorID:     req.ActorId,
		ActorName:   req.ActorName,
		ActorAvatar: req.ActorAvatar,
		ActionType:  req.ActionType,
		ActionName:  req.ActionName,
		Content:     req.Content,
		Metadata:    convertStructToMap(req.Metadata),
	}
}
