package domain

import (
	"bdspro/enums"
	_models "common/models"
	"time"
)

type Post struct {
	_models.BaseEntity
	ID              uint64               `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId       uint64               `json:"productId"`
	ExpiredAt       *time.Time           `json:"expiredAt"`
	Status          enums.PostStatus     `json:"status"`
	Visibility      enums.PostVisibility `json:"visibility"`
	Hidden          bool                 `json:"hidden"`
	Content         string               `json:"content"`
	TransactionType uint                 `json:"transactionType"`
	Title           string               `json:"title"`
	Type            uint                 `json:"type"` // 1: normal, 2: feed, 3: share
	Price           uint                 `json:"postPrice"`
	PriceType       uint                 `json:"priceType"`      // 1: chính thức, 2: thỏa thuận
	PackageVisible  uint                 `json:"packageVisible"` // 1: thường, 2: VIP
	Like            uint64               `json:"numLike"`
	Comment         uint64               `json:"numComment"`
	NumView         uint64               `json:"numView"`
}
