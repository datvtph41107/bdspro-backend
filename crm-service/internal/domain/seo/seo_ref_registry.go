package seo_domain

import (
	"strings"

	"crm/internal/enums"
)

const (
	SeoRefSourceManual = "manual"

	SeoRefSourceBdsproRegion     = "bdspro.region"
	SeoRefSourceBdsproAreaRegion = "bdspro.area_region"
	SeoRefSourceBdsproProject    = "bdspro.project"
	SeoRefSourceBdsproDocType    = "bdspro.doc_type"

	SeoRefSourceTqdParcel   = "tqd.parcel"
	SeoRefSourceTqdLayer    = "tqd.qh_layer"
	SeoRefSourceTqdDocument = "tqd.legal_document"
	SeoRefSourceTqdReport   = "tqd.report"

	SeoRefSourceCrmReport = "crm.report"
)

const (
	SeoSourceStatusReady      = "ready"
	SeoSourceStatusComingSoon = "coming_soon"
	SeoSourceStatusDisabled   = "disabled"
)

const (
	SeoSourceHealthHealthy        = "healthy"
	SeoSourceHealthDegraded       = "degraded"
	SeoSourceHealthUnavailable    = "unavailable"
	SeoSourceHealthNotImplemented = "not_implemented"
)

// SeoSourceCapabilities is the source-of-truth capability contract for FE/BE.
// Legacy flat fields such as Searchable/UniquePerRef/SupportsTemplate are kept
// in SeoRefTypeConfig for backward compatibility with existing handlers/mappers.
type SeoSourceCapabilities struct {
	Searchable       bool `json:"searchable"`
	Resolvable       bool `json:"resolvable"`
	UniquePerRef     bool `json:"uniquePerRef"`
	SupportsTemplate bool `json:"supportsTemplate"`
	SupportsPreview  bool `json:"supportsPreview"`
	SupportsGenerate bool `json:"supportsGenerate"`
	SupportsHealth   bool `json:"supportsHealth"`
}

// SeoRefTypeConfig is the SEO source registry item.
//
// Compatibility contract:
//   - Existing code can keep using Searchable, UniquePerRef, SupportsTemplate.
//   - New code should prefer Status, Resolvable, DisabledReason and Capabilities.
//   - Source registry is the only place deciding whether a source is ready.
//   - Frontend must not hard-code source readiness.
type SeoRefTypeConfig struct {
	Value            uint32 `json:"value"`
	Key              string `json:"key"`
	Label            string `json:"label"`
	Description      string `json:"description"`
	Family           string `json:"family"`
	SourceService    string `json:"sourceService"`
	ResolverKey      string `json:"resolverKey"`
	Searchable       bool   `json:"searchable"`
	UniquePerRef     bool   `json:"uniquePerRef"`
	SupportsTemplate bool   `json:"supportsTemplate"`
	DefaultPageType  string `json:"defaultPageType"`
	AdminURLPattern  string `json:"adminUrlPattern"`
	PublicURLPattern string `json:"publicUrlPattern"`

	// V2 source-of-truth fields.
	Status         string                `json:"status"`
	Resolvable     bool                  `json:"resolvable"`
	DisabledReason string                `json:"disabledReason"`
	Capabilities   SeoSourceCapabilities `json:"capabilities"`
}

func readyCapabilities(searchable bool, resolvable bool, uniquePerRef bool, supportsTemplate bool, supportsPreview bool, supportsGenerate bool, supportsHealth bool) SeoSourceCapabilities {
	return SeoSourceCapabilities{
		Searchable:       searchable,
		Resolvable:       resolvable,
		UniquePerRef:     uniquePerRef,
		SupportsTemplate: supportsTemplate,
		SupportsPreview:  supportsPreview,
		SupportsGenerate: supportsGenerate,
		SupportsHealth:   supportsHealth,
	}
}

func comingSoonCapabilities(uniquePerRef bool) SeoSourceCapabilities {
	return SeoSourceCapabilities{
		Searchable:       false,
		Resolvable:       false,
		UniquePerRef:     uniquePerRef,
		SupportsTemplate: false,
		SupportsPreview:  false,
		SupportsGenerate: false,
		SupportsHealth:   false,
	}
}

func disabledCapabilities(uniquePerRef bool) SeoSourceCapabilities {
	return SeoSourceCapabilities{
		Searchable:       false,
		Resolvable:       false,
		UniquePerRef:     uniquePerRef,
		SupportsTemplate: false,
		SupportsPreview:  false,
		SupportsGenerate: false,
		SupportsHealth:   false,
	}
}

// NormalizeSeoRefTypeConfig keeps legacy flat fields and V2 Capabilities in sync.
// This protects old code while new code migrates to capability-based decisions.
func NormalizeSeoRefTypeConfig(cfg SeoRefTypeConfig) SeoRefTypeConfig {
	if cfg.Status == "" {
		cfg.Status = SeoSourceStatusReady
	}

	if cfg.Capabilities == (SeoSourceCapabilities{}) {
		cfg.Capabilities = SeoSourceCapabilities{
			Searchable:       cfg.Searchable,
			Resolvable:       cfg.Resolvable,
			UniquePerRef:     cfg.UniquePerRef,
			SupportsTemplate: cfg.SupportsTemplate,
			SupportsPreview:  cfg.SupportsTemplate,
			SupportsGenerate: cfg.SupportsTemplate,
			SupportsHealth:   cfg.Searchable || cfg.Resolvable,
		}
	}

	cfg.Searchable = cfg.Capabilities.Searchable
	cfg.Resolvable = cfg.Capabilities.Resolvable
	cfg.UniquePerRef = cfg.Capabilities.UniquePerRef
	cfg.SupportsTemplate = cfg.Capabilities.SupportsTemplate

	if cfg.DefaultPageType == "" {
		if cfg.Key == "manual" || cfg.ResolverKey == SeoRefSourceManual {
			cfg.DefaultPageType = "manual"
		} else {
			cfg.DefaultPageType = "entity"
		}
	}

	return cfg
}

func SeoRefTypeConfigs() []SeoRefTypeConfig {
	items := []SeoRefTypeConfig{
		{
			Value:            0,
			Key:              "manual",
			Label:            "Trang thủ công",
			Description:      "Landing, hub hoặc guide không gắn dữ liệu nguồn",
			Family:           "manual",
			SourceService:    "crm",
			ResolverKey:      SeoRefSourceManual,
			Searchable:       false,
			UniquePerRef:     false,
			SupportsTemplate: true,
			DefaultPageType:  "manual",
			Status:           SeoSourceStatusReady,
			Resolvable:       false,
			Capabilities: readyCapabilities(
				false, // searchable
				false, // resolvable
				false, // uniquePerRef
				true,  // supportsTemplate
				true,  // supportsPreview
				true,  // supportsGenerate
				false, // supportsHealth
			),
		},
		{
			Value:            uint32(enums.ESEORefTypeAdmUnit),
			Key:              "adm_unit",
			Label:            "Đơn vị hành chính",
			Description:      "Tỉnh, huyện, xã hoặc đơn vị hành chính",
			Family:           "location",
			SourceService:    "bdspro",
			ResolverKey:      SeoRefSourceBdsproRegion,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
			AdminURLPattern:  "/admin/region/{id}",
			PublicURLPattern: "/dia-diem/{slug}",
			Status:           SeoSourceStatusReady,
			Resolvable:       true,
			Capabilities: readyCapabilities(
				true,
				true,
				true,
				true,
				true,
				true,
				true,
			),
		},
		{
			Value:            uint32(enums.ESEORefTypeRegion),
			Key:              "region",
			Label:            "Khu vực",
			Description:      "Khu vực/area-region phục vụ SEO địa phương",
			Family:           "location",
			SourceService:    "bdspro",
			ResolverKey:      SeoRefSourceBdsproAreaRegion,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
			AdminURLPattern:  "/admin/area-regions/{id}",
			PublicURLPattern: "/khu-vuc/{slug}",
			Status:           SeoSourceStatusReady,
			Resolvable:       true,
			Capabilities: readyCapabilities(
				true,
				true,
				true,
				true,
				true,
				true,
				true,
			),
		},
		{
			Value:            uint32(enums.ESEORefTypeParcel),
			Key:              "parcel",
			Label:            "Thửa đất",
			Description:      "Thửa đất hoặc parcel từ TQD",
			Family:           "planning",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdParcel,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
			AdminURLPattern:  "/admin/parcels/{id}",
			PublicURLPattern: "/quy-hoach/thua-dat/{id}-{slug}",
			Status:           SeoSourceStatusReady,
			Resolvable:       true,
			Capabilities: readyCapabilities(
				true,
				true,
				true,
				true,
				true,
				true,
				true,
			),
		},
		{
			Value:            uint32(enums.ESEORefTypeProject),
			Key:              "project",
			Label:            "Dự án",
			Description:      "Dự án bất động sản",
			Family:           "real_estate",
			SourceService:    "bdspro",
			ResolverKey:      SeoRefSourceBdsproProject,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
			AdminURLPattern:  "/admin/project/{id}",
			PublicURLPattern: "/du-an/{slug}",
			Status:           SeoSourceStatusReady,
			Resolvable:       true,
			Capabilities: readyCapabilities(
				true,
				true,
				true,
				true,
				true,
				true,
				true,
			),
		},
		{
			Value:            uint32(enums.ESEORefTypeDocument),
			Key:              "document",
			Label:            "Tài liệu",
			Description:      "Tài liệu pháp lý, quy hoạch hoặc hồ sơ liên quan",
			Family:           "document",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdDocument,
			Searchable:       false,
			UniquePerRef:     true,
			SupportsTemplate: false,
			DefaultPageType:  "entity",
			Status:           SeoSourceStatusComingSoon,
			Resolvable:       false,
			DisabledReason:   "Resolver TQD legal document chưa implement",
			Capabilities:     comingSoonCapabilities(true),
		},
		{
			Value:            uint32(enums.ESEORefTypeMap),
			Key:              "map",
			Label:            "Bản đồ / lớp quy hoạch",
			Description:      "Lớp bản đồ, QH layer hoặc workspace",
			Family:           "map",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdLayer,
			Searchable:       false,
			UniquePerRef:     true,
			SupportsTemplate: false,
			DefaultPageType:  "entity",
			Status:           SeoSourceStatusComingSoon,
			Resolvable:       false,
			DisabledReason:   "Resolver TQD QH layer chưa implement",
			Capabilities:     comingSoonCapabilities(true),
		},
		{
			Value:            uint32(enums.ESEORefTypeReport),
			Key:              "report",
			Label:            "Báo cáo",
			Description:      "Báo cáo, phân tích hoặc thống kê",
			Family:           "report",
			SourceService:    "crm",
			ResolverKey:      SeoRefSourceCrmReport,
			Searchable:       false,
			UniquePerRef:     true,
			SupportsTemplate: false,
			DefaultPageType:  "entity",
			Status:           SeoSourceStatusComingSoon,
			Resolvable:       false,
			DisabledReason:   "Resolver CRM report chưa implement",
			Capabilities:     comingSoonCapabilities(true),
		},
	}

	for i := range items {
		items[i] = NormalizeSeoRefTypeConfig(items[i])
	}

	return items
}

// GetSeoRefTypeConfig keeps the legacy lookup contract.
// Existing code passes refType + optional resolverKey/sourceService.
func GetSeoRefTypeConfig(value uint32, resolverKey string, sourceService string) (SeoRefTypeConfig, bool) {
	resolverKey = strings.TrimSpace(resolverKey)
	sourceService = strings.TrimSpace(sourceService)

	for _, item := range SeoRefTypeConfigs() {
		if item.Value != value {
			continue
		}
		if resolverKey != "" && item.ResolverKey != resolverKey {
			continue
		}
		if sourceService != "" && item.SourceService != sourceService {
			continue
		}
		return item, true
	}

	return SeoRefTypeConfig{}, false
}

// GetSeoRefTypeConfigByKey is the V2 lookup path.
// New APIs should prefer sourceKey over refType+resolverKey+sourceService.
func GetSeoRefTypeConfigByKey(sourceKey string) (SeoRefTypeConfig, bool) {
	sourceKey = strings.TrimSpace(sourceKey)
	if sourceKey == "" {
		return SeoRefTypeConfig{}, false
	}

	for _, item := range SeoRefTypeConfigs() {
		if item.Key == sourceKey {
			return item, true
		}
	}

	return SeoRefTypeConfig{}, false
}

func GetSeoRefTypeConfigByResolverKey(resolverKey string) (SeoRefTypeConfig, bool) {
	resolverKey = strings.TrimSpace(resolverKey)
	if resolverKey == "" {
		return SeoRefTypeConfig{}, false
	}

	for _, item := range SeoRefTypeConfigs() {
		if item.ResolverKey == resolverKey {
			return item, true
		}
	}

	return SeoRefTypeConfig{}, false
}

// IsSeoRefTypeReady tells whether a source can be used by admin UI.
// Manual is ready but not searchable/resolvable.
func IsSeoRefTypeReady(cfg SeoRefTypeConfig) bool {
	return cfg.Status == SeoSourceStatusReady
}

// IsSeoRefTypeSearchable tells whether source search API should call adapter/resolver.
func IsSeoRefTypeSearchable(cfg SeoRefTypeConfig) bool {
	cfg = NormalizeSeoRefTypeConfig(cfg)
	return cfg.Status == SeoSourceStatusReady && cfg.Capabilities.Searchable
}

// IsSeoRefTypeResolvable tells whether source resolve API should call adapter/resolver.
func IsSeoRefTypeResolvable(cfg SeoRefTypeConfig) bool {
	cfg = NormalizeSeoRefTypeConfig(cfg)
	return cfg.Status == SeoSourceStatusReady && cfg.Capabilities.Resolvable
}

// IsSeoResolverReady is useful while the codebase still uses the old resolver switch.
// Once adapters fully replace the switch, this helper can be removed.
func IsSeoResolverReady(resolverKey string) bool {
	switch strings.TrimSpace(resolverKey) {
	case SeoRefSourceBdsproProject,
		SeoRefSourceBdsproRegion,
		SeoRefSourceBdsproAreaRegion,
		SeoRefSourceTqdParcel:
		return true
	default:
		return false
	}
}

// SeoRefTypeSourceKeys returns all configured source keys.
// Useful for tests and frontend docs.
func SeoRefTypeSourceKeys() []string {
	items := SeoRefTypeConfigs()
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.Key)
	}
	return out
}
