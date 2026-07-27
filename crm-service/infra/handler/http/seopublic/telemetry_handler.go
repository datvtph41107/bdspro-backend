package seopublic

import (
	"crypto/hmac"
	"net/http"
	"strings"

	"crm/internal/dto"
	"crm/internal/usecase"

	"github.com/gin-gonic/gin"
)

const seoRuntimeTelemetryKeyHeader = "X-QHPro-Telemetry-Key"

type TelemetryHandler struct {
	usecase     *usecase.SeoRuntimeEventUsecase
	expectedKey string
}

func NewTelemetryHandler(uc *usecase.SeoRuntimeEventUsecase, expectedKey string) *TelemetryHandler {
	return &TelemetryHandler{usecase: uc, expectedKey: strings.TrimSpace(expectedKey)}
}

func (h *TelemetryHandler) Ingest(c *gin.Context) {
	if h == nil || h.usecase == nil || h.expectedKey == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": http.StatusServiceUnavailable, "message": "SEO runtime telemetry ingest chưa được cấu hình"})
		return
	}
	providedKey := strings.TrimSpace(c.GetHeader(seoRuntimeTelemetryKeyHeader))
	if providedKey == "" || !hmac.Equal([]byte(providedKey), []byte(h.expectedKey)) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "SEO runtime telemetry key không hợp lệ"})
		return
	}

	var req dto.SeoRuntimeEventIngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "SEO runtime event payload không hợp lệ"})
		return
	}
	response, err := h.usecase.Ingest(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	c.JSON(http.StatusAccepted, response)
}
