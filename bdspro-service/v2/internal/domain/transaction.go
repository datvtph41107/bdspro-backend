package domain

import (
	"time"

	"gorm.io/gorm"
)

// Bảng tỉnh/thành phố
type Transaction struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      int            `gorm:"unique;not null" json:"code"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
