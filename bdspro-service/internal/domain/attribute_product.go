package domain

import (
	"time"

	"gorm.io/gorm"
)

// Bảng liên kết sản phẩm và thuộc tính
type AttributeProduct struct {
	ID          uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint64           `gorm:"not null" json:"product_id"`
	Product     *Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	AttributeID uint64           `gorm:"not null" json:"attribute_id"`
	Attribute   *Attribute       `gorm:"foreignKey:AttributeID" json:"attribute,omitempty"`
	Values      []AttributeValue `gorm:"many2many:product_atb_val;" json:"values"`
	Name        string           `gorm:"size:255" json:"name"`
	Value       string           `json:"value"`
	Deleted     bool             `gorm:"-" json:"-"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

// func (AttributeProduct) RTable {

// }
