package handler_http

import (
	"net/http"
	"strconv"

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
	LayerID     uint64 `json:"layerId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
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
		LayerID:     req.LayerID,
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
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
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
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
// @Success      200  {array}  qh_domain.QHLabel
// @Router       /v2/tqd/qh/layers/{layerId}/labels [get]
func (h *QHLabelHandler) ListQHLabelsByLayerHTTP(c *gin.Context) {
	layerID, err := parseID(c, "layerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "layerId không hợp lệ"})
		return
	}

	result, err := h.labelUsecase.ListByLayerID(c.Request.Context(), layerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": len(result),
	})
}

// parseID là helper parse uint64 từ path param
func parseID(c *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}
