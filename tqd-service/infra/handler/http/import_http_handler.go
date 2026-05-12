package handler_http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	_utils "common/utils"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"github.com/gin-gonic/gin"
)

const multipartImportMaxBytes = 512 << 20 // 512MB tương tự ngưỡng ParseMultipartForm

// ImportHTTPHandler — REST multipart cho import (stream vào temp, không đọc hết RAM).
type ImportHTTPHandler struct {
	importUsecase usecase.ImportUsecase
}

func NewImportHTTPHandler(importUsecase usecase.ImportUsecase) *ImportHTTPHandler {
	return &ImportHTTPHandler{importUsecase: importUsecase}
}

// PostImportGeoJSONMultipart
// @Summary      Import GeoJSON/NDJSON (multipart form-data)
// @Description  Giống upload file: field **file** (bắt buộc). Form: file_format, layer_id, label_field, label_mappings (JSON string), default_label_id, source_file_name, validate_geometry, auto_fix_geometry, skip_invalid.
// @Tags         TQD Import
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File .json / .geojson / .ndjson"
// @Router       /v2/tqd/admin/import/geojson [post]
func (h *ImportHTTPHandler) PostImportGeoJSONMultipart(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(multipartImportMaxBytes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "multipart form: " + err.Error()})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field \"file\" is required"})
		return
	}

	layerID, err := parseUint64Form(c, "layer_id", "layerId")
	if err != nil || layerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "layer_id is required"})
		return
	}

	fileFormat := firstNonEmpty(c.PostForm("file_format"), c.PostForm("fileFormat"))
	labelField := firstNonEmpty(c.PostForm("label_field"), c.PostForm("labelField"))
	labelMappings, err := parseLabelMappingsForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label_mappings: " + err.Error()})
		return
	}

	var defaultLabelID *uint64
	if s := firstNonEmpty(c.PostForm("default_label_id"), c.PostForm("defaultLabelId")); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "default_label_id invalid"})
			return
		}
		defaultLabelID = &v
	}

	sourceName := firstNonEmpty(c.PostForm("source_file_name"), c.PostForm("sourceFileName"))
	if sourceName == "" {
		sourceName = fh.Filename
	}

	userID := _utils.GetOriginIdFromContext(c.Request.Context())
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	job := &dto.ImportRequest{
		LayerID:          layerID,
		LabelField:       labelField,
		LabelMappings:    labelMappings,
		DefaultLabelID:   defaultLabelID,
		SourceFileName:   sourceName,
		UserID:           userID,
		ValidateGeometry: parseBoolForm(c, "validate_geometry", "validateGeometry"),
		AutoFixGeometry:  parseBoolForm(c, "auto_fix_geometry", "autoFixGeometry"),
		SkipInvalid:      parseBoolForm(c, "skip_invalid", "skipInvalid"),
	}

	src, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// Giới hạn kích thước theo header (best-effort)
	if fh.Size > 0 && fh.Size > multipartImportMaxBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
		return
	}

	result, err := h.importUsecase.EnqueueImportFromReader(c.Request.Context(), src, fh.Filename, fileFormat, job)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already has an import") {
			status = http.StatusConflict
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    result.Code,
		"batchId": result.BatchID,
		"status":  result.Status,
		"message": result.Message,
	})
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}
	return strings.TrimSpace(b)
}

func parseUint64Form(c *gin.Context, snake, camel string) (uint64, error) {
	s := firstNonEmpty(c.PostForm(snake), c.PostForm(camel))
	if s == "" {
		return 0, nil
	}
	return strconv.ParseUint(s, 10, 64)
}

func parseBoolForm(c *gin.Context, snake, camel string) bool {
	s := strings.ToLower(firstNonEmpty(c.PostForm(snake), c.PostForm(camel)))
	return s == "1" || s == "true" || s == "yes"
}

func parseLabelMappingsForm(c *gin.Context) (map[string]uint64, error) {
	raw := firstNonEmpty(c.PostForm("label_mappings"), c.PostForm("labelMappings"))
	if raw == "" {
		return nil, nil
	}
	var m map[string]uint64
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	return m, nil
}
