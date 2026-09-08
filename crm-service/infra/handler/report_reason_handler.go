package handler

import (
	"context"
	"crm/internal/dto"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	_dto "common/domain/dto"
	_utils "common/utils"
)

type ReportReasonHandler struct {
	crmpb.UnimplementedReportReasonServiceServer
	reportReasonUsecase *usecase.ReportReasonUsecase
}

// @bind: crm/infra/handler.ReportReasonHandler
func NewReportReasonHandler(reportReasonUsecase *usecase.ReportReasonUsecase) *ReportReasonHandler {
	return &ReportReasonHandler{
		reportReasonUsecase: reportReasonUsecase,
	}
}

// CreateReportReason tạo lý do báo cáo mới
// @Summary Tạo lý do báo cáo mới
// @Description Tạo lý do báo cáo mới
// @Tags ReportReason
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ReportReasonCreateRequest true "Thông tin lý do báo cáo"
// @Router /report-reason [post]
func (h *ReportReasonHandler) CreateReportReason(ctx context.Context, req *crmpb.CreateReportReasonRequest) (*crmpb.CreateReportReasonResponse, error) {
	// Convert protobuf request to DTO
	createReq := &dto.ReportReasonCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	// Create reason
	response, err := h.reportReasonUsecase.CreateReportReason(ctx, createReq)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	reasonResponse := convertReportReasonToProtobuf(response)

	return &crmpb.CreateReportReasonResponse{
		Reason: reasonResponse,
	}, nil
}

// GetReportReasonList lấy danh sách lý do báo cáo
// @Summary Lấy danh sách lý do báo cáo
// @Description Lấy danh sách lý do báo cáo với phân trang
// @Tags ReportReason
// @Accept json
// @Produce json
// @Param is_active query bool false "Lọc theo trạng thái hoạt động"
// @Param page query int false "Số trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} crmpb.GetReportReasonListResponse
// @Router /report-reason [get]
func (h *ReportReasonHandler) GetReportReasonList(ctx context.Context, req *crmpb.GetReportReasonListRequest) (*crmpb.GetReportReasonListResponse, error) {
	// Convert protobuf request to DTO
	listReq := &dto.ReportReasonListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		IsActive: req.IsActive,
	}

	// Get reason list
	response, err := h.reportReasonUsecase.GetReportReasonList(ctx, listReq)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	reasons := make([]*crmpb.ReportReason, len(response.Data))
	for i, reason := range response.Data {
		reasons[i] = convertReportReasonToProtobuf(&reason)
	}

	return &crmpb.GetReportReasonListResponse{
		Data:  reasons,
		Total: response.Total,
	}, nil
}

// GetReportReason lấy lý do báo cáo theo ID
// @Summary Lấy lý do báo cáo theo ID
// @Description Lấy thông tin lý do báo cáo theo ID
// @Tags ReportReason
// @Accept json
// @Produce json
// @Param id path uint64 true "ID lý do báo cáo"
// @Success 200 {object} crmpb.GetReportReasonResponse
// @Router /report-reason/{id} [get]
func (h *ReportReasonHandler) GetReportReason(ctx context.Context, req *sharepb.IdRequest) (*crmpb.GetReportReasonResponse, error) {
	// Get reason by ID
	response, err := h.reportReasonUsecase.GetReportReasonByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	reasonResponse := convertReportReasonToProtobuf(response)

	return &crmpb.GetReportReasonResponse{
		Reason: reasonResponse,
	}, nil
}

// UpdateReportReason cập nhật lý do báo cáo
// @Summary Cập nhật lý do báo cáo
// @Description Cập nhật thông tin lý do báo cáo
// @Tags ReportReason
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID lý do báo cáo"
// @Param request body dto.ReportReasonUpdateRequest true "Thông tin cập nhật"
// @Success 200 {object} crmpb.UpdateReportReasonResponse
// @Router /report-reason/{id} [put]
func (h *ReportReasonHandler) UpdateReportReason(ctx context.Context, req *crmpb.UpdateReportReasonRequest) (*crmpb.UpdateReportReasonResponse, error) {
	// Convert protobuf request to DTO
	updateReq := &dto.ReportReasonUpdateRequest{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	// Update reason
	response, err := h.reportReasonUsecase.UpdateReportReason(ctx, updateReq)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	reasonResponse := convertReportReasonToProtobuf(response)

	return &crmpb.UpdateReportReasonResponse{
		Reason: reasonResponse,
	}, nil
}

// DeleteReportReason xóa lý do báo cáo
// @Summary Xóa lý do báo cáo
// @Description Xóa lý do báo cáo theo ID
// @Tags ReportReason
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID lý do báo cáo"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /report-reason/{id} [delete]
func (h *ReportReasonHandler) DeleteReportReason(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Delete reason
	err := h.reportReasonUsecase.DeleteReportReason(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Report reason deleted successfully",
	}, nil
}

// Helper functions to convert between protobuf and domain
func convertReportReasonToProtobuf(reason *dto.ReportReasonResponse) *crmpb.ReportReason {
	if reason == nil {
		return nil
	}

	return &crmpb.ReportReason{
		Id:          reason.ID,
		Name:        reason.Name,
		Description: reason.Description,
		IsActive:    reason.IsActive,
		CreatedAt:   _utils.FormatTimeToString(reason.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(reason.UpdatedAt),
	}
}