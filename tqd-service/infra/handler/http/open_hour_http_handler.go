package handler_http

import (
	"net/http"
	"strconv"

	_dto "common/domain/dto"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

// OpenHourHandler handles HTTP requests for open hour
type OpenHourHandler struct {
	usecase usecase.OpenHourUsecase
	mapper  *mapper.OpenHourMapper
}

// NewOpenHourHandler creates a new OpenHourHandler
func NewOpenHourHandler(
	usecase usecase.OpenHourUsecase,
	mapper *mapper.OpenHourMapper,
) *OpenHourHandler {
	return &OpenHourHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// CreateOpenHourHTTP handles HTTP POST request for creating open hour
// @Summary Create a new open hour
// @Description Create a new open hour entry
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateOpenHourRequest true "Open hour data"
// @Success 200 {object} map[string]interface{} "Open hour created successfully"
// @Router /v2/tqd/open-hours [post]
func (h *OpenHourHandler) CreateOpenHourHTTP(c *gin.Context) {
	var req dto.CreateOpenHourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), &req)
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
		"data":    result,
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
			"message": "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.GetByID(c.Request.Context(), id)
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
		"data":    result,
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
// @Param request body dto.UpdateOpenHourRequest true "Open hour data"
// @Success 200 {object} map[string]interface{} "Open hour updated successfully"
// @Router /v2/tqd/open-hours/{id} [put]
func (h *OpenHourHandler) UpdateOpenHourHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	var req dto.UpdateOpenHourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.Update(c.Request.Context(), id, &req)
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
		"data":    result,
	})
}

// DeleteOpenHourHTTP handles HTTP DELETE request for deleting open hour
// @Summary Delete open hour
// @Description Delete open hour by ID
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
			"message": "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
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
	})
}

// ListOpenHoursHTTP handles HTTP GET request for listing open hours
// @Summary List open hours
// @Description List open hours with pagination and filters
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param poiId query int false "Filter by POI ID"
// @Param dayOfWeek query int false "Filter by day of week (2-8)"
// @Param isOpen query bool false "Filter by is open status"
// @Success 200 {object} map[string]interface{} "Open hours list"
// @Router /v2/tqd/open-hours [get]
func (h *OpenHourHandler) ListOpenHoursHTTP(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Create filter
	filter := &dto.OpenHourFilter{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
	}

	// Parse poiId
	if poiIdStr := c.Query("poiId"); poiIdStr != "" {
		if poiId, err := strconv.ParseUint(poiIdStr, 10, 64); err == nil {
			filter.POIID = poiId
		}
	}

	// Parse dayOfWeek
	if dayStr := c.Query("dayOfWeek"); dayStr != "" {
		if day, err := strconv.Atoi(dayStr); err == nil && day >= 2 && day <= 8 {
			filter.DayOfWeek = &day
		}
	}

	// Parse isOpen
	if isOpenStr := c.Query("isOpen"); isOpenStr != "" {
		if isOpen, err := strconv.ParseBool(isOpenStr); err == nil {
			filter.IsOpen = &isOpen
		}
	}

	result, err := h.usecase.List(c.Request.Context(), filter)
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
		"data": gin.H{
			"items": result.Data,
			"total": result.Total,
			"page":  result.Page,
			"size":  result.Size,
			"pages": (result.Total + int64(result.Size) - 1) / int64(result.Size),
		},
	})
}

// GetOpenHourTimesByDayHTTP handles HTTP GET request for getting open hours grouped by day
// @Summary Get open hours grouped by day
// @Description Get open hours grouped by day of week
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dayOfWeek query string true "Day of week (Monday, Tuesday, etc.)"
// @Param poiId query int false "Filter by POI ID"
// @Success 200 {object} map[string]interface{} "Open hours by day"
// @Router /v2/tqd/open-hours/times/by-day [get]
func (h *OpenHourHandler) GetOpenHourTimesByDayHTTP(c *gin.Context) {
	dayOfWeek := c.Query("dayOfWeek")
	if dayOfWeek == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "dayOfWeek is required",
		})
		return
	}

	req := &dto.GetTimesByDayRequest{
		DayOfWeek: dayOfWeek,
	}

	if poiIdStr := c.Query("poiId"); poiIdStr != "" {
		if poiId, err := strconv.ParseUint(poiIdStr, 10, 64); err == nil {
			req.POIID = poiId
		}
	}

	result, err := h.usecase.GetTimesByDay(c.Request.Context(), req)
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

// BulkSaveOpenHoursHTTP handles HTTP POST request for bulk operations
// @Summary Bulk save/delete open hours
// @Description Create, update, and delete multiple open hours at once
// @Tags Open Hour
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BulkSaveOpenHourRequest true "Bulk save request"
// @Success 200 {object} map[string]interface{} "Bulk operation result"
// @Router /v2/tqd/open-hours/bulk [post]
func (h *OpenHourHandler) BulkSaveOpenHoursHTTP(c *gin.Context) {
	var req dto.BulkSaveOpenHourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.BulkSave(c.Request.Context(), &req)
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
