// Package controller registers all HTTP routes for modules/seo.
// Path parity with CRM public/admin API is maintained intentionally.
package controller

import (
	"net/http"

	"userv4/modules/seo/services"
)

// Handler is the modules/seo HTTP entry for all SEO capabilities.
type Handler struct {
	news      *NewsController
	domain    *SeoDomainController
	sitemap   *SitemapController
}

func NewHandler(
	newsSvc services.NewsService,
	domainSvc services.SeoDomainService,
) *Handler {
	return &Handler{
		news:    NewNewsController(newsSvc),
		domain:  NewSeoDomainController(domainSvc),
		sitemap: NewSitemapController(domainSvc),
	}
}

// RegisterHTTP mounts all SEO routes.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.news.RegisterHTTP(mux)
	h.domain.RegisterHTTP(mux)
	h.sitemap.RegisterHTTP(mux)
}
