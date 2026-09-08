package handler

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	_dto "common/domain/dto"
)

type ReportAdminHandler struct {
	crmpb.UnimplementedReportAdminServiceServer
	reportUsecase *usecase.ReportUsecase
	reportMapper  *mapper.ReportMapper
}

// @bind: crm/infra/handler.ReportAdminHandler
func NewReportAdminHandler(reportUsecase *usecase.ReportUsecase, reportMapper *mapper.ReportMapper) *ReportAdminHandler {
	return &ReportAdminHandler{
		reportUsecase: reportUsecase,
		reportMapper:  reportMapper,
	}
}

// GetAdminReportList lấy danh sách báo cáo cho admin
// @Summary Lấy danh sách báo cáo (Admin)
// @Description Lấy danh sách tất cả báo cáo theo registry cho admin
// @Tags Admin Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param registry path string true "Registry"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Param status query int false "Trạng thái báo cáo"
// @Param reason_id query int false "ID lý do báo cáo"
// @Param userId query int false "ID người dùng bị báo cáo"
// @Param ownerId query int false "ID đối tượng bị báo cáo"
// @Success 200 {object} crmpb.GetAdminReportListResponse
// @Router /admin/report/{registry} [get]
func (h *ReportAdminHandler) GetAdminReportList(ctx context.Context, req *crmpb.GetAdminReportListRequest) (*crmpb.GetAdminReportListResponse, error) {
	// Convert request to DTO
	searchReq := &dto.ReportListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Registry: req.Registry,
		Status:   req.Status,
		ReasonID: req.ReasonId,
	}

	if req.UserId != nil {
		userID := *req.UserId
		searchReq.OwnerId = &userID
	}

	if req.OwnerId != nil {
		ownerID := *req.OwnerId
		searchReq.OwnerId = &ownerID
	}

	// Get report list
	reports, total, err := h.reportUsecase.GetAdminReportList(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var reportResponses []*crmpb.ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, h.reportMapper.ReportToPb(&report))
	}

	return &crmpb.GetAdminReportListResponse{
		Data:  reportResponses,
		Total: total,
	}, nil
}

// ApproveReport duyệt báo cáo
// @Summary Duyệt báo cáo (Admin)
// @Description Admin duyệt báo cáo
// @Tags Admin Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Param request body dto.AdminReportActionRequest true "Thông tin duyệt"
// @Router /admin/report/{id}/approve [post]
func (h *ReportAdminHandler) ApproveReport(ctx context.Context, req *crmpb.ApproveReportRequest) (*sharepb.SubmitResponse, error) {
	// Convert to DTO
	actionReq := &dto.AdminReportActionRequest{
		Response:  req.Response,
		AdminNote: req.AdminNote,
	}
	if req.SendNotification != nil {
		actionReq.SendNotification = *req.SendNotification
	}

	// Approve report
	err := h.reportUsecase.ApproveReport(ctx, req.Id, actionReq)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Báo cáo đã được duyệt",
	}, nil
}

// RejectReport từ chối báo cáo
// @Summary Từ chối báo cáo (Admin)
// @Description Admin từ chối báo cáo
// @Tags Admin Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Param request body dto.AdminReportActionRequest true "Thông tin từ chối"
// @Router /admin/report/{id}/reject [post]
func (h *ReportAdminHandler) RejectReport(ctx context.Context, req *crmpb.RejectReportRequest) (*sharepb.SubmitResponse, error) {
	// Convert to DTO
	actionReq := &dto.AdminReportActionRequest{
		Response:  req.Response,
		AdminNote: req.AdminNote,
	}
	if req.SendNotification != nil {
		actionReq.SendNotification = *req.SendNotification
	}

	// Reject report
	err := h.reportUsecase.RejectReport(ctx, req.Id, actionReq)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Báo cáo đã bị từ chối",
	}, nil
}

// RemoveReport gỡ báo cáo (xóa hẳn)
// @Summary Gỡ báo cáo (Admin)
// @Description Admin gỡ báo cáo khỏi hệ thống
// @Tags Admin Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Router /admin/report/{id}/remove [delete]
func (h *ReportAdminHandler) RemoveReport(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Remove report (hard delete or mark as removed)
	err := h.reportUsecase.RemoveReport(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Báo cáo đã được gỡ khỏi hệ thống",
	}, nil
}