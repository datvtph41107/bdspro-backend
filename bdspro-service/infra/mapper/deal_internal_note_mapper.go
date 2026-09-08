package mapper

import (
	bdspropb "pb/types/bdspro"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/enums"
)

type InternalNoteMapper interface {
	TransformAddInternalNoteRequest(req *bdspropb.AddInternalNoteRequest) (string, enums.DealActionType)
	TransformAddInternalNoteResponse(note *domain.DealInternalNote) *bdspropb.AddInternalNoteResponse
	TransformGetDealHistoryRequest(req *bdspropb.GetDealHistoryRequest) ([]enums.DealActionType, uint32, uint32)
	TransformGetDealHistoryResponse(dealID uint64, notes []*domain.DealInternalNote, total, page, limit uint32) *bdspropb.GetDealHistoryResponse
}

type internalNoteTransformer struct{}

func NewInternalNoteTransformer() InternalNoteMapper {
	return &internalNoteTransformer{}
}

func (t *internalNoteTransformer) TransformAddInternalNoteRequest(req *bdspropb.AddInternalNoteRequest) (string, enums.DealActionType) {
	actionType := enums.DealActionType(req.ActionType)
	if actionType == 0 {
		actionType = enums.ActionTypeInternalNoteAdded // Default
	}
	return req.Content, actionType
}

func (t *internalNoteTransformer) TransformAddInternalNoteResponse(note *domain.DealInternalNote) *bdspropb.AddInternalNoteResponse {
	return &bdspropb.AddInternalNoteResponse{
		DealId:     note.DealID,
		NoteId:     note.ID,
		Content:    note.Content,
		ActionType: uint32(note.ActionType),
		CreatedAt:  note.CreatedAt.Format(time.RFC3339),
	}
}

func (t *internalNoteTransformer) TransformGetDealHistoryRequest(req *bdspropb.GetDealHistoryRequest) ([]enums.DealActionType, uint32, uint32) {
	// Transform action types
	actionTypes := make([]enums.DealActionType, len(req.ActionTypes))
	for i, actionType := range req.ActionTypes {
		actionTypes[i] = enums.DealActionType(actionType)
	}

	// Set default values for pagination
	page := req.Page
	if page == 0 {
		page = 1
	}
	size := req.Size
	if size == 0 {
		size = 10
	}

	return actionTypes, page, size
}

func (t *internalNoteTransformer) TransformGetDealHistoryResponse(dealID uint64, notes []*domain.DealInternalNote, total, page, limit uint32) *bdspropb.GetDealHistoryResponse {
	items := make([]*bdspropb.DealHistoryItem, len(notes))
	for i, note := range notes {
		items[i] = &bdspropb.DealHistoryItem{
			Id:          note.ID,
			DealId:      note.DealID,
			ActorId:     note.ActorID,
			ActorName:   "", // Sẽ được populate từ User service
			ActorAvatar: "", // Sẽ được populate từ User service
			ActionType:  uint32(note.ActionType),
			ActionName:  enums.DealActionTypeMap[note.ActionType],
			Content:     note.Content,
			CreatedAt:   note.CreatedAt.Format(time.RFC3339),
		}
	}

	return &bdspropb.GetDealHistoryResponse{
		Data:  items,
		Total: total,
	}
}
