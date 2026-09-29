package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/enums"
	"organization/internal/usecase"
)

type InternalNoteHandler struct {
	organizationpb.UnimplementedInternalNoteServiceServer
	internalNoteUsecase     usecase.InternalNoteUsecase
	internalNoteTransformer transformer.InternalNoteTransformer
	internalNoteValidator   validator.InternalNoteValidator
}

func NewInternalNoteHandler(
	internalNoteUsecase usecase.InternalNoteUsecase,
	internalNoteTransformer transformer.InternalNoteTransformer,
	internalNoteValidator validator.InternalNoteValidator,
) *InternalNoteHandler {
	return &InternalNoteHandler{
		internalNoteUsecase:     internalNoteUsecase,
		internalNoteTransformer: internalNoteTransformer,
		internalNoteValidator:   internalNoteValidator,
	}
}

// @Summary Thêm ghi chú nội bộ
// @Description Thêm ghi chú nội bộ cho deal
// @Tags Ghi chú nội bộ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param request body organizationpb.AddInternalNoteRequest true "Thông tin ghi chú"
// @Security BearerAuth
// @Success 200 {object} organizationpb.AddInternalNoteResponse "Thành công"
// @Router /deal/{dealId}/internal-note [post]
func (h *InternalNoteHandler) AddInternalNote(ctx context.Context, req *organizationpb.AddInternalNoteRequest) (*organizationpb.AddInternalNoteResponse, error) {
	// Validate request
	if err := h.internalNoteValidator.ValidateAddInternalNoteRequest(req); err != nil {
		return nil, err
	}

	// Transform request
	content, actionType := h.internalNoteTransformer.TransformAddInternalNoteRequest(req)

	// Add internal note
	note, err := h.internalNoteUsecase.AddInternalNote(ctx, req.DealId, content, actionType)
	if err != nil {
		return nil, err
	}

	// Transform response
	return h.internalNoteTransformer.TransformAddInternalNoteResponse(note), nil
}

// @Summary Lấy danh sách lịch sử
// @Description Lấy danh sách lịch sử của deal với phân trang và lọc theo loại hành động
// @Tags Ghi chú nội bộ
// @Accept json
// @Produce json
// @Param dealId path uint64 true "ID thương vụ"
// @Param actionTypes query []uint32 false "Lọc theo loại hành động"
// @Param page query uint32 false "Trang hiện tại" default(1)
// @Param limit query uint32 false "Số lượng item mỗi trang" default(20)
// @Security BearerAuth
// @Success 200 {object} organizationpb.GetDealHistoryResponse "Thành công"
// @Router /deal/{dealId}/history [get]
func (h *InternalNoteHandler) GetDealHistory(ctx context.Context, req *organizationpb.GetDealHistoryRequest) (*organizationpb.GetDealHistoryResponse, error) {
	// Validate request
	if err := h.internalNoteValidator.ValidateGetDealHistoryRequest(req); err != nil {
		return nil, err
	}
	page := req.Page
	if page == 0 {
		page = 0
	}
	size := req.Size
	if size == 0 {
		size = 10
	}
	actionTypes := make([]enums.ActionType, len(req.ActionTypes))
	for i, actionType := range req.ActionTypes {
		actionTypes[i] = enums.ActionType(actionType)
	}
	// Get deal history
	notes, total, err := h.internalNoteUsecase.GetDealHistory(ctx, req.DealId, actionTypes, page, size)
	if err != nil {
		return nil, err
	}

	// Transform response
	return h.internalNoteTransformer.TransformGetDealHistoryResponse(req.DealId, notes, total, req.Page, req.Size), nil
}