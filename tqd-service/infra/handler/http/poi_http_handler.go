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
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// POIHTTPHandler handles HTTP requests for POI operations
type POIHTTPHandler struct {
	usecase usecase.PoiUsecase
	mapper  *mapper.PoiMapper
}

// NewPOIHTTPHandler creates a new POI HTTP handler
func NewPOIHTTPHandler(
	usecase usecase.PoiUsecase,
	mapper *mapper.PoiMapper,
) *POIHTTPHandler {
	return &POIHTTPHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// CreatePOI handles HTTP POST request for creating POI
// @Summary Create a new POI
// @Description Create a new Point of Interest
// @Tags POI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePoiRequest true "POI creation request"
// @Success 201 {object} map[string]interface{} "POI created successfully"
// @Router /v2/tqd/pois [post]
func (h *POIHTTPHandler) CreatePOI(c *gin.Context) {
	var req dto.CreatePoiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Validation failed",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create POI",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "POI created successfully",
		"data":    result,
	})
}

// GetPOIByID handles HTTP GET request for getting POI by ID
// @Summary Get POI by ID
// @Description Get a Point of Interest by its ID
// @Tags POI
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Success 200 {object} map[string]interface{} "POI details"
// @Router /v2/tqd/pois/{id} [get]
func (h *POIHTTPHandler) GetPOIByID(c *gin.Context) {
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
			"message": "POI not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI retrieved successfully",
		"data":    result,
	})
}

// GetPOIByCode handles HTTP GET request for getting POI by code
// @Summary Get POI by code
// @Description Get a Point of Interest by its code
// @Tags POI
// @Produce json
// @Security BearerAuth
// @Param code path string true "POI Code"
// @Success 200 {object} map[string]interface{} "POI details"
// @Router /v2/tqd/pois/code/{code} [get]
func (h *POIHTTPHandler) GetPOIByCode(c *gin.Context) {
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
			"message": "POI not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI retrieved successfully",
		"data":    result,
	})
}

// UpdatePOI handles HTTP PUT request for updating POI
// @Summary Update a POI
// @Description Update an existing Point of Interest
// @Tags POI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Param request body dto.UpdatePoiRequest true "POI update request"
// @Success 200 {object} map[string]interface{} "POI updated successfully"
// @Router /v2/tqd/pois/{id} [put]
func (h *POIHTTPHandler) UpdatePOI(c *gin.Context) {
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

	var req dto.UpdatePoiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate request if needed
	if req.Name != nil {
		if err := validate.Var(*req.Name, "required,min=1,max=255"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid name",
				"details": err.Error(),
			})
			return
		}
	}

	result, err := h.usecase.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update POI",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI updated successfully",
		"data":    result,
	})
}

// DeletePOI handles HTTP DELETE request for deleting POI
// @Summary Delete a POI
// @Description Delete a Point of Interest
// @Tags POI
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Success 200 {object} map[string]interface{} "POI deleted successfully"
// @Router /v2/tqd/pois/{id} [delete]
func (h *POIHTTPHandler) DeletePOI(c *gin.Context) {
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
			"message": "Failed to delete POI",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI deleted successfully",
	})
}

// ListPOIs handles HTTP GET request for listing POIs
// @Summary List POIs
// @Description Get a paginated list of Points of Interest with optional filters
// @Tags POI
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param search query string false "Search by name or address or code"
// @Param categoryId query int false "Category ID"
// @Param categoryIds query string false "Category IDs (comma-separated)"
// @Param isActive query bool false "Filter by active status"
// @Param isVerified query bool false "Filter by verified status"
// @Param isFeatured query bool false "Filter by featured status"
// @Param minRating query float64 false "Minimum rating"
// @Param maxRating query float64 false "Maximum rating"
// @Param sortBy query string false "Sort by (rating, newest)"
// @Success 200 {object} map[string]interface{} "POIs list"
// @Router /v2/tqd/pois [get]
func (h *POIHTTPHandler) ListPOIs(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Create filter
	filter := &dto.PoiFilter{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Search: c.Query("search"),
		SortBy: c.Query("sortBy"),
	}

	// Parse categoryId
	if categoryIdStr := c.Query("categoryId"); categoryIdStr != "" {
		if categoryId, err := strconv.ParseUint(categoryIdStr, 10, 64); err == nil {
			filter.CategoryID = &categoryId
		}
	}

	// Parse categoryIds
	if categoryIdsStr := c.Query("categoryIds"); categoryIdsStr != "" {
		ids := strings.Split(categoryIdsStr, ",")
		for _, idStr := range ids {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64); err == nil {
				filter.CategoryIDs = append(filter.CategoryIDs, id)
			}
		}
	}

	// Parse boolean filters
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}
	if isVerifiedStr := c.Query("isVerified"); isVerifiedStr != "" {
		if isVerified, err := strconv.ParseBool(isVerifiedStr); err == nil {
			filter.IsVerified = &isVerified
		}
	}
	if isFeaturedStr := c.Query("isFeatured"); isFeaturedStr != "" {
		if isFeatured, err := strconv.ParseBool(isFeaturedStr); err == nil {
			filter.IsFeatured = &isFeatured
		}
	}

	// Parse rating filters
	if minRatingStr := c.Query("minRating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			filter.MinRating = &minRating
		}
	}
	if maxRatingStr := c.Query("maxRating"); maxRatingStr != "" {
		if maxRating, err := strconv.ParseFloat(maxRatingStr, 64); err == nil {
			filter.MaxRating = &maxRating
		}
	}

	result, err := h.usecase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to list POIs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POIs retrieved successfully",
		"data": gin.H{
			"items": result.Data,
			"total": result.Total,
			"page":  result.Page,
			"size":  result.Size,
			"pages": (result.Total + int64(result.Size) - 1) / int64(result.Size),
		},
	})
}

// GetPOIsByCategory handles HTTP GET request for getting POIs by category
// @Summary Get POIs by category
// @Description Get POIs filtered by category
// @Tags POI
// @Produce json
// @Security BearerAuth
// @Param categoryId path uint64 true "Category ID"
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} map[string]interface{} "POIs list"
// @Router /v2/tqd/pois/category/{categoryId} [get]
func (h *POIHTTPHandler) GetPOIsByCategory(c *gin.Context) {
	categoryIdStr := c.Param("categoryId")
	categoryId, err := strconv.ParseUint(categoryIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid category ID",
			"details": err.Error(),
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	results, err := h.usecase.ListByCategory(c.Request.Context(), categoryId, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get POIs by category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POIs retrieved successfully",
		"data":    results,
	})
}

// GetNearbyPOIs handles HTTP GET request for getting nearby POIs
// @Summary Get nearby POIs
// @Description Get POIs within a certain radius from given coordinates
// @Tags POI
// @Produce json
// @Security BearerAuth
// @Param lat query float64 true "Latitude"
// @Param lng query float64 true "Longitude"
// @Param radius query float64 false "Radius in kilometers" default(10)
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param categoryIds query string false "Category IDs (comma-separated)"
// @Success 200 {object} map[string]interface{} "Nearby POIs list"
// @Router /v2/tqd/pois/nearby [get]
func (h *POIHTTPHandler) GetNearbyPOIs(c *gin.Context) {
	// Parse coordinates
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid latitude",
			"details": err.Error(),
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid longitude",
			"details": err.Error(),
		})
		return
	}

	// Parse radius
	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "10"), 64)
	if radius <= 0 {
		radius = 10
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// Create nearby request
	req := &dto.NearbyPoiRequest{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		Latitude:  lat,
		Longitude: lng,
		Radius:    radius,
	}

	// Parse categoryIds
	if categoryIdsStr := c.Query("categoryIds"); categoryIdsStr != "" {
		ids := strings.Split(categoryIdsStr, ",")
		for _, idStr := range ids {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64); err == nil {
				req.CategoryIDs = append(req.CategoryIDs, id)
			}
		}
	}

	result, err := h.usecase.ListNearby(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get nearby POIs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Nearby POIs retrieved successfully",
		"data": gin.H{
			"items": result.Data,
			"total": result.Total,
			"page":  result.Page,
			"size":  result.Size,
			"pages": (result.Total + int64(result.Size) - 1) / int64(result.Size),
		},
	})
}

// UpdatePOIRating handles HTTP PUT request for updating POI rating
// @Summary Update POI rating
// @Description Update the rating information for a POI
// @Tags POI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Param request body UpdateRatingRequest true "Rating update request"
// @Success 200 {object} map[string]interface{} "POI updated"
// @Router /v2/tqd/pois/{id}/rating [put]
func (h *POIHTTPHandler) UpdatePOIRating(c *gin.Context) {
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
		Rating      float64 `json:"rating" binding:"required,min=0,max=5"`
		ReviewCount uint32  `json:"reviewCount" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.UpdateRating(c.Request.Context(), id, req.Rating, req.ReviewCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update POI rating",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI rating updated successfully",
		"data":    result,
	})
}

// IncrementPOIViewCount handles HTTP POST request for incrementing view count
// @Summary Increment POI view count
// @Description Increment the view count for a POI
// @Tags POI
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Success 200 {object} map[string]interface{} "View count incremented"
// @Router /v2/tqd/pois/{id}/view [post]
func (h *POIHTTPHandler) IncrementPOIViewCount(c *gin.Context) {
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

	if err := h.usecase.IncrementViewCount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to increment view count",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "View count incremented successfully",
	})
}

// TogglePOILike handles HTTP POST request for toggling like
// @Summary Toggle POI like
// @Description Increment or decrement like count for a POI
// @Tags POI
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Param action query string true "Action (like/unlike)"
// @Success 200 {object} map[string]interface{} "Like toggled"
// @Router /v2/tqd/pois/{id}/like [post]
func (h *POIHTTPHandler) TogglePOILike(c *gin.Context) {
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

	action := c.Query("action")
	var likeErr error

	switch action {
	case "like":
		likeErr = h.usecase.IncrementLikeCount(c.Request.Context(), id)
	case "unlike":
		likeErr = h.usecase.DecrementLikeCount(c.Request.Context(), id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid action. Use 'like' or 'unlike'",
		})
		return
	}

	if likeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to toggle like",
			"details": likeErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Like toggled successfully",
	})
}

// VerifyPOI handles HTTP PUT request for verifying POI
// @Summary Verify POI
// @Description Update verification status for a POI
// @Tags POI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "POI ID"
// @Param request body VerifyRequest true "Verification request"
// @Success 200 {object} map[string]interface{} "POI updated"
// @Router /v2/tqd/pois/{id}/verify [put]
func (h *POIHTTPHandler) VerifyPOI(c *gin.Context) {
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
		Verified bool `json:"verified"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.usecase.VerifyPOI(c.Request.Context(), id, req.Verified)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to verify POI",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "POI verification updated successfully",
		"data":    result,
	})
}
