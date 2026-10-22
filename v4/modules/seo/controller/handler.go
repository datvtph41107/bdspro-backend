// Package controller registers all HTTP routes for modules/seo.
// Path parity with CRM public/admin API is maintained intentionally.
package controller

import (
	"net/http"

	"mdn4/modules/seo/services"
)

// Handler is the modules/seo HTTP entry for all SEO capabilities.
type Handler struct {
	news      *NewsController
	domain    *SeoDomainController
	sitemap   *SitemapController
	domainSvc services.SeoDomainService
}

func NewHandler(
	newsSvc services.NewsService,
	domainSvc services.SeoDomainService,
) *Handler {
	return &Handler{
		news:      NewNewsController(newsSvc),
		domain:    NewSeoDomainController(domainSvc),
		sitemap:   NewSitemapController(domainSvc),
		domainSvc: domainSvc,
	}
}

// RegisterHTTP mounts all SEO routes.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.news.RegisterHTTP(mux)
	h.domain.RegisterHTTP(mux)
	h.sitemap.RegisterHTTP(mux)
}

// RegisterInternalHTTP mounts the renderer-only /v2/internal/seo/* routes guarded by secret.
func (h *Handler) RegisterInternalHTTP(mux *http.ServeMux, secret string) {
	NewInternalController(h.domainSvc, secret).RegisterHTTP(mux)
}
