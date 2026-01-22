package handler

import (
	"bdspro/infra/mapper"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	bdspropb "pb/types/bdspro"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ObservationHandler struct {
	bdspropb.UnimplementedObservationServiceServer
	ObservationMapper  *mapper.ObservationMapper
	ObservationUsecase *shared_usecase.ObservationUsecase
}

func NewObservationHandler(
	observationUsecase *shared_usecase.ObservationUsecase,
	observationMapper *mapper.ObservationMapper,
) *ObservationHandler {
	return &ObservationHandler{
		ObservationMapper:  observationMapper,
		ObservationUsecase: observationUsecase,
	}
}

func (h *ObservationHandler) TrackObservation(
	ctx context.Context,
	req *bdspropb.TrackObservationRequest,
) (*emptypb.Empty, error) {

	dtoReq := h.ObservationMapper.MapTrackObservationReq(req)

	err := h.ObservationUsecase.Track(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *ObservationHandler) GetProductStats(
	ctx context.Context,
	req *bdspropb.GetProductStatsRequest,
) (*bdspropb.ProductStatsResponse, error) {

	stats, err := h.ObservationUsecase.GetProductStats(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return h.ObservationMapper.MapProductStatsToPB(stats), nil
}
