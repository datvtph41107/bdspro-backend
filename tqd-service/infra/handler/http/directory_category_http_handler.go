package handler_http

import (
	"net/http"
	"strconv"

	_dto "common/domain/dto"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DirectoryCategoryHandler handles DirectoryCategory HTTP requests
type DirectoryCategoryHandler struct {
	directoryCategoryUsecase *usecase.DirectoryCategoryUsecase
	mapper                   *mapper.DirectoryCategoryMapper
	validator                *validator.DirectoryCategoryValidator
}

// NewDirectoryCategoryHandler creates a new DirectoryCategoryHandler
func NewDirectoryCategoryHandler(
	directoryCategoryUsecase *usecase.DirectoryCategoryUsecase,
	mapper *mapper.DirectoryCategoryMapper,
	validator *validator.DirectoryCategoryValidator,
) *DirectoryCategoryHandler {
	return &DirectoryCategoryHandler{
		directoryCategoryUsecase: directoryCategoryUsecase,
		mapper:                   mapper,
		validator:                validator,
	}
}

// CreateDirectoryCategoryHTTP handles HTTP POST request for creating directory category
// @Summary Create a new directory category
// @Description Create a new directory category
// @Tags Directory Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.DirectoryCategory true "Directory category"
// @Success 200 {object} domain.DirectoryCategory "Directory category"
// @Router /v2/tqd/directory/categories [post]
func (h *DirectoryCategoryHandler) CreateDirectoryCategoryHTTP(c *gin.Context) {
	var req domain.DirectoryCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validator.ValidateCreateRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create directory category
	err := h.directoryCategoryUsecase.CreateDirectoryCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// GetDirectoryCategoryHTTP handles HTTP GET request for getting directory category by ID
// @Summary Get a directory category by ID
// @Description Get a directory category by ID
// @Tags Directory Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Directory category ID"
// @Success 200 {object} domain.DirectoryCategory "Directory category"
// @Router /v2/tqd/directory/categories/{id} [get]
func (h *DirectoryCategoryHandler) GetDirectoryCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get directory category
	result, err := h.directoryCategoryUsecase.GetDirectoryCategoryByID(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Directory category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response using mapper
	response := h.mapper.ToProto(result)
	c.JSON(http.StatusOK, response)
}

// UpdateDirectoryCategoryHTTP handles HTTP PUT request for updating directory category
// @Summary Update a directory category
// @Description Update a directory category
// @Tags Directory Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Directory category ID"
// @Param request body domain.DirectoryCategory true "Directory category"
// @Success 200 {object} domain.DirectoryCategory "Directory category"
// @Router /v2/tqd/directory/categories/{id} [put]
func (h *DirectoryCategoryHandler) UpdateDirectoryCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req domain.DirectoryCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validator.ValidateUpdateRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update directory category
	result, err := h.directoryCategoryUsecase.UpdateDirectoryCategory(c.Request.Context(), id, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Directory category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response using mapper
	response := h.mapper.ToProto(result)
	c.JSON(http.StatusOK, response)
}

// DeleteDirectoryCategoryHTTP handles HTTP DELETE request for deleting directory category
// @Summary Delete a directory category
// @Description Delete a directory category
// @Tags Directory Category
// @Accept json
// @Produce json
// @Param id path string true "Directory category ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Router /v2/tqd/directory/categories/{id} [delete]
func (h *DirectoryCategoryHandler) DeleteDirectoryCategoryHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Delete directory category
	err = h.directoryCategoryUsecase.DeleteDirectoryCategory(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Directory category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Directory category deleted successfully",
	})
}

// ListDirectoryCategoriesHTTP handles HTTP GET request for listing directory categories
// @Summary List directory categories with pagination
// @Description List directory categories with pagination
// @Tags Directory Category
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search"
// @Param isActive query bool false "Is active"
// @Success 200 {object} usecase.ListDirectoryCategoriesResponse "List directory categories response"
// @Router /v2/tqd/directory/categories [get]
func (h *DirectoryCategoryHandler) ListDirectoryCategoriesHTTP(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 32)
	size, _ := strconv.ParseUint(c.DefaultQuery("size", "10"), 10, 32)
	search := c.Query("search")

	var isActive *bool
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	// Convert to request
	listReq := &dto.ListDirectoryCategoriesRequest{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Search:   search,
		IsActive: isActive,
	}

	// List directory categories
	result, total, err := h.directoryCategoryUsecase.ListDirectoryCategories(c.Request.Context(), listReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response using mapper
	response := gin.H{
		"code":          0,
		"message":       "success",
		"data":          h.mapper.ToProtoList(result),
		"totalElements": total,
	}

	// Return response
	c.JSON(http.StatusOK, response)
}
