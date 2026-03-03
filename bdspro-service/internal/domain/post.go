package domain

import (
	_models "common/models"
	"encoding/json"
	"time"

	"bdspro/internal/enums"
)

type Post struct {
	_models.BaseEntity
	ID              uint64                `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductId       uint64                `json:"productId"`
	ExpiredAt       *time.Time            `json:"expiredAt"`
	Visibility      enums.EVisibility     `json:"visibility"`
	Hidden          bool                  `gorm:"default:false" json:"hidden"`
	Content         string                `json:"content"`
	TransactionType enums.TransactionType `json:"transactionType"`
	Title           string                `json:"title"`
	NumDate         int                   `json:"numDate"`
	Price           float64               `json:"postPrice"`
	PriceType       uint32                `json:"priceType"`      // 1: chính thức, 2: thỏa thuận
	PackageVisible  uint                  `json:"packageVisible"` // 1: thường, 2: VIP
	Like            uint64                `json:"numLike"`
	Comment         uint64                `json:"numComment"`
	NumView         uint64                `json:"numView"`
	Product         *Product              `json:"product"`
	MediaList       []PostMediaEntity     `json:"mediaList"`
	OwnerID         uint64                `json:"ownerId"`
	OwnerOf         enums.EOwnerOf        `gorm:"column:owner_of;default:10" json:"ownerOf"`
	PublishedAt     *time.Time            `json:"publishedAt"`
	PublishedBy     uint64                `json:"publishedBy"`
	Code            string                `json:"code"`

	Status            enums.EPostStatus        `gorm:"default:20" json:"status"`   // Trạng thái cơ bản
	TransactionStatus enums.EPostVisibleStatus `gorm:"-" json:"transactionStatus"` // Trạng thái giao dịch
	VisibleStatus     enums.EPostVisibleStatus `gorm:"-" json:"visibleStatus"`
	// Asset           *Asset               `json:"asset"`
}

type PostItem struct {
	Post
	PostPrice         float64                 `gorm:"column:post_price" json:"postPrice"`
	Address           uint64                  `gorm:"column:address" json:"address"`
	ProvinceId        *uint64                 `gorm:"column:province_id" json:"provinceId"`
	DistrictId        *uint64                 `gorm:"column:district_id" json:"districtId"`
	WardId            *uint64                 `gorm:"column:ward_id" json:"wardId"`
	ProvinceName      string                  `gorm:"column:province_name" json:"provinceName"`
	DistrictName      string                  `gorm:"column:district_name" json:"districtName"`
	WardName          string                  `gorm:"column:ward_name" json:"wardName"`
	MediaListRaw      json.RawMessage         `gorm:"column:media_list_raw" json:"-"`
	Status            enums.EPostStatus       `gorm:"column:status" json:"status"`
	TransactionStatus enums.TransactionStatus `gorm:"column:transaction_status" json:"transactionStatus"`

	VisibleStatus         enums.EPostVisibleStatus `gorm:"-" json:"visibleStatus"`
	TransactionTypeName   string                   `gorm:"-" json:"transactionTypeName"`
	TransactionStatusName string                   `gorm:"-" json:"transactionStatusName"`
	VisibleStatusName     string                   `gorm:"-" json:"visibleStatusName"`

	Area float64 `gorm:"column:area" json:"area"`

	// is not column
	ProductUpdatedAt *time.Time `gorm:"column:product_updated_at;->;-:migration" json:"productUpdatedAt"`
}

func (PostItem) TableName() string {
	return "posts"
}

type CountPostByTime struct {
	Date  time.Time `json:"date"`
	Count uint32    `json:"count"`
}
