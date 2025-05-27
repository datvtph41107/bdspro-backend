package repo

import (
	"bdspro/internal/domain"
	"context"
)

type DepositeRepo interface {
	GetDepositeByID(ctx context.Context, id uint64) (*domain.Deposite, error)
}
