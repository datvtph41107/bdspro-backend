package domain

// InteractiveEvent lưu lại các sự kiện tương tác của người dùng trên hệ thống
type InteractiveEvent struct {
	StartTime int64   `gorm:"not null;index;comment:Thời gian bắt đầu (timestamp)"` // Thời gian bắt đầu (timestamp)
	EndTime   *int64  `gorm:"index;comment:Thời gian kết thúc (timestamp)"`         // Thời gian kết thúc (timestamp, optional)
	Event     string  `gorm:"type:varchar(30);not null;index;comment:Tên sự kiện"`  // Tên sự kiện
	RefID     uint64  `gorm:"index;comment:ID tham chiếu"`                          // ID tham chiếu (ref_id)
	Screen    *string `gorm:"type:varchar(40);comment:Màn hình/trang"`              // Màn hình/trang (optional)
	Duration  *int64  `gorm:"comment:Thời lượng (giây)"`                            // Thời lượng (giây, optional)
	ProfileID *uint64 `gorm:"index;comment:ID người dùng nếu đã đăng nhập"`         // ID user trong hệ thống, null nếu là guest
}

func (InteractiveEvent) TableName() string {
	return "interactive_events"
}
