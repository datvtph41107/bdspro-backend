package domain

import "time"

type ProductNote struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autuUpdateTime" json:"updatedAt"`
	PinnedAt  *time.Time `gorm:"column:pinned_at" json:"pinnedAt"`
	ProductID uint64     `gorm:"not null;index"`
	AuthorID  uint64     `gorm:"not null;index"`
	Content   string     `gorm:"type:text;not null"`
}

func (ProductNote) TableName() string {
	return "product_notes"
}
