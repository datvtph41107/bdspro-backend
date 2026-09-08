package entity

type DealProduct struct {
	DealID    uint64 `gorm:"primaryKey"`
	ProductID uint64 `gorm:"primaryKey"`
}

func (DealProduct) TableName() string {
	return "deal_products"
}
