package domain

// ClickURL lưu lại sự kiện click link trên hệ thống.
type ClickURL struct {
	EventTime  int64   `gorm:"not null;index;column:event_time;comment:Thời gian click (timestamp)"`
	DomainID   uint64  `gorm:"not null;index;column:domain_id;comment:ID domain"`
	DomainType uint32  `gorm:"not null;index;column:domain_type;comment:Loại domain"`
	ProfileID  *uint64 `gorm:"index;column:profile_id;comment:ID người dùng nếu đã đăng nhập"`
}

func (ClickURL) TableName() string {
	return "click_urls"
}
