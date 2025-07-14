package entity

import (
	_models "common/models"
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
}

func (d *Deal) TableName() string {
	return "group_deals"
}
