package handler

import (
	"context"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type ReportHandler struct {
	userpb.UnimplementedReportServiceServer
	reportUsecase *usecases.ReportUsecase
	reportMapper  *mapper.ReportMapper
}

func NewReportHandler(
	reportUsecase *usecases.ReportUsecase,
	reportMapper *mapper.ReportMapper,
) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
		reportMapper:  reportMapper,
	}
}

// CreateReport tạo báo cáo mới
// @Summary Tạo báo cáo
// @Description Tạo báo cáo người dùng mới
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.CreateReportRequest true "Thông tin báo cáo"
// @Success 200 {object} userpb.ReportResponse
// @Router /v2/user/report [post]
func (h *ReportHandler) CreateReport(ctx context.Context, req *userpb.CreateReportRequest) (*userpb.ReportResponse, error) {
	createReq := h.reportMapper.PbToCreateRequest(req)

	report, err := h.reportUsecase.CreateReport(ctx, createReq)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	_, err = h.reportUsecase.GetReportWithProofs(ctx, report.ID)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf from DTO which has complete data
	pbReport := h.reportMapper.EntityToPbReport(report)
	return &userpb.ReportResponse{Report: pbReport}, nil
}

// GetReports lấy danh sách báo cáo theo userId
// @Summary Lấy danh sách báo cáo
// @Description Lấy danh sách báo cáo theo userId
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId query uint64 false "ID người dùng"
// @Param status query uint32 false "Trạng thái báo cáo"
// @Param reasonId query uint64 false "ID nguyên nhân"
// @Param page query int32 false "Trang"
// @Param size query int32 false "Kích thước trang"
// @Success 200 {object} userpb.GetReportsResponse
// @Router /v2/user/report/list/{userId} [get]
func (h *ReportHandler) GetReports(ctx context.Context, req *userpb.GetReportsRequest) (*userpb.GetReportsResponse, error) {
	listReq := h.reportMapper.PbToListRequest(req)

	listResp, total, err := h.reportUsecase.GetReportsByUser(ctx, listReq)
	if err != nil {
		return nil, err
	}

	pbReports, err := h.reportMapper.ListResponseToPb(listResp, total)
	if err != nil {
		return nil, err
	}

	return &userpb.GetReportsResponse{
		Data:  pbReports,
		Total: total,
		Page:  int32(listReq.Page),
		Size:  int32(listReq.Size),
	}, nil
}

// UpdateReport cập nhật báo cáo (chỉ khi chưa duyệt)
// @Summary Cập nhật báo cáo
// @Description Cập nhật thông tin báo cáo (chỉ khi chưa duyệt)
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Param request body userpb.UpdateReportRequest true "Thông tin cập nhật"
// @Success 200 {object} userpb.ReportResponse
// @Router /v2/user/report/{id} [put]
func (h *ReportHandler) UpdateReport(ctx context.Context, req *userpb.UpdateReportRequest) (*userpb.ReportResponse, error) {
	updateReq := h.reportMapper.PbToUpdateRequest(req)

	report, err := h.reportUsecase.UpdateReport(ctx, req.Id, updateReq)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	_, err = h.reportUsecase.GetReportWithProofs(ctx, report.ID)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	pbReport := h.reportMapper.EntityToPbReport(report)
	return &userpb.ReportResponse{Report: pbReport}, nil
}

// DeleteReport xóa báo cáo
// @Summary Xóa báo cáo
// @Description Xóa báo cáo (chỉ khi chưa duyệt)
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Success 200 {object} userpb.DeleteReportResponse
// @Router /v2/user/report/{id} [delete]
func (h *ReportHandler) DeleteReport(ctx context.Context, req *userpb.DeleteReportRequest) (*userpb.DeleteReportResponse, error) {
	err := h.reportUsecase.DeleteReport(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.DeleteReportResponse{
		Message: "Xóa báo cáo thành công",
	}, nil
}

// RespondToReport phản hồi báo cáo (admin)
// @Summary Phản hồi báo cáo
// @Description Duyệt hoặc từ chối báo cáo (chỉ dành cho admin)
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Param request body userpb.RespondToReportRequest true "Phản hồi"
// @Success 200 {object} userpb.ReportResponse
// @Router /v2/user/report/{id}/respond [post]
func (h *ReportHandler) RespondToReport(ctx context.Context, req *userpb.RespondToReportRequest) (*userpb.ReportResponse, error) {
	respondReq := h.reportMapper.PbToResponseRequest(req)

	report, err := h.reportUsecase.RespondToReport(ctx, req.Id, respondReq)
	if err != nil {
		return nil, err
	}

	// Get complete report with proofs
	_, err = h.reportUsecase.GetReportWithProofs(ctx, report.ID)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf
	pbReport := h.reportMapper.EntityToPbReport(report)
	return &userpb.ReportResponse{Report: pbReport}, nil
}

// GetReport lấy thông tin chi tiết báo cáo
// @Summary Lấy thông tin chi tiết báo cáo
// @Description Lấy thông tin chi tiết báo cáo
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Success 200 {object} userpb.ReportResponse
// @Router /v2/user/report/detail/{id} [get]
func (h *ReportHandler) GetReport(ctx context.Context, req *userpb.GetReportRequest) (*userpb.Report, error) {
	report, err := h.reportUsecase.GetReport(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.reportMapper.EntityToPbReport(report), nil
}

// AdminNote cập nhật ghi chú admin cho báo cáo
// @Summary Cập nhật ghi chú admin
// @Description Cập nhật ghi chú nội bộ hoặc công khai cho báo cáo (chỉ dành cho admin)
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.AdminNoteRequest true "Thông tin ghi chú admin"
// @Success 200 {object} userpb.AdminNoteResponse
// @Router /v2/user/report/admin-note [put]
func (h *ReportHandler) AdminNote(ctx context.Context, req *userpb.AdminNoteRequest) (*userpb.AdminNoteResponse, error) {
	adminNoteReq := h.reportMapper.PbToAdminNoteRequest(req)

	report, err := h.reportUsecase.AdminNote(ctx, adminNoteReq)
	if err != nil {
		return nil, err
	}

	pbReport := h.reportMapper.EntityToPbReport(report)

	return &userpb.AdminNoteResponse{
		Success: true,
		Message: "Cập nhật ghi chú thành công",
		Report:  pbReport,
	}, nil
}
