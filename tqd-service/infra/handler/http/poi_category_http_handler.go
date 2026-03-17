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

// PoiCategoryHandler handles HTTP requests for POI category
type PoiCategoryHandler struct {
	usecase usecase.PoiCategoryUsecase
	mapper  *mapper.PoiCategoryMapper
}

// NewPoiCategoryHandler creates a new PoiCategoryHandler
func NewPoiCategoryHandler(
	usecase usecase.PoiCategoryUsecase,
	mapper *mapper.PoiCategoryMapper,
) *PoiCategoryHandler {
	return &PoiCategoryHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// CreatePoiCategoryHTTP handles HTTP POST request for creating POI category
// @Summary Create a new POI category
// @Description Create a new POI category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePoiCategoryRequest true "POI category"
// @Success 200 {object} map[string]interface{} "POI category created successfully"
// @Router /v2/tqd/poi/categories [post]
func (h *PoiCategoryHandler) CreatePoiCategoryHTTP(c *gin.Context) {
	var req dto.CreatePoiCategoryRequest
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
			"message": "Failed to create POI category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category created successfully",
		"data":    result,
	})
}

// GetPoiCategoryHTTP handles HTTP GET request for getting POI category by ID
// @Summary Get POI category by ID
// @Description Get POI category by ID
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Success 200 {object} map[string]interface{} "POI category"
// @Router /v2/tqd/poi/categories/{id} [get]
func (h *PoiCategoryHandler) GetPoiCategoryHTTP(c *gin.Context) {
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
			"message": "POI category not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category retrieved successfully",
		"data":    result,
	})
}

// GetPoiCategoryByCodeHTTP handles HTTP GET request for getting POI category by code
// @Summary Get POI category by code
// @Description Get POI category by code
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "POI category code"
// @Success 200 {object} map[string]interface{} "POI category"
// @Router /v2/tqd/poi/categories/code/{code} [get]
func (h *PoiCategoryHandler) GetPoiCategoryByCodeHTTP(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Code is required",
		})
		return
	}

	// Note: Need to implement GetByCode in usecase if needed
	// This would require adding GetByCode to PoiCategoryUsecase interface
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "Get by code not implemented yet",
	})
}

// UpdatePoiCategoryHTTP handles HTTP PUT request for updating POI category
// @Summary Update POI category
// @Description Update POI category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Param request body dto.UpdatePoiCategoryRequest true "POI category"
// @Success 200 {object} map[string]interface{} "POI category updated successfully"
// @Router /v2/tqd/poi/categories/{id} [put]
func (h *PoiCategoryHandler) UpdatePoiCategoryHTTP(c *gin.Context) {
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

	var req dto.UpdatePoiCategoryRequest
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
			"message": "Failed to update POI category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category updated successfully",
		"data":    result,
	})
}

// DeletePoiCategoryHTTP handles HTTP DELETE request for deleting POI category
// @Summary Delete POI category
// @Description Delete POI category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Success 200 {object} map[string]interface{} "POI category deleted successfully"
// @Router /v2/tqd/poi/categories/{id} [delete]
func (h *PoiCategoryHandler) DeletePoiCategoryHTTP(c *gin.Context) {
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
			"message": "Failed to delete POI category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category deleted successfully",
	})
}

// ListPoiCategoriesHTTP handles HTTP GET request for listing POI categories
// @Summary List POI categories
// @Description List POI categories with pagination and filters
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param search query string false "Search by name or code"
// @Param code query string false "Filter by code"
// @Param name query string false "Filter by name"
// @Param parentId query int false "Filter by parent ID"
// @Param isActive query bool false "Filter by active status"
// @Param level query int false "Filter by level"
// @Success 200 {object} map[string]interface{} "POI categories list"
// @Router /v2/tqd/poi/categories [get]
func (h *PoiCategoryHandler) ListPoiCategoriesHTTP(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Create filter
	filter := &dto.PoiCategoryFilter{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Search: c.Query("search"),
		Code:   c.Query("code"),
		Name:   c.Query("name"),
	}

	// Parse parentId
	if parentIdStr := c.Query("parentId"); parentIdStr != "" {
		if parentId, err := strconv.ParseUint(parentIdStr, 10, 64); err == nil {
			filter.ParentID = &parentId
		}
	}

	// Parse level
	if levelStr := c.Query("level"); levelStr != "" {
		if level, err := strconv.Atoi(levelStr); err == nil {
			filter.Level = &level
		}
	}

	// Parse isActive
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}

	result, err := h.usecase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list POI categories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI categories retrieved successfully",
		"data": gin.H{
			"items": result.Data,
			"total": result.Total,
			"page":  result.Page,
			"size":  result.Size,
			"pages": (result.Total + int64(result.Size) - 1) / int64(result.Size),
		},
	})
}

// GetPoiCategoryTreeHTTP handles HTTP GET request for getting POI categories in tree structure
// @Summary Get POI categories tree
// @Description Get POI categories in hierarchical tree structure (parent-children)
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "POI categories tree"
// @Router /v2/tqd/poi/categories/tree [get]
func (h *PoiCategoryHandler) GetPoiCategoryTreeHTTP(c *gin.Context) {
	result, err := h.usecase.GetTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get POI category tree",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI category tree retrieved successfully",
		"data":    result,
	})
}

// UpdatePoiCategoryStatsHTTP handles HTTP PUT request for updating category statistics
// @Summary Update POI category statistics
// @Description Update POI count for category
// @Tags POI Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "POI category ID"
// @Param request body UpdateStatsRequest true "Statistics"
// @Success 200 {object} map[string]interface{} "POI category updated"
// @Router /v2/tqd/poi/categories/{id}/stats [put]
func (h *PoiCategoryHandler) UpdatePoiCategoryStatsHTTP(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	var req struct {
		POICount int32 `json:"poiCount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Note: Need to implement UpdateStats in usecase
	// This would require adding methods to PoiCategoryUsecase interface
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "Update stats not implemented yet",
	})
}
