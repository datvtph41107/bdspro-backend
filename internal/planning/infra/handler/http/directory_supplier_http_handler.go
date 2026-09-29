package handler_http

import (
	"net/http"
	"strconv"
	"strings"

	_dto "common/domain/dto"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

// DirectorySupplierHandler handles HTTP requests for directory supplier
type DirectorySupplierHandler struct {
	usecase *usecase.DirectorySupplierUsecase
	mapper  *mapper.DirectorySupplierMapper
}

// NewDirectorySupplierHandler creates new handler
func NewDirectorySupplierHandler(
	usecase *usecase.DirectorySupplierUsecase,
	mapper *mapper.DirectorySupplierMapper,
) *DirectorySupplierHandler {
	return &DirectorySupplierHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// CreateDirectorySupplierHTTP handles HTTP POST request for creating directory supplier
// @Summary Create a new directory supplier
// @Description Create a new directory supplier
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateDirectorySupplierRequestDTO true "Directory supplier"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /v2/tqd/directory/suppliers [post]
func (h *DirectorySupplierHandler) CreateDirectorySupplierHTTP(c *gin.Context) {
	var req dto.CreateDirectorySupplierRequestDTO
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
			"message": "Failed to create directory supplier",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier created successfully",
		"data":    result,
	})
}

// GetDirectorySupplierHTTP handles HTTP GET request for getting directory supplier by ID
// @Summary Get directory supplier by ID
// @Description Get directory supplier by ID
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Directory supplier ID"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /v2/tqd/directory/suppliers/{id} [get]
func (h *DirectorySupplierHandler) GetDirectorySupplierHTTP(c *gin.Context) {
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
			"message": "Directory supplier not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier retrieved successfully",
		"data":    result,
	})
}

// GetDirectorySupplierByCodeHTTP handles HTTP GET request for getting directory supplier by code
// @Summary Get directory supplier by code
// @Description Get directory supplier by code
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Directory supplier code"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /v2/tqd/directory/suppliers/code/{code} [get]
func (h *DirectorySupplierHandler) GetDirectorySupplierByCodeHTTP(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Code is required",
		})
		return
	}

	result, err := h.usecase.GetByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Directory supplier not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier retrieved successfully",
		"data":    result,
	})
}

// UpdateDirectorySupplierHTTP handles HTTP PUT request for updating directory supplier
// @Summary Update directory supplier
// @Description Update directory supplier
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Directory supplier ID"
// @Param request body dto.UpdateDirectorySupplierRequestDTO true "Directory supplier"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /v2/tqd/directory/suppliers/{id} [put]
func (h *DirectorySupplierHandler) UpdateDirectorySupplierHTTP(c *gin.Context) {
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

	var req dto.UpdateDirectorySupplierRequestDTO
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
			"message": "Failed to update directory supplier",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier updated successfully",
		"data":    result,
	})
}

// DeleteDirectorySupplierHTTP handles HTTP DELETE request for deleting directory supplier
// @Summary Delete directory supplier
// @Description Delete directory supplier
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Directory supplier ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Router /v2/tqd/directory/suppliers/{id} [delete]
func (h *DirectorySupplierHandler) DeleteDirectorySupplierHTTP(c *gin.Context) {
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
			"message": "Failed to delete directory supplier",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier deleted successfully",
	})
}

// ListDirectorySuppliersHTTP handles HTTP GET request for listing directory suppliers
// @Summary List directory suppliers
// @Description List directory suppliers with pagination and filters
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Param categories query string false "Categories (comma-separated)"
// @Param isActive query bool false "Filter by active status"
// @Param minRating query number false "Minimum rating"
// @Param maxRating query number false "Maximum rating"
// @Param sortBy query string false "Sort by field (name, rating, created_at)"
// @Param sortOrder query string false "Sort order (asc, desc)"
// @Success 200 {object} map[string]interface{} "Directory suppliers list"
// @Router /v2/tqd/directory/suppliers [get]
func (h *DirectorySupplierHandler) ListDirectorySuppliersHTTP(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Create filter
	filter := &dto.DirectorySupplierFilterDTO{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Search:    c.Query("search"),
		SortBy:    c.Query("sortBy"),
		SortOrder: c.Query("sortOrder"),
	}

	// Parse categories
	if categoriesStr := c.Query("categories"); categoriesStr != "" {
		filter.Categories = strings.Split(categoriesStr, ",")
	}

	// Parse isActive
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &val
		}
	}

	// Parse minRating
	if minRatingStr := c.Query("minRating"); minRatingStr != "" {
		if val, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			filter.MinRating = &val
		}
	}

	// Parse maxRating
	if maxRatingStr := c.Query("maxRating"); maxRatingStr != "" {
		if val, err := strconv.ParseFloat(maxRatingStr, 64); err == nil {
			filter.MaxRating = &val
		}
	}

	result, err := h.usecase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list directory suppliers",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory suppliers retrieved successfully",
		"data": gin.H{
			"items": result.Data,
			"total": result.Total,
			"page":  result.Page,
			"size":  result.Size,
			"pages": (result.Total + int64(result.Size) - 1) / int64(result.Size),
		},
	})
}

// UpdateDirectorySupplierStatsHTTP handles HTTP PUT request for updating supplier stats
// @Summary Update directory supplier statistics
// @Description Update service count, contract count and total value
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Directory supplier ID"
// @Param request body UpdateStatsRequest true "Statistics"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /v2/tqd/directory/suppliers/{id}/stats [put]
func (h *DirectorySupplierHandler) UpdateDirectorySupplierStatsHTTP(c *gin.Context) {
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

	var req struct {
		ServiceCount  uint32 `json:"serviceCount"`
		ContractCount uint32 `json:"contractCount"`
		TotalValue    int64  `json:"totalValue"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.UpdateStats(c.Request.Context(), id, req.ServiceCount, req.ContractCount, req.TotalValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update supplier statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Supplier statistics updated successfully",
		"data":    result,
	})
}
