package handler

import (
	"context"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"

	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/enums"
	"organization/internal/usecase"
)

type GroupDealHandler struct {
	organizationpb.UnimplementedGroupDealServiceServer
	GroupDealUsecase     usecase.DealUsecase
	GroupDealTransformer transformer.GroupDealTransformer
	GroupDealValidator   validator.GroupDealValidator
	BdsproClient         *client.BdsproClient
	UserClient           *client.UserClient
}

func NewGroupDealHandler(
	groupDealUsecase usecase.DealUsecase,
	groupDealTransformer transformer.GroupDealTransformer,
	groupDealValidator validator.GroupDealValidator,
	bdsproClient *client.BdsproClient,
	userClient *client.UserClient,
) *GroupDealHandler {
	return &GroupDealHandler{
		GroupDealUsecase:     groupDealUsecase,
		GroupDealTransformer: groupDealTransformer,
		GroupDealValidator:   groupDealValidator,
		BdsproClient:         bdsproClient,
		UserClient:           userClient,
	}
}

// @Summary Tạo thương vụ
// @Description Tạo thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body organizationpb.CreateGroupDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.CreateGroupDealResponse "Thành công"
// @Router /deal [post]
func (h *GroupDealHandler) CreateGroupDeal(ctx context.Context, req *organizationpb.CreateGroupDealRequest) (*organizationpb.CreateGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateCreateGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal := h.GroupDealTransformer.CreateGroupDealRequestToEntity(req)

	deal, err := h.GroupDealUsecase.CreateDeal(ctx, deal)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToCreateGroupDealResponse(deal), nil
}

// @Summary Cập nhật thương vụ
// @Description Cập nhật thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request body organizationpb.UpdateGroupDealRequest true "Thông tin thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateGroupDealResponse "Thành công"
// @Router /deal/{id} [put]
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

// @Summary Xóa thương vụ
// @Description Xóa thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.DeleteGroupDealResponse "Thành công"
// @Router /deal/{id} [delete]
func (h *GroupDealHandler) DeleteGroupDeal(ctx context.Context, req *organizationpb.DeleteGroupDealRequest) (*organizationpb.DeleteGroupDealResponse, error) {
	if err := h.GroupDealValidator.ValidateDeleteGroupDealRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupDealUsecase.DeleteGroupDeal(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	return &organizationpb.DeleteGroupDealResponse{Id: req.Id}, nil
}

// @Summary Lấy thông tin thương vụ
// @Description Lấy thông tin thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetGroupDealResponse "Thành công"
// @Router /deal/{id} [get]
func (h *GroupDealHandler) GetGroupDeal(ctx context.Context, req *organizationpb.GetGroupDealRequest) (*organizationpb.GroupDeal, error) {
	if err := h.GroupDealValidator.ValidateGetGroupDealRequest(req); err != nil {
		return nil, err
	}

	deal, err := h.GroupDealUsecase.GetGroupDealByID(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	response := h.GroupDealTransformer.EntityToGetGroupDealResponse(deal)
	h.BdsproClient.MapProductToDealPb(ctx, []*organizationpb.GroupDeal{response})
	h.UserClient.MapMemberToDealPb(ctx, []*organizationpb.GroupDeal{response})
	return response, nil
}

// @Summary Lấy danh sách thương vụ
// @Description Lấy danh sách thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param request query organizationpb.GetGroupDealsRequest true "Thông tin lấy danh sách thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetGroupDealsResponse "Thành công"
// @Router /deal [get]
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
	deals, total, err := h.GroupDealUsecase.GetGroupDealsByGroupID(ctx, uint64(req.GroupId), page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupDealTransformer.EntityToGetGroupDealsResponse(deals, total), nil
}

// @Summary Cập nhật trạng thái thương vụ
// @Description Cập nhật trạng thái thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Param request body organizationpb.UpdateGroupDealStatusRequest true "Thông tin cập nhật trạng thái thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateGroupDealStatusResponse "Thành công"
// @Router /deal/{id}/status [put]
func (h *GroupDealHandler) UpdateGroupDealStatus(ctx context.Context, req *organizationpb.UpdateGroupDealStatusRequest) (*organizationpb.UpdateGroupDealStatusResponse, error) {
	if err := h.GroupDealValidator.ValidateUpdateGroupDealStatusRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupDealUsecase.UpdateGroupDealStatus(ctx, uint64(req.Id), enums.DealStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &organizationpb.UpdateGroupDealStatusResponse{Id: req.Id}, nil
}

// @Summary Lấy thông tin đầu tư thống kê khai báo thương vụ
// @Description Lấy thông tin đầu tư thống kê khai báo thương vụ
// @Tags Thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.InfoInvestmentResponse "Thành công"
// @Router /deal/{id}/investment [get]
func (h *GroupDealHandler) InfoInvestment(ctx context.Context, req *sharepb.IdDTO) (*organizationpb.InfoInvestmentResponse, error) {
	info, err := h.GroupDealUsecase.InfoInvestment(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &organizationpb.InfoInvestmentResponse{
		AmountPending:   info.AmountPending,
		AmountApproved:  info.AmountApproved,
		AmountRejected:  info.AmountRejected,
		AmountRest:      info.AmountRest,
		NumPending:      info.NumPending,
		NumApproved:     info.NumApproved,
		NumRejected:     info.NumRejected,
		NumInvestment:   info.NumInvestment,
		NumMemberSubmit: info.NumMemberSubmit,
		AmountTarget:    info.AmountTarget,
		PercentDone:     info.PercentDone,
	}, nil
}
