package domain

import (
	"time"

	"gorm.io/gorm"
)

// HistoryEntity lưu lịch sử thông báo đã gửi
type HistoryFcmEntity struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	Body      string    `gorm:"type:text" json:"body"`
	Image     string    `gorm:"type:varchar(255)" json:"image"`
	Target    string    `gorm:"type:varchar(255)" json:"target"`
	Error     string    `gorm:"type:text" json:"error,omitempty"`
	Response  string    `gorm:"type:text" json:"response,omitempty"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

func (n *HistoryFcmEntity) TableName() string {
	return "tb_history"
}

// BeforeCreate tự động gán thời gian khi tạo bản ghi
func (h *HistoryFcmEntity) BeforeCreate(tx *gorm.DB) (err error) {
	h.Timestamp = time.Now()
	return nil
}
