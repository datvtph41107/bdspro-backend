package domain

import (
	_models "common/models"
	"strings"
)

// type PrivateFieldEnum = [""]

type ProductPrivate struct {
	_models.BaseEntity
	ProductId     uint64   `json:"productId,omitempty"`
	ImportPrice   *float64 `json:"importPrice,omitempty"`                             // Giá nhập/giá chủ nhà muốn net - CHỈ NỘI BỘ
	OperatingCost *float64 `gorm:"type:DECIMAL(15,0)" json:"operatingCost,omitempty"` // Chi phí vận hành (sửa chữa, marketing...)
	InternalNote  string   `json:"internalNote,omitempty"`                            // Ghi chú nội bộ (thông tin nhạy cảm)
	TargetProfit  *float64 `json:"targetProfit,omitempty"`                            // Lợi nhuận mục tiêu
	PrivateDocs   string   `gorm:"type:text" json:"privateDocs,omitempty"`            // Danh sách tài liệu riêng tư (lưu dạng CSV)
}

// Chuyển đổi khi lưu vào DB
func (p *ProductPrivate) SetPrivateDocs(docs []string) {
	p.PrivateDocs = strings.Join(docs, ",")
}

// Chuyển đổi khi đọc từ DB
func (p *ProductPrivate) GetPrivateDocs() []string {
	if p.PrivateDocs == "" {
		return []string{}
	}
	return strings.Split(p.PrivateDocs, ",")
}

func (ProductPrivate) TableName() string {
	return "product_private"
}
