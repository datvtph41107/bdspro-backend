package usecase

import (
	"context"

	_err "common/domain/err"
	_utils "common/utils"
	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type IInteractiveEventUsecase interface {
	Track(ctx context.Context, req *dto.TrackInteractiveEventRequest) *_err.ErrorDTO
	TrackBatch(ctx context.Context, req *dto.TrackInteractiveEventBatchRequest) *_err.ErrorDTO
}

type InteractiveEventUsecase struct {
	InteractiveEventRepo repo.IInteractiveEventRepo
}

func NewInteractiveEventUsecase(
	interactiveEventRepo repo.IInteractiveEventRepo,
) IInteractiveEventUsecase {
	return &InteractiveEventUsecase{
		InteractiveEventRepo: interactiveEventRepo,
	}
}



func (u *InteractiveEventUsecase) Track(
	ctx context.Context,
	req *dto.TrackInteractiveEventRequest,
) *_err.ErrorDTO {
	if req == nil {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Request không hợp lệ",
		}
	}

	var profileID *uint64
	if profileIDVal := _utils.GetProfileIdWithContext(ctx); profileIDVal != 0 {
		profileID = &profileIDVal
	}

	event := &domain.InteractiveEvent{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Event:     req.Event,
		RefID:     req.RefID,
		Screen:    req.Screen,
		Duration:  req.Duration,
		ProfileID: profileID,
	}

	if err := u.InteractiveEventRepo.Insert(ctx, event); err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lưu interactive event: " + err.Error(),
		}
	}

	return nil
}

func (u *InteractiveEventUsecase) TrackBatch(
	ctx context.Context,
	req *dto.TrackInteractiveEventBatchRequest,
) *_err.ErrorDTO {
	if req == nil || len(req.Events) == 0 {
		return &_err.ErrorDTO{
			Code:    400,
			Message: "Request không hợp lệ hoặc danh sách events rỗng",
		}
	}

	var profileID *uint64
	if profileIDVal := _utils.GetProfileIdWithContext(ctx); profileIDVal != 0 {
		profileID = &profileIDVal
	}

	events := make([]*domain.InteractiveEvent, 0, len(req.Events))
	for _, dtoReq := range req.Events {
		event := &domain.InteractiveEvent{
			StartTime: dtoReq.StartTime,
			EndTime:   dtoReq.EndTime,
			Event:     dtoReq.Event,
			RefID:     dtoReq.RefID,
			Screen:    dtoReq.Screen,
			Duration:  dtoReq.Duration,
			ProfileID: profileID,
		}
		events = append(events, event)
	}

	if err := u.InteractiveEventRepo.InsertBatch(ctx, events); err != nil {
		return &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lưu batch interactive events: " + err.Error(),
		}
	}

	return nil
}
