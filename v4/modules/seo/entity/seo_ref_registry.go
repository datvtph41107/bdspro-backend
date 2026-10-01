package entity

import "strings"

// Source status constants for SEO ref registry.
const (
	SeoSourceStatusReady      = "ready"
	SeoSourceStatusComingSoon = "coming_soon"
	SeoSourceStatusDisabled   = "disabled"

	SeoRefSourceManual                = "manual"
	SeoRefSourceBdsproRegion          = "bdspro.region"
	SeoRefSourceBdsproAreaRegion      = "bdspro.area_region"
	SeoRefSourceBdsproProject         = "bdspro.project"
	SeoRefSourceTqdAdministrativeUnit = "tqd.administrative_unit"
	SeoRefSourceTqdParcel             = "tqd.parcel"
	SeoRefSourceTqdPlanningProject    = "tqd.planning_project"
	SeoRefSourceTqdLayer              = "tqd.qh_layer"
	SeoRefSourceTqdDocument           = "tqd.legal_document"
	SeoRefSourceCrmReport             = "crm.report"
)

// SeoSourceCapabilities describes operations a source supports.
type SeoSourceCapabilities struct {
	Searchable       bool `json:"searchable"`
	Resolvable       bool `json:"resolvable"`
	UniquePerRef     bool `json:"uniquePerRef"`
	SupportsTemplate bool `json:"supportsTemplate"`
	SupportsPreview  bool `json:"supportsPreview"`
	SupportsGenerate bool `json:"supportsGenerate"`
	SupportsHealth   bool `json:"supportsHealth"`
}

// SeoRefTypeConfig is the local SEO source registry entry (CRM path parity).
type SeoRefTypeConfig struct {
	Value            uint32                `json:"value"`
	Key              string                `json:"key"`
	Label            string                `json:"label"`
	Description      string                `json:"description"`
	Family           string                `json:"family"`
	SourceService    string                `json:"sourceService"`
	ResolverKey      string                `json:"resolverKey"`
	Searchable       bool                  `json:"searchable"`
	UniquePerRef     bool                  `json:"uniquePerRef"`
	SupportsTemplate bool                  `json:"supportsTemplate"`
	DefaultPageType  string                `json:"defaultPageType"`
	AdminURLPattern  string                `json:"adminUrlPattern"`
	PublicURLPattern string                `json:"publicUrlPattern"`
	Status           string                `json:"status"`
	Resolvable       bool                  `json:"resolvable"`
	DisabledReason   string                `json:"disabledReason,omitempty"`
	Capabilities     SeoSourceCapabilities `json:"capabilities"`
}

func readyCapabilities(searchable, resolvable, uniquePerRef, supportsTemplate, supportsPreview, supportsGenerate, supportsHealth bool) SeoSourceCapabilities {
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
	return SeoSourceCapabilities{UniquePerRef: uniquePerRef}
}

// NormalizeSeoRefTypeConfig syncs legacy fields with Capabilities.
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

// SeoRefTypeConfigs returns the configured SEO source catalog.
func SeoRefTypeConfigs() []SeoRefTypeConfig {
	items := []SeoRefTypeConfig{
		{
			Value: uint32(SEORefTypeLegacy), Key: "manual", Label: "Nội dung thủ công",
			Description: "Landing, hub hoặc guide không gắn dữ liệu nguồn", Family: "manual",
			SourceService: "crm", ResolverKey: SeoRefSourceManual,
			DefaultPageType: "manual", PublicURLPattern: "/tin-tuc/tin/{id}/{slug}",
			Status: SeoSourceStatusReady,
			Capabilities: readyCapabilities(false, false, false, true, true, true, false),
		},
		{
			Value: uint32(SEORefTypeAdmUnit), Key: "adm_unit", Label: "Đơn vị hành chính",
			Description: "Tỉnh, xã/phường hoặc đơn vị hành chính từ TQD projection", Family: "location",
			SourceService: "tqd", ResolverKey: SeoRefSourceTqdAdministrativeUnit,
			UniquePerRef: true, DefaultPageType: "entity",
			PublicURLPattern: "/tin-tuc/dia-ban/{id}/{slug}",
			Status: SeoSourceStatusComingSoon,
			DisabledReason: "Generic numeric refId cannot preserve typed province/ward identity",
			Capabilities: comingSoonCapabilities(true),
		},
		{
			Value: uint32(SEORefTypeRegion), Key: "region", Label: "Khu vực",
			Description: "Khu vực/area-region phục vụ SEO địa phương", Family: "location",
			SourceService: "bdspro", ResolverKey: SeoRefSourceBdsproAreaRegion,
			DefaultPageType: "entity", AdminURLPattern: "/admin/area-regions/{id}",
			PublicURLPattern: "/tin-tuc/vung-quy-hoach/{id}/{slug}",
			Status: SeoSourceStatusReady,
			Capabilities: readyCapabilities(true, true, true, true, true, true, true),
		},
		{
			Value: uint32(SEORefTypeParcel), Key: "parcel", Label: "Thửa đất",
			Description: "Thửa đất hoặc parcel từ TQD", Family: "planning",
			SourceService: "tqd", ResolverKey: SeoRefSourceTqdParcel,
			DefaultPageType: "entity", AdminURLPattern: "/admin/parcels/{id}",
			PublicURLPattern: "/tin-tuc/thua-dat/{id}/{slug}",
			Status: SeoSourceStatusReady,
			Capabilities: readyCapabilities(true, true, true, true, true, true, true),
		},
		{
			Value: uint32(SEORefTypeProject), Key: "project", Label: "Dự án",
			Description: "Dự án bất động sản", Family: "real_estate",
			SourceService: "bdspro", ResolverKey: SeoRefSourceBdsproProject,
			DefaultPageType: "entity", AdminURLPattern: "/admin/project/{id}",
			PublicURLPattern: "/du-an/{id}/{slug}",
			Status: SeoSourceStatusReady,
			Capabilities: readyCapabilities(true, true, true, true, true, true, true),
		},
		{
			Value: uint32(SEORefTypeProject), Key: "planning_project", Label: "Đồ án quy hoạch",
			Description: "Đồ án quy hoạch từ TQD Planning Project projection", Family: "planning",
			SourceService: "tqd", ResolverKey: SeoRefSourceTqdPlanningProject,
			DefaultPageType: "entity", PublicURLPattern: "/tin-tuc/do-an-quy-hoach/{id}/{slug}",
			Status: SeoSourceStatusReady,
			Capabilities: readyCapabilities(false, true, true, true, true, true, true),
		},
		{
			Value: uint32(SEORefTypeDocument), Key: "document", Label: "Tài liệu",
			Description: "Tài liệu pháp lý, quy hoạch hoặc hồ sơ liên quan", Family: "document",
			SourceService: "tqd", ResolverKey: SeoRefSourceTqdDocument,
			UniquePerRef: true, DefaultPageType: "entity",
			PublicURLPattern: "/tin-tuc/tai-lieu/{id}/{slug}",
			Status: SeoSourceStatusComingSoon,
			DisabledReason: "Resolver TQD legal document chưa implement",
			Capabilities: comingSoonCapabilities(true),
		},
		{
			Value: uint32(SEORefTypeMap), Key: "map", Label: "Bản đồ / lớp quy hoạch",
			Description: "Lớp bản đồ, QH layer hoặc workspace", Family: "map",
			SourceService: "tqd", ResolverKey: SeoRefSourceTqdLayer,
			UniquePerRef: true, DefaultPageType: "entity",
			PublicURLPattern: "/tin-tuc/lop/{id}/{slug}",
			Status: SeoSourceStatusComingSoon,
			DisabledReason: "Resolver TQD QH layer chưa implement",
			Capabilities: comingSoonCapabilities(true),
		},
		{
			Value: uint32(SEORefTypeReport), Key: "report", Label: "Báo cáo",
			Description: "Báo cáo, phân tích hoặc thống kê", Family: "report",
			SourceService: "crm", ResolverKey: SeoRefSourceCrmReport,
			UniquePerRef: true, DefaultPageType: "entity",
			PublicURLPattern: "/tin-tuc/bao-cao/{id}/{slug}",
			Status: SeoSourceStatusComingSoon,
			DisabledReason: "Resolver CRM report chưa implement",
			Capabilities: comingSoonCapabilities(true),
		},
	}
	for i := range items {
		items[i] = NormalizeSeoRefTypeConfig(items[i])
	}
	return items
}

// GetSeoRefTypeConfig finds config by legacy refType (+ optional resolver/service filters).
func GetSeoRefTypeConfig(value uint32, resolverKey, sourceService string) (SeoRefTypeConfig, bool) {
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

// GetSeoRefTypeConfigByKey finds config by sourceKey.
func GetSeoRefTypeConfigByKey(sourceKey string) (SeoRefTypeConfig, bool) {
	sourceKey = strings.TrimSpace(sourceKey)
	if sourceKey == "legacy" {
		sourceKey = "manual"
	}
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

// IsSeoRefTypeSearchable reports whether search is allowed for the source.
func IsSeoRefTypeSearchable(cfg SeoRefTypeConfig) bool {
	cfg = NormalizeSeoRefTypeConfig(cfg)
	return cfg.Status == SeoSourceStatusReady && cfg.Capabilities.Searchable
}

// IsSeoRefTypeResolvable reports whether resolve is allowed for the source.
func IsSeoRefTypeResolvable(cfg SeoRefTypeConfig) bool {
	cfg = NormalizeSeoRefTypeConfig(cfg)
	return cfg.Status == SeoSourceStatusReady && cfg.Capabilities.Resolvable
}
