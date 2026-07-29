package handler_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	"pb/clients"
	tqdpb "pb/types/tqd"
	"tqd/env"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
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

func (h *ReportGrpcHandler) auth() *clients.AuthGrpcClient {
	if h.authClient == nil {
		h.authClient = clients.MakeAuthGrpcClient(env.ENV_RUNTIME)
	}
	return h.authClient
}

func (h *ReportGrpcHandler) requirePerm(ctx context.Context, keys []string) error {
	if err := h.auth().HasPermissions(ctx, keys); err != nil {
		return _errors.ReturnError(403, err.Error())
	}
	return nil
}

// AdminListReports — GET /v2/tqd/admin/reports
func (h *ReportGrpcHandler) AdminListReports(ctx context.Context, req *tqdpb.AdminListReportsRequest) (*tqdpb.AdminListReportsResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	page := int(req.Page)
	size := int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	filter := repo.AdminReportFilter{
		Page:   page,
		Limit:  size,
		Q:      req.Q,
		UserID: req.UserId,
	}
	if req.ReportType != nil {
		filter.ReportType = req.ReportType
	}
	if req.Status != nil {
		filter.Status = req.Status
	}
	if req.ProblemReport != nil {
		filter.ProblemReport = req.ProblemReport
	}
	if req.QaStatus != nil {
		filter.QaStatus = req.QaStatus
	}
	if req.Severity != nil {
		filter.Severity = req.Severity
	}
	if req.AssigneeId != nil {
		filter.AssigneeID = req.AssigneeId
	}

	reports, total, err := h.reportUsecase.AdminList(ctx, filter)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	displays := h.reportUsecase.ResolveLinkageDisplays(ctx, reports)
	users := h.reportUsecase.ResolveUserDisplays(ctx, reports)
	items := make([]*tqdpb.AdminReportItem, 0, len(reports))
	for _, r := range reports {
		items = append(items, toAdminReportItem(r, displays[r.ID], users[r.ID], false))
	}
	return &tqdpb.AdminListReportsResponse{Data: items, Total: total}, nil
}

// AdminGetReport — GET /v2/tqd/admin/reports/{reportId}
func (h *ReportGrpcHandler) AdminGetReport(ctx context.Context, req *tqdpb.AdminGetReportRequest) (*tqdpb.AdminReportItem, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	report, err := h.reportUsecase.AdminGet(ctx, req.ReportId)
	if err != nil {
		return nil, _errors.ReturnError(404, err.Error())
	}
	displays := h.reportUsecase.ResolveLinkageDisplays(ctx, []domain.Report{*report})
	users := h.reportUsecase.ResolveUserDisplays(ctx, []domain.Report{*report})
	return toAdminReportItem(*report, displays[report.ID], users[report.ID], true), nil
}

// AdminReportSummary — GET /v2/tqd/admin/reports/summary
func (h *ReportGrpcHandler) AdminReportSummary(ctx context.Context, _ *tqdpb.AdminReportSummaryRequest) (*tqdpb.AdminReportSummaryResponse, error) {
	if err := h.requirePerm(ctx, permReportView); err != nil {
		return nil, err
	}
	m, err := h.reportUsecase.AdminSummary(ctx)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &tqdpb.AdminReportSummaryResponse{ByQaStatus: m}, nil
}

// AdminCreateInternalReport — POST /v2/tqd/admin/reports/internal
func (h *ReportGrpcHandler) AdminCreateInternalReport(ctx context.Context, req *tqdpb.AdminCreateInternalReportRequest) (*tqdpb.AdminReportItem, error) {
	if err := h.requirePerm(ctx, permReportCreate); err != nil {
		return nil, err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	report, err := h.reportUsecase.AdminCreateInternal(ctx, actorID, dto.AdminCreateReportRequest{
		Title:           req.Title,
		Description:     req.Description,
		ProblemReport:   req.ProblemReport,
		ReportType:      req.ReportType,
		Severity:        req.Severity,
		TargetID:        req.TargetId,
		SupportTicketID: req.SupportTicketId,
		Images:          req.Images,
	})
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	users := h.reportUsecase.ResolveUserDisplays(ctx, []domain.Report{*report})
	return toAdminReportItem(*report, usecase.LinkageDisplay{}, users[report.ID], true), nil
}

// AdminAssignReport — POST /v2/tqd/admin/reports/{reportId}/assign
func (h *ReportGrpcHandler) AdminAssignReport(ctx context.Context, req *tqdpb.AdminAssignReportRequest) (*tqdpb.AdminActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	if err := h.reportUsecase.AdminAssign(ctx, req.ReportId, req.AssigneeId); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &tqdpb.AdminActionResponse{Id: req.ReportId, Message: "assigned"}, nil
}

// AdminUpdateSeverity — PATCH /v2/tqd/admin/reports/{reportId}/severity
func (h *ReportGrpcHandler) AdminUpdateSeverity(ctx context.Context, req *tqdpb.AdminUpdateSeverityRequest) (*tqdpb.AdminActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	if err := h.reportUsecase.AdminUpdateSeverity(ctx, req.ReportId, req.Severity); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &tqdpb.AdminActionResponse{Id: req.ReportId, Message: "severity updated"}, nil
}

// AdminUpdateQaStatus — PATCH /v2/tqd/admin/reports/{reportId}/qa-status
func (h *ReportGrpcHandler) AdminUpdateQaStatus(ctx context.Context, req *tqdpb.AdminUpdateQaStatusRequest) (*tqdpb.AdminActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	if err := h.reportUsecase.AdminUpdateQaStatus(ctx, req.ReportId, req.QaStatus, req.Note); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &tqdpb.AdminActionResponse{Id: req.ReportId, Message: "qa status updated"}, nil
}

// AdminUpdateImages — POST /v2/tqd/admin/reports/{reportId}/images
func (h *ReportGrpcHandler) AdminUpdateImages(ctx context.Context, req *tqdpb.AdminUpdateImagesRequest) (*tqdpb.AdminActionResponse, error) {
	if err := h.requirePerm(ctx, permReportProcess); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	if err := h.reportUsecase.AdminUpdateImages(ctx, req.ReportId, req.Images); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &tqdpb.AdminActionResponse{Id: req.ReportId, Message: "images updated"}, nil
}

// AdminCloseReport — POST /v2/tqd/admin/reports/{reportId}/close
func (h *ReportGrpcHandler) AdminCloseReport(ctx context.Context, req *tqdpb.AdminCloseReportRequest) (*tqdpb.AdminActionResponse, error) {
	if err := h.requirePerm(ctx, permReportClose); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	if err := h.reportUsecase.AdminClose(ctx, req.ReportId, req.Reject, req.Note); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	return &tqdpb.AdminActionResponse{Id: req.ReportId, Message: "closed"}, nil
}

// AdminDeleteReport — DELETE /v2/tqd/admin/reports/{reportId}
func (h *ReportGrpcHandler) AdminDeleteReport(ctx context.Context, req *tqdpb.AdminDeleteReportRequest) (*emptypb.Empty, error) {
	if err := h.requirePerm(ctx, permReportClose); err != nil {
		return nil, err
	}
	if req.ReportId == 0 {
		return nil, _errors.ReturnError(400, "report_id is required")
	}
	if err := h.reportUsecase.AdminDelete(ctx, req.ReportId); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func toAdminReportItem(r domain.Report, display usecase.LinkageDisplay, user usecase.UserDisplay, full bool) *tqdpb.AdminReportItem {
	item := &tqdpb.AdminReportItem{
		Id:            r.ID,
		UserId:        r.UserID,
		Title:         r.Title,
		ProblemReport: uint32(r.ProblemReport),
		ReportType:    uint32(r.ReportType),
		Status:        uint32(r.Status),
		QaStatus:      uint32(r.QaStatus),
		Severity:      uint32(r.Severity),
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
	}
	if full {
		item.Description = r.Description
		item.Images = reportImagesProto(r.Images)
		if st, err := jsonToStruct(r.Linkage); err == nil {
			item.Linkage = st
		}
	}
	if r.AssigneeID != nil {
		item.AssigneeId = r.AssigneeID
	}
	if r.TargetID != nil {
		item.TargetId = r.TargetID
	}
	if r.SupportTicketID != nil {
		item.SupportTicketId = r.SupportTicketID
	}
	if r.ClosedAt != nil {
		s := r.ClosedAt.Format(time.RFC3339)
		item.ClosedAt = &s
	}
	applyLinkageSummaryProto(item, r.Linkage, r.TargetID, display)
	if user.UserName != "" {
		item.UserName = user.UserName
	}
	if user.UserFullName != "" {
		item.UserFullName = user.UserFullName
	}
	if user.UserUsername != "" {
		item.UserUsername = user.UserUsername
	}
	if user.AssigneeName != "" {
		item.AssigneeName = user.AssigneeName
	}
	if user.AssigneeFullName != "" {
		item.AssigneeFullName = user.AssigneeFullName
	}
	if user.AssigneeUsername != "" {
		item.AssigneeUsername = user.AssigneeUsername
	}
	return item
}

func reportImagesProto(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

func jsonToStruct(raw datatypes.JSON) (*structpb.Struct, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}

func applyLinkageSummaryProto(item *tqdpb.AdminReportItem, raw datatypes.JSON, targetID *uint64, display usecase.LinkageDisplay) {
	parcelOrRegion := ""
	if targetID != nil {
		parcelOrRegion = strconv.FormatUint(*targetID, 10)
	}
	var m map[string]interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &m)
	}
	if m != nil {
		if parcelOrRegion == "" {
			parcelOrRegion = firstStringProto(m, "parcelOrRegion", "parcelId", "regionId", "targetId")
		}
		if v := firstStringProto(m, "layerName", "layer"); v != "" {
			item.LayerName = v
		}
		if v := firstStringProto(m, "layerId"); v != "" {
			item.LayerId = v
		}
		if v := firstStringProto(m, "version", "layerVersion", "versionNo"); v != "" {
			item.Version = v
		}
		if v := firstStringProto(m, "legalStatusName", "legalStatus"); v != "" {
			item.LegalStatus = v
		}
		if v := firstStringProto(m, "parcelName"); v != "" {
			item.ParcelName = v
		}
	}
	if parcelOrRegion != "" {
		item.ParcelOrRegion = parcelOrRegion
	}
	if display.ParcelName != "" {
		item.ParcelName = display.ParcelName
	}
	if display.LayerName != "" {
		item.LayerName = display.LayerName
	}
}

func firstStringProto(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch t := v.(type) {
			case string:
				if strings.TrimSpace(t) != "" {
					return strings.TrimSpace(t)
				}
			case float64:
				return strconv.FormatInt(int64(t), 10)
			case json.Number:
				return t.String()
			default:
				s := strings.TrimSpace(fmt.Sprint(t))
				if s != "" && s != "<nil>" {
					return s
				}
			}
		}
	}
	return ""
}
