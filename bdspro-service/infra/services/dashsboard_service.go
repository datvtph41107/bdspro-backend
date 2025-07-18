package service

import (
	"bdspro/internal/usecases"
	_utils "common/utils"
	"context"
	pb_bdspro "pb/types/bdspro"
)

type DashboardService struct {
	usecase *usecases.DashboardUsecase
}

func NewDashboardService(usecase *usecases.DashboardUsecase) *DashboardService {
	return &DashboardService{usecase: usecase}
}

// @Summary Get cash profit
// @Description Get cash profit
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param fromDate query string false "From date"
// @Param toDate query string false "To date"
// @Success 200 {object} pb_bdspro.GetCashProfitResponse
// @Router /dashboard/cash-profit/type [get]
func (s *DashboardService) GetCashProfit(ctx context.Context, req *pb_bdspro.GetCashProfitRequest) (*pb_bdspro.GetCashProfitResponse, error) {
	fromDate := _utils.ParseStringToTime(req.FromDate)
	toDate := _utils.ParseStringToTime(req.ToDate)

	_, err := s.usecase.GetCashProfit(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return &pb_bdspro.GetCashProfitResponse{Data: nil}, nil
}

// @Summary Get cash profit time
// @Description Get cash profit time
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param fromDate query string false "From date"
// @Param toDate query string false "To date"
// @Success 200 {object} pb_bdspro.GetCashProfitTimeResponse
// @Router /dashboard/cash-profit/time [get]
func (s *DashboardService) GetCashProfitTime(ctx context.Context, req *pb_bdspro.GetCashProfitRequest) (*pb_bdspro.GetCashProfitTimeResponse, error) {
	// fromDate := _utils.ParseStringToTime(req.FromDate)
	// toDate := _utils.ParseStringToTime(req.ToDate)

	// _, err := s.usecase.GetCashProfitTime(ctx, fromDate, toDate)
	// if err != nil {
	// 	return nil, err
	// }

	return &pb_bdspro.GetCashProfitTimeResponse{Data: nil}, nil
}
