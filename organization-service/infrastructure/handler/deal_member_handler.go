package handler

// import (
// 	"context"
// 	organizationpb "pb/types/organization"

// 	"organization/infrastructure/client"
// 	"organization/infrastructure/transformer"
// 	"organization/infrastructure/validator"
// 	"organization/internal/usecase"
// )

// type DealInvitationHandler struct {
// 	organizationpb.UnimplementedDealMemberServiceServer
// 	dealInvitationUsecase     usecase.DealInvitationUsecase
// 	dealInvitationTransformer transformer.DealInvitationTransformer
// 	dealInvitationValidator   validator.DealInvitationValidator
// 	userClient                *client.UserClient
// }

// func NewDealInvitationHandler(
// 	dealInvitationUsecase usecase.DealInvitationUsecase,
// 	dealInvitationTransformer transformer.DealInvitationTransformer,
// 	dealInvitationValidator validator.DealInvitationValidator,
// 	userClient *client.UserClient,
// ) *DealInvitationHandler {
// 	return &DealInvitationHandler{
// 		dealInvitationUsecase:     dealInvitationUsecase,
// 		dealInvitationTransformer: dealInvitationTransformer,
// 		dealInvitationValidator:   dealInvitationValidator,
// 		userClient:                userClient,
// 	}
// }

// // @Summary Gửi lời mời tham gia thương vụ
// // @Description Gửi lời mời tham gia thương vụ
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.SendInvitationRequest true "Thông tin lời mời"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.SendInvitationResponse "Thành công"
// // @Router /deal-invitation [post]
// func (h *DealInvitationHandler) SendInvitation(ctx context.Context, req *organizationpb.SendDealInvitationRequest) (*organizationpb.SendDealInvitationResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateSendInvitationRequest(req); err != nil {
// 		return nil, err
// 	}

// 	invitation := h.dealInvitationTransformer.SendInvitationRequestToEntity(req)

// 	createdInvitation, err := h.dealInvitationUsecase.SendInvitation(ctx, invitation)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.dealInvitationTransformer.EntityToSendInvitationResponse(createdInvitation), nil
// }

// // @Summary Xác nhận lời mời
// // @Description Xác nhận lời mời tham gia thương vụ
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.AcceptInvitationRequest true "Thông tin xác nhận"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.AcceptInvitationResponse "Thành công"
// // @Router /deal-invitation/accept [post]
// func (h *DealInvitationHandler) AcceptInvitation(ctx context.Context, req *organizationpb.AcceptDealInvitationRequest) (*organizationpb.AcceptDealInvitationResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateAcceptInvitationRequest(req); err != nil {
// 		return nil, err
// 	}

// 	err := h.dealInvitationUsecase.AcceptInvitation(ctx, req.InvitationId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &organizationpb.AcceptDealInvitationResponse{
// 		InvitationId: req.InvitationId,
// 		Success:      true,
// 	}, nil
// }

// // @Summary Từ chối lời mời
// // @Description Từ chối lời mời tham gia thương vụ
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.RejectInvitationRequest true "Thông tin từ chối"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.RejectInvitationResponse "Thành công"
// // @Router /deal-invitation/reject [post]
// func (h *DealInvitationHandler) RejectInvitation(ctx context.Context, req *organizationpb.RejectDealInvitationRequest) (*organizationpb.RejectDealInvitationResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateRejectInvitationRequest(req); err != nil {
// 		return nil, err
// 	}

// 	err := h.dealInvitationUsecase.RejectInvitation(ctx, req.InvitationId, req.Reason)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &organizationpb.RejectDealInvitationResponse{
// 		InvitationId: req.InvitationId,
// 		Success:      true,
// 	}, nil
// }

// // @Summary Gửi lại lời mời
// // @Description Gửi lại lời mời tham gia thương vụ
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.ResendInvitationRequest true "Thông tin gửi lại"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.ResendInvitationResponse "Thành công"
// // @Router /deal-invitation/resend [post]
// func (h *DealInvitationHandler) ResendInvitation(ctx context.Context, req *organizationpb.ResendDealInvitationRequest) (*organizationpb.ResendDealInvitationResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateResendInvitationRequest(req); err != nil {
// 		return nil, err
// 	}

// 	updatedInvitation, err := h.dealInvitationUsecase.ResendInvitation(ctx, req.InvitationId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.dealInvitationTransformer.EntityToResendInvitationResponse(updatedInvitation), nil
// }

// // @Summary Lấy danh sách thành viên đã accept
// // @Description Lấy danh sách thành viên đã chấp nhận lời mời tham gia thương vụ
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param deal_id path int true "ID của thương vụ"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.GetAcceptedMembersResponse "Thành công"
// // @Router /deal/{deal_id}/accepted-members [get]
// func (h *DealInvitationHandler) GetAcceptedMembers(ctx context.Context, req *organizationpb.GetAcceptedMembersRequest) (*organizationpb.GetAcceptedMembersResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateGetAcceptedMembersRequest(req); err != nil {
// 		return nil, err
// 	}

// 	members, err := h.dealInvitationUsecase.GetAcceptedMembers(ctx, req.DealId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.dealInvitationTransformer.EntitiesToGetAcceptedMembersResponse(members), nil
// }

// // @Summary Rút khỏi thương vụ
// // @Description Rút khỏi thương vụ (dành cho thành viên)
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.WithdrawFromDealRequest true "Thông tin rút khỏi"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.WithdrawFromDealResponse "Thành công"
// // @Router /deal-invitation/withdraw [post]
// func (h *DealInvitationHandler) WithdrawFromDeal(ctx context.Context, req *organizationpb.WithdrawFromDealRequest) (*organizationpb.WithdrawFromDealResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateWithdrawFromDealRequest(req); err != nil {
// 		return nil, err
// 	}

// 	err := h.dealInvitationUsecase.WithdrawFromDeal(ctx, req.InvitationId, req.Reason)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &organizationpb.WithdrawFromDealResponse{
// 		InvitationId: req.InvitationId,
// 		Success:      true,
// 	}, nil
// }

// // @Summary Gỡ khỏi thương vụ
// // @Description Gỡ khỏi thương vụ (dành cho admin/phụ trách)
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.RemoveFromDealRequest true "Thông tin gỡ khỏi"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.RemoveFromDealResponse "Thành công"
// // @Router /deal-invitation/remove [post]
// func (h *DealInvitationHandler) RemoveFromDeal(ctx context.Context, req *organizationpb.RemoveFromDealRequest) (*organizationpb.RemoveFromDealResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateRemoveFromDealRequest(req); err != nil {
// 		return nil, err
// 	}

// 	err := h.dealInvitationUsecase.RemoveFromDeal(ctx, req.InvitationId, req.Reason)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &organizationpb.RemoveFromDealResponse{
// 		InvitationId: req.InvitationId,
// 		Success:      true,
// 	}, nil
// }

// // @Summary Xác nhận rút khỏi thương vụ
// // @Description Xác nhận rút khỏi thương vụ (dành cho admin/phụ trách)
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request body organizationpb.ConfirmWithdrawalRequest true "Thông tin xác nhận"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.ConfirmWithdrawalResponse "Thành công"
// // @Router /deal-invitation/confirm-withdrawal [post]
// func (h *DealInvitationHandler) ConfirmWithdrawal(ctx context.Context, req *organizationpb.ConfirmWithdrawalRequest) (*organizationpb.ConfirmWithdrawalResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateConfirmWithdrawalRequest(req); err != nil {
// 		return nil, err
// 	}

// 	err := h.dealInvitationUsecase.ConfirmWithdrawal(ctx, req.InvitationId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &organizationpb.ConfirmWithdrawalResponse{
// 		InvitationId: req.InvitationId,
// 		Success:      true,
// 	}, nil
// }

// // @Summary Tìm kiếm thành viên
// // @Description Tìm kiếm thành viên + sort theo thương vụ chung
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request query organizationpb.SearchMembersRequest true "Thông tin tìm kiếm"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.SearchMembersResponse "Thành công"
// // @Router /deal-invitation/search [get]
// func (h *DealInvitationHandler) SearchMembers(ctx context.Context, req *organizationpb.SearchDealMembersRequest) (*organizationpb.SearchDealMembersResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateSearchMembersRequest(req); err != nil {
// 		return nil, err
// 	}

// 	page := int(req.Page)
// 	size := int(req.Size)
// 	if page <= 0 {
// 		page = 1
// 	}
// 	if size <= 0 {
// 		size = 10
// 	}

// 	members, total, err := h.dealInvitationUsecase.SearchMembers(ctx, req.DealId, req.Keyword, page, size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.dealInvitationTransformer.EntitiesToSearchMembersResponse(members, total), nil
// }

// // @Summary Lấy lời mời theo ID
// // @Description Lấy thông tin lời mời theo ID
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param id path int true "ID của lời mời"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.GetInvitationResponse "Thành công"
// // @Router /deal-invitation/{id} [get]
// func (h *DealInvitationHandler) GetInvitation(ctx context.Context, req *organizationpb.GetDealInvitationRequest) (*organizationpb.GetDealInvitationResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateGetInvitationRequest(req); err != nil {
// 		return nil, err
// 	}

// 	invitation, err := h.dealInvitationUsecase.GetInvitationByID(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.dealInvitationTransformer.EntityToGetInvitationResponse(invitation), nil
// }

// // @Summary Lấy danh sách lời mời đang chờ
// // @Description Lấy danh sách lời mời đang chờ xác nhận
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request query organizationpb.GetPendingInvitationsRequest true "Thông tin lấy danh sách"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.GetPendingInvitationsResponse "Thành công"
// // @Router /deal-invitation/pending [get]
// func (h *DealInvitationHandler) GetPendingInvitations(ctx context.Context, req *organizationpb.GetPendingInvitationsRequest) (*organizationpb.GetPendingInvitationsResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateGetPendingInvitationsRequest(req); err != nil {
// 		return nil, err
// 	}

// 	page := 0
// 	size := 10

// 	if req.Page != 0 {
// 		page = int(req.Page)
// 	}
// 	if req.Size != 0 {
// 		size = int(req.Size)
// 	}

// 	invitations, total, err := h.dealInvitationUsecase.GetPendingInvitations(ctx, req.InviteeId, page, size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h.dealInvitationTransformer.EntitiesToGetPendingInvitationsResponse(invitations, total), nil
// }

// // @Summary Lấy danh sách lời mời theo thương vụ
// // @Description Lấy danh sách lời mời theo thương vụ
// // @Tags Lời mời thương vụ
// // @Accept json
// // @Produce json
// // @Param request query organizationpb.GetInvitationsByDealRequest true "Thông tin lấy danh sách"
// // @Security BearerAuth
// // @Success 200 {object} organizationpb.GetInvitationsByDealResponse "Thành công"
// // @Router /deal/{deal_id}/invitations [get]
// func (h *DealInvitationHandler) GetInvitationsByDeal(ctx context.Context, req *organizationpb.GetDealInvitationsByDealRequest) (*organizationpb.GetDealInvitationsByDealResponse, error) {
// 	if err := h.dealInvitationValidator.ValidateGetInvitationsByDealRequest(req); err != nil {
// 		return nil, err
// 	}

// 	page := 0
// 	size := 10

// 	if req.Page != 0 {
// 		page = int(req.Page)
// 	}
// 	if req.Size != 0 {
// 		size = int(req.Size)
// 	}

// 	invitations, total, err := h.dealInvitationUsecase.GetInvitationsByDealID(ctx, req.DealId, page, size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pbs := h.dealInvitationTransformer.EntitiesToGetInvitationsByDealResponse(invitations, total)
// 	h.userClient.MapToDealInvitationPb(ctx, pbs.Data)

// 	return pbs, nil
// }
