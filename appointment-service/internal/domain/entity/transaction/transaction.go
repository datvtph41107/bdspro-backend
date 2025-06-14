package transaction

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	Id                  uint32           `json:"id" gorm:"type:uint32;primary_key"`
	OwnerId            uint32           `json:"owner_id" gorm:"type:uint32"`
	OrganizationId     *uint32          `json:"organization_id" gorm:"type:uint32"`
	Amount             float64          `json:"amount" gorm:"type:decimal(20,2)"`
	Currency           string           `json:"currency" gorm:"type:varchar(3)"`
	TransactionName    string           `json:"transaction_name" gorm:"type:varchar(255)"`
	Description        string           `json:"description" gorm:"type:text"`
	TransactionType    TransactionType  `json:"transaction_type" gorm:"type:varchar(20)"`
	CategoryId         uint32           `json:"category_id" gorm:"type:uint32"`
	PaymentMethodId    uint32           `json:"payment_method_id" gorm:"type:uint32"`
	RelatedDealId      *uint32          `json:"related_deal_id" gorm:"type:uint32"`
	RelatedTransactionId *uint32        `json:"related_transaction_id" gorm:"type:uint32"`
	ApprovalStatus     ApprovalStatus   `json:"approval_status" gorm:"type:varchar(20)"`
	ApprovedBy         *uint32          `json:"approved_by" gorm:"type:uint32"`
	ApprovalDate       *time.Time       `json:"approval_date" gorm:"type:datetime"`
	TransactionStatus  TransactionStatus `json:"transaction_status" gorm:"type:varchar(20)"`
	TransactionDate    time.Time        `json:"transaction_date" gorm:"type:date"`
	IsDeleted          bool             `json:"is_deleted" gorm:"type:boolean;default:false"`
	CreatedAt          time.Time        `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"type:timestamp"`
} 