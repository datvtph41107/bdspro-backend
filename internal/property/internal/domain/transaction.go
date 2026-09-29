package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

// Transaction represents a property transaction
type Transaction struct {
	_models.BaseEntity
	OwnerId              uint64                  `gorm:"type:int8;not null" json:"ownerId"`
	OwnerType            enums.EOwnerOf          `gorm:"type:smallint;not null;default:10" json:"ownerType"` // 10: Member, 20: Group, 30: Organization
	OrganizationId       *uint64                 `gorm:"type:int8" json:"organizationId,omitempty"`
	Amount               float64                 `gorm:"type:decimal(20,2);not null" json:"amount"`
	Currency             string                  `gorm:"type:varchar(10);default:'VND'" json:"currency"`
	TransactionName      string                  `gorm:"type:varchar(255);not null" json:"transactionName"`
	Description          string                  `gorm:"type:text" json:"description"`
	TransactionType      enums.TransactionType   `gorm:"type:smallint;not null;default:10" json:"transactionType"` // 10: Bán, 20: Thuê
	TransactionStatus    enums.TransactionStatus `gorm:"type:smallint;not null;default:10" json:"transactionStatus"`
	ApprovalStatus       enums.ApprovalStatus    `gorm:"type:varchar(20);default:'pending'" json:"approvalStatus"`
	TransactionDate      time.Time               `gorm:"type:timestamp" json:"transactionDate"`
	CategoryId           uint32                  `gorm:"type:int4" json:"categoryId"`
	PaymentMethodId      uint32                  `gorm:"type:int4" json:"paymentMethodId"`
	RelatedDealId        *uint32                 `gorm:"type:int4" json:"relatedDealId,omitempty"`
	RelatedTransactionId *uint32                 `gorm:"type:int4" json:"relatedTransactionId,omitempty"`
	ApprovedBy           *uint32                 `gorm:"type:int4" json:"approvedBy,omitempty"`
	ApprovalDate         *time.Time              `gorm:"type:timestamp" json:"approvalDate,omitempty"`
	IsDeleted            bool                    `gorm:"default:false" json:"isDeleted"`
	Status               uint32                  `gorm:"type:smallint;default:10" json:"status"`
	DepositeAmount       *float64                `gorm:"type:decimal(20,2)" json:"depositeAmount,omitempty"`
	DepositeNote         *string                 `gorm:"type:text" json:"depositeNote,omitempty"`
	ContactId            *uint64                 `gorm:"type:int8" json:"contactId,omitempty"`
	ProductId            *uint64                 `gorm:"type:int8" json:"productId,omitempty"`
	Product              *Product                `gorm:"foreignKey:ProductId;references:ID" json:"product,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

// Helper methods để get tên trạng thái
func (t *Transaction) GetTransactionTypeName() string {
	switch t.TransactionType {
	case enums.TransactionTypeSale:
		return "Bán"
	case enums.TransactionTypeRent:
		return "Thuê"
	default:
		return "Không xác định"
	}
}

func (t *Transaction) GetTransactionStatusName() string {
	switch t.TransactionStatus {
	case enums.TransactionStatusDraft:
		return "Nháp"
	case enums.TransactionStatusDeposite:
		return "Đặt cọc"
	case enums.TransactionStatusSigned:
		return "Đã ký"
	case enums.TransactionStatusCompleted:
		return "Hoàn thành"
	case enums.TransactionStatusCanceled:
		return "Hủy"
	default:
		return "Không xác định"
	}
}
