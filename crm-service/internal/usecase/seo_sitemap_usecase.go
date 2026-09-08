package usecase

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SeoSitemapUsecase struct {
	seoDomainRepo repo.SeoDomainRepo
}

func NewSeoSitemapUsecase(seoDomainRepo repo.SeoDomainRepo) *SeoSitemapUsecase {
	return &SeoSitemapUsecase{
		seoDomainRepo: seoDomainRepo,
	}
}

type seoSitemapXMLURLSet struct {
	XMLName xml.Name           `xml:"urlset"`
	Xmlns   string             `xml:"xmlns,attr"`
	URLs    []seoSitemapXMLURL `xml:"url"`
}

type seoSitemapXMLURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

func (u *SeoSitemapUsecase) BuildSitemapXML(ctx context.Context, baseURL string) (*dto.SeoSitemapXMLResponse, error) {
	items, err := u.seoDomainRepo.GetPublishedSitemapItems(ctx)
	if err != nil {
		return nil, err
	}

	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")

	urls := make([]seoSitemapXMLURL, 0, len(items))
	for i := range items {
		page := &items[i]

		if !page.IsSitemapEligible() ||
			page.RenderStatus != seo_domain.SeoRenderStatusSuccess ||
			strings.TrimSpace(page.CanonicalURL) == "" {
			continue
		}

		loc := buildSeoSitemapLoc(page.CanonicalURL, baseURL)
		if strings.TrimSpace(loc) == "" {
			continue
		}

		urls = append(urls, seoSitemapXMLURL{
			Loc:        loc,
			LastMod:    resolveSeoSitemapLastMod(page),
			ChangeFreq: strings.TrimSpace(page.SitemapChangeFreq),
			Priority:   formatSeoSitemapPriority(page.SitemapPriority),
		})
	}

	out := seoSitemapXMLURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	raw, err := xml.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, err
	}

	return &dto.SeoSitemapXMLResponse{
		XML:         xml.Header + string(raw),
		ContentType: "application/xml; charset=utf-8",
		Total:       len(urls),
	}, nil
}

func buildSeoSitemapLoc(canonicalURL string, baseURL string) string {
	loc := strings.TrimSpace(canonicalURL)
	if loc == "" {
		return ""
	}

	if strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
		return loc
	}

	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return "/" + strings.TrimLeft(loc, "/")
	}

	return baseURL + "/" + strings.TrimLeft(loc, "/")
}

func resolveSeoSitemapLastMod(page *seo_domain.SeoDomain) string {
	if page == nil {
		return ""
	}

	if page.SiteMapLastedAt != nil && !page.SiteMapLastedAt.IsZero() {
		return page.SiteMapLastedAt.Format(time.RFC3339)
	}

	if !page.UpdatedAt.IsZero() {
		return page.UpdatedAt.Format(time.RFC3339)
	}

	if page.GeneratedAt != nil && !page.GeneratedAt.IsZero() {
		return page.GeneratedAt.Format(time.RFC3339)
	}

	return ""
}

func formatSeoSitemapPriority(priority float32) string {
	if priority <= 0 {
		priority = 0.5
	}

	if priority > 1 {
		priority = 1
	}

	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", priority), "0"), ".")
}
