package handler

import (
	"context"

	"hub/infra/mapper"
	"hub/internal/dto"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ClickURLHandler struct {
	hubpb.UnimplementedClickURLServiceServer
	ClickURLMapper  *mapper.ClickURLMapper
	ClickURLUsecase _usecase.IClickURLUsecase
}

func NewClickURLHandler(
	clickURLUsecase _usecase.IClickURLUsecase,
	clickURLMapper *mapper.ClickURLMapper,
) *ClickURLHandler {
	return &ClickURLHandler{
		ClickURLMapper:  clickURLMapper,
		ClickURLUsecase: clickURLUsecase,
	}
}

// @Summary Track click URL events
// @Description Track một hoặc nhiều sự kiện click link
// @Tags ClickURL
// @Accept json
// @Produce json
// @Param body body hubpb.TrackClickURLBatchRequest true "Danh sách click url cần track"
// @Success 200 {object} google.protobuf.Empty
// @Router /v2/hub/v2/click-url [post]
func (h *ClickURLHandler) TrackClickURL(
	ctx context.Context,
	req *hubpb.TrackClickURLBatchRequest,
) (*emptypb.Empty, error) {
	if req == nil || len(req.Events) == 0 {
		return &emptypb.Empty{}, nil
	}

	dtoReqs := make([]dto.TrackClickURLRequest, 0, len(req.Events))
	for _, pbReq := range req.Events {
		dtoReq := h.ClickURLMapper.MapTrackClickURLReq(pbReq)
		if dtoReq != nil {
			dtoReqs = append(dtoReqs, *dtoReq)
		}
	}

	if len(dtoReqs) == 0 {
		return &emptypb.Empty{}, nil
	}

	err := h.ClickURLUsecase.TrackBatch(ctx, &dto.TrackClickURLBatchRequest{
		Events: dtoReqs,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
