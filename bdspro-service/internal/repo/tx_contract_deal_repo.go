package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
)

type TxContractDealRepository interface {
	Create(ctx context.Context, transaction *domain.TxContractDeal) error
	Update(ctx context.Context, transaction *domain.TxContractDeal) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.TxContractDeal, error)
	GetByIds(ctx context.Context, ids []uint64) ([]*domain.TxContractDeal, error)
	GetByDealID(ctx context.Context, dealID uint64, transactionType *enums.TxTransactionType, page, size int) ([]*domain.TxContractDeal, int64, error)
	Approve(ctx context.Context, id uint64, approvedBy uint64, status enums.TxApprovedStatus, rejectReason string) error
	GetTotalAmountByType(ctx context.Context, dealID uint64, transactionType enums.TxTransactionType) (float64, error)
	GetListByDealID(ctx context.Context, dealID uint64, req dto.TxDealSearch) ([]*dto.TxContractDealListResponse, int64, error)
	GetListByProfileID(ctx context.Context, profileID uint64, req dto.TxDealSearch) ([]*dto.TxContractDealListResponse, int64, error)
	GetListByTransactionIds(ctx context.Context, transactionIds []uint64) ([]*dto.TxContractDealListResponse, error)
	GetTotalAmountByDealID(ctx context.Context, dealID uint64) (*dto.TxDealCostStatisticsDTO, error)
	GetDealContractProcess(ctx context.Context, dealID uint64) ([]dto.TxProcessDTO, error)
}
