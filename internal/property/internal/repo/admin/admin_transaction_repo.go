package admin_repo

import (
	"bdspro/internal/domain"
	"context"
)

// @bind: infra.postgres.admin.AdminTransactionPostgres
type IAdminTransactionRepo interface {
	GetList(ctx context.Context, filter *TransactionFilter) ([]*domain.Transaction, int64, error)
	GetDetail(ctx context.Context, id uint64) (*domain.Transaction, error)
	UpdateApprovalStatus(ctx context.Context, id uint64, status string, approvedBy *uint32) error
	Delete(ctx context.Context, id uint64) error
}

type TransactionFilter struct {
	Page              int
	Size              int
	Keyword           *string
	TransactionType   *uint32
	TransactionStatus *uint32
	ApprovalStatus    *string
	StartDate         *string
	EndDate           *string
	ProductId         *uint64
	ContactId         *uint64
	OwnerId           *uint64
	OwnerType         *uint32
}

