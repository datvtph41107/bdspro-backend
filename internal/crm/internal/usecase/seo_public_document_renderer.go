package usecase

import (
	"crm/internal"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	_errors "common/errors"
	"crm/internal/dto"
)

type seoDocumentRenderOptions struct {
	CanonicalURL string
	TemplateKey  string
	ModuleID     string
	Environment  string
	GeneratedAt  time.Time
}

type seoDocumentTemplateView struct {
	Document       *dto.SeoPublicDocument
	CanonicalURL   string
	Robots         string
	TemplateKey    string
	ModuleID       string
	Environment    string
	GeneratedAtISO string
	StructuredData template.JS
	PrimaryMedia   *dto.SeoMedia
	HasRelated     bool
}

func renderSeoPublicDocumentHTML(document *dto.SeoPublicDocument, options seoDocumentRenderOptions) (string, error) {
	if document == nil {
		return "", _errors.ReturnError(service.SEOPublicDocumentRequired)
	}
	if strings.TrimSpace(document.Heading.Title) == "" {
		return "", _errors.ReturnError(service.SEOTitleRequiredForRender)
	}

	canonicalURL := firstNonEmptyRender(options.CanonicalURL, document.SEO.CanonicalPath)
	robots := "noindex,follow"
	if document.SEO.Indexable {
		robots = "index,follow"
	}

	structuredData, err := buildSeoStructuredDataJSON(document, canonicalURL)
	if err != nil {
		return "", fmt.Errorf("create SEO JSON-LD: %w", err)
	}

	var primaryMedia *dto.SeoMedia
	if document.Media != nil {
		primaryMedia = document.Media.Primary
	}

	view := seoDocumentTemplateView{
		Document:       document,
		CanonicalURL:   canonicalURL,
		Robots:         robots,
		TemplateKey:    firstNonEmptyRender(options.TemplateKey, document.ModuleID, "default"),
		ModuleID:       firstNonEmptyRender(options.ModuleID, document.ModuleID, "default"),
		Environment:    firstNonEmptyRender(options.Environment, "crm-rendered"),
		GeneratedAtISO: options.GeneratedAt.Format(time.RFC3339),
		StructuredData: template.JS(structuredData),
		PrimaryMedia:   primaryMedia,
		HasRelated:     len(document.RelatedLinks) > 0 || len(document.RelatedEntities) > 0,
	}

	tpl, err := template.New("seo_public_document").Parse(seoPublicDocumentHTMLTemplate)
	if err != nil {
		return "", fmt.Errorf("parse public SEO template: %w", err)
	}

	var output strings.Builder
	if err := tpl.Execute(&output, view); err != nil {
		return "", fmt.Errorf("render public SEO document: %w", err)
	}
	return output.String(), nil
}

func buildSeoStructuredDataJSON(document *dto.SeoPublicDocument, canonicalURL string) (string, error) {
	if document == nil {
		return "[]", nil
	}

	mainEntity := map[string]any{
		"@context":     "https://schema.org",
		"@type":        schemaTypeForModule(document.ModuleID),
		"@id":          canonicalURL + "#primary",
		"url":          canonicalURL,
		"name":         document.Heading.Title,
		"description":  document.Heading.Summary,
		"dateModified": document.Source.UpdatedAt,
		"isPartOf": map[string]any{
			"@type": "WebSite",
			"name":  "QHPro",
			"url":   siteOriginFromCanonical(canonicalURL),
		},
	}
	if document.Source.SourceName != "" {
		mainEntity["provider"] = map[string]any{
			"@type": "Organization",
			"name":  document.Source.SourceName,
		}
	}
	if document.Source.PublishedAt != "" {
		mainEntity["datePublished"] = document.Source.PublishedAt
	}
	if document.SEO.ImageURL != "" {
		mainEntity["image"] = document.SEO.ImageURL
	}
	if len(document.Facts) > 0 {
		properties := make([]map[string]any, 0, len(document.Facts))
		for _, fact := range document.Facts {
			properties = append(properties, map[string]any{
				"@type": "PropertyValue",
				"name":  fact.Label,
				"value": fact.Value,
			})
		}
		mainEntity["additionalProperty"] = properties
	}

	breadcrumbItems := make([]map[string]any, 0, len(document.Breadcrumbs))
	for index, item := range document.Breadcrumbs {
		breadcrumbItems = append(breadcrumbItems, map[string]any{
			"@type":    "ListItem",
			"position": index + 1,
			"name":     item.Label,
			"item":     absolutePublicURL(canonicalURL, item.Path),
		})
	}
	breadcrumb := map[string]any{
		"@context":        "https://schema.org",
		"@type":           "BreadcrumbList",
		"itemListElement": breadcrumbItems,
	}

	payload := []any{mainEntity, breadcrumb}
	if document.SEO.Indexable && document.FAQ != nil && len(document.FAQ.Items) > 0 {
		faqItems := make([]map[string]any, 0, len(document.FAQ.Items))
		for _, item := range document.FAQ.Items {
			faqItems = append(faqItems, map[string]any{
				"@type": "Question",
				"name":  item.Question,
				"acceptedAnswer": map[string]any{
					"@type": "Answer",
					"text":  item.Answer,
				},
			})
		}
		payload = append(payload, map[string]any{
			"@context":   "https://schema.org",
			"@type":      "FAQPage",
			"mainEntity": faqItems,
		})
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func schemaTypeForModule(moduleID string) string {
	switch normalizeSeoModuleID(moduleID) {
	case "administrative-unit", "parcel", "planning-region":
		return "Place"
	case "planning-news":
		return "Article"
	case "planning-report":
		return "Report"
	default:
		return "Project"
	}
}

func siteOriginFromCanonical(canonicalURL string) string {
	canonicalURL = strings.TrimSpace(canonicalURL)
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(canonicalURL, prefix) {
			remainder := strings.TrimPrefix(canonicalURL, prefix)
			if slash := strings.IndexByte(remainder, '/'); slash >= 0 {
				return prefix + remainder[:slash]
			}
			return canonicalURL
		}
	}
	return "https://qhpro.vn"
}

func absolutePublicURL(canonicalURL string, path string) string {
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		return path
	}
	return strings.TrimRight(siteOriginFromCanonical(canonicalURL), "/") + "/" + strings.TrimLeft(path, "/")
}

const seoPublicDocumentHTMLTemplate = `<!doctype html>
<html lang="vi">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Document.SEO.Title}}</title>
  <meta name="description" content="{{.Document.SEO.Description}}">
  <link rel="canonical" href="{{.CanonicalURL}}">
  <meta name="robots" content="{{.Robots}}">
  <meta property="og:type" content="website">
  <meta property="og:title" content="{{.Document.SEO.Title}}">
  <meta property="og:description" content="{{.Document.SEO.Description}}">
  <meta property="og:url" content="{{.CanonicalURL}}">
  {{if .Document.SEO.ImageURL}}<meta property="og:image" content="{{.Document.SEO.ImageURL}}">{{end}}
  <meta name="x-seo-template" content="{{.TemplateKey}}">
  <meta name="x-seo-generated-at" content="{{.GeneratedAtISO}}">
  <script type="application/ld+json">{{.StructuredData}}</script>
  <style>
    :root{color-scheme:light;--ink:#0f172a;--muted:#64748b;--line:#e2e8f0;--paper:#fff;--surface:#f8fafc;--brand:#2563eb;--brand-dark:#172554;--warning:#fffbeb;--warning-line:#fde68a;--info:#eff6ff;--info-line:#bfdbfe}
    *{box-sizing:border-box}html{scroll-behavior:smooth}body{margin:0;background:var(--surface);color:var(--ink);font-family:Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;line-height:1.6}a{color:inherit}.container{width:min(1180px,calc(100% - 32px));margin-inline:auto}.site-header{background:#fff;border-bottom:1px solid var(--line)}.site-header__inner{min-height:64px;display:flex;align-items:center;justify-content:space-between;gap:20px}.brand{font-weight:900;letter-spacing:-.03em;text-decoration:none;font-size:22px}.brand span{color:var(--brand)}.site-header__nav{display:flex;gap:18px;font-size:14px;font-weight:700}.site-header__nav a{text-decoration:none;color:#334155}.breadcrumbs{background:#fff;border-bottom:1px solid var(--line)}.breadcrumbs ol{list-style:none;padding:14px 0;margin:0;display:flex;flex-wrap:wrap;gap:8px;font-size:14px;color:var(--muted)}.breadcrumbs li{display:flex;align-items:center;gap:8px}.breadcrumbs li+li:before{content:"/";color:#cbd5e1}.breadcrumbs a{text-decoration:none}.hero{background:linear-gradient(135deg,#172554 0%,#0f172a 58%,#111827 100%);color:#fff}.hero__grid{display:grid;grid-template-columns:minmax(0,1fr) minmax(280px,420px);gap:48px;padding:72px 0}.eyebrow{display:inline-flex;border:1px solid rgba(147,197,253,.28);background:rgba(59,130,246,.12);color:#bfdbfe;border-radius:999px;padding:6px 12px;font-size:12px;font-weight:900;text-transform:uppercase;letter-spacing:.14em}.hero h1{font-size:clamp(38px,6vw,64px);line-height:1.05;letter-spacing:-.045em;margin:20px 0 0}.hero__summary{font-size:18px;line-height:1.8;color:#cbd5e1;max-width:780px;margin:24px 0 0}.button{display:inline-flex;align-items:center;justify-content:center;min-height:48px;margin-top:28px;padding:12px 20px;border-radius:12px;background:#60a5fa;color:#082f49;font-weight:900;text-decoration:none}.media-card{border:1px solid rgba(255,255,255,.12);background:rgba(255,255,255,.06);padding:14px;border-radius:24px;align-self:start}.media-card img{display:block;width:100%;aspect-ratio:4/3;object-fit:contain;background:#0f172a;border-radius:16px}.media-card figcaption{font-size:13px;color:#cbd5e1;margin-top:10px}.facts{padding:36px 0}.facts dl{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:14px;margin:0}.fact{background:#fff;border:1px solid var(--line);border-radius:18px;padding:20px;box-shadow:0 10px 30px rgba(15,23,42,.04)}.fact dt{font-size:11px;text-transform:uppercase;letter-spacing:.12em;color:#94a3b8;font-weight:900}.fact dd{margin:8px 0 0;font-size:18px;font-weight:850}.fact p{margin:8px 0 0;color:var(--muted);font-size:13px}.content-grid{display:grid;grid-template-columns:minmax(0,1fr) 320px;gap:28px;padding-bottom:64px}.article,.aside{display:flex;flex-direction:column;gap:22px}.card{background:#fff;border:1px solid var(--line);border-radius:22px;padding:28px;box-shadow:0 10px 30px rgba(15,23,42,.04)}.card h2{margin:0;font-size:25px;line-height:1.25;letter-spacing:-.02em}.blocks{display:flex;flex-direction:column;gap:18px;margin-top:18px}.blocks p{white-space:pre-line;margin:0;color:#334155;line-height:1.8}.blocks ul{margin:0;padding-left:22px;color:#334155}.blocks li+li{margin-top:8px}.notice{border-radius:16px;padding:18px}.notice--warning{background:var(--warning);border:1px solid var(--warning-line);color:#78350f}.notice--info{background:var(--info);border:1px solid var(--info-line);color:#1e3a8a}.notice h3{margin:0 0 6px}.timeline{list-style:none;padding:0 0 0 22px;margin:0;border-left:2px solid #dbeafe}.timeline li{position:relative;padding:0 0 24px 12px}.timeline li:last-child{padding-bottom:0}.timeline li:before{content:"";position:absolute;left:-29px;top:6px;width:12px;height:12px;border-radius:999px;background:var(--brand);box-shadow:0 0 0 4px #dbeafe}.timeline h3{margin:0;font-size:16px}.timeline time{font-size:13px;color:var(--muted)}.timeline p{margin-top:7px}.faq details{border-top:1px solid var(--line);padding:16px 0}.faq details:first-of-type{border-top:0}.faq summary{cursor:pointer;font-weight:850}.faq p{margin:10px 0 0;color:#334155}.source dl{margin:16px 0 0;display:grid;gap:14px;font-size:14px}.source dt{color:#94a3b8}.source dd{margin:2px 0 0;font-weight:750;overflow-wrap:anywhere}.limitations{margin:18px 0 0;padding:16px 16px 16px 32px;border-radius:14px;background:var(--warning);color:#78350f;font-size:13px}.related-item{border:1px solid var(--line);border-radius:14px;padding:15px}.related-item+.related-item{margin-top:10px}.related-item strong{display:block}.related-item p{font-size:13px;color:var(--muted);margin:5px 0 0}.related-actions{display:flex;flex-wrap:wrap;gap:12px;margin-top:10px}.related-actions a{font-size:13px;color:#1d4ed8;font-weight:850;text-decoration:none}.link-list{list-style:none;padding:0;margin:16px 0 0}.link-list li+li{border-top:1px solid var(--line)}.link-list a{display:block;padding:12px 0;color:#1d4ed8;font-weight:800;text-decoration:none}.link-list small{display:block;color:var(--muted);font-weight:400}.readiness{display:inline-flex;margin-top:16px;padding:5px 10px;border-radius:999px;font-size:12px;font-weight:900;background:#ecfdf5;color:#166534}.readiness--not-ready{background:#fff7ed;color:#9a3412}.site-footer{border-top:1px solid var(--line);background:#fff}.site-footer__inner{padding:28px 0;color:var(--muted);font-size:13px;display:flex;justify-content:space-between;gap:20px;flex-wrap:wrap}@media(max-width:900px){.hero__grid,.content-grid{grid-template-columns:1fr}.hero__grid{padding:52px 0}.facts dl{grid-template-columns:repeat(2,minmax(0,1fr))}.site-header__nav{display:none}}@media(max-width:560px){.container{width:min(100% - 24px,1180px)}.facts dl{grid-template-columns:1fr}.card{padding:21px}.hero h1{font-size:39px}}
  </style>
</head>
<body data-qhpro-seo-document data-qhpro-seo-module="{{.ModuleID}}" data-qhpro-seo-environment="{{.Environment}}" data-qhpro-seo-indexable="{{.Document.SEO.Indexable}}">
  <header class="site-header"><div class="container site-header__inner"><a class="brand" href="/">QH<span>Pro</span></a><nav class="site-header__nav" aria-label="Điều hướng chính"><a href="/ban-do">Bản đồ</a><a href="/quy-hoach/tra-cuu">Tra cứu</a><a href="/quy-hoach/tai-lieu">Tài liệu</a></nav></div></header>
  <nav class="breadcrumbs" aria-label="Breadcrumb"><ol class="container">{{range .Document.Breadcrumbs}}<li><a href="{{.Path}}">{{.Label}}</a></li>{{end}}</ol></nav>
  <main data-seo-slug="{{.Document.Identity.CanonicalPath}}">
    <section class="hero"><div class="container hero__grid"><div>{{if .Document.Heading.Eyebrow}}<span class="eyebrow">{{.Document.Heading.Eyebrow}}</span>{{end}}<h1>{{.Document.Heading.Title}}</h1><p class="hero__summary">{{.Document.Heading.Summary}}</p>{{if .Document.MapTarget}}<a class="button" href="{{.Document.MapTarget.Path}}">{{.Document.MapTarget.Label}}</a>{{end}}</div>{{if .PrimaryMedia}}<figure class="media-card">{{if .PrimaryMedia.URL}}<img src="{{.PrimaryMedia.URL}}" alt="{{.PrimaryMedia.Alt}}"{{if .PrimaryMedia.Width}} width="{{.PrimaryMedia.Width}}"{{end}}{{if .PrimaryMedia.Height}} height="{{.PrimaryMedia.Height}}"{{end}}>{{else}}<div role="img" aria-label="{{.PrimaryMedia.Alt}}" style="display:grid;place-items:center;aspect-ratio:4/3;border-radius:16px;background:#0f172a;color:#cbd5e1;padding:28px;text-align:center">{{.PrimaryMedia.Alt}}</div>{{end}}<figcaption>{{if .Document.Media.Note}}{{.Document.Media.Note}}{{else}}{{.PrimaryMedia.Alt}}{{end}}</figcaption></figure>{{end}}</div></section>
    <section class="facts"><div class="container"><dl>{{range .Document.Facts}}<div class="fact" data-seo-fact-id="{{.ID}}"><dt>{{.Label}}</dt><dd>{{.Value}}</dd>{{if .Description}}<p>{{.Description}}</p>{{end}}</div>{{end}}</dl></div></section>
    <div class="container content-grid"><article class="article">{{range .Document.Sections}}<section class="card" id="{{.ID}}" data-seo-section-id="{{.ID}}"><h2>{{.Title}}</h2><div class="blocks">{{range .Blocks}}{{if eq .Type "paragraph"}}<p>{{.Text}}</p>{{else if eq .Type "list"}}<ul>{{range .Items}}<li>{{.}}</li>{{end}}</ul>{{else if eq .Type "notice"}}<div class="notice notice--{{.Tone}}">{{if .Title}}<h3>{{.Title}}</h3>{{end}}<p>{{.Text}}</p></div>{{else if eq .Type "timeline"}}<ol class="timeline">{{range .TimelineItems}}<li data-seo-timeline-id="{{.ID}}"><h3>{{.Title}}</h3>{{if .Date}}<time datetime="{{.Date}}">{{.Date}}</time>{{end}}{{if .Description}}<p>{{.Description}}</p>{{end}}</li>{{end}}</ol>{{end}}{{end}}</div></section>{{end}}{{if .Document.FAQ}}<section class="card faq" id="faq"><h2>Giải đáp thắc mắc</h2>{{range .Document.FAQ.Items}}<details data-seo-faq-id="{{.ID}}"><summary>{{.Question}}</summary><p>{{.Answer}}</p></details>{{end}}</section>{{end}}</article>
      <aside class="aside"><section class="card source"><h2>Nguồn dữ liệu</h2><dl><div><dt>Nguồn</dt><dd>{{if .Document.Source.SourceURL}}<a href="{{.Document.Source.SourceURL}}" rel="nofollow noopener">{{if .Document.Source.SourceName}}{{.Document.Source.SourceName}}{{else}}{{.Document.Source.SourceURL}}{{end}}</a>{{else}}{{if .Document.Source.SourceName}}{{.Document.Source.SourceName}}{{else}}Nội dung biên tập{{end}}{{end}}</dd></div><div><dt>Cập nhật</dt><dd><time datetime="{{.Document.Source.UpdatedAt}}">{{.Document.Source.UpdatedAt}}</time></dd></div></dl></section>
      {{if .Document.RelatedEntities}}<section class="card"><h2>Dữ liệu liên quan</h2>{{range .Document.RelatedEntities}}<div class="related-item"><strong>{{.Title}}</strong>{{if .Description}}<p>{{.Description}}</p>{{end}}<div class="related-actions">{{if .PublicPath}}<a href="{{.PublicPath}}">Xem chi tiết</a>{{end}}{{if .MapPath}}<a href="{{.MapPath}}">Xem bản đồ</a>{{end}}</div></div>{{end}}</section>{{end}}
      {{if .Document.RelatedLinks}}<section class="card"><h2>Liên kết hữu ích</h2><ul class="link-list">{{range .Document.RelatedLinks}}<li><a href="{{.Path}}">{{.Label}}{{if .Description}}<small>{{.Description}}</small>{{end}}</a></li>{{end}}</ul></section>{{end}}</aside></div>
  </main>
  <footer class="site-footer"><div class="container site-footer__inner"><span>QHPro · Nền tảng thông tin quy hoạch</span><span>Dữ liệu tham khảo; vui lòng đối chiếu hồ sơ chính thức.</span></div></footer>
</body>
</html>`

func formatRenderError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
