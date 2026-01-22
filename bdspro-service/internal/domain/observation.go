package domain

import (
	_models "common/domain/entity"
)

// Observation lưu lại mọi hành vi người dùng trên hệ thống
// Dùng chung cho: product, asset, post, blog, project, agency...
type Observation struct {
	_models.BaseEntity
	// ID của đối tượng
	EntityID uint64 `gorm:"not null;index;comment:ID của entity được theo dõi"`
	// Loại đối tượng: product | asset | post | blog | project | agency | profile
	EntityType *string `gorm:"type:varchar(32);not null;index;comment:Loại entity được theo dõi"`
	// ID user trong hệ thống, null nếu là guest
	ProfileID *uint64 `gorm:"index;comment:ID người dùng nếu đã đăng nhập"`
	// ID thiết bị/trình duyệt, dùng để phân biệt unique viewer
	ClientID *string `gorm:"type:varchar(64);index;comment:ID thiết bị hoặc browser (cookie/localStorage)"`
	// ID phiên truy cập
	SessionID *string `gorm:"type:varchar(64);index;comment:ID session để gom các hành vi cùng phiên"`
	// Vai trò user: buyer | seller | agent | admin
	Role *string `gorm:"type:varchar(32);comment:Vai trò người dùng trong hệ thống"`
	// ID tổ chức/sàn (nếu user thuộc sàn)
	OrgID *uint64 `gorm:"index;comment:ID sàn hoặc tổ chức của user"`
	// Loại hành vi: view | save | click_phone | submit_lead | share...
	Action string `gorm:"type:varchar(32);not null;index;comment:Loại hành vi người dùng"`
	// Thời gian thực hiện hành vi (giây), chủ yếu dùng cho action=view
	Duration uint32 `gorm:"default:0;comment:Thời gian thực hiện hành vi (giây)"`
	// URL hiện tại
	PageURL *string `gorm:"type:varchar(255);comment:URL trang user đang xem"`
	// Trang nguồn (từ đâu tới)
	Referrer *string `gorm:"type:varchar(255);comment:Trang nguồn (referrer)"`
	// Khu vực địa lý (ví dụ: HCM, HN)
	Geo *string `gorm:"type:varchar(64);comment:Vị trí địa lý người dùng"`
	// Thiết bị: mobile | desktop | tablet
	Device *string `gorm:"type:varchar(32);comment:Loại thiết bị"`
	// User-Agent trình duyệt
	UserAgent *string `gorm:"type:varchar(255);comment:User agent trình duyệt"`
	// Có phải bot không
	// IsBot bool `gorm:"default:false;comment:Đánh dấu bot/crawler"`
	// // Giá trị số (tuỳ action), ví dụ: scroll depth, score...
	// Value float64 `gorm:"default:0;comment:Giá trị bổ sung cho hành vi (tuỳ use case)"`
	// // ID dự án (nếu asset thuộc dự án)
	// ProjectID *uint64 `gorm:"index;comment:ID dự án bất động sản"`
	// // ID sàn đăng tin
	// AgencyID *uint64 `gorm:"index;comment:ID sàn bất động sản"`
	// // Khoảng giá user quan tâm
	// PriceMin *uint64 `gorm:"comment:Giá thấp nhất user quan tâm"`
	// PriceMax *uint64 `gorm:"comment:Giá cao nhất user quan tâm"`
	// // Khoảng diện tích
	// AreaMin *uint64 `gorm:"comment:Diện tích nhỏ nhất user quan tâm"`
	// AreaMax *uint64 `gorm:"comment:Diện tích lớn nhất user quan tâm"`
	// // Loại tài sản: apartment | land | house | office
	// PropertyType string `gorm:"type:varchar(32);comment:Loại hình bất động sản"`
}

func (Observation) TableName() string {
	return "observations"
}
