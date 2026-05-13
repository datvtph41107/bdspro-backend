package handler_http

import (
	"net/http"
	"strconv"
	"strings"

	_utils "common/utils"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

// QHLabelHandler xử lý các request HTTP cho nhãn quy hoạch (QHLabel)
type QHLabelHandler struct {
	labelUsecase usecase.QHLabelUsecase
}

func NewQHLabelHandler(labelUsecase usecase.QHLabelUsecase) *QHLabelHandler {
	return &QHLabelHandler{labelUsecase: labelUsecase}
}

// createLabelRequest là body cho tạo / cập nhật nhãn
type createLabelRequest struct {
	LayerID      uint64  `json:"layerId"`
	Name         string  `json:"name"`
	DisplayName  string  `json:"displayName"`
	Description  string  `json:"description"`
	Color        string  `json:"color"`
	FillOpacity  float64 `json:"fillOpacity"`
	StrokeColor  string  `json:"strokeColor"`
	StrokeWidth  int     `json:"strokeWidth"`
	DisplayOrder int     `json:"displayOrder"`
	IsVisible    bool    `json:"isVisible"`
	MinZoom      int     `json:"minZoom"`
	MaxZoom      int     `json:"maxZoom"`
	StandardAt   *string `json:"standardAt,omitempty"` // RFC3339; omit = không đổi (update) / không gán (create)
}

// CreateQHLabelHTTP godoc
// @Summary      Tạo nhãn mới
// @Tags         QH Label
// @Accept       json
// @Produce      json
// @Param        request body createLabelRequest true "Thông tin nhãn"
// @Success      200  {object}  qh_domain.QHLabel
// @Router       /v2/tqd/qh/labels [post]
func (h *QHLabelHandler) CreateQHLabelHTTP(c *gin.Context) {
	var req createLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label := &qh_domain.QHLabel{
		LayerID:      req.LayerID,
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Color:        req.Color,
		FillOpacity:  req.FillOpacity,
		StrokeColor:  req.StrokeColor,
		StrokeWidth:  req.StrokeWidth,
		DisplayOrder: req.DisplayOrder,
		IsVisible:    req.IsVisible,
		MinZoom:      req.MinZoom,
		MaxZoom:      req.MaxZoom,
		Status:       10,
	}

	// Set defaults if empty
	if label.Color == "" {
		label.Color = "#CCCCCC"
	}
	if label.FillOpacity == 0 {
		label.FillOpacity = 0.6
	}
	if label.StrokeColor == "" {
		label.StrokeColor = "#000000"
	}
	if label.StrokeWidth == 0 {
		label.StrokeWidth = 1
	}

	if req.StandardAt != nil {
		s := strings.TrimSpace(*req.StandardAt)
		if s != "" {
			if t := _utils.ParseStringToTime(s); t != nil {
				label.StandardAt = t
			}
		}
	}

	result, err := h.labelUsecase.Create(c.Request.Context(), label)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetQHLabelHTTP godoc
// @Summary      Lấy nhãn theo ID
// @Tags         QH Label
// @Produce      json
// @Param        id   path  int  true  "Label ID"
// @Success      200  {object}  qh_domain.QHLabel
// @Router       /v2/tqd/qh/labels/{id} [get]
func (h *QHLabelHandler) GetQHLabelHTTP(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	result, err := h.labelUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateQHLabelHTTP godoc
// @Summary      Cập nhật nhãn
// @Tags         QH Label
// @Accept       json
// @Produce      json
// @Param        id      path  int                 true  "Label ID"
// @Param        request body  createLabelRequest  true  "Thông tin nhãn"
// @Success      200  {object}  qh_domain.QHLabel
// @Router       /v2/tqd/qh/labels/{id} [put]
func (h *QHLabelHandler) UpdateQHLabelHTTP(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req createLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label := &qh_domain.QHLabel{
		ID:           id,
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Color:        req.Color,
		FillOpacity:  req.FillOpacity,
		StrokeColor:  req.StrokeColor,
		StrokeWidth:  req.StrokeWidth,
		DisplayOrder: req.DisplayOrder,
		IsVisible:    req.IsVisible,
		MinZoom:      req.MinZoom,
		MaxZoom:      req.MaxZoom,
	}
	if req.StandardAt != nil {
		s := strings.TrimSpace(*req.StandardAt)
		if s == "" {
			label.ClearStandardAt = true
		} else if t := _utils.ParseStringToTime(s); t != nil {
			label.StandardAt = t
		}
	}

	result, err := h.labelUsecase.Update(c.Request.Context(), id, label)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteQHLabelHTTP godoc
// @Summary      Xóa nhãn
// @Tags         QH Label
// @Produce      json
// @Param        id  path  int  true  "Label ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /v2/tqd/qh/labels/{id} [delete]
func (h *QHLabelHandler) DeleteQHLabelHTTP(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := h.labelUsecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "message": "Xóa nhãn thành công"})
}

// ListQHLabelsByLayerHTTP godoc
// @Summary      Danh sách nhãn theo layer
// @Tags         QH Label
// @Produce      json
// @Param        layerId  path  int  true  "Layer ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /v2/tqd/qh/layers/{layerId}/labels [get]
func (h *QHLabelHandler) ListQHLabelsByLayerHTTP(c *gin.Context) {
	layerID, err := parseID(c, "layerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "layerId không hợp lệ"})
		return
	}

	// SỬA: ListByLayerID trả về 3 giá trị: labels, total, error
	labels, total, err := h.labelUsecase.ListByLayerID(c.Request.Context(), layerID, nil, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  labels,
		"total": total,
	})
}

// parseID là helper parse uint64 từ path param
func parseID(c *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}
