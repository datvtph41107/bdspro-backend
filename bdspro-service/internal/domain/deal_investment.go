package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type DealInvestment struct {
	_models.BaseEntity
	DealID             uint64 `gorm:"not null;index"`
	MemberID           uint64
	InvestType         enums.DealMemberType `gorm:"not null;default:10"`
	Amount             float64              `gorm:"not null"`
	TransferTime       *time.Time           ``
	Note               string               `gorm:"not null"`
	TransferProofImage string               `gorm:"column:transfer_proof_image"` // URL to transfer proof image
	AttachDocument     []AttachDocument     `gorm:"-"`

	Status    enums.TxApprovedStatus `gorm:"not null;default:10"`
	ChangedAt *time.Time             ``
	Confirmed bool                   `gorm:"not null;default:false"` // For tracking confirmation status

	// For proxy investments (when admin creates on behalf of user)
	ProxyUserID    *uint64 `gorm:"index"` // ID of admin/manager who created on behalf
	DoneInvestment bool    `gorm:"-"`     // Đánh dấu đã đủ vốn góp
}

func (DealInvestment) TableName() string {
	return "deal_investments"
}
