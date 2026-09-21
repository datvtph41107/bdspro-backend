package handler_grpc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"pb/clients"
	tqdpb "pb/types/tqd"
	"tqd/internal"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	permReportView = []string{
		"ADMIN_BC_GIS_XEM",
		"ADMIN_BC_GIS_XU_LY", "ADMIN_BC_GIS_SUA",
		"ADMIN_BC_GIS_DONG", "ADMIN_BC_GIS_XOA",
	}
	permReportCreate  = []string{"ADMIN_BC_GIS_TAO"}
	permReportProcess = []string{"ADMIN_BC_GIS_XU_LY", "ADMIN_BC_GIS_SUA"}
	permReportClose   = []string{"ADMIN_BC_GIS_DONG", "ADMIN_BC_GIS_XOA"}
)

// GisReportAdminGrpcHandler — admin GIS reports qua gateway (grpc-gateway).
type GisReportAdminGrpcHandler struct {
	tqdpb.UnimplementedGisReportAdminServiceServer
	usecase    usecase.ReportUsecase
	authClient *clients.AuthGrpcClient
}

func NewGisReportAdminGrpcHandler(
	reportUsecase usecase.ReportUsecase,
	authClient *clients.AuthGrpcClient,
) *GisReportAdminGrpcHandler {
	return &GisReportAdminGrpcHandler{
		usecase:    reportUsecase,
		authClient: authClient,
	}
}

func (h *GisReportAdminGrpcHandler) HasPermissions(ctx context.Context, keys []string) error {
	if h == nil || h.authClient == nil {
		return _errors.ReturnError(service.PermissionAuthorityUnavailable)
	}
	return h.authClient.HasPermissions(ctx, keys)
}

func (h *GisReportAdminGrpcHandler) requirePerm(ctx context.Context, keys []string) error {
	if err := h.authClient.HasPermissions(ctx, keys); err != nil {
		return err
	}
	return nil
}

func (h *GisReportAdminGrpcHandler) List(ctx context.Context, req *tqdpb.AdminGisReportListRequest) (*tqdpb.AdminGisReportListResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	filter := repo.AdminReportFilter{
		Page:  int(pagable.GetPage()),
		Limit: pagable.GetLimit(),
	}
	if req.Q != nil {
		filter.Q = strings.TrimSpace(*req.Q)
	}
	filter.ProblemReport = req.ProblemReport
	filter.QaStatus = req.QaStatus
	filter.Severity = req.Severity
	filter.AssigneeID = req.AssigneeId

	reports, total, err := h.usecase.AdminList(ctx, filter)
	if err != nil {
		return nil, err
	}
	displays := h.usecase.ResolveLinkageDisplays(ctx, reports)
	users := h.usecase.ResolveUserDisplays(ctx, reports)
	return adminGisReportListResponse(reports, total, int32(pagable.GetPage()), int32(pagable.GetSize()), "", displays, users), nil
}

func (h *GisReportAdminGrpcHandler) Summary(ctx context.Context, _ *emptypb.Empty) (*tqdpb.AdminGisReportSummaryResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	m, err := h.usecase.AdminSummary(ctx)
	if err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportSummaryResponse{ByQaStatus: m}, nil
}

func (h *GisReportAdminGrpcHandler) Queue(ctx context.Context, req *tqdpb.AdminGisReportQueueRequest) (*tqdpb.AdminGisReportListResponse, error) {
	start := time.Now()
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("[GisReportAdminGrpcHandler.Queue] Permission check took %s", time.Since(start)))

	bucket := strings.TrimSpace(req.GetBucket())
	if bucket == "" {
		bucket = "mine"
	}
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	var severity *uint32
	if req.Severity != nil {
		severity = req.Severity
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	slog.InfoContext(ctx, fmt.Sprintf("[GisReportAdminGrpcHandler.Queue] ActorID: %d", actorID))

	queryStart := time.Now()
	reports, total, err := h.usecase.AdminQueue(ctx, bucket, severity, int(pagable.GetPage()), pagable.GetLimit(), actorID)
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("[GisReportAdminGrpcHandler.Queue] AdminQueue query took %s, found %d reports", time.Since(queryStart), len(reports)))

	resolveStart := time.Now()
	displays := h.usecase.ResolveLinkageDisplays(ctx, reports)
	users := h.usecase.ResolveUserDisplays(ctx, reports)
	slog.InfoContext(ctx, fmt.Sprintf("[GisReportAdminGrpcHandler.Queue] ResolveLinkageDisplays and ResolveUserDisplays took %s", time.Since(resolveStart)))

	return adminGisReportListResponse(reports, total, int32(pagable.GetPage()), int32(pagable.GetSize()), bucket, displays, users), nil
}

func (h *GisReportAdminGrpcHandler) QueueSummary(ctx context.Context, _ *emptypb.Empty) (*tqdpb.AdminGisReportQueueSummaryResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	sum, err := h.usecase.AdminQueueSummary(ctx, actorID)
	if err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportQueueSummaryResponse{
		Mine:       sum.Mine,
		Unassigned: sum.Unassigned,
		DataFix:    sum.DataFix,
	}, nil
}

func (h *GisReportAdminGrpcHandler) ListEvents(ctx context.Context, req *tqdpb.AdminGisReportListEventsRequest) (*tqdpb.AdminGisReportEventListResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	filter := repo.AdminEventFilter{
		Page:  int(pagable.GetPage()),
		Limit: pagable.GetLimit(),
	}
	filter.ReportID = req.ReportId
	filter.ActorID = req.ActorId
	if req.Action != nil {
		filter.Action = strings.TrimSpace(*req.Action)
	}
	filter.From = req.From
	filter.To = req.To

	events, total, err := h.usecase.AdminListEvents(ctx, filter)
	if err != nil {
		return nil, err
	}
	names := h.actorNames(ctx, events)
	return toAdminGisReportEvents(events, total, names), nil
}

func (h *GisReportAdminGrpcHandler) ListReportEvents(ctx context.Context, req *tqdpb.AdminGisReportListReportEventsRequest) (*tqdpb.AdminGisReportEventListResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 50
	}
	events, err := h.usecase.AdminListReportEvents(ctx, req.GetId(), limit)
	if err != nil {
		return nil, err
	}
	names := h.actorNames(ctx, events)
	return toAdminGisReportEvents(events, int64(len(events)), names), nil
}

func (h *GisReportAdminGrpcHandler) Get(ctx context.Context, req *tqdpb.AdminGisReportGetRequest) (*tqdpb.AdminGisReportDetail, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	if req.GetId() == 0 {
		return nil, _errors.ReturnError(service.AdminReportIDInvalid)
	}
	report, err := h.usecase.AdminGet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	displays := h.usecase.ResolveLinkageDisplays(ctx, []domain.Report{*report})
	users := h.usecase.ResolveUserDisplays(ctx, []domain.Report{*report})
	return toAdminGisReportDetail(*report, displays[report.ID], users[report.ID]), nil
}

func (h *GisReportAdminGrpcHandler) CreateInternal(ctx context.Context, req *tqdpb.AdminGisReportCreateInternalRequest) (*tqdpb.AdminGisReportDetail, error) {
	if err := h.requirePerm(ctx, permReportCreate); err != nil {
		return nil, err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	createReq := dto.AdminCreateReportRequest{
		Title:           req.GetTitle(),
		Description:     req.GetDescription(),
		ProblemReport:   req.GetProblemReport(),
		ReportType:      req.GetReportType(),
		Severity:        req.GetSeverity(),
		TargetID:        req.TargetId,
		SupportTicketID: req.SupportTicketId,
		Images:          req.GetImages(),
	}
	report, err := h.usecase.AdminCreateInternal(ctx, actorID, createReq)
	if err != nil {
		return nil, err
	}
	users := h.usecase.ResolveUserDisplays(ctx, []domain.Report{*report})
	return toAdminGisReportDetail(*report, usecase.LinkageDisplay{}, users[report.ID]), nil
}

func (h *GisReportAdminGrpcHandler) Assign(ctx context.Context, req *tqdpb.AdminGisReportAssignRequest) (*tqdpb.AdminGisReportActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if err := h.usecase.AdminAssign(ctx, req.GetId(), req.GetAssigneeId()); err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportActionResponse{Id: req.GetId(), Message: "assigned"}, nil
}

func (h *GisReportAdminGrpcHandler) UpdateSeverity(ctx context.Context, req *tqdpb.AdminGisReportUpdateSeverityRequest) (*tqdpb.AdminGisReportActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if err := h.usecase.AdminUpdateSeverity(ctx, req.GetId(), req.GetSeverity(), req.GetReason()); err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportActionResponse{Id: req.GetId(), Message: "severity updated"}, nil
}

func (h *GisReportAdminGrpcHandler) UpdateQaStatus(ctx context.Context, req *tqdpb.AdminGisReportUpdateQaStatusRequest) (*tqdpb.AdminGisReportActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if err := h.usecase.AdminUpdateQaStatus(ctx, req.GetId(), req.GetQaStatus(), req.GetNote()); err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportActionResponse{Id: req.GetId(), Message: "qa status updated"}, nil
}

func (h *GisReportAdminGrpcHandler) UpdateLinkage(ctx context.Context, req *tqdpb.AdminGisReportUpdateLinkageRequest) (*tqdpb.AdminGisReportActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	linkage := structToMap(req.GetLinkage())
	if err := h.usecase.AdminUpdateLinkage(ctx, req.GetId(), linkage); err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportActionResponse{Id: req.GetId(), Message: "linkage updated"}, nil
}

func (h *GisReportAdminGrpcHandler) UpdateImages(ctx context.Context, req *tqdpb.AdminGisReportUpdateImagesRequest) (*tqdpb.AdminGisReportActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if err := h.usecase.AdminUpdateImages(ctx, req.GetId(), req.GetImages()); err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportActionResponse{Id: req.GetId(), Message: "images updated"}, nil
}

func (h *GisReportAdminGrpcHandler) Close(ctx context.Context, req *tqdpb.AdminGisReportCloseRequest) (*tqdpb.AdminGisReportActionResponse, error) {
	if err := h.requirePerm(ctx, permReportClose); err != nil {
		return nil, err
	}
	if err := h.usecase.AdminClose(ctx, req.GetId(), req.GetReject(), req.GetNote()); err != nil {
		return nil, err
	}
	return &tqdpb.AdminGisReportActionResponse{Id: req.GetId(), Message: "closed"}, nil
}

func (h *GisReportAdminGrpcHandler) actorNames(ctx context.Context, events []domain.ReportEvent) map[uint64]string {
	actorIDs := make([]uint64, 0, len(events))
	for _, e := range events {
		actorIDs = append(actorIDs, e.ActorID)
	}
	return h.usecase.ResolveActorNames(ctx, actorIDs)
}
