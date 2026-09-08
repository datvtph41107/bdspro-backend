package handler_http

import (
	"tqd/internal/usecase"
)

type RegionHttpHandler struct {
	regionUsecase usecase.RegionUsecase
}

func NewRegionHttpHandler(regionUsecase usecase.RegionUsecase) *RegionHttpHandler {
	return &RegionHttpHandler{regionUsecase: regionUsecase}
}

// CreateRegionBatch godoc
// POST /v2/tqd/admin/regions/batch
// Body: { "items": [{ "geometry": "<GeoJSON>", "layerId": 1, ... }] }
// func (h *RegionHttpHandler) CreateRegionBatch(c *gin.Context) {
// 	var req dto.CreateRegionBatchRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if len(req.Items) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "items must not be empty"})
// 		return
// 	}

// 	resp, err := h.regionUsecase.CreateBatch(c.Request.Context(), &req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }
