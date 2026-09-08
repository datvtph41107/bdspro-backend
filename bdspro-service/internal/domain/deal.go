package domain

import (
	_models "common/models"
	"time"

	"bdspro/internal/enums"
)

type Deal struct {
	_models.BaseEntity

	OwnerId   uint64
	OwnerType enums.EOwnerOf

	// GroupId        uint64
	Name           string
	TargetProfit   float64          `gorm:"column:target_profit;default:0"`
	Status         enums.DealStatus `gorm:"column:status;default:10"`
	DealType       enums.DealType   `gorm:"column:deal_type;default:10"`
	Note           string           `gorm:"column:note;type:text"`
	NumMember      uint32           `gorm:"-"`
	ChargePersonID *uint64
	CancelReason   *string          `gorm:"type:text"`     // Lý do hủy thương vụ
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
	FromDate                *time.Time
	ToDate                  *time.Time
	AllowSharing            bool
	MemberCanAddTransaction bool
	OnlyOwnerGetCommission  bool
	InternalNote            string
	AllowManualInput        bool
	Description             string            `gorm:"column:description;type:text"`
	ShareVisibility         enums.EVisibility `gorm:"column:share_visibility;default:10"` // Phạm vi chia sẻ
}

func (d *Deal) TableName() string {
	return "deals"
}

func (d *Deal) IsModifiable() bool {
	switch d.Status {
	case enums.DealStatusCanceled, enums.DealStatusCompleted:
		return false
	default:
		return true
	}
}

func (d *Deal) IsOwner(userID uint64) bool {
	return d.OwnerId == userID
}

func (d *Deal) IsChargePerson(userID uint64) bool {
	return d.ChargePersonID != nil && *d.ChargePersonID == userID
}
