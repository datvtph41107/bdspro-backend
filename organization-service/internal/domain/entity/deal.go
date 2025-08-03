package entity

import (
	_models "common/models"
	"time"

	"organization/internal/enums"
)

type Deal struct {
	_models.BaseEntity

	OwnerId   uint64
	OwnerType enums.OwnerType

	// GroupId        uint64
	Name           string
	TargetProfit   float64          // kỳ vọng lợi nhuận
	Status         enums.DealStatus `gorm:"default:10"`
	DealType       enums.DealType
	Note           string
	NumMember      uint32 `gorm:"-"`
	ChargePersonID *uint64
	CancelReason   *string          // Lý do hủy thương vụ
	IsUnilateral   bool             `gorm:"default:false"` // Đơn phương (gỡ khỏi thương vụ)
	IsManual       bool             `gorm:"default:false"` // Nhập tay (dành cho doanh thu)
	Products       []DealProduct    `gorm:"foreignKey:DealID;references:ID"`
	ProductIds     []uint64         `gorm:"-"`
	Documents      []AttachDocument `gorm:"foreignKey:OwnerId;references:ID;-:migration"`
	BankAccountID  *uint64
	BankAccount    *BankAccount `gorm:"foreignKey:BankAccountID;references:ID"`
	Members        []DealMember `gorm:"foreignKey:DealID;references:ID"`
	Customers      []DealMember `gorm:"foreignKey:DealID;references:ID"`
	Partners       []DealMember `gorm:"foreignKey:DealID;references:ID"`
	MemberIds      []uint64     `gorm:"-"`
	CustomerIds    []uint64     `gorm:"-"`
	PartnerIds     []uint64     `gorm:"-"`

	// Setting
	FromDate                time.Time
	ToDate                  time.Time
	AllowSharing            bool
	MemberCanAddTransaction bool
	OnlyOwnerGetCommission  bool
	InternalNote            string
}

func (d *Deal) TableName() string {
	return "group_deals"
}
