package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"
)

type DealMilestoneHandler struct {
	bdspropb.UnimplementedDealMilestoneServiceServer
	milestoneUsecase     usecases.DealMilestoneUsecase
	milestoneTransformer mapper.DealMilestoneTransformer
}

func NewDealMilestoneHandler(
	milestoneUsecase usecases.DealMilestoneUsecase,
	milestoneTransformer mapper.DealMilestoneTransformer,
) *DealMilestoneHandler {
	return &DealMilestoneHandler{
		milestoneUsecase:     milestoneUsecase,
		milestoneTransformer: milestoneTransformer,
	}
}

func (h *DealMilestoneHandler) CreateDealMilestone(ctx context.Context, req *bdspropb.CreateDealMilestoneRequest) (*bdspropb.DealMilestone, error) {
	milestone := h.milestoneTransformer.CreateRequestToEntity(req)
	createdMilestone, err := h.milestoneUsecase.CreateMilestone(ctx, milestone)
	if err != nil {
		return nil, err
	}
	return h.milestoneTransformer.EntityToProto(createdMilestone), nil
}

func (h *DealMilestoneHandler) UpdateDealMilestone(ctx context.Context, req *bdspropb.UpdateDealMilestoneRequest) (*bdspropb.DealMilestone, error) {
	milestone := h.milestoneTransformer.UpdateRequestToEntity(req)
	updatedMilestone, err := h.milestoneUsecase.UpdateMilestone(ctx, req.Id, milestone)
	if err != nil {
		return nil, err
	}
	return h.milestoneTransformer.EntityToProto(updatedMilestone), nil
}

func (h *DealMilestoneHandler) DeleteDealMilestone(ctx context.Context, req *bdspropb.DeleteDealMilestoneRequest) (*bdspropb.DeleteDealMilestoneResponse, error) {
	err := h.milestoneUsecase.DeleteMilestone(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.DeleteDealMilestoneResponse{Id: req.Id}, nil
}

func (h *DealMilestoneHandler) GetDealMilestones(ctx context.Context, req *bdspropb.GetDealMilestonesRequest) (*bdspropb.GetDealMilestonesResponse, error) {
	milestones, err := h.milestoneUsecase.GetMilestonesByDealID(ctx, req.DealId)
	if err != nil {
		return nil, err
	}

	response := &bdspropb.GetDealMilestonesResponse{
		Milestones: make([]*bdspropb.DealMilestone, len(milestones)),
	}
	for i, milestone := range milestones {
		response.Milestones[i] = h.milestoneTransformer.EntityToProto(milestone)
	}
	return response, nil
}

func (h *DealMilestoneHandler) UpdateMilestoneOrder(ctx context.Context, req *bdspropb.UpdateMilestoneOrderRequest) (*bdspropb.UpdateMilestoneOrderResponse, error) {
	err := h.milestoneUsecase.UpdateMilestoneOrder(ctx, req.DealId, req.MilestoneIds)
	if err != nil {
		return nil, err
	}
	return &bdspropb.UpdateMilestoneOrderResponse{Success: true}, nil
}
