package entity

import (
	_models "common/models"
	"organization/internal/enums"
	"time"
)

type Investment struct {
	_models.BaseEntity
	DealID             uint64 `gorm:"not null;index"`
	MemberID           uint64
	InvestType         enums.DealMemberType `gorm:"not null;default:10"`
	Amount             float64              `gorm:"not null"`
	TransferTime       *time.Time           ``
	Note               string               `gorm:"not null"`
	TransferProofImage string               `gorm:"column:transfer_proof_image"` // URL to transfer proof image
	AttachDocument     []AttachDocument     `gorm:"-"`

	Status    enums.ApprovedStatus `gorm:"not null;default:10"`
	ChangedAt *time.Time           ``
	Confirmed bool                 `gorm:"not null;default:false"` // For tracking confirmation status

	// For proxy investments (when admin creates on behalf of user)
	ProxyUserID    *uint64 `gorm:"index"` // ID of admin/manager who created on behalf
	DoneInvestment bool    `gorm:"-"`     // Đánh dấu đã đủ vốn góp
}

func (Investment) TableName() string {
	return "investments"
}
