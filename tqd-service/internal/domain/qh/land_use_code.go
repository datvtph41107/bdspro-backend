package qh_domain

import "time"

// LandUseCode - Mã loại đất (Master Data) standard cho toàn hệ thống, tham chiếu bởi các label cụ thể (ví dụ: ODT, DAT_O_DT_THI, Đất ở đô thị đều trỏ về ODT)
// Lý do: Mỗi loại đất chỉ 1 định nghĩa, các label tham chiếu đến
// Dẫn chứng: ODT chỉ xuất hiện 1 lần, DAT_O_DT_THI và Đất ở đô thị đều trỏ về ODT
type LandUseCode struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code           string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	GroupID        uint64    `gorm:"not null;index" json:"groupId"`
	CanBuild       bool      `gorm:"default:false" json:"canBuild"`
	Priority       int       `gorm:"default:50" json:"priority"`
	BuildCondition string    `gorm:"type:text" json:"buildCondition"`
	IsActive       bool      `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdatedBy      string    `gorm:"type:varchar(100)" json:"updatedBy"`

	Group *LandUseGroup `gorm:"-" json:"group,omitempty"`
}

func (LandUseCode) TableName() string {
	return "land_use_codes"
}
