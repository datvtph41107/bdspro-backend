package entity

import (
	"time"

	"gorm.io/gorm"
)

// Deal maps to the deals table.
type Deal struct {
	ID                      uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt               *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt               *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt               gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy               *uint64        `gorm:"column:created_by"`
	UpdatedBy               *uint64        `gorm:"column:updated_by"`
	Name                    string         `gorm:"column:name;size:255"`
	TargetProfit            float64        `gorm:"column:target_profit;default:0"`
	Status                  uint32         `gorm:"column:status;default:10"`
	DealType                uint32         `gorm:"column:deal_type;default:10"`
	Note                    string         `gorm:"column:note;type:text"`
	ChargePersonID          *uint64        `gorm:"column:charge_person_id"`
	CancelReason            *string        `gorm:"column:cancel_reason;type:text"`
	IsUnilateral            bool           `gorm:"column:is_unilateral;default:false"`
	IsManual                bool           `gorm:"column:is_manual;default:false"`
	BankAccountID           *uint64        `gorm:"column:bank_account_id"`
	FromDate                *time.Time     `gorm:"column:from_date"`
	ToDate                  *time.Time     `gorm:"column:to_date"`
	AllowSharing            bool           `gorm:"column:allow_sharing;default:false"`
	MemberCanAddTransaction bool           `gorm:"column:member_can_add_transaction;default:false"`
	OnlyOwnerGetCommission  bool           `gorm:"column:only_owner_get_commission;default:false"`
	InternalNote            string         `gorm:"column:internal_note;type:text"`
	AllowManualInput        bool           `gorm:"column:allow_manual_input;default:false"`
	Description             string         `gorm:"column:description;type:text"`
	ShareVisibility         uint32         `gorm:"column:share_visibility;default:10"`

	// hiện tại đang trỏ vào UserProfile.ProfileID
	// mục tiêu là đổi sang OriginProfile.OriginID
	OwnerID uint64 `gorm:"column:owner_id;index"`
	// deprecated: sau có thể bỏ do thông qua OwnerID -> OriginProfile.OriginId là xác định được
	// ownerType rồi
	OwnerType uint32 `gorm:"column:owner_type;default:10"`

	Members []DealMember `gorm:"foreignKey:DealID;references:ID"`
}

func (Deal) TableName() string { return "deals" }
