package entity

import (
	"time"

	"gorm.io/gorm"
)

// DealInvestment maps to deal_investments.
type DealInvestment struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	CreatedBy          *uint64        `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy          *uint64        `gorm:"column:updated_by" json:"updatedBy"`
	DealID             uint64         `gorm:"column:deal_id;index;not null" json:"dealId"`
	MemberID           uint64         `gorm:"column:member_id;index" json:"memberId"`
	InvestType         uint32         `gorm:"column:invest_type;default:10"`
	Amount             float64        `gorm:"column:amount;not null;default:0" json:"amount"`
	TransferTime       *time.Time     `gorm:"column:transfer_time" json:"transferTime"`
	Note               string         `gorm:"column:note;type:text" json:"note"`
	TransferProofImage string         `gorm:"column:transfer_proof_image;type:text" json:"transferProofImage"`
	Status             uint32         `gorm:"column:status;default:10" json:"status"`
	ChangedAt          *time.Time     `gorm:"column:changed_at" json:"changedAt"`
	Confirmed          bool           `gorm:"column:confirmed;default:false" json:"confirmed"`
	ProxyUserID        *uint64        `gorm:"column:proxy_user_id" json:"proxyUserId"`
}

func (DealInvestment) TableName() string { return "deal_investments" }

// TxCostType maps to tx_cost_types.
type TxCostType struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt   *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	TypeName    string         `gorm:"column:type_name"`
	IsCustom    bool           `gorm:"column:is_custom;default:false"`
	CostType    uint32         `gorm:"column:cost_type;default:10"`
	Description string         `gorm:"column:description;type:text"`
}

func (TxCostType) TableName() string { return "tx_cost_types" }

// TxDealCost maps to tx_deal_costs.
type TxDealCost struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt   *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CostName    string         `gorm:"column:cost_name"`
	DealID      uint64         `gorm:"column:deal_id;index;not null"`
	Amount      float64        `gorm:"column:amount;default:0"`
	ConfirmedAt *time.Time     `gorm:"column:confirmed_at"`
	CostTypeID  uint64         `gorm:"column:cost_type_id;index"`
	Status      uint32         `gorm:"column:status;default:10"`
	Note        string         `gorm:"column:note;type:text"`
	CostType    *TxCostType    `gorm:"foreignKey:CostTypeID;references:ID"`
}

func (TxDealCost) TableName() string { return "tx_deal_costs" }
