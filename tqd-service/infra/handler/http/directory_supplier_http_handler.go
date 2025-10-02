package handler_http

import (
	"net/http"
	"strconv"
	"strings"

	"tqd/infra/mapper"
	"tqd/internal/domain"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

type DirectorySupplierHandler struct {
	directorySupplierUsecase *usecase.DirectorySupplierUsecase
	mapper                   *mapper.DirectorySupplierMapper
}

func NewDirectorySupplierHandler(
	directorySupplierUsecase *usecase.DirectorySupplierUsecase,
	mapper *mapper.DirectorySupplierMapper,
) *DirectorySupplierHandler {
	return &DirectorySupplierHandler{
		directorySupplierUsecase: directorySupplierUsecase,
		mapper:                   mapper,
	}
}

// CreateDirectorySupplierHTTP handles HTTP POST request for creating directory supplier
// @Summary Create a new directory supplier
// @Description Create a new directory supplier
// @Tags Directory Supplier
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.DirectorySupplier true "Directory supplier"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /v2/tqd/directory/suppliers [post]
func (h *DirectorySupplierHandler) CreateDirectorySupplierHTTP(c *gin.Context) {
	var supplier domain.DirectorySupplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	createdSupplier, err := h.directorySupplierUsecase.Create(c.Request.Context(), &supplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create directory supplier",
			"details": err.Error(),
		})
		return
	}

	response := h.mapper.ToResponse(createdSupplier)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier created successfully",
		"data":    response,
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

	supplier, err := h.directorySupplierUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Directory supplier not found",
			"details": err.Error(),
		})
		return
	}

	response := h.mapper.ToResponse(supplier)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier retrieved successfully",
		"data":    response,
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
// @Param request body domain.DirectorySupplier true "Directory supplier"
// @Success 200 {object} map[string]interface{} "Directory supplier"
// @Router /directory/suppliers/{id} [put]
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

	var supplier domain.DirectorySupplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	supplier.ID = id
	updatedSupplier, err := h.directorySupplierUsecase.Update(c.Request.Context(), &supplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update directory supplier",
			"details": err.Error(),
		})
		return
	}

	response := h.mapper.ToResponse(updatedSupplier)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Directory supplier updated successfully",
		"data":    response,
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
// @Success 200 {object} _dto.SuccessResponse "Success response"
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

	err = h.directorySupplierUsecase.Delete(c.Request.Context(), id)
	if err != nil {
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
// @Success 200 {object} map[string]interface{} "Directory suppliers list"
// @Router /v2/tqd/directory/suppliers [get]
func (h *DirectorySupplierHandler) ListDirectorySuppliersHTTP(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Parse search parameter
	search := c.Query("search")

	// Parse categories parameter
	var categories []string
	if categoriesStr := c.Query("categories"); categoriesStr != "" {
		categories = strings.Split(categoriesStr, ",")
	}

	// Parse isActive parameter
	var isActive *bool
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	// Parse rating parameters
	var minRating, maxRating *float64
	if minRatingStr := c.Query("minRating"); minRatingStr != "" {
		if val, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			minRating = &val
		}
	}
	if maxRatingStr := c.Query("maxRating"); maxRatingStr != "" {
		if val, err := strconv.ParseFloat(maxRatingStr, 64); err == nil {
			maxRating = &val
		}
	}

	suppliers, total, err := h.directorySupplierUsecase.List(c.Request.Context(), page, size, search, categories, isActive, minRating, maxRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list directory suppliers",
			"details": err.Error(),
		})
		return
	}

	response := h.mapper.ToListResponse(suppliers, total)
	response["page"] = page
	response["size"] = size

	c.JSON(http.StatusOK, response)
}
