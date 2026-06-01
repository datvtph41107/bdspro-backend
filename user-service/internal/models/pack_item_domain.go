package models

type PackItemDomain struct {
	ID                 uint64           `gorm:"primaryKey" json:"id"`
	Domain             string           `gorm:"not null" json:"domain"`
	Price              float64          `gorm:"" json:"price"`
	Currency           string           `gorm:"" json:"currency"`
	PriceTableDomainID uint64           `gorm:"not null" json:"priceTableDomainId"`
	PriceTableDomain   PriceTableDomain `gorm:"foreignKey:PriceTableDomainID;references:ID" json:"priceTableDomain"`
}

func (PackItemDomain) TableName() string {
	return "pack_items"
}
