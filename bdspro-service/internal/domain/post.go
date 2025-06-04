package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type Post struct {
	_models.BaseEntity
	ID              uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId       uint64            `json:"productId"`
	ExpiredAt       *time.Time        `json:"expiredAt"`
	Status          enums.EPostStatus `json:"status"`
	Visibility      enums.EVisibility `json:"visibility"`
	Hidden          bool              `gorm:"default:false" json:"hidden"`
	Content         string            `json:"content"`
	TransactionType uint              `json:"transactionType"`
	Title           string            `json:"title"`
	NumDate         int               `json:"numDate"`
	Type            uint              `json:"type"` // 1: normal, 2: feed, 3: share
	Price           float64           `json:"postPrice"`
	PriceType       uint              `json:"priceType"`      // 1: chính thức, 2: thỏa thuận
	PackageVisible  uint              `json:"packageVisible"` // 1: thường, 2: VIP
	Like            uint64            `json:"numLike"`
	Comment         uint64            `json:"numComment"`
	NumView         uint64            `json:"numView"`
	Product         *Product          `json:"product"`
	MediaList       []PostMediaEntity `json:"mediaList"`
	OwnerID         uint64            `json:"ownerId"`
	OwnerType       enums.EOwnerType  `gorm:"column:owner_type;default:10" json:"ownerType"`
	// Asset           *Asset               `json:"asset"`
}

type PostItem struct {
	_models.BaseEntity
	ID              uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId       uint64            `json:"productId"`
	ExpiredAt       *time.Time        `json:"expiredAt"`
	Status          enums.EPostStatus `json:"status"`
	Visibility      enums.EVisibility `json:"visibility"`
	Hidden          bool              `json:"hidden"`
	Content         string            `json:"content"`
	TransactionType uint              `json:"transactionType"`
	Title           string            `json:"title"`
	Type            uint              `json:"type"` // 1: normal, 2: feed, 3: share
	Price           uint              `json:"postPrice"`
	PriceType       uint              `json:"priceType"`      // 1: chính thức, 2: thỏa thuận
	PackageVisible  uint              `json:"packageVisible"` // 1: thường, 2: VIP
	Like            uint64            `json:"numLike"`
	Comment         uint64            `json:"numComment"`
	NumView         uint64            `json:"numView"`
	Address         uint64            `gorm:"column:address" json:"address"`
	ProvinceName    string            `gorm:"column:province_name" json:"provinceName"`
	DistrictName    string            `gorm:"column:district_name" json:"districtName"`
	WardName        string            `gorm:"column:ward_name" json:"wardName"`
}

func (PostItem) TableName() string {
	return "posts"
}
