package discovery

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	commonhttp "common/httpresponse"

	"tqd/internal/domain/discovery/model"
	"tqd/internal/usecase/discovery/application"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *application.Service }

func NewHandler(service *application.Service) *Handler { return &Handler{service: service} }

func (h *Handler) Identify(c *gin.Context) {
	var req domain.IdentifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
		return
	}
	response, err := h.service.Identify(c.Request.Context(), req)
	if err != nil {
		statusCode, payload := identifyHTTPProblem(err)
		c.JSON(statusCode, payload)
		return
	}
	c.JSON(http.StatusOK, response)
}

func identifyHTTPProblem(err error) (int, gin.H) {
	problem := commonhttp.ProblemFromError(err)
	return problem.Status, gin.H{
		"code":       "DISCOVERY_IDENTIFY_FAILED",
		"message":    problem.Detail,
		"error_code": problem.Code,
	}
}

func parseViewport(raw string) (*domain.Bounds, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("viewport must contain minLongitude,minLatitude,maxLongitude,maxLatitude")
	}

	values := make([]float64, 4)
	for i, part := range parts {
		value, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("viewport value %q is invalid", part)
		}
		values[i] = value
	}

	bounds := &domain.Bounds{
		MinLongitude: values[0],
		MinLatitude:  values[1],
		MaxLongitude: values[2],
		MaxLatitude:  values[3],
	}

	if bounds.MinLongitude < -180 || bounds.MaxLongitude > 180 ||
		bounds.MinLatitude < -90 || bounds.MaxLatitude > 90 ||
		bounds.MinLongitude >= bounds.MaxLongitude ||
		bounds.MinLatitude >= bounds.MaxLatitude {
		return nil, fmt.Errorf("viewport bounds are outside the valid coordinate range")
	}

	return bounds, nil
}

func (h *Handler) Search(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	req := domain.SearchRequest{
		Query:              c.Query("q"),
		Mode:               c.DefaultQuery("mode", "suggest"),
		Limit:              limit,
		Cursor:             strings.TrimSpace(c.Query("cursor")),
		AdministrativeCode: strings.TrimSpace(c.Query("administrativeCode")),
	}

	if kinds := strings.TrimSpace(c.Query("kinds")); kinds != "" {
		for _, kind := range strings.Split(kinds, ",") {
			k := domain.EntityKind(strings.TrimSpace(kind))
			if k == "" {
				continue
			}
			if !k.IsValid() {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "INVALID_ENTITY_KIND",
					"message": "Unsupported entity kind: " + string(k),
				})
				return
			}
			req.EntityKinds = append(req.EntityKinds, k)
		}
	}

	if mode := strings.TrimSpace(c.Query("mapMode")); mode != "" {
		req.MapContext.MapMode = mode
	}

	if layerIDs := strings.TrimSpace(c.Query("activeLayerIds")); layerIDs != "" {
		for _, layerID := range strings.Split(layerIDs, ",") {
			if value := strings.TrimSpace(layerID); value != "" {
				req.MapContext.ActiveLayerIDs = append(req.MapContext.ActiveLayerIDs, value)
			}
		}
	}

	viewport, err := parseViewport(c.Query("viewport"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_VIEWPORT",
			"message": err.Error(),
		})
		return
	}
	req.MapContext.Viewport = viewport

	response, err := h.service.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DISCOVERY_SEARCH_FAILED",
			"message": err.Error(),
		})
		return
	}

	c.Header("X-QHPro-Discovery-Source", "tqd-discovery-v2")
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetEntity(c *gin.Context) {
	key := c.Param("key")
	entity, err := h.service.GetEntity(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_ENTITY_KEY", "message": err.Error()})
		return
	}
	if entity == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "ENTITY_NOT_FOUND", "message": "Không tìm thấy đối tượng"})
		return
	}
	c.JSON(http.StatusOK, entity)
}
