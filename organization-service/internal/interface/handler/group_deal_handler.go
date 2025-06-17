package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

type GroupDealHandler struct {
	organizationpb.UnimplementedGroupDealServiceServer
	GroupDealUsecase     usecase.GroupDealUsecase
	GroupDealTransformer transformer.GroupDealTransformer
	GroupDealValidator   validator.GroupDealValidator
}

func NewGroupDealHandler(
	groupDealUsecase usecase.GroupDealUsecase,
	groupDealTransformer transformer.GroupDealTransformer,
	groupDealValidator validator.GroupDealValidator,
) *GroupDealHandler {
	return &GroupDealHandler{
		GroupDealUsecase:     groupDealUsecase,
		GroupDealTransformer: groupDealTransformer,
		GroupDealValidator:   groupDealValidator,
	}
}

func (h *GroupDealHandler) CreateGroupDeal(ctx context.Context, req *organizationpb.CreateGroupDealRequest) (*organizationpb.CreateGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateCreateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.GroupDealTransformer.CreateGroupDealRequestToEntity(req)

	deal, err := h.GroupDealUsecase.CreateGroupDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToCreateGroupDealResponse(deal), nil
}

func (h *GroupDealHandler) UpdateGroupDeal(ctx context.Context, req *organizationpb.UpdateGroupDealRequest) (*organizationpb.UpdateGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateUpdateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.GroupDealTransformer.UpdateGroupDealRequestToEntity(req)

	deal, err := h.GroupDealUsecase.UpdateGroupDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToUpdateGroupDealResponse(deal), nil
}

func (h *GroupDealHandler) DeleteGroupDeal(ctx context.Context, req *organizationpb.DeleteGroupDealRequest) (*organizationpb.DeleteGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateDeleteGroupDealRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupDealUsecase.DeleteGroupDeal(ctx, uint32(req.Id))
	if err != nil {
		return nil, err
	}

	return &organizationpb.DeleteGroupDealResponse{Id: req.Id}, nil
}

func (h *GroupDealHandler) GetGroupDeal(ctx context.Context, req *organizationpb.GetGroupDealRequest) (*organizationpb.GetGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateGetGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal, err := h.GroupDealUsecase.GetGroupDealByID(ctx, uint32(req.Id))
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToGetGroupDealResponse(deal), nil
}

func (h *GroupDealHandler) GetGroupDeals(ctx context.Context, req *organizationpb.GetGroupDealsRequest) (*organizationpb.GetGroupDealsResponse, error) {
	if err := h.GroupDealValidator.ValidateGetGroupDealsRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}
	deals, total, err := h.GroupDealUsecase.GetGroupDealsByGroupID(ctx, req.GroupId, page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToGetGroupDealsResponse(deals, total), nil
}

func (h *GroupDealHandler) UpdateGroupDealStatus(ctx context.Context, req *organizationpb.UpdateGroupDealStatusRequest) (*organizationpb.UpdateGroupDealStatusResponse, error) {
	if err := h.GroupDealValidator.ValidateUpdateGroupDealStatusRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupDealUsecase.UpdateGroupDealStatus(ctx, uint32(req.Id), entity.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &organizationpb.UpdateGroupDealStatusResponse{Id: req.Id}, nil
}
