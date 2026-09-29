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

type ReportHandler struct {
	crmpb.UnimplementedReportServiceServer
	reportUsecase *usecase.ReportUsecase
	reportMapper  *mapper.ReportMapper
}

// @bind: crm/infra/handler.ReportHandler
func NewReportHandler(reportUsecase *usecase.ReportUsecase, reportMapper *mapper.ReportMapper) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
		reportMapper:  reportMapper,
	}
}

// CreateReport tạo báo cáo mới
// @Summary Tạo báo cáo mới
// @Description Tạo báo cáo mới cho một đối tượng
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param registry path string true "Registry"
// @Param request body dto.ReportCreateRequest true "Thông tin báo cáo"
// @Router /report/{registry} [post]
func (h *ReportHandler) CreateReport(ctx context.Context, req *crmpb.CreateReportRequest) (*crmpb.ReportResponse, error) {
	// Convert protobuf request to DTO
	reportReq := &dto.ReportCreateRequest{
		ReasonID:     req.ReasonId,
		UserID:       req.UserId,
		OwnerID:      req.OwnerId,
		OwnerOf:      req.OwnerOf,
		RegistryName: req.Registry,
		Content:      req.Content,
		ProofDocs:    req.ProofDocs,
	}

	// Create report
	createdReport, err := h.reportUsecase.CreateReport(ctx, reportReq)
	if err != nil {
		return nil, err
	}

	return h.reportMapper.ReportToPb(createdReport), nil
}

// GetReportList lấy danh sách báo cáo
// @Summary Lấy danh sách báo cáo
// @Description Lấy danh sách báo cáo của người dùng
// @Tags Report
// @Accept json
// @Produce json
// @Param registry path string true "Registry"
// @Param ownerId path uint64 true "Owner ID"
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} crmpb.GetReportListResponse
// @Router /report/{registry}/list/{ownerId} [get]
func (h *ReportHandler) GetReportList(ctx context.Context, req *crmpb.GetReportListRequest) (*crmpb.GetReportListResponse, error) {
	// Convert request to DTO
	searchReq := &dto.ReportListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		OwnerId:  &req.OwnerId,
		Status:   req.Status,
		ReasonID: req.ReasonId,
		Registry: req.Registry,
	}

	// Get report list
	reports, total, err := h.reportUsecase.GetReportList(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var reportResponses []*crmpb.ReportResponse
	for _, report := range reports {
		reportResponses = append(reportResponses, h.reportMapper.ReportToPb(&report))
	}

	return &crmpb.GetReportListResponse{
		Data:  reportResponses,
		Total: total,
	}, nil
}

// GetReportReasons lấy danh sách lý do báo cáo
// @Summary Lấy danh sách lý do báo cáo
// @Description Lấy danh sách lý do báo cáo có sẵn
// @Tags Report
// @Accept json
// @Produce json
// @Success 200 {object} crmpb.GetReportReasonsResponse
// @Router /report/reasons [get]
func (h *ReportHandler) GetReportReasons(ctx context.Context, req *crmpb.GetReportReasonsRequest) (*crmpb.GetReportReasonsResponse, error) {
	// Get report reasons
	reasons, err := h.reportUsecase.GetReportReasons(ctx)
	if err != nil {
		return nil, err
	}

	return &crmpb.GetReportReasonsResponse{
		Data:  h.reportMapper.ReportReasonListToReportReasonResponseList(reasons),
		Total: int64(len(reasons)),
	}, nil
}

// UpdateReport cập nhật báo cáo
// @Summary Cập nhật báo cáo
// @Description Cập nhật báo cáo của người dùng
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Param request body dto.ReportUpdateRequest true "Thông tin cập nhật"
// @Router /report/{id} [put]
func (h *ReportHandler) UpdateReport(ctx context.Context, req *crmpb.UpdateReportRequest) (*crmpb.UpdateReportResponse, error) {
	// Convert protobuf request to DTO
	reportReq := &dto.ReportUpdateRequest{
		ProofDocs: req.ProofDocs,
	}

	// Update report
	updatedReport, err := h.reportUsecase.UpdateReport(ctx, req.Id, reportReq)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	pbResponse := h.reportMapper.ReportToPb(updatedReport)

	return &crmpb.UpdateReportResponse{
		Report: pbResponse,
	}, nil
}

// DeleteReport xóa báo cáo
// @Summary Xóa báo cáo
// @Description Xóa báo cáo của người dùng
// @Tags Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID báo cáo"
// @Router /report/{id} [delete]
func (h *ReportHandler) DeleteReport(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Delete report
	err := h.reportUsecase.DeleteReport(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Report deleted successfully",
	}, nil
}