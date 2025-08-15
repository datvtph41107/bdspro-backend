package dto

import (
	"organization/internal/enums"
	"time"
)

// Request DTOs
type CreateDealTransactionRequest struct {
	DealID          uint64                `json:"deal_id" validate:"required"`
	Type            enums.TransactionType `json:"type" validate:"required"`
	Amount          float64               `json:"amount" validate:"required,gt=0"`
	Description     string                `json:"description" validate:"required"`
	TransactionDate *time.Time            `json:"transaction_date" validate:"required"`
	Category        string                `json:"category" validate:"required"`
	Note            string                `json:"note"`
	IsManual        bool                  `json:"is_manual"` // Đánh dấu doanh thu được nhập tay
	DocumentIds     []uint64              `json:"document_ids"`
	TransactionType enums.TransactionType `json:"transaction_type"`
}

type UpdateDealTransactionRequest struct {
	ID              uint64     `json:"id" validate:"required"`
	Amount          float64    `json:"amount" validate:"required,gt=0"`
	Description     string     `json:"description" validate:"required"`
	TransactionDate *time.Time `json:"transaction_date" validate:"required"`
	Category        string     `json:"category" validate:"required"`
	Note            string     `json:"note"`
	IsManual        bool       `json:"is_manual"` // Đánh dấu doanh thu được nhập tay
	DocumentIds     []uint64   `json:"document_ids"`
}

type ApproveDealTransactionRequest struct {
	ID           uint64 `json:"id" validate:"required"`
	RejectReason string `json:"reject_reason"` // Chỉ cần khi reject
}

type GetDealTransactionsRequest struct {
	DealID uint64                 `json:"deal_id" validate:"required"`
	Type   *enums.TransactionType `json:"type"` // Optional filter
	Page   int                    `json:"page" validate:"min=1"`
	Size   int                    `json:"size" validate:"min=1,max=100"`
}

// Response DTOs
type DealTransactionResponse struct {
	ID              uint64                `json:"id"`
	DealID          uint64                `json:"deal_id"`
	Type            enums.TransactionType `json:"type"`
	TypeName        string                `json:"type_name"`
	Amount          float64               `json:"amount"`
	Description     string                `json:"description"`
	TransactionDate *time.Time            `json:"transaction_date"`
	Category        string                `json:"category"`
	Note            string                `json:"note"`
	IsManual        bool                  `json:"is_manual"` // Đánh dấu doanh thu được nhập tay
	CreatedBy       uint64                `json:"created_by"`
	Status          enums.ApprovedStatus  `json:"status"`
	StatusName      string                `json:"status_name"`
	ApprovedBy      *uint64               `json:"approved_by"`
	ApprovedAt      *time.Time            `json:"approved_at"`
	RejectReason    string                `json:"reject_reason"`
	DocumentIds     []uint64              `json:"document_ids"`
	CreatedAt       *time.Time            `json:"created_at"`
	UpdatedAt       *time.Time            `json:"updated_at"`
	CustomerID      uint64                `json:"customer_id"`
	CustomerName    string                `json:"customer_name"`
	TransactionID   *uint64               `json:"transaction_id"`
}

type DealTransactionListResponse struct {
	Transactions []DealTransactionResponse `json:"transactions"`
	Total        int64                     `json:"total"`
	Page         int                       `json:"page"`
	Limit        int                       `json:"limit"`
}

// Statistics DTOs
type DealFinancialSummary struct {
	DealID          uint64   `json:"deal_id"`
	TotalRevenue    float64  `json:"total_revenue"`
	TotalExpense    float64  `json:"total_expense"`
	NetProfit       float64  `json:"net_profit"`
	ProfitMargin    float64  `json:"profit_margin"` // Tỷ lệ lợi nhuận (%)
	TargetProfit    float64  `json:"target_profit"`
	ProfitStatus    string   `json:"profit_status"` // "on_track", "below_target", "above_target"
	WarningMessages []string `json:"warning_messages"`
}

type SignContractRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

type TransferContractRequest struct {
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

type DealContractResponse struct {
	TransactionID uint64 `json:"transaction_id"`
}
