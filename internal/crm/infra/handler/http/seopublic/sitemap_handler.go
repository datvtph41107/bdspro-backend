package seopublic

import (
	"net/http"
	"strings"

	"crm/internal/usecase"

	"github.com/gin-gonic/gin"
)

const unavailableSitemapXML = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`

// SitemapHandler serves the public sitemap as raw XML.
//
// CRM remains the owner of sitemap eligibility and XML rendering.
// Gateway and Next only transport this response.
type SitemapHandler struct {
	usecase       *usecase.SeoSitemapUsecase
	publicBaseURL string
}

func NewSitemapHandler(
	uc *usecase.SeoSitemapUsecase,
	publicBaseURL string,
) *SitemapHandler {
	return &SitemapHandler{
		usecase:       uc,
		publicBaseURL: strings.TrimRight(strings.TrimSpace(publicBaseURL), "/"),
	}
}

func (h *SitemapHandler) Get(c *gin.Context) {
	response, err := h.usecase.BuildSitemapXML(
		c.Request.Context(),
		h.publicBaseURL,
	)
	if err != nil {
		c.Header("Content-Type", "application/xml; charset=utf-8")
		c.Header("Cache-Control", "no-store, max-age=0")
		c.Header("Retry-After", "60")
		c.Header("X-Content-Type-Options", "nosniff")

		if c.Request.Method == http.MethodHead {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Data(
			http.StatusServiceUnavailable,
			"application/xml; charset=utf-8",
			[]byte(unavailableSitemapXML),
		)
		return
	}

	contentType := strings.TrimSpace(response.ContentType)
	if contentType == "" {
		contentType = "application/xml; charset=utf-8"
	}

	c.Header("Content-Type", contentType)
	c.Header(
		"Cache-Control",
		"public, max-age=300, stale-while-revalidate=3600",
	)
	c.Header("X-Content-Type-Options", "nosniff")

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
		return
	}

	c.Data(
		http.StatusOK,
		contentType,
		[]byte(response.XML),
	)
}
