package handler_http

import (
	"net/http"
	sharepb "pb/types/shared"
	"strconv"

	_dto "common/domain/dto"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AmenityHandler handles Amenity HTTP requests
type AmenityHandler struct {
	amenityUsecase *usecase.AmenityUsecase
	amenityMapper  *mapper.AmenityMapper
}

// NewAmenityHandler creates a new AmenityHandler
func NewAmenityHandler(
	amenityUsecase *usecase.AmenityUsecase,
	amenityMapper *mapper.AmenityMapper,
) *AmenityHandler {
	return &AmenityHandler{
		amenityUsecase: amenityUsecase,
		amenityMapper:  amenityMapper,
	}
}

// CreateAmenityHTTP handles HTTP POST request for creating amenity
// @Summary Create a new amenity
// @Description Create a new amenity
// @Tags Amenity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAmenityRequestDTO true "Amenity"
// @Success 200 {object} dto.AmenityDTO "Amenity"
// @Router /v2/tqd/amenities [post]
func (h *AmenityHandler) CreateAmenityHTTP(c *gin.Context) {
	var createReq dto.CreateAmenityRequestDTO
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to domain entity
	amenity := h.amenityMapper.ToDomainFromCreateRequest(&createReq)

	// Create amenity
	result, err := h.amenityUsecase.Create(c.Request.Context(), amenity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert domain to DTO
	amenityDTO := h.amenityMapper.ToDTO(result)

	c.JSON(http.StatusOK, amenityDTO)
}

// GetAmenityHTTP handles HTTP GET request for getting amenity by ID
// @Summary Get an amenity by ID
// @Description Get an amenity by ID
// @Tags Amenity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Amenity ID"
// @Success 200 {object} dto.AmenityDTO "Amenity"
// @Router /v2/tqd/amenities/{id} [get]
func (h *AmenityHandler) GetAmenityHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get amenity
	amenity, err := h.amenityUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Amenity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert domain to DTO
	amenityDTO := h.amenityMapper.ToDTO(amenity)

	c.JSON(http.StatusOK, amenityDTO)
}

// UpdateAmenityHTTP handles HTTP PUT request for updating amenity
// @Summary Update an amenity
// @Description Update an amenity
// @Tags Amenity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Amenity ID"
// @Param request body dto.UpdateAmenityRequestDTO true "Amenity"
// @Success 200 {object} dto.AmenityDTO "Amenity"
// @Router /v2/tqd/amenities/{id} [put]
func (h *AmenityHandler) UpdateAmenityHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updateReq dto.UpdateAmenityRequestDTO
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to domain entity
	amenity := h.amenityMapper.ToDomainFromUpdateRequest(&updateReq, id)

	// Update amenity
	result, err := h.amenityUsecase.Update(c.Request.Context(), id, amenity)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Amenity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert domain to DTO
	amenityDTO := h.amenityMapper.ToDTO(result)

	c.JSON(http.StatusOK, amenityDTO)
}

// DeleteAmenityHTTP handles HTTP DELETE request for deleting amenity
// @Summary Delete an amenity
// @Description Delete an amenity
// @Tags Amenity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Amenity ID"
// @Success 200 {object} sharepb.SubmitResponse "Submit response"
// @Router /v2/tqd/amenities/{id} [delete]
func (h *AmenityHandler) DeleteAmenityHTTP(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Delete amenity
	_, err = h.amenityUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Amenity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, &sharepb.SubmitResponse{
		Id:      id,
		Message: "Amenity deleted successfully",
	})
}

// ListAmenitiesHTTP handles HTTP GET request for listing amenities
// @Summary List amenities with pagination
// @Description List amenities with pagination
// @Tags Amenity
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search"
// @Success 200 {object} dto.ListAmenitiesResponseDTO "List amenities response"
// @Router /v2/tqd/amenities [get]
func (h *AmenityHandler) ListAmenitiesHTTP(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 32)
	size, _ := strconv.ParseUint(c.DefaultQuery("size", "10"), 10, 32)
	// search := c.Query("search")

	// Convert to DTO
	listReq := &dto.AmenityFilterDTO{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
	}

	// List amenities
	amenitiesDomain, total, err := h.amenityUsecase.GetList(c.Request.Context(), &listReq.Pagable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, &dto.ListAmenitiesResponseDTO{
		Data:  amenitiesDomain,
		Total: total,
	})
}
