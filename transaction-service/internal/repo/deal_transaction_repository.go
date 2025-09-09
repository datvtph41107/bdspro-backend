package repo

import (
	"context"
	"tx/internal/domain/entity"
	"tx/internal/dto"
	"tx/internal/enums"
)

type DealTransactionRepository interface {
	Create(ctx context.Context, transaction *entity.ContractDeal) error
	Update(ctx context.Context, transaction *entity.ContractDeal) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.ContractDeal, error)
	GetByIds(ctx context.Context, ids []uint64) ([]*entity.ContractDeal, error)
	GetByDealID(ctx context.Context, dealID uint64, transactionType *enums.TransactionType, page, size int) ([]*entity.ContractDeal, int64, error)
	Approve(ctx context.Context, id uint64, approvedBy uint64, status enums.ApprovedStatus, rejectReason string) error
	GetTotalAmountByType(ctx context.Context, dealID uint64, transactionType enums.TransactionType) (float64, error)
	GetListByDealID(ctx context.Context, dealID uint64, req dto.DealSearch) ([]*dto.ContractDealListResponse, int64, error)
	GetTotalAmountByDealID(ctx context.Context, dealID uint64) (*dto.DealCostStatisticsDTO, error)
	GetDealContractProcess(ctx context.Context, dealID uint64) ([]dto.ProcessDTO, error)
}
