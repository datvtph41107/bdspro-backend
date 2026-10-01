package controller

import (
	"net/http"

	"userv4/modules/seo/services"
)

const unavailableSitemapXML = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`

// SitemapController serves the sitemap XML for public search engines.
type SitemapController struct {
	svc           services.SeoDomainService
	publicBaseURL string
}

func NewSitemapController(svc services.SeoDomainService) *SitemapController {
	return &SitemapController{svc: svc, publicBaseURL: "https://qhpro.vn"}
}

func (c *SitemapController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/crm/public/content/sitemap.xml", c.serve)
	mux.HandleFunc("GET /v2/crm/sitemap.xml", c.serve)
	mux.HandleFunc("GET /v2/seo/sitemap.xml", c.serve)
}

func (c *SitemapController) serve(w http.ResponseWriter, r *http.Request) {
	xml, err := c.svc.ListForSitemap(r.Context(), c.publicBaseURL)
	if err != nil {
		writeXML(w, unavailableSitemapXML, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=3600")
	writeXML(w, xml, http.StatusOK)
}
