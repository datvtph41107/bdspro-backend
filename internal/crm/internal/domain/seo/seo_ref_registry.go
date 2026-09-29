package seo_domain

import (
	"strings"

	"crm/internal/enums"
)

const (
	// SeoRefSourceManual là nguồn nội dung nhập trực tiếp trong CRM.
	SeoRefSourceManual = "manual"

	SeoRefSourceBdsproRegion     = "bdspro.region"
	SeoRefSourceBdsproAreaRegion = "bdspro.area_region"
	SeoRefSourceBdsproProject    = "bdspro.project"
	SeoRefSourceBdsproDocType    = "bdspro.doc_type"

	SeoRefSourceTqdAdministrativeUnit = "tqd.administrative_unit"
	SeoRefSourceTqdParcel             = "tqd.parcel"
	SeoRefSourceTqdPlanningProject    = "tqd.planning_project"
	SeoRefSourceTqdLayer              = "tqd.qh_layer"
	SeoRefSourceTqdDocument           = "tqd.legal_document"
	SeoRefSourceTqdReport             = "tqd.report"

	SeoRefSourceCrmReport = "crm.report"
)

// Các trạng thái cho biết nguồn có sẵn để sử dụng hay chưa.
const (
	SeoSourceStatusReady      = "ready"
	SeoSourceStatusComingSoon = "coming_soon"
	SeoSourceStatusDisabled   = "disabled"
)

// Các trạng thái kết nối kỹ thuật tới nguồn dữ liệu.
const (
	SeoSourceHealthHealthy        = "healthy"
	SeoSourceHealthDegraded       = "degraded"
	SeoSourceHealthUnavailable    = "unavailable"
	SeoSourceHealthNotImplemented = "not_implemented"
)

// SeoSourceCapabilities mô tả các thao tác kỹ thuật mà một nguồn hỗ trợ.
type SeoSourceCapabilities struct {
	// Searchable cho phép tìm bản ghi nguồn.
	Searchable bool `json:"searchable"`
	// Resolvable cho phép đọc chi tiết một bản ghi nguồn.
	Resolvable bool `json:"resolvable"`
	// UniquePerRef chỉ cho phép một trang SEO trên mỗi bản ghi nguồn.
	UniquePerRef bool `json:"uniquePerRef"`
	// SupportsTemplate cho phép tạo nội dung từ mẫu.
	SupportsTemplate bool `json:"supportsTemplate"`
	// SupportsPreview cho phép xem trước HTML.
	SupportsPreview bool `json:"supportsPreview"`
	// SupportsGenerate cho phép tạo HTML tĩnh.
	SupportsGenerate bool `json:"supportsGenerate"`
	// SupportsHealth giữ thông tin tương thích về kiểm tra kết nối nguồn.
	SupportsHealth bool `json:"supportsHealth"`
}

// SeoRefTypeConfig là cấu hình kết nối một loại dữ liệu nguồn với trang SEO.
type SeoRefTypeConfig struct {
	// Value là mã số loại nguồn để tương thích dữ liệu cũ.
	Value uint32 `json:"value"`
	// Key là tên ngắn của loại dữ liệu.
	Key string `json:"key"`
	// Label là tên hiển thị trong Admin.
	Label string `json:"label"`
	// Description mô tả ngắn dữ liệu nguồn.
	Description string `json:"description"`
	// Family là nhóm dữ liệu nguồn.
	Family string `json:"family"`
	// SourceService là dịch vụ đang sở hữu dữ liệu.
	SourceService string `json:"sourceService"`
	// ResolverKey chọn bộ đọc dữ liệu nguồn.
	ResolverKey string `json:"resolverKey"`
	// Searchable là trường tương thích cho khả năng tìm kiếm.
	Searchable bool `json:"searchable"`
	// UniquePerRef là trường tương thích cho giới hạn một trang trên mỗi nguồn.
	UniquePerRef bool `json:"uniquePerRef"`
	// SupportsTemplate là trường tương thích cho khả năng dùng mẫu.
	SupportsTemplate bool `json:"supportsTemplate"`
	// DefaultPageType là loại trang mặc định được tạo.
	DefaultPageType string `json:"defaultPageType"`
	// AdminURLPattern là mẫu đường dẫn quản trị của bản ghi nguồn.
	AdminURLPattern string `json:"adminUrlPattern"`
	// PublicURLPattern là mẫu đường dẫn công khai của trang.
	PublicURLPattern string `json:"publicUrlPattern"`

	// Status là trạng thái sử dụng hiện tại của nguồn.
	Status string `json:"status"`
	// Resolvable cho biết có thể đọc chi tiết bản ghi nguồn.
	Resolvable bool `json:"resolvable"`
	// DisabledReason giải thích ngắn khi nguồn chưa dùng được.
	DisabledReason string `json:"disabledReason"`
	// Capabilities gom các thao tác kỹ thuật nguồn hỗ trợ.
	Capabilities SeoSourceCapabilities `json:"capabilities"`
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

// NormalizeSeoRefTypeConfig đồng bộ các trường tương thích với nhóm Capabilities.
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

// SeoRefTypeConfigs trả về danh sách nguồn dữ liệu SEO đã cấu hình.
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
				false, // không tìm kiếm
				false, // không đọc từ nguồn
				false, // không giới hạn theo nguồn
				true,  // hỗ trợ mẫu
				true,  // hỗ trợ xem trước
				true,  // hỗ trợ tạo HTML
				false, // không kiểm tra kết nối
			),
		},
		{
			Value:            uint32(enums.ESEORefTypeAdmUnit),
			Key:              "adm_unit",
			Label:            "Đơn vị hành chính",
			Description:      "Tỉnh, xã/phường hoặc đơn vị hành chính từ TQD projection",
			Family:           "location",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdAdministrativeUnit,
			Searchable:       false,
			UniquePerRef:     true,
			SupportsTemplate: false,
			DefaultPageType:  "entity",
			PublicURLPattern: "/dia-ban/{slug}",
			Status:           SeoSourceStatusComingSoon,
			Resolvable:       false,
			DisabledReason:   "Generic numeric refId cannot preserve typed province/ward identity; dedicated canonical provisioning is required",
			Capabilities:     comingSoonCapabilities(true),
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
			Value:            uint32(enums.ESEORefTypeProject),
			Key:              "planning_project",
			Label:            "Đồ án quy hoạch",
			Description:      "Đồ án quy hoạch từ TQD Planning Project projection",
			Family:           "planning",
			SourceService:    "tqd",
			ResolverKey:      SeoRefSourceTqdPlanningProject,
			Searchable:       false,
			UniquePerRef:     true,
			SupportsTemplate: true,
			DefaultPageType:  "entity",
			PublicURLPattern: "/do-an-quy-hoach/{slug}",
			Status:           SeoSourceStatusReady,
			Resolvable:       true,
			Capabilities: readyCapabilities(
				false, // không tìm kiếm
				true,  // đọc được từ nguồn
				true,  // một trang trên mỗi nguồn
				true,  // hỗ trợ mẫu
				true,  // hỗ trợ xem trước
				true,  // hỗ trợ tạo HTML
				true,  // hỗ trợ kiểm tra kết nối
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

// GetSeoRefTypeConfig tìm cấu hình theo mã loại và thông tin kết nối cũ.
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

// GetSeoRefTypeConfigByKey tìm cấu hình nguồn theo một khóa duy nhất.
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

// GetSeoRefTypeConfigByResolverKey tìm cấu hình theo khóa bộ đọc dữ liệu.
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

// IsSeoRefTypeReady kiểm tra nguồn đã sẵn sàng cho Admin hay chưa.
func IsSeoRefTypeReady(cfg SeoRefTypeConfig) bool {
	return cfg.Status == SeoSourceStatusReady
}

// IsSeoRefTypeSearchable kiểm tra API có thể tìm kiếm nguồn hay không.
func IsSeoRefTypeSearchable(cfg SeoRefTypeConfig) bool {
	cfg = NormalizeSeoRefTypeConfig(cfg)
	return cfg.Status == SeoSourceStatusReady && cfg.Capabilities.Searchable
}

// IsSeoRefTypeResolvable kiểm tra API có thể đọc chi tiết nguồn hay không.
func IsSeoRefTypeResolvable(cfg SeoRefTypeConfig) bool {
	cfg = NormalizeSeoRefTypeConfig(cfg)
	return cfg.Status == SeoSourceStatusReady && cfg.Capabilities.Resolvable
}

// IsSeoResolverReady kiểm tra bộ đọc dữ liệu nguồn đã được triển khai.
func IsSeoResolverReady(resolverKey string) bool {
	switch strings.TrimSpace(resolverKey) {
	case SeoRefSourceBdsproProject,
		SeoRefSourceBdsproRegion,
		SeoRefSourceBdsproAreaRegion,
		SeoRefSourceTqdParcel,
		SeoRefSourceTqdPlanningProject:
		return true
	default:
		return false
	}
}

// SeoRefTypeSourceKeys trả về toàn bộ khóa nguồn đã cấu hình.
func SeoRefTypeSourceKeys() []string {
	items := SeoRefTypeConfigs()
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.Key)
	}
	return out
}
