package usecase

import (
	"strings"
	"testing"
	"time"
)

func TestBuildMetadataScriptProducesValidJSONLD(t *testing.T) {
	got := buildMetadataScript(`{"@context":"https://schema.org","name":"A < B"}`)
	if got == "" {
		t.Fatal("expected JSON-LD script")
	}
	if strings.Contains(got, "&#34;") {
		t.Fatalf("JSON-LD must not be HTML entity encoded: %s", got)
	}
	if !strings.Contains(got, `"@context":"https://schema.org"`) {
		t.Fatalf("missing JSON-LD payload: %s", got)
	}
	if !strings.Contains(got, `A \u003c B`) {
		t.Fatalf("expected script-safe JSON escaping: %s", got)
	}
}

func TestRenderSeoHTMLIncludesTelemetryAndRobotsPolicy(t *testing.T) {
	htmlBody, err := renderSeoHTML(renderSeoInput{
		Slug:            "/do-an-quy-hoach/qh2",
		Title:           `Quy hoạch <H2-1>`,
		Description:     "Mô tả",
		CanonicalURL:    "https://qhpro.vn/do-an-quy-hoach/qh2",
		Content:         "<p>Nội dung</p>",
		Summary:         "Tóm tắt",
		TemplateKey:     "planning-project",
		ModuleID:        "planning-project",
		Environment:     "crm-rendered",
		Indexable:       true,
		TelemetryScript: "/static/qhpro/seo-m8-client.js",
		GeneratedAt:     time.Unix(0, 0).UTC(),
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	checks := []string{
		`<meta name="robots" content="index,follow">`,
		`data-qhpro-seo-module="planning-project"`,
		`data-qhpro-seo-environment="crm-rendered"`,
		`src="/static/qhpro/seo-m8-client.js"`,
		`<title>Quy hoạch &lt;H2-1&gt;</title>`,
	}
	for _, check := range checks {
		if !strings.Contains(htmlBody, check) {
			t.Fatalf("missing %q in rendered HTML", check)
		}
	}
}
