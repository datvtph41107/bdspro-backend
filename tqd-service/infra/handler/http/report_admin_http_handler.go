package handler_http

import (
	"net/http"
	"strconv"
	"time"

	_utils "common/utils"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ReportAdminHTTPHandler struct {
	usecase usecase.ReportUsecase
}

func NewReportAdminHTTPHandler(usecase usecase.ReportUsecase) *ReportAdminHTTPHandler {
	return &ReportAdminHTTPHandler{usecase: usecase}
}

func (h *ReportAdminHTTPHandler) List(c *gin.Context) {
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
	reports, total, err := h.usecase.AdminList(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := make([]gin.H, 0, len(reports))
	for _, r := range reports {
		data = append(data, reportToJSON(r))
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "total": total})
}

func (h *ReportAdminHTTPHandler) Summary(c *gin.Context) {
	m, err := h.usecase.AdminSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"byQaStatus": m})
}

func (h *ReportAdminHTTPHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	report, err := h.usecase.AdminGet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reportToJSON(*report))
}

func (h *ReportAdminHTTPHandler) CreateInternal(c *gin.Context) {
	var req dto.AdminCreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actorID := _utils.GetProfileIdWithContext(c.Request.Context())
	report, err := h.usecase.AdminCreateInternal(c.Request.Context(), actorID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reportToJSON(*report))
}

func (h *ReportAdminHTTPHandler) Assign(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminAssign(c.Request.Context(), id, req.AssigneeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "assigned"})
}

func (h *ReportAdminHTTPHandler) UpdateSeverity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminSeverityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateSeverity(c.Request.Context(), id, req.Severity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "severity updated"})
}

func (h *ReportAdminHTTPHandler) UpdateQaStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminQaStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateQaStatus(c.Request.Context(), id, req.QaStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "qa status updated"})
}

func (h *ReportAdminHTTPHandler) UpdateLinkage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var linkage map[string]interface{}
	if err := c.ShouldBindJSON(&linkage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.usecase.AdminUpdateLinkage(c.Request.Context(), id, linkage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "linkage updated"})
}

func (h *ReportAdminHTTPHandler) Close(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.AdminCloseRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.usecase.AdminClose(c.Request.Context(), id, req.Reject, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "closed"})
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
