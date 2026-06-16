package seo_domain

import "crm/internal/enums"

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
}

func SeoRefTypeConfigs() []SeoRefTypeConfig {
	return []SeoRefTypeConfig{
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
			AdminURLPattern:  "/admin/project/{id}",
			PublicURLPattern: "/du-an/{slug}",
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
		},
		{
			Value:            uint32(enums.ESEORefTypeDocument),
			Key:              "document",
			Label:            "Tài liệu",
			Description:      "Tài liệu pháp lý, quy hoạch hoặc hồ sơ liên quan",
			Family:           "document",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdDocument,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
		},
		{
			Value:            uint32(enums.ESEORefTypeMap),
			Key:              "map",
			Label:            "Bản đồ / lớp quy hoạch",
			Description:      "Lớp bản đồ, QH layer hoặc workspace",
			Family:           "map",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdLayer,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
		},
		{
			Value:            uint32(enums.ESEORefTypeReport),
			Key:              "report",
			Label:            "Báo cáo",
			Description:      "Báo cáo, phân tích hoặc thống kê",
			Family:           "report",
			SourceService:    "crm",
			ResolverKey:      SeoRefSourceCrmReport,
			Searchable:       true,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
		},
	}
}

func GetSeoRefTypeConfig(value uint32, resolverKey string, sourceService string) (SeoRefTypeConfig, bool) {
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
