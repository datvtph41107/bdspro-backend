package transformer

import (
	organizationpb "pb/types/organization"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/enums"
)

type InternalNoteTransformer interface {
	TransformAddInternalNoteRequest(req *organizationpb.AddInternalNoteRequest) (string, enums.ActionType)
	TransformAddInternalNoteResponse(note *entity.InternalNote) *organizationpb.AddInternalNoteResponse
	TransformGetDealHistoryRequest(req *organizationpb.GetDealHistoryRequest) ([]enums.ActionType, uint32, uint32)
	TransformGetDealHistoryResponse(dealID uint64, notes []*entity.InternalNote, total, page, limit uint32) *organizationpb.GetDealHistoryResponse
}

type internalNoteTransformer struct{}

func NewInternalNoteTransformer() InternalNoteTransformer {
	return &internalNoteTransformer{}
}

func (t *internalNoteTransformer) TransformAddInternalNoteRequest(req *organizationpb.AddInternalNoteRequest) (string, enums.ActionType) {
	actionType := enums.ActionType(req.ActionType)
	if actionType == 0 {
		actionType = enums.ActionTypeInternalNoteAdded // Default
	}
	return req.Content, actionType
}

func (t *internalNoteTransformer) TransformAddInternalNoteResponse(note *entity.InternalNote) *organizationpb.AddInternalNoteResponse {
	return &organizationpb.AddInternalNoteResponse{
		DealId:     note.DealID,
		NoteId:     note.ID,
		Content:    note.Content,
		ActionType: uint32(note.ActionType),
		CreatedAt:  note.CreatedAt.Format(time.RFC3339),
	}
}

func (t *internalNoteTransformer) TransformGetDealHistoryRequest(req *organizationpb.GetDealHistoryRequest) ([]enums.ActionType, uint32, uint32) {
	// Transform action types
	actionTypes := make([]enums.ActionType, len(req.ActionTypes))
	for i, actionType := range req.ActionTypes {
		actionTypes[i] = enums.ActionType(actionType)
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

func (t *internalNoteTransformer) TransformGetDealHistoryResponse(dealID uint64, notes []*entity.InternalNote, total, page, limit uint32) *organizationpb.GetDealHistoryResponse {
	items := make([]*organizationpb.DealHistoryItem, len(notes))
	for i, note := range notes {
		items[i] = &organizationpb.DealHistoryItem{
			Id:          note.ID,
			DealId:      note.DealID,
			ActorId:     note.ActorID,
			ActorName:   "", // Sẽ được populate từ User service
			ActorAvatar: "", // Sẽ được populate từ User service
			ActionType:  uint32(note.ActionType),
			ActionName:  enums.ActionTypeMap[note.ActionType],
			Content:     note.Content,
			CreatedAt:   note.CreatedAt.Format(time.RFC3339),
		}
	}

	return &organizationpb.GetDealHistoryResponse{
		Data:  items,
		Total: total,
	}
}
