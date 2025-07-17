package entity

import (
	_models "common/models"
	"organization/internal/enums"
	"time"
)

type DealTransaction struct {
	_models.BaseEntity
	DealID         uint64               `gorm:"not null;index"`
	Type           enums.TransactionType `gorm:"not null"` // Chi phí hoặc doanh thu
	Amount         float64              `gorm:"not null"`
	Description    string               `gorm:"not null"`
	TransactionDate *time.Time          `gorm:"not null"`
	Category       string               `gorm:"not null"` // Danh mục chi phí/doanh thu
	Note           string               `gorm:"type:text"`
	IsManual       bool                 `gorm:"default:false"` // Đánh dấu doanh thu được nhập tay
	
	// Người tạo giao dịch
	CreatedBy      uint64               `gorm:"not null"`
	
	// Trạng thái phê duyệt
	Status         enums.ApprovedStatus `gorm:"not null;default:10"`
	ApprovedBy     *uint64              `gorm:"index"`
	ApprovedAt     *time.Time
	RejectReason   string               `gorm:"type:text"`
	
	// Tài liệu đính kèm
	AttachDocuments []AttachDocument    `gorm:"foreignKey:OwnerId;references:ID"`
	
	// Thông tin deal (relationship)
	Deal           *Deal                `gorm:"foreignKey:DealID;references:ID"`
}

func (DealTransaction) TableName() string {
	return "deal_transactions"
} 