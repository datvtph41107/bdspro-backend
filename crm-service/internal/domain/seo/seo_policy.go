package seo_domain

import "strings"

type SeoPolicyDecision struct {
	Allowed  bool
	Code     string
	Message  string
	Reasons  []string
	Warnings []string
}

func SeoPolicyAllow(warnings ...string) SeoPolicyDecision {
	return SeoPolicyDecision{Allowed: true, Code: "allowed", Warnings: compactStrings(warnings)}
}

func SeoPolicyDeny(code string, message string, reasons ...string) SeoPolicyDecision {
	if code == "" {
		code = "denied"
	}
	return SeoPolicyDecision{
		Allowed: false,
		Code:    code,
		Message: message,
		Reasons: compactStrings(reasons),
	}
}

func (d SeoPolicyDecision) WithWarning(warnings ...string) SeoPolicyDecision {
	d.Warnings = append(d.Warnings, compactStrings(warnings)...)
	return d
}

func compactStrings(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func CanUseSourceForSeo(cfg SeoRefTypeConfig) SeoPolicyDecision {
	cfg = NormalizeSeoRefTypeConfig(cfg)

	if cfg.Status != SeoSourceStatusReady {
		return SeoPolicyDeny("source_not_ready", "Nguồn SEO chưa sẵn sàng", cfg.DisabledReason)
	}
	if !cfg.Capabilities.Resolvable && cfg.Key != SeoRefSourceManual && cfg.ResolverKey != SeoRefSourceManual {
		return SeoPolicyDeny("source_not_resolvable", "Nguồn SEO chưa hỗ trợ resolve entity")
	}
	return SeoPolicyAllow()
}

func CanSearchSource(cfg SeoRefTypeConfig) SeoPolicyDecision {
	cfg = NormalizeSeoRefTypeConfig(cfg)

	if cfg.Status != SeoSourceStatusReady {
		return SeoPolicyDeny("source_not_ready", "Nguồn SEO chưa sẵn sàng", cfg.DisabledReason)
	}
	if !cfg.Capabilities.Searchable {
		return SeoPolicyDeny("source_not_searchable", "Nguồn SEO không hỗ trợ tìm kiếm")
	}
	return SeoPolicyAllow()
}

func CanPublish(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}

	page.NormalizeLifecycle()

	if page.PageStatus == SeoPageStatusArchived {
		return SeoPolicyDeny("page_archived", "Không thể publish page đã archived")
	}
	if strings.TrimSpace(page.CanonicalURL) == "" {
		return SeoPolicyDeny("missing_canonical", "Thiếu canonical URL")
	}
	if strings.TrimSpace(page.Slug) == "" {
		return SeoPolicyDeny("missing_slug", "Thiếu slug")
	}
	if strings.TrimSpace(page.Title) == "" {
		return SeoPolicyDeny("missing_title", "Thiếu title")
	}
	if page.SourceStatus == SeoSourceStatusMissing || page.RefMissing {
		return SeoPolicyDeny("source_missing", "Source gốc không còn tồn tại")
	}
	if page.SourceStatus == SeoSourceStatusPermissionDenied {
		return SeoPolicyDeny("source_permission_denied", "Không đủ quyền truy cập source gốc")
	}

	warnings := make([]string, 0, 4)
	if strings.TrimSpace(page.Description) == "" {
		warnings = append(warnings, "missing_description")
	}
	if strings.TrimSpace(page.Content) == "" {
		warnings = append(warnings, "missing_content")
	}
	if strings.TrimSpace(page.TemplateKey) == "" {
		warnings = append(warnings, "missing_template_key")
	}
	if page.QualityScore > 0 && page.QualityScore < 50 {
		return SeoPolicyDeny("quality_score_too_low", "Quality score chưa đủ để publish")
	}
	if page.QualityScore >= 50 && page.QualityScore < 70 {
		warnings = append(warnings, "quality_score_medium")
	}

	return SeoPolicyAllow(warnings...)
}

func CanUnpublish(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}
	page.NormalizeLifecycle()

	if page.PageStatus == SeoPageStatusArchived {
		return SeoPolicyDeny("page_archived", "Page archived không cần unpublish")
	}
	return SeoPolicyAllow()
}

func CanArchive(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}
	return SeoPolicyAllow()
}

func CanRestore(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}
	page.NormalizeLifecycle()

	if page.PageStatus != SeoPageStatusArchived {
		return SeoPolicyDeny("page_not_archived", "Chỉ restore page đang archived")
	}
	return SeoPolicyAllow()
}

func CanGenerate(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}

	page.NormalizeLifecycle()

	if page.PageStatus != SeoPageStatusPublished || !page.Published {
		return SeoPolicyDeny("page_not_published", "Chỉ generate page đã published")
	}
	if strings.TrimSpace(page.CanonicalURL) == "" {
		return SeoPolicyDeny("missing_canonical", "Thiếu canonical URL")
	}
	if strings.TrimSpace(page.Title) == "" {
		return SeoPolicyDeny("missing_title", "Thiếu title")
	}
	if page.SourceStatus == SeoSourceStatusMissing || page.RefMissing {
		return SeoPolicyDeny("source_missing", "Source gốc không còn tồn tại")
	}
	if page.SourceStatus == SeoSourceStatusPermissionDenied {
		return SeoPolicyDeny("source_permission_denied", "Không đủ quyền truy cập source gốc")
	}

	warnings := make([]string, 0, 3)
	if strings.TrimSpace(page.TemplateKey) == "" {
		warnings = append(warnings, "missing_template_key_use_default")
	}
	if strings.TrimSpace(page.Description) == "" {
		warnings = append(warnings, "missing_description")
	}
	if strings.TrimSpace(page.Content) == "" {
		warnings = append(warnings, "missing_content")
	}

	return SeoPolicyAllow(warnings...)
}

func CanServePublic(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}

	page.NormalizeLifecycle()

	if page.PageStatus == SeoPageStatusArchived {
		return SeoPolicyDeny("page_archived", "Page đã archived")
	}
	if page.PageStatus != SeoPageStatusPublished || !page.Published {
		return SeoPolicyDeny("page_not_published", "Page chưa published")
	}
	if page.Scope != "" && page.Scope != SeoScopePublic {
		return SeoPolicyDeny("not_public_scope", "Page không thuộc public scope")
	}
	if page.RenderStatus != SeoRenderStatusSuccess && strings.TrimSpace(page.RenderedHTML) == "" {
		return SeoPolicyDeny("html_not_ready", "HTML public chưa sẵn sàng")
	}
	if strings.TrimSpace(page.RenderedHTML) == "" {
		return SeoPolicyDeny("empty_html", "Rendered HTML rỗng")
	}

	warnings := make([]string, 0, 2)
	if page.RenderStatus != SeoRenderStatusSuccess {
		warnings = append(warnings, "serve_stale_html")
	}
	if !page.IsIndex {
		warnings = append(warnings, "noindex_page")
	}

	return SeoPolicyAllow(warnings...)
}

func CanIncludeInSitemap(page *SeoDomain) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}

	page.NormalizeLifecycle()

	if !page.IsSitemapEligible() {
		return SeoPolicyDeny("not_sitemap_eligible", "Page chưa đủ điều kiện sitemap")
	}
	if page.RenderStatus != SeoRenderStatusSuccess {
		return SeoPolicyDeny("render_not_success", "Page chưa render thành công")
	}
	if strings.TrimSpace(page.CanonicalURL) == "" {
		return SeoPolicyDeny("missing_canonical", "Thiếu canonical URL")
	}
	if page.SourceStatus == SeoSourceStatusMissing || page.SourceStatus == SeoSourceStatusPermissionDenied {
		return SeoPolicyDeny("bad_source_status", "Source không hợp lệ để sitemap")
	}
	if page.QualityScore > 0 && page.QualityScore < 50 {
		return SeoPolicyDeny("quality_score_too_low", "Quality score chưa đủ để sitemap")
	}

	return SeoPolicyAllow()
}

func NeedRegenerate(page *SeoDomain, reason string) SeoPolicyDecision {
	if page == nil {
		return SeoPolicyDeny("page_nil", "SEO page không tồn tại")
	}

	page.NormalizeLifecycle()

	if page.PageStatus != SeoPageStatusPublished || !page.Published {
		return SeoPolicyDeny("page_not_published", "Page chưa published nên không cần regenerate public HTML")
	}
	if page.NeedGenerate {
		return SeoPolicyAllow("need_generate_true")
	}
	if page.RenderStatus == SeoRenderStatusFailed || page.RenderStatus == SeoRenderStatusNone {
		return SeoPolicyAllow("render_status_requires_generation")
	}
	if strings.TrimSpace(reason) != "" {
		return SeoPolicyAllow(reason)
	}

	return SeoPolicyDeny("no_regenerate_needed", "HTML hiện tại chưa có dấu hiệu cần regenerate")
}
