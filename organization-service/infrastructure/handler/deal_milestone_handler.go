package handler

import (
	"context"
	"organization/infrastructure/transformer"
	"organization/internal/usecase"
	organizationpb "pb/types/organization"
)

type DealMilestoneHandler struct {
	organizationpb.UnimplementedDealMilestoneServiceServer
	milestoneUsecase    usecase.DealMilestoneUsecase
	milestoneTransformer transformer.DealMilestoneTransformer
}

func NewDealMilestoneHandler(
	milestoneUsecase usecase.DealMilestoneUsecase,
	milestoneTransformer transformer.DealMilestoneTransformer,
) *DealMilestoneHandler {
	return &DealMilestoneHandler{
		milestoneUsecase:    milestoneUsecase,
		milestoneTransformer: milestoneTransformer,
	}
}

func (h *DealMilestoneHandler) CreateDealMilestone(ctx context.Context, req *organizationpb.CreateDealMilestoneRequest) (*organizationpb.DealMilestone, error) {
	milestone := h.milestoneTransformer.CreateRequestToEntity(req)
	createdMilestone, err := h.milestoneUsecase.CreateMilestone(ctx, milestone)
	if err != nil {
		return nil, err
	}
	return h.milestoneTransformer.EntityToProto(createdMilestone), nil
}

func (h *DealMilestoneHandler) UpdateDealMilestone(ctx context.Context, req *organizationpb.UpdateDealMilestoneRequest) (*organizationpb.DealMilestone, error) {
	milestone := h.milestoneTransformer.UpdateRequestToEntity(req)
	updatedMilestone, err := h.milestoneUsecase.UpdateMilestone(ctx, req.Id, milestone)
	if err != nil {
		return nil, err
	}
	return h.milestoneTransformer.EntityToProto(updatedMilestone), nil
}

func (h *DealMilestoneHandler) DeleteDealMilestone(ctx context.Context, req *organizationpb.DeleteDealMilestoneRequest) (*organizationpb.DeleteDealMilestoneResponse, error) {
	err := h.milestoneUsecase.DeleteMilestone(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeleteDealMilestoneResponse{Id: req.Id}, nil
}

func (h *DealMilestoneHandler) GetDealMilestones(ctx context.Context, req *organizationpb.GetDealMilestonesRequest) (*organizationpb.GetDealMilestonesResponse, error) {
	milestones, err := h.milestoneUsecase.GetMilestonesByDealID(ctx, req.DealId)
	if err != nil {
		return nil, err
	}

	response := &organizationpb.GetDealMilestonesResponse{
		Milestones: make([]*organizationpb.DealMilestone, len(milestones)),
	}
	for i, milestone := range milestones {
		response.Milestones[i] = h.milestoneTransformer.EntityToProto(milestone)
	}
	return response, nil
}

func (h *DealMilestoneHandler) UpdateMilestoneOrder(ctx context.Context, req *organizationpb.UpdateMilestoneOrderRequest) (*organizationpb.UpdateMilestoneOrderResponse, error) {
	err := h.milestoneUsecase.UpdateMilestoneOrder(ctx, req.DealId, req.MilestoneIds)
	if err != nil {
		return nil, err
	}
	return &organizationpb.UpdateMilestoneOrderResponse{Success: true}, nil
} 