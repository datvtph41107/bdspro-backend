package handler

import (
	"context"

	"hub/infra/mapper"
	"hub/internal/dto"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"

	"google.golang.org/protobuf/types/known/emptypb"
)

type InteractiveEventHandler struct {
	hubpb.UnimplementedInteractiveEventServiceServer
	InteractiveEventMapper  *mapper.InteractiveEventMapper
	InteractiveEventUsecase _usecase.IInteractiveEventUsecase
}

func NewInteractiveEventHandler(
	interactiveEventUsecase _usecase.IInteractiveEventUsecase,
	interactiveEventMapper *mapper.InteractiveEventMapper,
) *InteractiveEventHandler {
	return &InteractiveEventHandler{
		InteractiveEventMapper:  interactiveEventMapper,
		InteractiveEventUsecase: interactiveEventUsecase,
	}
}

// @Summary Track interactive events
// @Description Track một hoặc nhiều interactive events
// @Tags InteractiveEvent
// @Accept json
// @Produce json
// @Param body body hubpb.TrackInteractiveEventBatchRequest true "Danh sách events cần track"
// @Success 200 {object} google.protobuf.Empty
// @Router /v2/hub/v2/interactive-event [post]
func (h *InteractiveEventHandler) TrackInteractiveEvent(
	ctx context.Context,
	req *hubpb.TrackInteractiveEventBatchRequest,
) (*emptypb.Empty, error) {
	if req == nil || len(req.Events) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Convert protobuf requests to DTOs
	dtoReqs := make([]*dto.TrackInteractiveEventRequest, 0, len(req.Events))
	for _, pbReq := range req.Events {
		dtoReq := h.InteractiveEventMapper.MapTrackInteractiveEventReq(pbReq)
		if dtoReq != nil {
			dtoReqs = append(dtoReqs, dtoReq)
		}
	}

	if len(dtoReqs) == 0 {
		return &emptypb.Empty{}, nil
	}

	// Create batch request
	batchReq := &dto.TrackInteractiveEventBatchRequest{
		Events: make([]dto.TrackInteractiveEventRequest, 0, len(dtoReqs)),
	}
	for _, dtoReq := range dtoReqs {
		batchReq.Events = append(batchReq.Events, *dtoReq)
	}

	// Track batch
	err := h.InteractiveEventUsecase.TrackBatch(ctx, batchReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
