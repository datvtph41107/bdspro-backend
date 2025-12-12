package handler_http

import (
	"context"
	"strconv"

	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// POIHTTPHandler handles HTTP requests for POI operations
type POIHTTPHandler struct {
	poiUsecase usecase.IPOIUsecase
}

// NewPOIHTTPHandler creates a new POI HTTP handler
// @bind: infra/handler/http/poi_http_handler.go
func NewPOIHTTPHandler(poiUsecase usecase.IPOIUsecase) *POIHTTPHandler {
	return &POIHTTPHandler{
		poiUsecase: poiUsecase,
	}
}

// CreatePOI creates a new POI
// @Summary Create a new POI
// @Description Create a new Point of Interest
// @Tags POI
// @Accept json
// @Produce json
// @Param request body dto.CreatePOIRequestDTO true "POI creation request"
// @Success 201 {object} dto.POIDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois [post]
func (h *POIHTTPHandler) CreatePOI(c *gin.Context) {
	var req dto.CreatePOIRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "details": err.Error()})
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "details": err.Error()})
		return
	}

	ctx := context.Background()
	poiDTO, errorDTO := h.poiUsecase.CreatePOI(ctx, &req)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(201, poiDTO)
}

// UpdatePOI updates an existing POI
// @Summary Update a POI
// @Description Update an existing Point of Interest
// @Tags POI
// @Accept json
// @Produce json
// @Param id path uint64 true "POI ID"
// @Param request body dto.UpdatePOIRequestDTO true "POI update request"
// @Success 200 {object} dto.POIDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/{id} [put]
func (h *POIHTTPHandler) UpdatePOI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req dto.UpdatePOIRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "details": err.Error()})
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "details": err.Error()})
		return
	}

	ctx := context.Background()
	poiDTO, errorDTO := h.poiUsecase.UpdatePOI(ctx, id, &req)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, poiDTO)
}

// DeletePOI deletes a POI
// @Summary Delete a POI
// @Description Delete a Point of Interest
// @Tags POI
// @Param id path uint64 true "POI ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/{id} [delete]
func (h *POIHTTPHandler) DeletePOI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	ctx := context.Background()
	errorDTO := h.poiUsecase.DeletePOI(ctx, id)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.Status(204)
}

// GetPOIByID retrieves a POI by ID
// @Summary Get POI by ID
// @Description Get a Point of Interest by its ID
// @Tags POI
// @Param id path uint64 true "POI ID"
// @Success 200 {object} dto.POIDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/{id} [get]
func (h *POIHTTPHandler) GetPOIByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	ctx := context.Background()
	poiDTO, errorDTO := h.poiUsecase.GetPOIByID(ctx, id)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, poiDTO)
}

// GetPOIByCode retrieves a POI by code
// @Summary Get POI by code
// @Description Get a Point of Interest by its code
// @Tags POI
// @Param code path string true "POI Code"
// @Success 200 {object} dto.POIDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/code/{code} [get]
func (h *POIHTTPHandler) GetPOIByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "Mã POI không được để trống"})
		return
	}

	ctx := context.Background()
	poiDTO, errorDTO := h.poiUsecase.GetPOIByCode(ctx, code)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, poiDTO)
}

// GetPOIList retrieves a paginated list of POIs
// @Summary Get POI list
// @Description Get a paginated list of Points of Interest with optional filters
// @Tags POI
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(20)
// @Param text query string false "Search by name or address or code"
// @Param categoryId query int false "Category ID"
// @Param isActive query bool false "Filter by active status"
// @Param minLat query float64 false "Minimum latitude"
// @Param maxLat query float64 false "Maximum latitude"
// @Param minLng query float64 false "Minimum longitude"
// @Param maxLng query float64 false "Maximum longitude"
// @Success 200 {object} dto.ListPOIsResponseDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois [get]
func (h *POIHTTPHandler) GetPOIList(c *gin.Context) {
	var req dto.ListPOIsRequestDTO

	// Bind pagination parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "Tham số không hợp lệ", "details": err.Error()})
		return
	}

	// Parse optional filters
	if categoryIdStr := c.Query("categoryId"); categoryIdStr != "" {
		if categoryId, err := strconv.ParseUint(categoryIdStr, 10, 64); err == nil {
			req.CategoryID = &categoryId
		}
	}

	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	if minLatStr := c.Query("minLat"); minLatStr != "" {
		if minLat, err := strconv.ParseFloat(minLatStr, 64); err == nil {
			req.MinLat = &minLat
		}
	}

	if maxLatStr := c.Query("maxLat"); maxLatStr != "" {
		if maxLat, err := strconv.ParseFloat(maxLatStr, 64); err == nil {
			req.MaxLat = &maxLat
		}
	}

	if minLngStr := c.Query("minLng"); minLngStr != "" {
		if minLng, err := strconv.ParseFloat(minLngStr, 64); err == nil {
			req.MinLng = &minLng
		}
	}

	if maxLngStr := c.Query("maxLng"); maxLngStr != "" {
		if maxLng, err := strconv.ParseFloat(maxLngStr, 64); err == nil {
			req.MaxLng = &maxLng
		}
	}

	ctx := c.Request.Context()
	response, errorDTO := h.poiUsecase.GetPOIList(ctx, &req)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, response)
}

// GetPOIsByCategory retrieves POIs by category
// @Summary Get POIs by category
// @Description Get POIs filtered by category
// @Tags POI
// @Param categoryId path uint64 true "Category ID"
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} dto.POIDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/category/{categoryId} [get]
func (h *POIHTTPHandler) GetPOIsByCategory(c *gin.Context) {
	categoryIdStr := c.Param("categoryId")
	categoryId, err := strconv.ParseUint(categoryIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID danh mục không hợp lệ"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	ctx := context.Background()
	poiDTOs, errorDTO := h.poiUsecase.GetPOIsByCategory(ctx, categoryId, limit)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, poiDTOs)
}

// GetNearbyPOIs retrieves POIs within a certain radius
// @Summary Get nearby POIs
// @Description Get POIs within a certain radius from given coordinates
// @Tags POI
// @Param lat query float64 true "Latitude"
// @Param lng query float64 true "Longitude"
// @Param radius query float64 false "Radius in kilometers" default(10)
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} dto.POIDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/nearby [get]
func (h *POIHTTPHandler) GetNearbyPOIs(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Vĩ độ không hợp lệ"})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Kinh độ không hợp lệ"})
		return
	}

	radiusStr := c.DefaultQuery("radius", "10")
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		radius = 10
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	ctx := context.Background()
	poiDTOs, errorDTO := h.poiUsecase.GetNearbyPOIs(ctx, lat, lng, radius, limit)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, poiDTOs)
}

// UpdatePOIRating updates the rating for a POI
// @Summary Update POI rating
// @Description Update the rating information for a POI
// @Tags POI
// @Accept json
// @Produce json
// @Param id path uint64 true "POI ID"
// @Param request body map[string]interface{} true "Rating update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/pois/{id}/rating [put]
func (h *POIHTTPHandler) UpdatePOIRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req struct {
		RatingPoint float64 `json:"ratingPoint" validate:"required,min=0,max=5"`
		RatingCount int     `json:"ratingCount" validate:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "details": err.Error()})
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ", "details": err.Error()})
		return
	}

	ctx := context.Background()
	errorDTO := h.poiUsecase.UpdatePOIRating(ctx, id, req.RatingPoint, req.RatingCount)
	if errorDTO != nil {
		c.JSON(errorDTO.Code, gin.H{"error": errorDTO.Message})
		return
	}

	c.JSON(200, gin.H{"message": "Cập nhật đánh giá thành công"})
}
