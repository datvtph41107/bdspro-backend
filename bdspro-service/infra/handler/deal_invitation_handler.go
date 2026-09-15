package handler

import (
	"context"
	bdspropb "pb/types/bdspro"

	"common/fault"

	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/usecases"
)

type DealInvitationHandler struct {
	bdspropb.UnimplementedDealMemberServiceServer
	dealInvitationUsecase     usecases.DealInvitationUsecase
	dealInvitationTransformer mapper.DealInvitationTransformer
	dealInvitationValidator   validator.DealInvitationValidator
	userClient                *client.UserClient
	authClient                *client.AuthClient
}

func NewDealInvitationHandler(
	dealInvitationUsecase usecases.DealInvitationUsecase,
	dealInvitationTransformer mapper.DealInvitationTransformer,
	dealInvitationValidator validator.DealInvitationValidator,
	userClient *client.UserClient,
	authClient *client.AuthClient,
) *DealInvitationHandler {
	return &DealInvitationHandler{
		dealInvitationUsecase:     dealInvitationUsecase,
		dealInvitationTransformer: dealInvitationTransformer,
		dealInvitationValidator:   dealInvitationValidator,
		userClient:                userClient,
		authClient:                authClient,
	}
}

// @Summary Gửi lời mời tham gia thương vụ
// @Description Gửi lời mời tham gia thương vụ
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.SendDealInvitationRequest true "Thông tin lời mời"
// @Security BearerAuth
// @Success 200 {object} bdspropb.SendDealInvitationResponse "Thành công"
// @Router /deal-invitation [post]
func (h *DealInvitationHandler) SendDealInvitation(ctx context.Context, req *bdspropb.SendDealInvitationRequest) (*bdspropb.SendDealInvitationResponse, error) {
	if err := h.dealInvitationValidator.ValidateSendInvitationRequest(req); err != nil {
		return nil, err
	}

	invitation := h.dealInvitationTransformer.SendInvitationRequestToEntity(req)

	createdInvitation, err := h.dealInvitationUsecase.SendInvitation(ctx, invitation)
	if err != nil {
		return nil, mapDealInvitationSendError(err)
	}

	return h.dealInvitationTransformer.EntityToSendInvitationResponse(createdInvitation), nil
}

func mapDealInvitationSendError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := fault.As(err); !ok {
		return err
	}
	return fault.ToGRPC(err)
}

// @Summary Xác nhận lời mời
// @Description Xác nhận lời mời tham gia thương vụ
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.AcceptDealInvitationRequest true "Thông tin xác nhận"
// @Security BearerAuth
// @Success 200 {object} bdspropb.AcceptDealInvitationResponse "Thành công"
// @Router /deal-invitation/accept [post]
func (h *DealInvitationHandler) AcceptDealInvitation(ctx context.Context, req *bdspropb.AcceptDealInvitationRequest) (*bdspropb.AcceptDealInvitationResponse, error) {
	if err := h.dealInvitationValidator.ValidateAcceptInvitationRequest(req); err != nil {
		return nil, err
	}

	err := h.dealInvitationUsecase.AcceptInvitation(ctx, req.InvitationId)
	if err != nil {
		return nil, err
	}

	return &bdspropb.AcceptDealInvitationResponse{
		InvitationId: req.InvitationId,
		Success:      true,
	}, nil
}

// @Summary Từ chối lời mời
// @Description Từ chối lời mời tham gia thương vụ
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.RejectDealInvitationRequest true "Thông tin từ chối"
// @Security BearerAuth
// @Success 200 {object} bdspropb.RejectDealInvitationResponse "Thành công"
// @Router /deal-invitation/reject [post]
func (h *DealInvitationHandler) RejectDealInvitation(ctx context.Context, req *bdspropb.RejectDealInvitationRequest) (*bdspropb.RejectDealInvitationResponse, error) {
	if err := h.dealInvitationValidator.ValidateRejectInvitationRequest(req); err != nil {
		return nil, err
	}

	err := h.dealInvitationUsecase.RejectInvitation(ctx, req.InvitationId, req.Reason)
	if err != nil {
		return nil, err
	}

	return &bdspropb.RejectDealInvitationResponse{
		InvitationId: req.InvitationId,
		Success:      true,
	}, nil
}

// @Summary Gửi lại lời mời
// @Description Gửi lại lời mời tham gia thương vụ
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.ResendDealInvitationRequest true "Thông tin gửi lại"
// @Security BearerAuth
// @Success 200 {object} bdspropb.ResendDealInvitationResponse "Thành công"
// @Router /deal-invitation/resend [post]
func (h *DealInvitationHandler) ResendDealInvitation(ctx context.Context, req *bdspropb.ResendDealInvitationRequest) (*bdspropb.ResendDealInvitationResponse, error) {
	if err := h.dealInvitationValidator.ValidateResendInvitationRequest(req); err != nil {
		return nil, err
	}

	updatedInvitation, err := h.dealInvitationUsecase.ResendInvitation(ctx, req.InvitationId)
	if err != nil {
		return nil, err
	}

	return h.dealInvitationTransformer.EntityToResendInvitationResponse(updatedInvitation), nil
}

// @Summary Lấy danh sách thành viên đã accept
// @Description Lấy danh sách thành viên đã chấp nhận lời mời tham gia thương vụ
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param deal_id path int true "ID của thương vụ"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetAcceptedMembersResponse "Thành công"
// @Router /deal/{deal_id}/accepted-members [get]
func (h *DealInvitationHandler) GetAcceptedMembers(ctx context.Context, req *bdspropb.GetAcceptedMembersRequest) (*bdspropb.GetAcceptedMembersResponse, error) {
	if err := h.dealInvitationValidator.ValidateGetAcceptedMembersRequest(req); err != nil {
		return nil, err
	}

	members, err := h.dealInvitationUsecase.GetAcceptedMembers(ctx, req.DealId)
	if err != nil {
		return nil, err
	}

	return h.dealInvitationTransformer.EntitiesToGetAcceptedMembersResponse(members), nil
}

// @Summary Rút khỏi thương vụ
// @Description Rút khỏi thương vụ (dành cho thành viên)
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.WithdrawFromDealRequest true "Thông tin rút khỏi"
// @Security BearerAuth
// @Success 200 {object} bdspropb.WithdrawFromDealResponse "Thành công"
// @Router /deal-invitation/withdraw [post]
func (h *DealInvitationHandler) WithdrawFromDeal(ctx context.Context, req *bdspropb.WithdrawFromDealRequest) (*bdspropb.WithdrawFromDealResponse, error) {
	if err := h.dealInvitationValidator.ValidateWithdrawFromDealRequest(req); err != nil {
		return nil, err
	}

	err := h.dealInvitationUsecase.WithdrawFromDeal(ctx, req.InvitationId, req.Reason)
	if err != nil {
		return nil, err
	}

	return &bdspropb.WithdrawFromDealResponse{
		InvitationId: req.InvitationId,
		Success:      true,
	}, nil
}

// @Summary Gỡ khỏi thương vụ
// @Description Gỡ khỏi thương vụ (dành cho admin/phụ trách)
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.RemoveFromDealRequest true "Thông tin gỡ khỏi"
// @Security BearerAuth
// @Success 200 {object} bdspropb.RemoveFromDealResponse "Thành công"
// @Router /deal-invitation/remove [post]
func (h *DealInvitationHandler) RemoveFromDeal(ctx context.Context, req *bdspropb.RemoveFromDealRequest) (*bdspropb.RemoveFromDealResponse, error) {
	if err := h.dealInvitationValidator.ValidateRemoveFromDealRequest(req); err != nil {
		return nil, err
	}

	err := h.dealInvitationUsecase.RemoveFromDeal(ctx, req.InvitationId, req.Reason)
	if err != nil {
		return nil, err
	}

	return &bdspropb.RemoveFromDealResponse{
		InvitationId: req.InvitationId,
		Success:      true,
	}, nil
}

// @Summary Xác nhận rút khỏi thương vụ
// @Description Xác nhận rút khỏi thương vụ (dành cho admin/phụ trách)
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request body bdspropb.ConfirmWithdrawalRequest true "Thông tin xác nhận"
// @Security BearerAuth
// @Success 200 {object} bdspropb.ConfirmWithdrawalResponse "Thành công"
// @Router /deal-invitation/confirm-withdrawal [post]
func (h *DealInvitationHandler) ConfirmWithdrawal(ctx context.Context, req *bdspropb.ConfirmWithdrawalRequest) (*bdspropb.ConfirmWithdrawalResponse, error) {
	if err := h.dealInvitationValidator.ValidateConfirmWithdrawalRequest(req); err != nil {
		return nil, err
	}

	err := h.dealInvitationUsecase.ConfirmWithdrawal(ctx, req.InvitationId)
	if err != nil {
		return nil, err
	}

	return &bdspropb.ConfirmWithdrawalResponse{
		InvitationId: req.InvitationId,
		Success:      true,
	}, nil
}

// @Summary Tìm kiếm thành viên
// @Description Tìm kiếm thành viên + sort theo thương vụ chung
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request query bdspropb.SearchDealMembersRequest true "Thông tin tìm kiếm"
// @Security BearerAuth
// @Success 200 {object} bdspropb.SearchDealMembersResponse "Thành công"
// @Router /deal-invitation/search [get]
func (h *DealInvitationHandler) SearchMembers(ctx context.Context, req *bdspropb.SearchDealMembersRequest) (*bdspropb.SearchDealMembersResponse, error) {
	if err := h.dealInvitationValidator.ValidateSearchMembersRequest(req); err != nil {
		return nil, err
	}

	page := int(req.Page)
	size := int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	members, total, err := h.dealInvitationUsecase.SearchMembers(ctx, req.DealId, req.Keyword, page, size)
	if err != nil {
		return nil, err
	}

	return h.dealInvitationTransformer.EntitiesToSearchMembersResponse(members, total), nil
}

// @Summary Lấy lời mời theo ID
// @Description Lấy thông tin lời mời theo ID
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param id path int true "ID của lời mời"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealInvitationResponse "Thành công"
// @Router /deal-invitation/{id} [get]
func (h *DealInvitationHandler) GetDealInvitation(ctx context.Context, req *bdspropb.GetDealInvitationRequest) (*bdspropb.GetDealInvitationResponse, error) {
	if err := h.dealInvitationValidator.ValidateGetInvitationRequest(req); err != nil {
		return nil, err
	}

	invitation, err := h.dealInvitationUsecase.GetInvitationByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.dealInvitationTransformer.EntityToGetInvitationResponse(invitation), nil
}

// @Summary Lấy danh sách lời mời đang chờ
// @Description Lấy danh sách lời mời đang chờ xác nhận
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request query bdspropb.GetPendingInvitationsRequest true "Thông tin lấy danh sách"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetPendingInvitationsResponse "Thành công"
// @Router /deal-invitation/pending [get]
func (h *DealInvitationHandler) GetPendingInvitations(ctx context.Context, req *bdspropb.GetPendingInvitationsRequest) (*bdspropb.GetPendingInvitationsResponse, error) {
	if err := h.dealInvitationValidator.ValidateGetPendingInvitationsRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10

	if req.Page != 0 {
		page = int(req.Page)
	}
	if req.Size != 0 {
		size = int(req.Size)
	}

	invitations, total, err := h.dealInvitationUsecase.GetPendingInvitations(ctx, req.MemberId, page, size)
	if err != nil {
		return nil, err
	}

	return h.dealInvitationTransformer.EntitiesToGetPendingInvitationsResponse(invitations, total), nil
}

// @Summary Lấy danh sách lời mời theo thương vụ
// @Description Lấy danh sách lời mời theo thương vụ
// @Tags Lời mời thương vụ
// @Accept json
// @Produce json
// @Param request query bdspropb.GetDealInvitationsByDealRequest true "Thông tin lấy danh sách"
// @Security BearerAuth
// @Success 200 {object} bdspropb.GetDealInvitationsByDealResponse "Thành công"
// @Router /deal/{deal_id}/invitations [get]
func (h *DealInvitationHandler) GetDealInvitationsByDeal(ctx context.Context, req *bdspropb.GetDealInvitationsByDealRequest) (*bdspropb.GetDealInvitationsByDealResponse, error) {
	if err := h.dealInvitationValidator.ValidateGetInvitationsByDealRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10

	if req.Page != 0 {
		page = int(req.Page)
	}
	if req.Size != 0 {
		size = int(req.Size)
	}

	invitations, total, err := h.dealInvitationUsecase.GetInvitationsByDealID(ctx, req.DealId, page, size)
	if err != nil {
		return nil, err
	}

	pbs := h.dealInvitationTransformer.EntitiesToGetInvitationsByDealResponse(invitations, total)
	h.userClient.MapToDealInvitationPb(ctx, pbs.Data)
	h.authClient.MapRoleToDealInvitationPb(ctx, pbs.Data)
	return pbs, nil
}
