package repository

import (
	"context"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
)

type WithdrawalRequestsRepository interface {
	CreateWithdrawalRequest(ctx context.Context, request *entity.WithdrawalRequest) error
	GetWithdrawalRequestById(ctx context.Context, id uint32) (*entity.WithdrawalRequest, error)
	UpdateWithdrawalRequest(ctx context.Context, request *entity.WithdrawalRequest) error
}