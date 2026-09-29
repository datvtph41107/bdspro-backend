package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
)

type InvestmentRepository interface {
	CreateInvestment(ctx context.Context, investment *entity.Investment) (*entity.Investment, error)
	UpdateInvestment(ctx context.Context, investment *entity.Investment) (*entity.Investment, error)
	DeleteInvestment(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.Investment, error)
	GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]entity.Investment, int64, error)
	ChangeStatus(ctx context.Context, id uint64, status enums.ApprovedStatus) (uint64, error)
	ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) (uint64, error)
	CountDashboard(ctx context.Context, summary *dto.SummaryRequest) (*dto.SummaryResponse, error)
	GetConfirmedInvestmentsByDealId(ctx context.Context, dealId uint64, memberId uint64) ([]entity.Investment, error)
}
