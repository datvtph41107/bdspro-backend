package usecases

import (
	"context"
	"time"

	pb_bdspro "pb/types/bdspro"
)

type DashboardUsecase struct {
}

func NewDashboardUsecase() *DashboardUsecase {
	return &DashboardUsecase{}
}

func (u *DashboardUsecase) GetCashProfit(ctx context.Context, fromDate *time.Time, toDate *time.Time) ([]*pb_bdspro.CashProfit, error) {
	return nil, nil
}
