package handler_http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	_utils "common/utils"
	"pb/clients"
	"tqd/env"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
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
	data := make([]gin.H, 0, len(reports))
	for _, r := range reports {
		data = append(data, reportToListJSON(r))
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
	c.JSON(http.StatusOK, reportToJSON(*report))
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
	c.JSON(http.StatusOK, reportToJSON(*report))
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
	if err := h.usecase.AdminUpdateQaStatus(ctx, id, req.QaStatus); err != nil {
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

// reportToListJSON — payload nhẹ cho bảng danh sách (không description/images/linkage).
func reportToListJSON(r domain.Report) gin.H {
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
	if r.ClosedAt != nil {
		out["closedAt"] = r.ClosedAt.Format(time.RFC3339)
	}
	return out
}

func reportToJSON(r domain.Report) gin.H {
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
	return out
}
