package dto

import (
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"time"
)

// Request DTOs
type TxDealContractPaymentRequest struct {
	Id              uint64     `json:"id,omitempty"`
	Amount          float64    `json:"amount,omitempty"`
	CustomerId      uint64     `json:"customerId,omitempty"`
	DealId          uint64     `json:"dealId,omitempty"`
	Note            string     `json:"note,omitempty"`
	ProductId       uint64     `json:"productId,omitempty"`
	Timestamp       *time.Time `json:"timestamp,omitempty"`
	TransactionStep uint32     `json:"transactionStep,omitempty"`
	TransactionType uint32     `json:"transactionType,omitempty"`
}

type TxUpdateDealTransactionRequest struct {
	ID              uint64     `json:"id" validate:"required"`
	Amount          float64    `json:"amount" validate:"required,gt=0"`
	Description     string     `json:"description" validate:"required"`
	TransactionDate *time.Time `json:"transaction_date" validate:"required"`
	Category        string     `json:"category" validate:"required"`
	Note            string     `json:"note"`
	IsManual        bool       `json:"is_manual"` // Đánh dấu doanh thu được nhập tay
	DocumentIds     []uint64   `json:"document_ids"`
}

type TxApproveDealTransactionRequest struct {
	ID           uint64 `json:"id" validate:"required"`
	RejectReason string `json:"reject_reason"` // Chỉ cần khi reject
}

type TxGetDealTransactionsRequest struct {
	DealID uint64                   `json:"deal_id" validate:"required"`
	Type   *enums.TxTransactionType `json:"type"` // Optional filter
	Page   int                      `json:"page" validate:"min=1"`
	Size   int                      `json:"size" validate:"min=1,max=100"`
}

// Response DTOs
type TxDealTransactionResponse struct {
	ID              uint64                  `json:"id"`
	DealID          uint64                  `json:"deal_id"`
	Type            enums.TxTransactionType `json:"type"`
	TypeName        string                  `json:"type_name"`
	Amount          float64                 `json:"amount"`
	Description     string                  `json:"description"`
	TransactionDate *time.Time              `json:"transaction_date"`
	Category        string                  `json:"category"`
	Note            string                  `json:"note"`
	IsManual        bool                    `json:"is_manual"` // Đánh dấu doanh thu được nhập tay
	CreatedBy       uint64                  `json:"created_by"`
	Status          enums.TxApprovedStatus  `json:"status"`
	StatusName      string                  `json:"status_name"`
	ApprovedBy      *uint64                 `json:"approved_by"`
	ApprovedAt      *time.Time              `json:"approved_at"`
	RejectReason    string                  `json:"reject_reason"`
	DocumentIds     []uint64                `json:"document_ids"`
	CreatedAt       *time.Time              `json:"created_at"`
	UpdatedAt       *time.Time              `json:"updated_at"`
	CustomerID      uint64                  `json:"customer_id"`
	CustomerName    string                  `json:"customer_name"`
	TransactionID   *uint64                 `json:"transaction_id"`
}

type TxDealTransactionListResponse struct {
	Transactions []TxDealTransactionResponse `json:"transactions"`
	Total        int64                       `json:"total"`
	Page         int                         `json:"page"`
	Limit        int                         `json:"limit"`
}

// Statistics DTOs
type TxDealFinancialSummary struct {
	DealID          uint64   `json:"deal_id"`
	TotalRevenue    float64  `json:"total_revenue"`
	TotalExpense    float64  `json:"total_expense"`
	NetProfit       float64  `json:"net_profit"`
	ProfitMargin    float64  `json:"profit_margin"` // Tỷ lệ lợi nhuận (%)
	TargetProfit    float64  `json:"target_profit"`
	ProfitStatus    string   `json:"profit_status"` // "on_track", "below_target", "above_target"
	WarningMessages []string `json:"warning_messages"`
	TotalPayment    float64  `json:"total_payment"`
}

type TxSignContractRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

type TxTransferContractRequest struct {
	ID              uint64  `json:"id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	CustomerID      uint64  `json:"customer_id" validate:"required"`
	DealID          uint64  `json:"deal_id" validate:"required"`
	Note            string  `json:"note"`
	ProductIDs      uint64  `json:"product_ids" validate:"required"`
	Timestamp       string  `json:"timestamp" validate:"required"`
	TransactionStep uint32  `json:"transaction_step" validate:"required"`
	TransactionType uint32  `json:"transaction_type" validate:"required"`
}

type TxDealContractResponse struct {
	TransactionID uint64 `json:"transaction_id"`
}

type TxDealSearch struct {
	_dto.Pagable
}

type TxActionPaymentRequest struct {
	ID            uint64     `json:"id,omitempty"`
	Amount        float64    `json:"amount,omitempty"`
	FromId        uint64     `json:"fromId,omitempty"`
	DealId        uint64     `json:"dealId,omitempty"`
	Note          string     `json:"note,omitempty"`
	ProductId     uint64     `json:"productId,omitempty"`
	Timestamp     *time.Time `json:"timestamp,omitempty"`
	Action        uint32     `json:"action,omitempty"`
	TransactionId uint64     `json:"transactionId,omitempty"`
	ContractId    uint64     `json:"contractId,omitempty"`
}

type TxContractDealListResponse struct {
	ContractID      uint64     `gorm:"column:contract_id" json:"id"`
	DealID          uint64     `gorm:"column:deal_id" json:"deal_id"`
	TransactionName string     `gorm:"column:transaction_name" json:"transaction_name"`
	Method          uint       `gorm:"column:method" json:"method"`
	CustomerID      uint64     `gorm:"column:customer_id" json:"customer_id"`
	CustomerName    string     `gorm:"-" json:"customer_name"`
	Status          uint32     `gorm:"column:status" json:"status"`
	StatusName      string     `gorm:"-" json:"status_name"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	Timestamp       *time.Time `gorm:"column:timestamp" json:"timestamp"`
	Code            string     `gorm:"-" json:"code"`
}

// DealAmountResponse represents the total amount response for a deal
type TxDealAmountResponse struct {
	DealID     uint64  `json:"deal_id"`
	Payment    TxStats `json:"payment"`
	NumPayment TxStats `json:"num_payment"`
}

type TxStats struct {
	Pending  float64 `json:"pending"`
	Approved float64 `json:"approved"`
	Rejected float64 `json:"rejected"`
	Total    float64 `json:"total"`
}
