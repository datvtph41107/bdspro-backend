package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
)

type InvestmentRepository interface {
	CreateInvestment(ctx context.Context, investment *domain.DealInvestment) (*domain.DealInvestment, error)
	UpdateInvestment(ctx context.Context, investment *domain.DealInvestment) (*domain.DealInvestment, error)
	DeleteInvestment(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.DealInvestment, error)
	GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]domain.DealInvestment, int64, error)
	ChangeStatus(ctx context.Context, id uint64, status enums.TxApprovedStatus) (uint64, error)
	ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) (uint64, error)
	CountDashboard(ctx context.Context, summary *dto.SummaryRequest) (*dto.SummaryResponse, error)
	GetConfirmedInvestmentsByDealId(ctx context.Context, dealId uint64, memberId uint64) ([]domain.DealInvestment, error)
}
