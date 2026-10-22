package services

import (
	"testing"

	"mdn4/modules/seo/entity"
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

func TestResolvePublicURLPattern_FiveGenerateSources(t *testing.T) {
	cases := []struct {
		refType entity.SEORefType
		source  string
		want    string
	}{
		{entity.SEORefTypeParcel, entity.SeoRefSourceTqdParcel, entity.PublicURLPatternThuaDat},
		{entity.SEORefTypeMap, entity.SeoRefSourceTqdLayer, entity.PublicURLPatternLop},
		{entity.SEORefTypeProject, entity.SeoRefSourceTqdPlanningProject, entity.PublicURLPatternDoAnQuyHoach},
		{entity.SEORefTypeDocument, entity.SeoRefSourceTqdDocument, entity.PublicURLPatternTaiLieu},
		{entity.SEORefTypeAdmUnit, entity.SeoRefSourceTqdAdministrativeUnit, entity.PublicURLPatternDiaBan},
		{entity.SEORefTypeProject, "tqd_planning_project", entity.PublicURLPatternDoAnQuyHoach},
		{entity.SEORefTypeLegacy, "manual", entity.PublicURLPatternTin},
	}
	for _, tc := range cases {
		got := entity.ResolvePublicURLPattern(tc.refType, tc.source)
		if got != tc.want {
			t.Fatalf("ResolvePublicURLPattern(%d,%q)=%q want %q", tc.refType, tc.source, got, tc.want)
		}
	}
}

func TestFormatPublicURL_PlanningProject(t *testing.T) {
	got := entity.FormatPublicURL(entity.PublicURLPatternDoAnQuyHoach, 11, "da-001")
	want := "/tin-tuc/do-an-quy-hoach/11/da-001"
	if got != want {
		t.Fatalf("FormatPublicURL=%q want %q", got, want)
	}
	if compound := entity.CompoundSeoSlug(11, "da-001"); compound != "11/da-001" {
		t.Fatalf("CompoundSeoSlug=%q want 11/da-001", compound)
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
