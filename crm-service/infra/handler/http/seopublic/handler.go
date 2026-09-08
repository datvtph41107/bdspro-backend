package seopublic

import (
	"net/http"
	"strings"

	"crm/internal/usecase"
	sharedpb "pb/types/shared"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type Handler struct {
	usecase *usecase.SeoPublicUsecase
}

func NewHandler(uc *usecase.SeoPublicUsecase) *Handler {
	return &Handler{usecase: uc}
}

func (h *Handler) GetDocument(c *gin.Context) {
	response, err := h.usecase.GetPublicSeoDocument(
		c.Request.Context(),
		c.Param("module"),
		c.Param("slug"),
	)
	if err != nil {
		writeError(c, err)
		return
	}

	if cacheControl := strings.TrimSpace(response.CacheControl); cacheControl != "" {
		c.Header("Cache-Control", cacheControl)
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-Content-Type-Options", "nosniff")
	if response.Document != nil && !response.Document.SEO.Indexable {
		c.Header("X-Robots-Tag", "noindex, follow")
	}

	etag := normalizeETag(response.ETag)
	if etag != "" {
		c.Header("ETag", etag)
		if matchesETag(c.GetHeader("If-None-Match"), etag) {
			c.Status(http.StatusNotModified)
			return
		}
	}

	statusCode := response.StatusCode
	if statusCode <= 0 {
		statusCode = http.StatusOK
	}
	if c.Request.Method == http.MethodHead {
		c.Status(statusCode)
		return
	}
	c.JSON(statusCode, response.Document)
}

func (h *Handler) GetPageBySlug(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))

	response, err := h.usecase.GetPublicSeoPageBySlug(
		c.Request.Context(),
		slug,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	statusCode := response.StatusCode
	if statusCode <= 0 {
		statusCode = http.StatusOK
	}

	contentType := strings.TrimSpace(response.ContentType)
	if contentType == "" {
		contentType = "text/html; charset=utf-8"
	}

	if cacheControl := strings.TrimSpace(response.CacheControl); cacheControl != "" {
		c.Header("Cache-Control", cacheControl)
	}

	if response.NoIndex {
		c.Header("X-Robots-Tag", "noindex, follow")
	}

	c.Header("Content-Type", contentType)
	c.Header("X-Content-Type-Options", "nosniff")

	etag := normalizeETag(response.ETag)
	if etag != "" {
		c.Header("ETag", etag)

		if matchesETag(c.GetHeader("If-None-Match"), etag) {
			c.Status(http.StatusNotModified)
			return
		}
	}

	if c.Request.Method == http.MethodHead {
		c.Status(statusCode)
		return
	}

	c.Data(statusCode, contentType, []byte(response.HTML))
}

func normalizeETag(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	if strings.HasPrefix(value, `W/"`) {
		return value
	}

	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return value
	}

	return `"` + strings.Trim(value, `"`) + `"`
}

func matchesETag(headerValue string, expected string) bool {
	headerValue = strings.TrimSpace(headerValue)
	if headerValue == "" {
		return false
	}

	for _, candidate := range strings.Split(headerValue, ",") {
		candidate = strings.TrimSpace(candidate)

		if candidate == "*" || candidate == expected {
			return true
		}
	}

	return false
}

func writeError(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"

	grpcStatus, ok := status.FromError(err)
	if ok {
		if strings.TrimSpace(grpcStatus.Message()) != "" {
			message = grpcStatus.Message()
		}

		for _, detail := range grpcStatus.Details() {
			errorResponse, ok := detail.(*sharedpb.ErrorResponse)
			if !ok {
				continue
			}

			if errorResponse.Code >= 400 && errorResponse.Code <= 599 {
				statusCode = int(errorResponse.Code)
			}

			if strings.TrimSpace(errorResponse.Message) != "" {
				message = errorResponse.Message
			}

			break
		}
	}

	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"message": message,
	})
}
