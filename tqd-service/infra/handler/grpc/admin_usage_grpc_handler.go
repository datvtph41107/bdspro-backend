package handler_grpc

import (
	commonmetering "common/metering"
	"common/request"
	"context"
	"strconv"
	"time"

	tqdpb "pb/types/tqd"
	"tqd/internal/access"
	reportadmin "tqd/internal/usecase/generatedreport/admin"
	quotausecase "tqd/internal/usecase/quota"
	"tqd/internal/usecase/quota/reconciliation/usageprojection"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UsagePermissionAuthorizer interface {
	HasPermissions(context.Context, []string) error
}

type AdminUsageGrpcHandler struct {
	tqdpb.UnimplementedAdminUsageServiceServer
	service    *quotausecase.AdminUsageService
	reports    *reportadmin.Service
	authorizer UsagePermissionAuthorizer
	reconciler *usageprojection.Service
}

func NewAdminUsageGrpcHandler(service *quotausecase.AdminUsageService, reports *reportadmin.Service, authorizer UsagePermissionAuthorizer, reconciler *usageprojection.Service) *AdminUsageGrpcHandler {
	return &AdminUsageGrpcHandler{service: service, reports: reports, authorizer: authorizer, reconciler: reconciler}
}

func (h *AdminUsageGrpcHandler) ReconcileUsage(ctx context.Context, req *tqdpb.ReconcileAdminUsageRequest) (*tqdpb.ReconcileAdminUsageResponse, error) {
	if h == nil || h.authorizer == nil {
		return nil, status.Error(codes.Unavailable, "usage permission authority unavailable")
	}
	if err := h.authorizer.HasPermissions(ctx, []string{"COMMERCIAL_USAGE_RECONCILE"}); err != nil {
		return nil, status.Error(codes.PermissionDenied, "COMMERCIAL_USAGE_RECONCILE permission is required")
	}
	if _, ok := request.IdempotencyKeyFromContext(ctx); !ok {
		return nil, status.Error(codes.InvalidArgument, "Idempotency-Key is required")
	}
	if h.reconciler == nil {
		return nil, status.Error(codes.Unavailable, "usage reconciliation is unavailable")
	}
	if req == nil || req.GetProfileId() == 0 || req.GetMeterCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "profile id and meter code are required")
	}
	periodStart, err := adminUsageTimestamp(req.GetPeriodStart(), "period_start")
	if err != nil {
		return nil, err
	}
	periodEnd, err := adminUsageTimestamp(req.GetPeriodEnd(), "period_end")
	if err != nil {
		return nil, err
	}
	if !periodEnd.After(periodStart) {
		return nil, status.Error(codes.InvalidArgument, "period_end must be after period_start")
	}
	result, err := h.reconciler.CheckAndFix(ctx, access.Subject{
		Type: access.SubjectProfile, ID: strconv.FormatUint(req.GetProfileId(), 10),
	}, commonmetering.Code(req.GetMeterCode()), periodStart, periodEnd)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "usage reconciliation failed")
	}
	health := "in_sync"
	if !result.Fixed && result.SavedUsed != result.RuntimeUsed {
		health = "drift"
	}
	return &tqdpb.ReconcileAdminUsageResponse{
		DurableUsed: result.SavedUsed, PreviousRuntimeUsed: result.RuntimeUsed,
		RuntimeReserved: result.Reserved, Fixed: result.Fixed, ReconciliationHealth: health,
	}, nil
}

func adminUsageTimestamp(value *timestamppb.Timestamp, name string) (time.Time, error) {
	if value == nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s is required", name)
	}
	if err := value.CheckValid(); err != nil {
		return time.Time{}, status.Errorf(codes.InvalidArgument, "%s is invalid", name)
	}
	return value.AsTime().UTC(), nil
}

func (h *AdminUsageGrpcHandler) ListUsage(ctx context.Context, req *tqdpb.ListAdminUsageRequest) (*tqdpb.ListAdminUsageResponse, error) {
	if h == nil || h.authorizer == nil {
		return nil, status.Error(codes.Unavailable, "usage permission authority unavailable")
	}
	if err := h.authorizer.HasPermissions(ctx, []string{quotausecase.PermissionUsageView}); err != nil {
		return nil, status.Error(codes.PermissionDenied, "COMMERCIAL_USAGE_VIEW permission is required")
	}
	if h.service == nil {
		return nil, status.Error(codes.Unavailable, "admin usage projection unavailable")
	}
	if req == nil {
		req = &tqdpb.ListAdminUsageRequest{}
	}
	page, err := h.service.List(ctx, quotausecase.AdminUsageQuery{Page: req.GetPage(), PageSize: req.GetPageSize(), ProfileID: req.GetProfileId(), Operation: req.GetOperation(), MeterCode: req.GetMeterCode()})
	if err != nil {
		if err.Error() == "page size must be between 1 and 100" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "admin usage projection failed")
	}
	response := &tqdpb.ListAdminUsageResponse{Total: page.Total, Page: page.Page, PageSize: page.PageSize}
	for _, item := range page.Usage {
		projection := &tqdpb.AdminUsageProjection{
			SubjectType: item.SubjectType, SubjectId: item.SubjectID, Operation: item.Operation, MeterCode: item.MeterCode,
			DurableUsed: item.DurableUsed, RuntimeUsed: item.RuntimeUsed, RuntimeReserved: item.RuntimeReserved,
			RemainingProjection: item.Remaining, LimitKnown: item.LimitKnown, LimitSnapshot: item.LimitSnapshot,
			SubscriptionId: item.SubscriptionID, PlanCode: item.PlanCode, PlanVersion: item.PlanVersion,
			PolicyVersion: item.PolicyVersion, RuntimeProjectionAvailable: item.RuntimeAvailable,
			ReconciliationHealth: item.Reconciliation,
		}
		if !item.PeriodStart.IsZero() {
			projection.PeriodStart = timestamppb.New(item.PeriodStart)
		}
		if !item.PeriodEnd.IsZero() {
			projection.PeriodEnd = timestamppb.New(item.PeriodEnd)
		}
		if !item.LastUsageAt.IsZero() {
			projection.LastUsageAt = timestamppb.New(item.LastUsageAt)
		}
		response.Usage = append(response.Usage, projection)
	}
	return response, nil
}

func (h *AdminUsageGrpcHandler) ListGeneratedReports(ctx context.Context, req *tqdpb.ListAdminGeneratedReportsRequest) (*tqdpb.ListAdminGeneratedReportsResponse, error) {
	if h == nil || h.authorizer == nil {
		return nil, status.Error(codes.Unavailable, "report permission authority unavailable")
	}
	if err := h.authorizer.HasPermissions(ctx, []string{reportadmin.PermissionReportView}); err != nil {
		return nil, status.Error(codes.PermissionDenied, "COMMERCIAL_USAGE_VIEW permission is required")
	}
	if h.reports == nil {
		return nil, status.Error(codes.Unavailable, "generated report projection unavailable")
	}
	if req == nil {
		req = &tqdpb.ListAdminGeneratedReportsRequest{}
	}
	page, err := h.reports.List(ctx, reportadmin.Query{
		Page: req.GetPage(), PageSize: req.GetPageSize(), ProfileID: req.GetProfileId(),
		Status: req.GetStatus(), JobStatus: req.GetJobStatus(),
	})
	if err != nil {
		switch err.Error() {
		case "page size must be between 1 and 100", "invalid report job status":
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "generated report projection failed")
		}
	}
	response := &tqdpb.ListAdminGeneratedReportsResponse{
		Total: page.Total, Page: page.Page, PageSize: page.PageSize,
		Reports: make([]*tqdpb.AdminGeneratedReportProjection, 0, len(page.Reports)),
	}
	for _, item := range page.Reports {
		projection := &tqdpb.AdminGeneratedReportProjection{
			ReportId: item.ReportID, ProfileId: item.ProfileID,
			ReportType: item.ReportType, Status: item.Status,
			Title: item.Title, Subtitle: item.Subtitle,
			ParcelId: item.ParcelID, RegionId: item.RegionID, Format: item.Format,
			ThumbnailUrl: item.ThumbnailURL, ImageUrl: item.ImageURL,
			PdfUrl: item.PDFURL, ShareUrl: item.ShareURL,
			FileSize: item.FileSize, ErrorMessage: item.ErrorMessage,
			JobId: item.JobID, JobStatus: item.JobStatus,
			JobAttempts: item.JobAttempts, JobLockedBy: item.JobLockedBy,
			JobClaimVersion: item.JobClaimVersion, JobLastError: item.JobLastError,
		}
		if !item.CreatedAt.IsZero() {
			projection.CreatedAt = timestamppb.New(item.CreatedAt)
		}
		if !item.UpdatedAt.IsZero() {
			projection.UpdatedAt = timestamppb.New(item.UpdatedAt)
		}
		if item.ExpiresAt != nil {
			projection.ExpiresAt = timestamppb.New(*item.ExpiresAt)
		}
		if item.JobAvailableAt != nil {
			projection.JobAvailableAt = timestamppb.New(*item.JobAvailableAt)
		}
		response.Reports = append(response.Reports, projection)
	}
	return response, nil
}
