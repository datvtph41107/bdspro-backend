package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type DealCommissionHandler struct {
	organizationpb.UnimplementedDealCommissionServiceServer
	dealCommissionUsecase     usecase.DealCommissionUsecase
	dealCommissionTransformer transformer.DealCommissionTransformer
	dealCommissionValidator   validator.DealCommissionValidator
}

func NewDealCommissionHandler(
	dealCommissionUsecase usecase.DealCommissionUsecase,
	dealCommissionTransformer transformer.DealCommissionTransformer,
	dealCommissionValidator validator.DealCommissionValidator,
) *DealCommissionHandler {
	return &DealCommissionHandler{
		dealCommissionUsecase:     dealCommissionUsecase,
		dealCommissionTransformer: dealCommissionTransformer,
		dealCommissionValidator:   dealCommissionValidator,
	}
}

// @Summary Cập nhật hoa hồng cho deal member
// @Description Cập nhật hoa hồng cho deal member
// @Tags Hoa hồng
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param memberId path uint64 true "ID thành viên"
// @Param request body organizationpb.UpdateDealMemberCommissionRequest true "Thông tin hoa hồng"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateDealMemberCommissionResponse "Thành công"
// @Router /deal/{dealId}/member/{memberId}/commission [put]
func (h *DealCommissionHandler) UpdateDealMemberCommission(ctx context.Context, req *organizationpb.UpdateDealMemberCommissionRequest) (*organizationpb.UpdateDealMemberCommissionResponse, error) {
	// Validate request
	if err := h.dealCommissionValidator.ValidateUpdateDealMemberCommissionRequest(req); err != nil {
		return nil, err
	}

	// Transform request
	commissionValue, commissionType, note := h.dealCommissionTransformer.TransformUpdateCommissionRequest(req)

	// Update commission
	err := h.dealCommissionUsecase.UpdateDealMemberCommission(ctx, req.DealId, req.MemberId, commissionValue, commissionType, note)
	if err != nil {
		return nil, err
	}

	// Transform response
	return h.dealCommissionTransformer.TransformUpdateCommissionResponse(req), nil
}

// @Summary Cập nhật ghi chú cho deal member
// @Description Cập nhật ghi chú cho deal member
// @Tags Hoa hồng
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param memberId path uint64 true "ID thành viên"
// @Param request body organizationpb.UpdateDealMemberNoteRequest true "Thông tin ghi chú"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateDealMemberNoteResponse "Thành công"
// @Router /deal/{dealId}/member/{memberId}/note [put]
func (h *DealCommissionHandler) UpdateDealMemberNote(ctx context.Context, req *organizationpb.UpdateDealMemberNoteRequest) (*organizationpb.UpdateDealMemberNoteResponse, error) {
	// Validate request
	if err := h.dealCommissionValidator.ValidateUpdateDealMemberNoteRequest(req); err != nil {
		return nil, err
	}

	// Update note
	err := h.dealCommissionUsecase.UpdateDealMemberNote(ctx, req.DealId, req.MemberId, req.Note)
	if err != nil {
		return nil, err
	}

	// Transform response
	return h.dealCommissionTransformer.TransformUpdateNoteResponse(req), nil
}

// @Summary Cập nhật hoa hồng cho nhiều deal members
// @Description Cập nhật hoa hồng cho nhiều deal members cùng lúc
// @Tags Hoa hồng
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param request body organizationpb.UpdateDealMembersCommissionRequest true "Thông tin hoa hồng cho nhiều members"
// @Security BearerAuth
// @Success 200 {object} organizationpb.UpdateDealMembersCommissionResponse "Thành công"
// @Router /deal/{dealId}/members/commission [put]
func (h *DealCommissionHandler) UpdateDealMembersCommission(ctx context.Context, req *organizationpb.UpdateDealMembersCommissionRequest) (*organizationpb.UpdateDealMembersCommissionResponse, error) {
	// Validate request
	if err := h.dealCommissionValidator.ValidateUpdateDealMembersCommissionRequest(req); err != nil {
		return nil, err
	}

	// Transform request
	members := h.dealCommissionTransformer.TransformUpdateMembersCommissionRequest(req)

	// Update commission for all members
	err := h.dealCommissionUsecase.UpdateDealMembersCommission(ctx, req.DealId, members)
	if err != nil {
		return nil, err
	}

	// Transform response
	return h.dealCommissionTransformer.TransformUpdateMembersCommissionResponse(req), nil
}

// @Summary Lấy thống kê lợi nhuận và hoa hồng
// @Description Lấy thống kê về lợi nhuận mục tiêu, tổng tiền đã chia, số người đã chia
// @Tags Hoa hồng
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetDealCommissionStatsResponse "Thành công"
// @Router /deal/{dealId}/commission/stats [get]
func (h *DealCommissionHandler) GetDealCommissionStats(ctx context.Context, req *organizationpb.GetDealCommissionStatsRequest) (*organizationpb.GetDealCommissionStatsResponse, error) {
	// Validate request
	if err := h.dealCommissionValidator.ValidateGetDealCommissionStatsRequest(req); err != nil {
		return nil, err
	}

	// Get commission stats
	stats, err := h.dealCommissionUsecase.GetDealCommissionStats(ctx, req.DealId)
	if err != nil {
		return nil, err
	}

	// Transform response
	return h.dealCommissionTransformer.TransformCommissionStatsResponse(stats), nil
} 