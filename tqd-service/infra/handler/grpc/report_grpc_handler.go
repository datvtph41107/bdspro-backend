package handler_grpc

import (
	"context"
	"encoding/json"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ReportGrpcHandler struct {
	tqdpb.UnimplementedReportServiceServer
	reportUsecase usecase.ReportUsecase
}

func NewReportGrpcHandler(reportUsecase usecase.ReportUsecase) *ReportGrpcHandler {
	return &ReportGrpcHandler{reportUsecase: reportUsecase}
}

func (h *ReportGrpcHandler) CreateReportAsync(
	ctx context.Context,
	req *tqdpb.CreateReportAsyncRequest,
) (*tqdpb.CreateReportAsyncResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	var targetData interface{}
	if req.TargetData != nil && *req.TargetData != "" {
		if err := json.Unmarshal([]byte(*req.TargetData), &targetData); err != nil {
			return nil, _errors.ReturnError(400, "invalid targetData json")
		}
	}

	reportID, err := h.reportUsecase.CreateAsync(
		ctx,
		userID,
		req.ReportType,
		req.TargetId,
		targetData,
		req.Profile,
		req.Format,
		req.GetIncludeMapImage(),
	)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return &tqdpb.CreateReportAsyncResponse{
		ReportId:  reportID,
		Status:    10,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (h *ReportGrpcHandler) GetReportStatus(
	ctx context.Context,
	req *tqdpb.GetReportStatusRequest,
) (*tqdpb.ReportStatusResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	report, err := h.reportUsecase.GetStatus(ctx, userID, req.ReportId)
	if err != nil {
		return nil, _errors.ReturnError(404, err.Error())
	}

	resp := &tqdpb.ReportStatusResponse{
		Id:     report.ID,
		Status: uint32(report.Status),
	}

	if report.FileURL != nil {
		resp.FileUrl = report.FileURL
	}
	if report.ErrorMessage != nil {
		resp.ErrorMessage = report.ErrorMessage
	}
	if report.CompletedAt != nil {
		completed := report.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &completed
	}
	if report.ExpiresAt != nil {
		expires := report.ExpiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &expires
	}

	return resp, nil
}

func (h *ReportGrpcHandler) DownloadReport(
	ctx context.Context,
	req *tqdpb.DownloadReportRequest,
) (*tqdpb.DownloadReportResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	_, err := h.reportUsecase.Download(ctx, userID, req.ReportId)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	return &tqdpb.DownloadReportResponse{
		FileContent: []byte{},
		ContentType: "application/pdf",
		FileName:    "report.pdf",
	}, nil
}

func (h *ReportGrpcHandler) ListReports(
	ctx context.Context,
	req *tqdpb.ListReportsRequest,
) (*tqdpb.ListReportsResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	reports, total, err := h.reportUsecase.ListUser(ctx, userID, req.ReportType, req.Status, int(page), int(limit))
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.ReportSummary, 0, len(reports))
	for _, r := range reports {
		item := &tqdpb.ReportSummary{
			Id:         r.ID,
			ReportType: uint32(r.ReportType),
			Profile:    r.Profile,
			Status:     uint32(r.Status),
			CreatedAt:  r.CreatedAt.Format(time.RFC3339),
		}
		if r.CompletedAt != nil {
			completed := r.CompletedAt.Format(time.RFC3339)
			item.CompletedAt = &completed
		}
		if r.FileURL != nil {
			item.FileUrl = r.FileURL
		}
		items = append(items, item)
	}

	return &tqdpb.ListReportsResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  limit,
	}, nil
}

func (h *ReportGrpcHandler) DeleteReport(
	ctx context.Context,
	req *tqdpb.DeleteReportRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.reportUsecase.Delete(ctx, userID, req.ReportId); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *ReportGrpcHandler) AdminListReports(
	ctx context.Context,
	req *tqdpb.AdminListReportsRequest,
) (*tqdpb.ListReportsResponse, error) {
	// Admin check tạm thời bỏ qua
	page := int32(1)
	limit := int32(20)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	reports, total, err := h.reportUsecase.AdminList(ctx, req.UserId, req.ReportType, req.Status, int(page), int(limit))
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.ReportSummary, 0, len(reports))
	for _, r := range reports {
		item := &tqdpb.ReportSummary{
			Id:         r.ID,
			ReportType: uint32(r.ReportType),
			Profile:    r.Profile,
			Status:     uint32(r.Status),
			CreatedAt:  r.CreatedAt.Format(time.RFC3339),
		}
		if r.CompletedAt != nil {
			completed := r.CompletedAt.Format(time.RFC3339)
			item.CompletedAt = &completed
		}
		if r.FileURL != nil {
			item.FileUrl = r.FileURL
		}
		items = append(items, item)
	}

	return &tqdpb.ListReportsResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  limit,
	}, nil
}

func (h *ReportGrpcHandler) AdminDeleteReport(
	ctx context.Context,
	req *tqdpb.AdminDeleteReportRequest,
) (*emptypb.Empty, error) {
	if err := h.reportUsecase.AdminDelete(ctx, req.ReportId); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}
