package handler_grpc

import (
	"context"
	"log"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"pb/clients"
	tqdpb "pb/types/tqd"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ReportGrpcHandler struct {
	tqdpb.UnimplementedReportServiceServer
	reportUsecase usecase.ReportUsecase
	SyncProvider  *_utils.SyncUtil
	authClient    *clients.AuthGrpcClient
}

func NewReportGrpcHandler(
	reportUsecase usecase.ReportUsecase,
	syncProvider *_utils.SyncUtil,
	authClient *clients.AuthGrpcClient,
) *ReportGrpcHandler {
	return &ReportGrpcHandler{
		reportUsecase: reportUsecase,
		SyncProvider:  syncProvider,
		authClient:    authClient,
	}
}

// CreateReportAsync - Tạo báo cáo bất đồng bộ
func (h *ReportGrpcHandler) CreateReportAsync(
	ctx context.Context,
	req *tqdpb.CreateReportAsyncRequest,
) (*tqdpb.CreateReportAsyncResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	log.Printf("[CreateReportAsync] UserID: %d, ReportType: %d, Profile: %d, Format: %d",
		userID, req.ReportType, req.Profile, req.Format)

	// Parse targetData if exists
	// if req.TargetData != nil && *req.TargetData != "" {
	// 	var targetData interface{}
	// 	if err := json.Unmarshal([]byte(*req.TargetData), &targetData); err != nil {
	// 		log.Printf("[CreateReportAsync] Invalid targetData: %v", err)
	// 		return nil, _errors.ReturnError(400, "invalid targetData json")
	// 	}
	// }

	report := dto.ReportDTO{
		ProblemReport: enums.ProblemReport(req.ProblemReport),
		Description:   req.Description,
		TargetId:      req.TargetId,
		ReportType:    enums.ReportType(req.ReportType),
	}

	// Create report
	reportID, err := h.reportUsecase.CreateAsync(
		ctx,
		userID,
		report,
	)
	if err != nil {
		log.Printf("[CreateReportAsync] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	log.Printf("[CreateReportAsync] Report created: %d", reportID)

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDReportList, userID), t.UnixMilli())

	return &tqdpb.CreateReportAsyncResponse{
		ReportId:  reportID,
		Status:    10, // Processing
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// GetReportStatus - Lấy trạng thái báo cáo
func (h *ReportGrpcHandler) GetReportStatus(
	ctx context.Context,
	req *tqdpb.GetReportStatusRequest,
) (*tqdpb.ReportStatusResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}

	log.Printf("[GetReportStatus] UserID: %d, ReportID: %d", userID, req.ReportId)

	report, err := h.reportUsecase.GetStatus(ctx, userID, req.ReportId)
	if err != nil {
		log.Printf("[GetReportStatus] Error: %v", err)
		return nil, _errors.ReturnError(404, err.Error())
	}

	resp := &tqdpb.ReportStatusResponse{
		Id:     report.ID,
		Status: uint32(report.Status),
	}

	// Convert string to *string for proto
	if report.FileURL != nil && *report.FileURL != "" {
		fileUrl := *report.FileURL
		resp.FileUrl = &fileUrl
	}
	if report.ErrorMessage != nil && *report.ErrorMessage != "" {
		errMsg := *report.ErrorMessage
		resp.ErrorMessage = &errMsg
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

// =====================================================
// DOWNLOAD REPORT
// =====================================================

// DownloadReport - Tải báo cáo
func (h *ReportGrpcHandler) DownloadReport(
	ctx context.Context,
	req *tqdpb.DownloadReportRequest,
) (*tqdpb.DownloadReportResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}

	log.Printf("[DownloadReport] UserID: %d, ReportID: %d", userID, req.ReportId)

	report, err := h.reportUsecase.Download(ctx, userID, req.ReportId)
	if err != nil {
		log.Printf("[DownloadReport] Error: %v", err)
		return nil, _errors.ReturnError(400, err.Error())
	}

	if report == nil {
		return nil, _errors.ReturnError(404, "report not found")
	}

	// TODO: Implement actual file download
	// Currently returns placeholder

	return &tqdpb.DownloadReportResponse{
		FileContent: []byte{},
		ContentType: "application/pdf",
		FileName:    "report.pdf",
	}, nil
}

// =====================================================
// LIST REPORTS (User)
// =====================================================

// ListReports - Danh sách báo cáo của user
func (h *ReportGrpcHandler) ListReports(ctx context.Context, req *tqdpb.ListReportsRequest) (*tqdpb.ListReportsResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDReportList, userID)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListReportsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	page := int(pagable.GetPage())
	limit := pagable.GetLimit()

	log.Printf("[ListReports] UserID: %d, Page: %d, Limit: %d, ReportType: %v, Status: %v",
		userID, page, limit, req.ReportType, req.Status)

	reports, total, err := h.reportUsecase.ListUser(ctx, userID, req.ReportType, req.Status, page, limit)
	if err != nil {
		log.Printf("[ListReports] Error: %v", err)
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.ReportSummary, 0, len(reports))
	for _, r := range reports {
		item := &tqdpb.ReportSummary{
			Id:         r.ID,
			ReportType: uint32(r.ProblemReport),
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
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// =====================================================
// DELETE REPORT (User)
// =====================================================

// DeleteReport - Xóa báo cáo của user
func (h *ReportGrpcHandler) DeleteReport(
	ctx context.Context,
	req *tqdpb.DeleteReportRequest,
) (*emptypb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}

	log.Printf("[DeleteReport] UserID: %d, ReportID: %d", userID, req.ReportId)

	if err := h.reportUsecase.Delete(ctx, userID, req.ReportId); err != nil {
		log.Printf("[DeleteReport] Error: %v", err)
		return nil, _errors.ReturnError(400, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDReportList, userID), t.UnixMilli())

	return &emptypb.Empty{}, nil
}

// =====================================================
// ADMIN DELETE
// =====================================================
