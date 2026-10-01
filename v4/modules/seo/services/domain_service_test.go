package services

import (
	"testing"

	"userv4/modules/seo/entity"
)

func TestSEORefTypeName(t *testing.T) {
	cases := []struct {
		ref  entity.SEORefType
		want string
	}{
		{entity.SEORefTypeAdmUnit, "Đơn vị hành chính"},
		{entity.SEORefTypeParcel, "Thửa đất"},
		{entity.SEORefTypeRegion, "Vùng"},
		{entity.SEORefTypeProject, "Đồ án"},
		{entity.SEORefTypeDocument, "Tài liệu"},
		{entity.SEORefTypeMap, "Bản đồ"},
		{entity.SEORefTypeReport, "Báo cáo"},
		{entity.SEORefTypeLegacy, "Nội dung thủ công"},
		{entity.SEORefType(0), ""},
	}
	for _, tc := range cases {
		if got := tc.ref.Name(); got != tc.want {
			t.Fatalf("SEORefType(%d).Name()=%q want %q", tc.ref, got, tc.want)
		}
	}
}

func TestSeoDomainNormalizeLifecycle(t *testing.T) {
	seo := &entity.SeoDomain{Published: true, NeedGenerate: true}
	seo.NormalizeLifecycle()
	if seo.PageStatus != entity.PageStatusPublished {
		t.Fatalf("page status=%q want published", seo.PageStatus)
	}
	if seo.RenderStatus != entity.RenderStatusPending {
		t.Fatalf("render status=%q want pending", seo.RenderStatus)
	}
}

func TestSeoDomainIsSitemapEligible(t *testing.T) {
	seo := &entity.SeoDomain{
		PageStatus: entity.PageStatusPublished,
		Published:  true,
		IsIndex:    true,
		IsSiteMap:  true,
		Scope:      entity.ScopePublic,
	}
	if !seo.IsSitemapEligible() {
		t.Fatal("expected sitemap eligible")
	}
	seo.Scope = entity.ScopeInternal
	if seo.IsSitemapEligible() {
		t.Fatal("internal scope must not be sitemap eligible")
	}
}

func TestNormalizeLastDates(t *testing.T) {
	if got := normalizeLastDates(0); got != 3 {
		t.Fatalf("default lastDates=%d want 3", got)
	}
	if got := normalizeLastDates(100); got != 30 {
		t.Fatalf("capped lastDates=%d want 30", got)
	}
	if got := normalizeLastDates(7); got != 7 {
		t.Fatalf("lastDates=%d want 7", got)
	}
}
