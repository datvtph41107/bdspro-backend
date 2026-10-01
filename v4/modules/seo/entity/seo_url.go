package entity

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Default public URL patterns for /tin-tuc/<loại>/<id>/<slug>.
const (
	PublicURLPatternTin             = "/tin-tuc/tin/{id}/{slug}"
	PublicURLPatternDiaBan          = "/tin-tuc/dia-ban/{id}/{slug}"
	PublicURLPatternVungQuyHoach    = "/tin-tuc/vung-quy-hoach/{id}/{slug}"
	PublicURLPatternThuaDat         = "/tin-tuc/thua-dat/{id}/{slug}"
	PublicURLPatternDoAnQuyHoach    = "/tin-tuc/do-an-quy-hoach/{id}/{slug}"
	PublicURLPatternTaiLieu         = "/tin-tuc/tai-lieu/{id}/{slug}"
	PublicURLPatternLop             = "/tin-tuc/lop/{id}/{slug}"
	PublicURLPatternBaoCao          = "/tin-tuc/bao-cao/{id}/{slug}"
	PublicURLPatternDuAn            = "/du-an/{id}/{slug}"
)

// NormalizeSeoResolverKey maps generate RefSource aliases onto registry ResolverKey values.
func NormalizeSeoResolverKey(refSource string) string {
	switch strings.ToLower(strings.TrimSpace(refSource)) {
	case "tqd_parcel", "tqd.parcel":
		return SeoRefSourceTqdParcel
	case "tqd_layer", "tqd.qh_layer", "tqd_qh_layer":
		return SeoRefSourceTqdLayer
	case "tqd_planning_project", "tqd.planning_project":
		return SeoRefSourceTqdPlanningProject
	case "tqd_planning_document", "tqd.legal_document", "tqd_legal_document", "tqd_document":
		return SeoRefSourceTqdDocument
	case "tqd_adm_unit", "tqd.administrative_unit", "tqd_administrative_unit":
		return SeoRefSourceTqdAdministrativeUnit
	case "bdspro_project", "bdspro.project":
		return SeoRefSourceBdsproProject
	case "bdspro_area_region", "bdspro.area_region", "bdspro_region", "bdspro.region":
		return SeoRefSourceBdsproAreaRegion
	case "crm_report", "crm.report":
		return SeoRefSourceCrmReport
	case "manual", "":
		return SeoRefSourceManual
	default:
		return strings.TrimSpace(refSource)
	}
}

// ResolvePublicURLPattern picks /tin-tuc/<loại>/{id}/{slug} (or /du-an/...) for a ref.
func ResolvePublicURLPattern(refType SEORefType, refSource string) string {
	resolverKey := NormalizeSeoResolverKey(refSource)
	if cfg, ok := GetSeoRefTypeConfig(uint32(refType), resolverKey, ""); ok && strings.TrimSpace(cfg.PublicURLPattern) != "" {
		return cfg.PublicURLPattern
	}
	if cfg, ok := GetSeoRefTypeConfigByKey(resolverKey); ok && cfg.Value == uint32(refType) && strings.TrimSpace(cfg.PublicURLPattern) != "" {
		return cfg.PublicURLPattern
	}
	// Fallbacks when RefSource is generic (e.g. ref-400).
	switch refType {
	case SEORefTypeAdmUnit:
		return PublicURLPatternDiaBan
	case SEORefTypeParcel:
		return PublicURLPatternThuaDat
	case SEORefTypeRegion:
		return PublicURLPatternVungQuyHoach
	case SEORefTypeProject:
		if resolverKey == SeoRefSourceBdsproProject {
			return PublicURLPatternDuAn
		}
		return PublicURLPatternDoAnQuyHoach
	case SEORefTypeDocument:
		return PublicURLPatternTaiLieu
	case SEORefTypeMap:
		return PublicURLPatternLop
	case SEORefTypeReport:
		return PublicURLPatternBaoCao
	case SEORefTypeLegacy:
		return PublicURLPatternTin
	default:
		return PublicURLPatternTin
	}
}

// FormatPublicURL fills {id}/{slug} into a PublicURLPattern.
func FormatPublicURL(pattern string, id uint64, slug string) string {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		pattern = PublicURLPatternTin
	}
	slug = SeoSlugSegment(slug, id)
	return strings.NewReplacer(
		"{id}", strconv.FormatUint(id, 10),
		"{slug}", slug,
	).Replace(pattern)
}

// CompoundSeoSlug stores slug field as "{id}/{slug}" (CRM parity).
func CompoundSeoSlug(id uint64, slug string) string {
	return fmt.Sprintf("%d/%s", id, SeoSlugSegment(slug, id))
}

// SeoSlugSegment normalizes a label into a URL slug segment.
func SeoSlugSegment(label string, fallbackID uint64) string {
	slug := slugifySEO(label)
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.Trim(slug, "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	if slug == "" || slug == "seo-page" {
		return strconv.FormatUint(fallbackID, 10)
	}
	return slug
}

// ExtractSeoSlugSegment pulls the trailing slug from compound "id/slug" or a path.
func ExtractSeoSlugSegment(raw string) string {
	raw = strings.Trim(strings.TrimSpace(raw), "/")
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, "/")
	seg := parts[len(parts)-1]
	slug := slugifySEO(seg)
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.Trim(slug, "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	if slug == "" || slug == "seo-page" {
		return ""
	}
	return slug
}

// IsManualSEORefType reports manual / legacy SEO pages (path uses seo_domain.id).
func IsManualSEORefType(refType SEORefType) bool {
	return refType == 0 || refType == SEORefTypeLegacy
}

func slugifySEO(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	s, _, err := transform.String(t, value)
	if err != nil {
		s = value
	}
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_' || r == '.':
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}
