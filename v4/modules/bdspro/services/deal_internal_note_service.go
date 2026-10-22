package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
	"mdn4/modules/bdspro/repository"
	bdsutils "mdn4/modules/bdspro/utils"
)

// DealInternalNoteService adds internal notes on a deal (migrated from DealInternalNoteService).
type DealInternalNoteService interface {
	Add(ctx context.Context, req dto.AddDealInternalNoteRequest) (*bdspropb.AddInternalNoteResponse, error)
}

// @bind: mdn4/modules/bdspro/services.DealInternalNoteService
type LocalDealInternalNoteService struct {
	deal    repository.DealRepository
	notes   repository.DealInternalNoteRepository
	history *bdsutils.DealHistoryLogger
}

func NewLocalDealInternalNoteService(
	deal repository.DealRepository,
	notes repository.DealInternalNoteRepository,
	history *bdsutils.DealHistoryLogger,
) *LocalDealInternalNoteService {
	return &LocalDealInternalNoteService{deal: deal, notes: notes, history: history}
}

func (s *LocalDealInternalNoteService) Add(ctx context.Context, req dto.AddDealInternalNoteRequest) (*bdspropb.AddInternalNoteResponse, error) {
	dealID := req.DealID
	if dealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if req.Content == "" {
		return nil, _errors.BadRequestException("content is required")
	}
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ForbiddenException("profile is required")
	}

	actionType := enums.DealActionType(req.ActionType)
	if actionType == 0 {
		actionType = enums.ActionTypeInternalNoteAdded
	}
	if actionType != enums.ActionTypeInternalNoteAdded && actionType != enums.ActionTypeInternalNoteEdited {
		return nil, _errors.BadRequestException("invalid action type for internal note")
	}

	deal, err := s.deal.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, _errors.NotFoundException("deal not found")
	}
	ok, err := s.deal.IsAcceptedMember(ctx, dealID, profileID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, _errors.ForbiddenException("you are not allowed to access this deal")
	}

	created, err := s.notes.Create(ctx, &entity.DealInternalNote{
		DealID:     dealID,
		ActorID:    profileID,
		Content:    req.Content,
		ActionType: actionType,
	})
	if err != nil {
		return nil, err
	}

	historyEvent := enums.DealHistoryEventAddInternalNote
	if actionType == enums.ActionTypeInternalNoteEdited {
		historyEvent = enums.DealHistoryEventEditInternalNote
	}
	if s.history != nil {
		actionName := enums.DealActionTypeMap[actionType]
		if actionName == "" {
			actionName = historyEvent.GetEventName()
		}
		s.history.LogDealHistory(ctx, dealID, historyEvent,
			fmt.Sprintf("%s: %s", actionName, req.Content),
			map[string]interface{}{
				"deal_id":     dealID,
				"deal_name":   deal.Name,
				"note_id":     created.ID,
				"action_type": uint32(actionType),
				"content":     req.Content,
			})
	}

	return &bdspropb.AddInternalNoteResponse{
		DealId:     created.DealID,
		NoteId:     created.ID,
		Content:    created.Content,
		ActionType: uint32(created.ActionType),
		CreatedAt:  created.CreatedAt.Format(time.RFC3339),
	}, nil
}
