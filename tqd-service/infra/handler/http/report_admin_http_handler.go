package handler_http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_utils "common/utils"
	"pb/clients"
	"tqd/env"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// Quyền theo hành động. Key ADMIN_BC_GIS_SUA/ADMIN_BC_GIS_XOA là key cũ trước khi
// tách Xử lý / Đóng, vẫn nhận cho tới khi mọi môi trường chạy script rename permission.
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

type ReportAdminHTTPHandler struct {
	usecase    usecase.ReportUsecase
	authClient *clients.AuthGrpcClient
}

func NewReportAdminHTTPHandler(usecase usecase.ReportUsecase) *ReportAdminHTTPHandler {
	return &ReportAdminHTTPHandler{
		usecase:    usecase,
		authClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

// ctxFromGin: JWT middleware gắn profileId trên Gin keys, không có trên Request.Context().
// HasPermissions (gRPC) cần ProfileIDKey trên context — dùng CreateContextFromGin.
func ctxFromGin(c *gin.Context) context.Context {
	return _utils.CreateContextFromGin(c)
}

func (h *ReportAdminHTTPHandler) List(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	filter := repo.AdminReportFilter{Page: page, Limit: size, Q: c.Query("q")}
	if v := c.Query("problemReport"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 32)
		u := uint32(n)
		filter.ProblemReport = &u
	}
	if v := c.Query("qaStatus"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 32)
		u := uint32(n)
		filter.QaStatus = &u
	}
	if v := c.Query("severity"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 32)
		u := uint32(n)
		filter.Severity = &u
	}
	if v := c.Query("assigneeId"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 64)
		filter.AssigneeID = &n
	}
	reports, total, err := h.usecase.AdminList(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	displays := h.usecase.ResolveLinkageDisplays(ctx, reports)
	users := h.usecase.ResolveUserDisplays(ctx, reports)
	data := make([]gin.H, 0, len(reports))
	for _, r := range reports {
		data = append(data, reportToListJSON(r, displays[r.ID], users[r.ID]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "total": total})
}

func (h *ReportAdminHTTPHandler) Summary(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	m, err := h.usecase.AdminSummary(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"byQaStatus": m})
}

func (h *ReportAdminHTTPHandler) Queue(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	bucket := c.DefaultQuery("bucket", "mine")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	var severity *uint32
	if v := c.Query("severity"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 32)
		u := uint32(n)
		severity = &u
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	reports, total, err := h.usecase.AdminQueue(ctx, bucket, severity, page, size, actorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	displays := h.usecase.ResolveLinkageDisplays(ctx, reports)
	users := h.usecase.ResolveUserDisplays(ctx, reports)
	data := make([]gin.H, 0, len(reports))
	for _, r := range reports {
		data = append(data, reportToListJSON(r, displays[r.ID], users[r.ID]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "total": total, "bucket": bucket})
}

func (h *ReportAdminHTTPHandler) QueueSummary(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	sum, err := h.usecase.AdminQueueSummary(ctx, actorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"mine":       sum.Mine,
		"unassigned": sum.Unassigned,
		"dataFix":    sum.DataFix,
	})
}

func (h *ReportAdminHTTPHandler) ListEvents(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	filter := repo.AdminEventFilter{}
	filter.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("size", "20"))
	if v := c.Query("reportId"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 64)
		filter.ReportID = &n
	}
	if v := c.Query("actorId"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 64)
		filter.ActorID = &n
	}
	filter.Action = strings.TrimSpace(c.Query("action"))
	if v := c.Query("from"); v != "" {
		filter.From = &v
	}
	if v := c.Query("to"); v != "" {
		filter.To = &v
	}
	events, total, err := h.usecase.AdminListEvents(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": eventsToJSON(h.usecase, ctx, events), "total": total})
}

func (h *ReportAdminHTTPHandler) ListReportEvents(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	events, err := h.usecase.AdminListReportEvents(ctx, id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": eventsToJSON(h.usecase, ctx, events)})
}

func eventsToJSON(u usecase.ReportUsecase, ctx context.Context, events []domain.ReportEvent) []gin.H {
	actorIDs := make([]uint64, 0, len(events))
	for _, e := range events {
		actorIDs = append(actorIDs, e.ActorID)
	}
	names := u.ResolveActorNames(ctx, actorIDs)
	out := make([]gin.H, 0, len(events))
	for _, e := range events {
		item := gin.H{
			"id":        e.ID,
			"reportId":  e.ReportID,
			"actorId":   e.ActorID,
			"action":    e.Action,
			"note":      e.Note,
			"createdAt": e.CreatedAt.Format(time.RFC3339),
		}
		if e.BeforeJSON != "" {
			item["beforeJson"] = e.BeforeJSON
		}
		if e.AfterJSON != "" {
			item["afterJson"] = e.AfterJSON
		}
		if n := names[e.ActorID]; n != "" {
			item["actorName"] = n
		}
		out = append(out, item)
	}
	return out
}

func (h *ReportAdminHTTPHandler) Get(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportView); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	report, err := h.usecase.AdminGet(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	displays := h.usecase.ResolveLinkageDisplays(ctx, []domain.Report{*report})
	users := h.usecase.ResolveUserDisplays(ctx, []domain.Report{*report})
	c.JSON(http.StatusOK, reportToJSON(*report, displays[report.ID], users[report.ID]))
}

func (h *ReportAdminHTTPHandler) CreateInternal(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportCreate); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	var req dto.AdminCreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	report, err := h.usecase.AdminCreateInternal(ctx, actorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	users := h.usecase.ResolveUserDisplays(ctx, []domain.Report{*report})
	c.JSON(http.StatusOK, reportToJSON(*report, usecase.LinkageDisplay{}, users[report.ID]))
}

func (h *ReportAdminHTTPHandler) Assign(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportProcess); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminAssign(ctx, id, req.AssigneeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "assigned"})
}

func (h *ReportAdminHTTPHandler) UpdateSeverity(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportProcess); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminSeverityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateSeverity(ctx, id, req.Severity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "severity updated"})
}

func (h *ReportAdminHTTPHandler) UpdateQaStatus(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportProcess); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminQaStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateQaStatus(ctx, id, req.QaStatus, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "qa status updated"})
}

func (h *ReportAdminHTTPHandler) UpdateLinkage(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportProcess); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var linkage map[string]interface{}
	if err := c.ShouldBindJSON(&linkage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateLinkage(ctx, id, linkage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "linkage updated"})
}

func (h *ReportAdminHTTPHandler) UpdateImages(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportProcess); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateImages(ctx, id, req.Images); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "images updated"})
}

func (h *ReportAdminHTTPHandler) Close(c *gin.Context) {
	ctx := ctxFromGin(c)
	if err := h.authClient.HasPermissions(ctx, permReportClose); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminCloseRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.usecase.AdminClose(ctx, id, req.Reject, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "closed"})
}

func reportImages(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

// reportToListJSON — payload nhẹ cho bảng danh sách (không description/images đầy đủ).
func reportToListJSON(r domain.Report, display usecase.LinkageDisplay, user usecase.UserDisplay) gin.H {
	out := gin.H{
		"id":            r.ID,
		"userId":        r.UserID,
		"title":         r.Title,
		"problemReport": uint32(r.ProblemReport),
		"reportType":    uint32(r.ReportType),
		"status":        uint32(r.Status),
		"qaStatus":      uint32(r.QaStatus),
		"severity":      uint32(r.Severity),
		"createdAt":     r.CreatedAt.Format(time.RFC3339),
	}
	if r.AssigneeID != nil {
		out["assigneeId"] = *r.AssigneeID
	}
	if r.SupportTicketID != nil {
		out["supportTicketId"] = *r.SupportTicketID
	}
	if r.TargetID != nil {
		out["targetId"] = *r.TargetID
	}
	if r.ClosedAt != nil {
		out["closedAt"] = r.ClosedAt.Format(time.RFC3339)
	}
	summary := linkageSummary(r.Linkage, r.TargetID, display)
	for k, v := range summary {
		out[k] = v
	}
	applyUserDisplay(out, user)
	return out
}

// linkageSummary extracts light fields for list columns (II.4).
func linkageSummary(raw datatypes.JSON, targetID *uint64, display usecase.LinkageDisplay) gin.H {
	out := gin.H{}
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
			parcelOrRegion = firstString(m, "parcelOrRegion", "parcelId", "regionId", "targetId")
		}
		if v := firstString(m, "layerName", "layer"); v != "" {
			out["layerName"] = v
		}
		if v := firstString(m, "layerId"); v != "" {
			out["layerId"] = v
		}
		if v := firstString(m, "version", "layerVersion", "versionNo"); v != "" {
			out["version"] = v
		}
		if v := firstString(m, "legalStatusName", "legalStatus"); v != "" {
			out["legalStatus"] = v
		}
		if v := firstString(m, "parcelName"); v != "" {
			out["parcelName"] = v
		}
	}
	if parcelOrRegion != "" {
		out["parcelOrRegion"] = parcelOrRegion
	}
	if display.ParcelName != "" {
		out["parcelName"] = display.ParcelName
	}
	if display.LayerName != "" {
		out["layerName"] = display.LayerName
	}
	return out
}

func firstString(m map[string]interface{}, keys ...string) string {
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

func reportToJSON(r domain.Report, display usecase.LinkageDisplay, user usecase.UserDisplay) gin.H {
	out := gin.H{
		"id":            r.ID,
		"userId":        r.UserID,
		"title":         r.Title,
		"description":   r.Description,
		"problemReport": uint32(r.ProblemReport),
		"reportType":    uint32(r.ReportType),
		"status":        uint32(r.Status),
		"qaStatus":      uint32(r.QaStatus),
		"severity":      uint32(r.Severity),
		"createdAt":     r.CreatedAt.Format(time.RFC3339),
		"linkage":       r.Linkage,
		"images":        reportImages(r.Images),
	}
	if r.AssigneeID != nil {
		out["assigneeId"] = *r.AssigneeID
	}
	if r.TargetID != nil {
		out["targetId"] = *r.TargetID
	}
	if r.SupportTicketID != nil {
		out["supportTicketId"] = *r.SupportTicketID
	}
	if r.ClosedAt != nil {
		out["closedAt"] = r.ClosedAt.Format(time.RFC3339)
	}
	summary := linkageSummary(r.Linkage, r.TargetID, display)
	for k, v := range summary {
		out[k] = v
	}
	applyUserDisplay(out, user)
	return out
}

func applyUserDisplay(out gin.H, user usecase.UserDisplay) {
	if user.UserName != "" {
		out["userName"] = user.UserName
	}
	if user.UserFullName != "" {
		out["userFullName"] = user.UserFullName
	}
	if user.UserUsername != "" {
		out["userUsername"] = user.UserUsername
	}
	if user.AssigneeName != "" {
		out["assigneeName"] = user.AssigneeName
	}
	if user.AssigneeFullName != "" {
		out["assigneeFullName"] = user.AssigneeFullName
	}
	if user.AssigneeUsername != "" {
		out["assigneeUsername"] = user.AssigneeUsername
	}
}
