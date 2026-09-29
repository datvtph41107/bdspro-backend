package handler

import (
	"net/http"
	"strconv"

	"map/internal/dto"
	"map/internal/usecase"

	"github.com/gin-gonic/gin"
)

// MapPointHandler handles map point related HTTP requests
type MapPointHandler struct {
	mapPointUsecase *usecase.MapPointUsecase
}

// NewMapPointHandler creates a new map point handler
func NewMapPointHandler(mapPointUsecase *usecase.MapPointUsecase) *MapPointHandler {
	return &MapPointHandler{
		mapPointUsecase: mapPointUsecase,
	}
}

// CreateMapPoint creates a new map point
// @Summary Create a new map point
// @Description Create a new point on the map
// @Tags Map Points
// @Accept json
// @Produce json
// @Param request body dto.CreateMapPointRequest true "Map point data"
// @Success 200 {object} dto.MapPointResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /map/locations/new [post]
// @Security BearerAuth
func (h *MapPointHandler) CreateMapPoint(c *gin.Context) {
	var req dto.CreateMapPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := h.mapPointUsecase.CreateMapPoint(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// FindNearbyLocations finds locations near a given point
// @Summary Find nearby locations
// @Description Find locations within a specified radius of a point
// @Tags Map Points
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude" example(10.762622)
// @Param lng query float64 true "Longitude" example(106.660172)
// @Param radius query float64 false "Search radius in meters" example(5000)
// @Success 200 {object} dto.NearbyLocationsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /map/locations/nearby [get]
func (h *MapPointHandler) FindNearbyLocations(c *gin.Context) {
	var req dto.NearbyLocationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	response, err := h.mapPointUsecase.FindNearbyLocations(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// FindLocationsByPolygon finds locations within a polygon
// @Summary Find locations by polygon
// @Description Find locations within a specified polygon area
// @Tags Map Points
// @Accept json
// @Produce json
// @Param request body dto.FindByPolygonRequest true "Polygon coordinates"
// @Success 200 {object} dto.FindByPolygonResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /map/locations/find-by-polygon [post]
func (h *MapPointHandler) FindLocationsByPolygon(c *gin.Context) {
	var req dto.FindByPolygonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := h.mapPointUsecase.FindLocationsByPolygon(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMapPointByID gets a map point by ID
// @Summary Get map point by ID
// @Description Get a specific map point by its ID
// @Tags Map Points
// @Accept json
// @Produce json
// @Param id path int true "Map point ID"
// @Success 200 {object} dto.MapPointResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /map/locations/{id} [get]
func (h *MapPointHandler) GetMapPointByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	response, err := h.mapPointUsecase.GetMapPointByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMapPoint updates an existing map point
// @Summary Update map point
// @Description Update an existing map point
// @Tags Map Points
// @Accept json
// @Produce json
// @Param id path int true "Map point ID"
// @Param request body dto.CreateMapPointRequest true "Updated map point data"
// @Success 200 {object} dto.MapPointResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /map/locations/{id} [put]
// @Security BearerAuth
func (h *MapPointHandler) UpdateMapPoint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.CreateMapPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := h.mapPointUsecase.UpdateMapPoint(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteMapPoint deletes a map point by ID
// @Summary Delete map point
// @Description Delete a map point by its ID
// @Tags Map Points
// @Accept json
// @Produce json
// @Param id path int true "Map point ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /map/locations/{id} [delete]
// @Security BearerAuth
func (h *MapPointHandler) DeleteMapPoint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.mapPointUsecase.DeleteMapPoint(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Map point deleted successfully"})
}
