package handler_http

import (
	"net/http"
	"strconv"

	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OpenHourHandler struct {
	openHourUsecase usecase.OpenHourUsecase
}

// NewOpenHourHandler creates a new OpenHourHandler
// @bind: *handler_http.OpenHourHandler
func NewOpenHourHandler(
	openHourUsecase usecase.OpenHourUsecase,
) *OpenHourHandler {
	return &OpenHourHandler{
		openHourUsecase: openHourUsecase,
	}
}

// CreateOpenHourHTTP handles HTTP POST request for creating open hour
// @Summary Create a new open hour
// @Description Create a new open hour entry
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.OpenHourCreateRequest true "Open hour data"
// @Success 200 {object} map[string]interface{} "Open hour created successfully"
// @Router /v2/tqd/open-hours [post]
func (h *OpenHourHandler) CreateOpenHourHTTP(c *gin.Context) {
	var req dto.OpenHourCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Map request to domain
	openHour := &domain.OpenHour{
		Name:      req.Name,
		POIID:     req.POIID,
		DayOfWeek: req.DayOfWeek,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
		Note:      req.Note,
		IsOpen:    req.IsOpen,
		Status:    enums.OpenHourStatus(req.Status),
		Type:      enums.OpenHourType(req.Type),
	}

	createdOpenHour, err := h.openHourUsecase.Create(c.Request.Context(), openHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create open hour",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Open hour created successfully",
		"data":    createdOpenHour,
	})
}

// GetOpenHourHTTP handles HTTP GET request for getting open hour by ID
// @Summary Get open hour by ID
// @Description Get open hour details by ID
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Open hour ID"
// @Success 200 {object} map[string]interface{} "Open hour details"
// @Router /v2/tqd/open-hours/{id} [get]
func (h *OpenHourHandler) GetOpenHourHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID",
			"details": err.Error(),
		})
		return
	}

	openHour, err := h.openHourUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Open hour not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Open hour retrieved successfully",
		"data":    openHour,
	})
}

// UpdateOpenHourHTTP handles HTTP PUT request for updating open hour
// @Summary Update open hour
// @Description Update open hour details
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Open hour ID"
// @Param request body dto.OpenHourUpdateRequest true "Open hour data"
// @Success 200 {object} map[string]interface{} "Open hour updated successfully"
// @Router /v2/tqd/open-hours/{id} [put]
func (h *OpenHourHandler) UpdateOpenHourHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID",
			"details": err.Error(),
		})
		return
	}

	var req dto.OpenHourUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing open hour
	existingOpenHour, err := h.openHourUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Open hour not found",
			"details": err.Error(),
		})
		return
	}

	// Update fields if provided
	if req.Name != nil {
		existingOpenHour.Name = *req.Name
	}
	if req.POIID != nil {
		existingOpenHour.POIID = req.POIID
	}
	if req.DayOfWeek != nil {
		existingOpenHour.DayOfWeek = *req.DayOfWeek
	}
	if req.OpenTime != nil {
		existingOpenHour.OpenTime = *req.OpenTime
	}
	if req.CloseTime != nil {
		existingOpenHour.CloseTime = *req.CloseTime
	}
	if req.Note != nil {
		existingOpenHour.Note = *req.Note
	}
	if req.IsOpen != nil {
		existingOpenHour.IsOpen = *req.IsOpen
	}
	if req.Status != nil {
		existingOpenHour.Status = enums.OpenHourStatus(*req.Status)
	}
	if req.Type != nil {
		existingOpenHour.Type = enums.OpenHourType(*req.Type)
	}

	updatedOpenHour, err := h.openHourUsecase.Update(c.Request.Context(), id, existingOpenHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update open hour",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Open hour updated successfully",
		"data":    updatedOpenHour,
	})
}

// DeleteOpenHourHTTP handles HTTP DELETE request for deleting open hour
// @Summary Delete open hour
// @Description Delete open hour by ID (soft delete)
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Open hour ID"
// @Success 200 {object} map[string]interface{} "Open hour deleted successfully"
// @Router /v2/tqd/open-hours/{id} [delete]
func (h *OpenHourHandler) DeleteOpenHourHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID",
			"details": err.Error(),
		})
		return
	}

	ok, err := h.openHourUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete open hour",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Open hour deleted successfully",
		"data":    ok,
	})
}

// ListOpenHoursHTTP handles HTTP GET request for listing open hours
// @Summary List open hours
// @Description List open hours with pagination and filters
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param poiId query int false "Filter by POI ID"
// @Param dayOfWeek query int false "Filter by day of week (2-8)"
// @Param isOpen query bool false "Filter by is open status"
// @Param status query int false "Filter by status"
// @Param type query int false "Filter by type"
// @Success 200 {object} map[string]interface{} "Open hours list"
// @Router /v2/tqd/open-hours [get]
func (h *OpenHourHandler) ListOpenHoursHTTP(c *gin.Context) {
	var req dto.OpenHourListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, total, err := h.openHourUsecase.GetList(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list open hours",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Open hours retrieved successfully",
		"data":    response,
		"total":   total,
	})
}

// GetOpenHourTimesByDayHTTP handles HTTP GET request for getting open hours grouped by day of week
// @Summary Get open hours grouped by day of week
// @Description Get unique days with their time slots
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param poiId query int false "Filter by POI ID"
// @Success 200 {object} map[string]interface{} "Open hours grouped by day"
// @Router /v2/tqd/open-hours/times [get]
func (h *OpenHourHandler) GetOpenHourTimesByDayHTTP(c *gin.Context) {
	var poiID *uint64
	if poiIDStr := c.Query("poiId"); poiIDStr != "" {
		id, err := strconv.ParseUint(poiIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid POI ID",
				"details": err.Error(),
			})
			return
		}
		poiID = &id
	}

	result, err := h.openHourUsecase.GetTimesByDayOfWeek(c.Request.Context(), poiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get open hours by day",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Open hours by day retrieved successfully",
		"data":    result,
	})
}

// BulkSaveOpenHoursHTTP handles HTTP POST request for bulk operations (create, update, delete)
// @Summary Bulk save/delete open hours
// @Description Create, update, and delete multiple open hours at once. Items with ID will be updated, without ID will be created, and IDs in deleteIds will be deleted
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.OpenHourBulkSaveRequest true "Bulk save request"
// @Success 200 {object} map[string]interface{} "Bulk operation result"
// @Router /v2/tqd/open-hours/bulk [post]
func (h *OpenHourHandler) BulkSaveOpenHoursHTTP(c *gin.Context) {
	var req dto.OpenHourBulkSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.openHourUsecase.BulkSave(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to bulk save open hours",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Bulk operation completed successfully",
		"data":    result,
	})
}
