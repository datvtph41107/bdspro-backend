package handler_http

import (
	"net/http"
	"strconv"

	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DirectorySourceHandler struct {
	// tqdpb.UnimplementedDirectorySourceServiceServer
	directorySourceUsecase *usecase.DirectorySourceUsecase
	mapper                 *mapper.DirectorySourceMapper
	validator              *validator.DirectorySourceValidator
}

func NewDirectorySourceHandler(
	directorySourceUsecase *usecase.DirectorySourceUsecase,
	mapper *mapper.DirectorySourceMapper,
	validator *validator.DirectorySourceValidator,
) *DirectorySourceHandler {
	return &DirectorySourceHandler{
		directorySourceUsecase: directorySourceUsecase,
		mapper:                 mapper,
		validator:              validator,
	}
}

// CreateDirectorySourceHTTP handles HTTP POST request for creating directory source
// @Summary Create a new directory source
// @Description Create a new directory source
// @Tags Directory Source
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.DirectorySource true "Directory source"
// @Success 200 {object} tqdpb.DirectorySource "Directory source"
// @Router /directory/sources [post]
func (h *DirectorySourceHandler) CreateDirectorySourceHTTP(c *gin.Context) {
	var createReq *dto.DirectorySourceDTO
	if err := c.ShouldBindJSON(createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validator.ValidateCreateRequest(createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create directory source
	result, err := h.directorySourceUsecase.Create(c.Request.Context(), createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to proto
	response := h.mapper.ToProto(result)

	c.JSON(http.StatusOK, response)
}

// GetDirectorySourceHTTP handles HTTP GET request for getting directory source by ID
// @Summary Get a directory source by ID
// @Description Get a directory source by ID
// @Tags Directory Source
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Directory source ID"
// @Success 200 {object} tqdpb.DirectorySource "Directory source"
// @Router /directory/sources/{id} [get]
func (h *DirectorySourceHandler) GetDirectorySourceHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get directory source
	result, err := h.directorySourceUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Directory source not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to proto
	response := h.mapper.ToProto(result)

	c.JSON(http.StatusOK, response)
}

// UpdateDirectorySourceHTTP handles HTTP PUT request for updating directory source
// @Summary Update a directory source
// @Description Update a directory source
// @Tags Directory Source
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Directory source ID"
// @Param request body tqdpb.DirectorySource true "Directory source"
// @Success 200 {object} tqdpb.DirectorySource "Directory source"
// @Router /directory/sources/{id} [put]
func (h *DirectorySourceHandler) UpdateDirectorySourceHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req domain.DirectorySource
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set ID from URL parameter
	req.ID = id

	// Validate request
	// if err := h.validator.ValidateUpdateRequest(req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// Update directory source
	result, err := h.directorySourceUsecase.Update(c.Request.Context(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Directory source not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to proto
	response := h.mapper.ToProto(result)

	c.JSON(http.StatusOK, response)
}

// DeleteDirectorySourceHTTP handles HTTP DELETE request for deleting directory source
// @Summary Delete a directory source
// @Description Delete a directory source
// @Tags Directory Source
// @Accept json
// @Produce json
// @Param id path string true "Directory source ID"
// @Success 200 {object} sharepb.SubmitResponse "Submit response"
// @Router /directory/sources/{id} [delete]
func (h *DirectorySourceHandler) DeleteDirectorySourceHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Delete directory source
	err = h.directorySourceUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Directory source not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, &sharepb.SubmitResponse{
		Id:      id,
		Message: "Directory source deleted successfully",
	})
}

// ListDirectorySourcesHTTP handles HTTP GET request for listing directory sources
// @Summary List directory sources with pagination
// @Description List directory sources with pagination
// @Tags Directory Source
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search"
// @Param category query string false "Category"
// @Param type query string false "Type"
// @Param isActive query bool false "Is active"
// @Param isRecurring query bool false "Is recurring"
// @Success 200 {object} tqdpb.ListDirectorySourcesResponse "List directory sources response"
// @Router /directory/sources [get]
func (h *DirectorySourceHandler) ListDirectorySourcesHTTP(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 32)
	size, _ := strconv.ParseUint(c.DefaultQuery("size", "10"), 10, 32)
	search := c.Query("search")
	category := c.Query("category")
	typeStr := c.Query("type")

	var isActive *bool
	var isRecurring *bool

	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	if isRecurringStr := c.Query("isRecurring"); isRecurringStr != "" {
		if val, err := strconv.ParseBool(isRecurringStr); err == nil {
			isRecurring = &val
		}
	}

	// Convert to DTO
	listReq := &dto.ListDirectorySourcesRequestDTO{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Search:      search,
		Category:    &category,
		Type:        &typeStr,
		IsActive:    isActive,
		IsRecurring: isRecurring,
	}

	// List directory sources
	result, total, err := h.directorySourceUsecase.List(c.Request.Context(), listReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DTOs to protos
	response := gin.H{
		"code":          0,
		"message":       "success",
		"data":          result,
		"totalElements": total,
	}
	// Return response
	c.JSON(http.StatusOK, response)
}
